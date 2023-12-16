package lua_module

import (
	"github.com/go-resty/resty/v2"
	lua "github.com/yuin/gopher-lua"
	"os"
)

type TLuaHttpClient struct {
	client *resty.Client
	trace  bool
}

func (this *TLuaHttpClient) Request() *resty.Request {
	if this.trace {
		return this.client.R().EnableTrace()
	}
	return this.client.R()
}

/**
 * @Author: icy
 * @Description: http 请求返回的内容处理
 * @Param lua 状态 回调函数的索引 http 返回的内容
 * @return 字符串  如果不返回为空 侧将引向调用请求方法的返回内容
 * @Date: 12/15/23 3:29 PM
 */

func (this *TLuaHttpClient) callResponseFunc(L *lua.LState, cbkIndex int, res *resty.Response) string {
	cbk := L.ToFunction(cbkIndex)
	if cbk != nil {
		table := L.NewTable()
		table.RawSetString(`statusCode`, lua.LNumber(res.StatusCode()))
		table.RawSetString(`statusText`, lua.LString(res.Status()))
		table.RawSetString(`size`, lua.LNumber(res.Size()))

		THeader := L.NewTable()
		list := res.Header()
		if len(list) > 0 {
			for key, val := range list {
				THeader.RawSetString(key, lua.LString(val[0]))
			}
		}
		table.RawSetString(`header`, THeader)

		TCookie := L.NewTable()
		cookies := res.Cookies()
		if len(cookies) > 0 {
			for _, P := range cookies {
				TCookie.RawSetString(P.Name, lua.LString(P.Value))
			}
		}

		table.RawSetString(`cookies`, TCookie)
		table.RawSetString(`content`, lua.LString(res.Body()))
		table.RawSetString(`time`, lua.LNumber(res.Time()))

		if err := L.CallByParam(lua.P{
			Fn:      cbk,
			NRet:    0,
			Protect: true,
		}, table); err == nil {
			r := L.Get(-1)
			L.Pop(1)
			if r.Type() == lua.LTString {
				return lua.LVAsString(r)
			}
		}
	}

	return ``
}

func (this *TLuaHttpClient) Trace(L *lua.LState) int {
	this.trace = L.ToBool(1)
	return 0
}

func (this *TLuaHttpClient) Get(L *lua.LState) int {
	res := ``
	url := L.ToString(1)
	if url != "" {
		if result, err := this.Request().Get(url); err == nil {
			res = this.callResponseFunc(L, 2, result)
			if res == "" {
				res = result.String()
			}
		}
	}

	L.Push(lua.LString(res))

	return 1
}

func (this *TLuaHttpClient) Delete(L *lua.LState) int {
	res := ``
	url := L.ToString(1)
	if url != "" {
		if r, err := this.Request().Delete(url); err == nil {
			res = this.callResponseFunc(L, 2, r)
			if res == "" {
				res = r.String()
			}
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
			if r, err := this.Request().SetFormData(LuaTableToMapString(table)).Put(url); err == nil {
				res = this.callResponseFunc(L, 3, r)
				if res == "" {
					res = r.String()
				}
			}
		} else {
			if r, err := this.Request().SetBody([]byte(L.ToString(2))).Put(url); err == nil {
				res = this.callResponseFunc(L, 3, r)
				if res == "" {
					res = r.String()
				}
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
			if r, err := this.Request().SetFormData(LuaTableToMapString(table)).Post(url); err == nil {
				res = this.callResponseFunc(L, 3, r)
				if res == "" {
					res = r.String()
				}
			}
		} else {
			if r, err := this.Request().SetBody([]byte(L.ToString(2))).Post(url); err == nil {
				res = this.callResponseFunc(L, 3, r)
				if res == "" {
					res = r.String()
				}
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
				this.Request().SetHeader(k, v)
			}
		}
	} else {
		k := L.ToString(1)
		v := L.ToString(2)
		if k != "" && v != "" {
			this.Request().SetHeader(k, v)
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
		L.Push(lua.LString(this.Request().Header.Get(k)))
	} else {
		table := L.NewTable()
		for k, v := range this.Request().Header {
			table.RawSetString(k, lua.LString(v[0]))
		}
		L.Push(table)
	}

	return 1
}
func (this *TLuaHttpClient) Upload(L *lua.LState) int {
	completed := false
	url := L.ToString(1)
	table := L.ToTable(2)
	if url != "" && table != nil {
		if r, err := this.Request().SetFiles(LuaTableToMapString(table)).Post(url); err == nil {
			completed = true
			this.callResponseFunc(L, 3, r)
			L.Push(lua.LString(r.String()))
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
		if r, err := this.Request().Get(url); err == nil {
			if f, err := os.Create(fileName); err == nil {
				defer f.Close()
				_, err = f.Write(r.Body())
				completed = err == nil
				this.callResponseFunc(L, 3, r)
			}
		}
	}

	L.Push(lua.LBool(completed))
	return 1
}

func (this *TLuaHttpClient) Proxy(L *lua.LState) int {
	res := this.client.SetProxy(L.ToString(1)) == nil
	L.Push(lua.LBool(res))
	return 1
}

func HttpClientPreload(L *lua.LState) {
	L.PreloadModule("httpClient", func(L *lua.LState) int {
		HttpClient := &TLuaHttpClient{client: resty.New()}

		t := L.NewTable()
		L.SetFuncs(t, map[string]lua.LGFunction{`get`: HttpClient.Get,
			`put`:       HttpClient.Put,
			`post`:      HttpClient.Post,
			`delete`:    HttpClient.Delete,
			`upload`:    HttpClient.Upload,
			`download`:  HttpClient.Download,
			`setHeader`: HttpClient.SetHeader,
			`getHeader`: HttpClient.GetHeader,
			`proxy`:     HttpClient.Proxy,
			`trace`:     HttpClient.Trace,
		})

		L.Push(t)
		return 1
	})
}
