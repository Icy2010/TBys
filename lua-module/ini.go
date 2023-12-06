package lua_module

import (
	"fmt"
	"github.com/go-ini/ini"
	lua "github.com/yuin/gopher-lua"
)

type TLuaIni struct {
	file *ini.File
}

func (this *TLuaIni) Load(L *lua.LState) int {
	fileName := L.ToString(1)
	var err error = nil
	this.file, err = ini.Load(fileName)
	L.Push(lua.LBool(err == nil))
	return 1
}

func (this *TLuaIni) Keys(L *lua.LState) int {
	name := L.ToString(1)
	table := L.NewTable()
	if this.file != nil {
		if name != "" {
			sec := this.file.Section(name)
			for _, key := range sec.Keys() {
				table.Append(lua.LString(key.String()))
			}
		}
	}

	L.Push(table)
	return 1
}

func (this *TLuaIni) Sections(L *lua.LState) int {
	table := L.NewTable()
	if this.file != nil {
		list := this.file.Sections()
		if len(list) > 0 {
			for _, sec := range list {
				table.Append(lua.LString(sec.Name()))
			}
		}
	}

	L.Push(table)
	return 1
}

func (this *TLuaIni) getValue(name, key string) (*ini.Key, error) {
	if this.file != nil {
		if name != "" && key != "" {
			if sec, err := this.file.GetSection(name); err == nil {
				return sec.GetKey(key)
			}
		}
	}
	return nil, fmt.Errorf(`无法获取`)
}

func (this *TLuaIni) GetString(L *lua.LState) int {
	val := ``
	if key, err := this.getValue(L.ToString(1), L.ToString(2)); err == nil {
		val = key.Value()
	}

	L.Push(lua.LString(val))
	return 1
}

func (this *TLuaIni) GetInt(L *lua.LState) int {
	var val int64 = 0
	if key, err := this.getValue(L.ToString(1), L.ToString(2)); err == nil {
		val, err = key.Int64()
	}

	L.Push(lua.LNumber(val))
	return 1
}

func (this *TLuaIni) GetFloat(L *lua.LState) int {
	var val float64 = 0
	if key, err := this.getValue(L.ToString(1), L.ToString(2)); err == nil {
		val, err = key.Float64()
	}

	L.Push(lua.LNumber(val))
	return 1
}

func (this *TLuaIni) GetBool(L *lua.LState) int {
	var val bool = false
	if key, err := this.getValue(L.ToString(1), L.ToString(2)); err == nil {
		val, err = key.Bool()
	}

	L.Push(lua.LBool(val))
	return 1
}

func (this *TLuaIni) SetValue(L *lua.LState) int {
	if key, err := this.getValue(L.ToString(1), L.ToString(2)); err == nil {
		key.SetValue(L.ToString(3))
	}

	return 0
}

func (this *TLuaIni) Save(L *lua.LState) int {
	fileName := L.ToString(1)
	err := this.file.SaveTo(fileName)
	L.Push(lua.LBool(err == nil))
	return 1
}

/*--------------------------------------------------------------------------------------------------------------------*/

func IniPreload(L *lua.LState) {
	L.PreloadModule(`ini`, func(L *lua.LState) int {
		ini := &TLuaIni{}
		t := L.NewTable()
		L.SetFuncs(t, map[string]lua.LGFunction{
			`load`:      ini.Load,
			`save`:      ini.Save,
			`sections`:  ini.Sections,
			`keys`:      ini.Keys,
			`getString`: ini.GetString,
			`getInt`:    ini.GetInt,
			`getFloat`:  ini.GetFloat,
			`getBool`:   ini.GetBool,
			`setValue`:  ini.SetValue,
		})

		L.Push(t)
		return 1
	})
}
