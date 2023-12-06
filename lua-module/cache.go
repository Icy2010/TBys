package lua_module

import (
	"github.com/patrickmn/go-cache"
	lua "github.com/yuin/gopher-lua"
	"time"
)

type TLuaCache struct {
	ch *cache.Cache
}

func (this *TLuaCache) set(L *lua.LState) int {
	key := L.ToString(1)
	table := L.ToTable(2)
	if key != "" && table != nil {
		dur := L.ToInt64(3)
		if dur == 0 {
			dur = int64(time.Minute)
		}
		this.ch.Set(key, table, time.Duration(dur))
	}

	return 0
}

func (this *TLuaCache) get(L *lua.LState) int {
	key := L.ToString(1)
	if key != "" {
		if table, ok := this.ch.Get(key); ok {
			L.Push(table.(*lua.LTable))
			L.Push(lua.LTrue)
			return 2
		}
	}

	L.Push(lua.LFalse)
	return 1
}

func (this *TLuaCache) delete(L *lua.LState) int {
	key := L.ToString(1)
	if key != "" {
		this.ch.Delete(key)
	}

	this.ch.ItemCount()
	return 0
}

func (this *TLuaCache) count(L *lua.LState) int {
	c := this.ch.ItemCount()
	L.Push(lua.LNumber(c))
	return 1
}

/*--------------------------------------------------------------------------------------------------------------------*/

func CacheRegisterModule(L *lua.LState) {
	ca := &TLuaCache{}
	ca.ch = cache.New(60*time.Minute, 10*time.Minute)
	L.RegisterModule(`cache`, map[string]lua.LGFunction{
		`set`:    ca.set,
		`get`:    ca.get,
		`delete`: ca.delete,
		`count`:  ca.count,
	})
}
