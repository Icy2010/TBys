package lua_module

import (
	. "github.com/Icy2010/TBys"
	lua "github.com/yuin/gopher-lua"
	"os"
)

type TLuaHttpClient struct {
	http *THttp
}

func (this *TLuaHttpClient) Get(L *lua.LState) int {
	res := ``
	url := L.ToString(1)
	if url != "" {
		if buff, err := this.http.Get(url); err == nil {
			res = string(buff)
		}
	}

	L.Push(lua.LString(res))

	return 1
}

func (this *TLuaHttpClient) Delete(L *lua.LState) int {
	res := ``
	url := L.ToString(1)
	if url != "" {
		if buff, err := this.http.Delete(url); err == nil {
			res = string(buff)
		}
	}

	L.Push(lua.LString(res))

	return 1
}

func (this *TLuaHttpClient) Put(L *lua.LState) int {
	res := ``
	url := L.ToString(1)
	if url != "" {
		table := L.ToTable(2)
		if table != nil {
			if buff, err := this.http.Put(url, LuaTableToMap(table)); err == nil {
				res = string(buff)
			}
		} else {
			if buff, err := this.http.PutBuff(url, []byte(L.ToString(2))); err == nil {
				res = string(buff)
			}
		}
	}

	L.Push(lua.LString(res))
	return 1
}

func (this *TLuaHttpClient) Post(L *lua.LState) int {
	res := ``
	url := L.ToString(1)
	if url != "" {
		table := L.ToTable(2)
		if table != nil {
			if buff, err := this.http.Post(url, LuaTableToMap(table)); err == nil {
				res = string(buff)
			}
		} else {
			if buff, err := this.http.PostBuff(url, []byte(L.ToString(2))); err == nil {
				res = string(buff)
			}
		}
	}

	L.Push(lua.LString(res))

	return 1
}

/*
   设置可以直接全部一次性设置完成  传入表即可如果
   传入 key ， value 两个属值  就是设置单个
*/

func (this *TLuaHttpClient) SetHeader(L *lua.LState) int {
	table := L.ToTable(1)
	if table != nil {
		data := LuaTableToMapString(table)
		if len(data) > 0 {
			for k, v := range data {
				this.http.SetHeader(k, v)
			}
		}
	} else {
		k := L.ToString(1)
		v := L.ToString(2)
		if k != "" && v != "" {
			this.http.SetHeader(k, v)
		}
	}

	return 0
}

/*
传入 key 就返回单个
如果什么都没有传入 返回全部 (table)
*/

func (this *TLuaHttpClient) GetHeader(L *lua.LState) int {
	k := L.ToString(1)
	if k != "" {
		L.Push(lua.LString(this.http.GetHeader(k)))
	} else {
		table := L.NewTable()
		for k, v := range this.http.Header() {
			table.RawSetString(k, lua.LString(v))
		}
		L.Push(table)
	}

	return 1
}

func (this *TLuaHttpClient) HasHeader(L *lua.LState) int {
	has := false
	k := L.ToString(1)
	if k != "" {
		_, has = this.http.HasHeader(k)
	}

	L.Push(lua.LBool(has))
	return 1
}

func (this *TLuaHttpClient) GZip(L *lua.LState) int {
	if v, ok := L.Get(1).(lua.LBool); ok {
		this.http.GZip = bool(v)
	}

	L.Push(lua.LBool(this.http.GZip))
	return 1
}

func (this *TLuaHttpClient) Upload(L *lua.LState) int {
	completed := false

	url := L.ToString(1)
	fileName := L.ToString(2)
	if url != "" && PathExist(fileName) {
		table := L.ToTable(3)
		form := make(map[string]any)
		if table != nil {
			form = LuaTableToMap(table)
		}
		if buff, err := this.http.Upload(url, fileName, form); err == nil {
			completed = true
			L.Push(lua.LString(buff))
		} else {
			L.Push(lua.LNil)
		}

	} else {
		L.Push(lua.LNil)
	}

	L.Push(lua.LBool(completed))
	return 2
}

func (this *TLuaHttpClient) Download(L *lua.LState) int {
	completed := false

	url := L.ToString(1)
	fileName := L.ToString(2)
	if url != "" {
		if buff, err := this.http.Get(url); err == nil {
			if f, err := os.Create(fileName); err == nil {
				defer f.Close()
				_, err = f.Write(buff)

				completed = err == nil
			}
		}
	}

	L.Push(lua.LBool(completed))
	return 1
}

func (this *TLuaHttpClient) Proxy(L *lua.LState) int {
	res := this.http.Proxy(L.ToString(1)) == nil
	L.Push(lua.LBool(res))
	return 1
}

func HttpClientPreload(L *lua.LState) {
	L.PreloadModule("httpClient", func(L *lua.LState) int {
		HttpClient := &TLuaHttpClient{http: NewHTTP(nil)}

		t := L.NewTable()
		L.SetFuncs(t, map[string]lua.LGFunction{`get`: HttpClient.Get,
			`put`:       HttpClient.Put,
			`post`:      HttpClient.Post,
			`delete`:    HttpClient.Delete,
			`upload`:    HttpClient.Upload,
			`download`:  HttpClient.Download,
			`setHeader`: HttpClient.SetHeader,
			`getHeader`: HttpClient.GetHeader,
			`hasHeader`: HttpClient.HasHeader,
			`gZip`:      HttpClient.GZip,
			`proxy`:     HttpClient.Proxy,
		})

		L.Push(t)
		return 1
	})
}
