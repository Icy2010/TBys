package lua_module

import (
	"errors"
	htmlquery "github.com/antchfx/xquery/html"
	"github.com/dop251/goja"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

type TLuaHtmlParser struct {
	TBaseBysLua
}

func (this *TLuaHtmlParser) getScript(content string, js *goja.Runtime, onScript func(value string, size uint) (string, bool)) (string, bool) {
	if onScript == nil {
		return "", false
	}

	root, _ := htmlquery.Parse(strings.NewReader(content))
	list := htmlquery.Find(root, `//script`)
	if len(list) > 0 {
		for _, v := range list {
			text := htmlquery.InnerText(v)
			if value, next := onScript(text, uint(len(text))); next {
				if val, err := js.RunString(value); err == nil {
					return val.String(), true
				}
			}
		}
	}

	return "", false
}

func (this *TLuaHtmlParser) getContent(src, expr string, OnText func(val string) bool) error {
	if OnText != nil {
		return errors.New(`请传入回调函数`)
	}

	root, _ := htmlquery.Parse(strings.NewReader(src))
	list := htmlquery.Find(root, expr)
	if len(list) > 0 {
		for _, v := range list {
			text := htmlquery.InnerText(v)
			if OnText(text) {
				return nil
			}
		}
	} else {
		return errors.New(`未查找到内容`)
	}

	return nil
}

/*--------------------------------------------------------------------------------------------------------------------*/

func (this *TLuaHtmlParser) Find(L *lua.LState) int {
	ok := false
	src := L.ToString(1)
	expr := L.ToString(2)
	cbk := L.ToFunction(3)
	if src != "" && expr != "" && cbk != nil {
		err := this.getContent(src, expr, func(val string) bool {
			if err := L.CallByParam(lua.P{
				Fn:      cbk,
				NRet:    1,
				Protect: true,
			}, lua.LString(val)); err == nil {
				res := L.Get(-1)
				L.Pop(1)
				if res.Type() == lua.LTBool {
					if lua.LVAsBool(res) {
						return true
					}
				}
			}

			return false
		})

		ok = err == nil
	}

	L.Push(lua.LBool(ok))
	return 1
}

func (this *TLuaHtmlParser) ReadScript(L *lua.LState) int {
	has := false
	Script := ""
	content := L.ToString(1)
	doScript := L.ToFunction(2)
	if content != "" && doScript != nil {
		js := goja.New()
		Script, has = this.getScript(content, js, func(value string, size uint) (string, bool) {
			if err := L.CallByParam(lua.P{
				Fn:      doScript,
				NRet:    1,
				Protect: true,
			}, lua.LString(value), lua.LNumber(size)); err == nil {
				res := L.Get(-1)
				L.Pop(1)

				if res.Type() == lua.LTString {
					return res.String(), true
				}
			}

			return ``, false
		})
	} else {
		LuaLoggedError(this.Errorf(`必要参数缺少,HTML 内容  脚本回调`))
	}

	L.Push(lua.LBool(has))
	L.Push(lua.LString(Script))
	return 2
}

/*--------------------------------------------------------------------------------------------------------------------*/

func HtmlParserPreload(L *lua.LState) {
	L.PreloadModule("htmlParser", HtmlParserLoader)
}

func HtmlParserLoader(L *lua.LState) int {
	p := &TLuaHtmlParser{}
	t := L.NewTable()

	L.SetFuncs(t, map[string]lua.LGFunction{
		`readScript`: p.ReadScript,
		`find`:       p.Find,
	})

	L.Push(t)
	return 1
}
