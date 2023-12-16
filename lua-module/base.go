package lua_module

import (
	"crypto/rand"
	"fmt"
	u "github.com/Icy2010/TBys"
	lua "github.com/yuin/gopher-lua"
	"math/big"
	"os"
	"strings"
	"time"
)

func logger(val ...any) *u.TLogger {
	return u.Logger()
}

type TBaseBysLua struct {
}

func (TBaseBysLua) Errorf(format string, val ...any) error {
	return fmt.Errorf(format, val...)
}

func (TBaseBysLua) LuaTableValueString(table *lua.LTable, key string, value *string) {
	val := table.RawGetString(key)
	if val.Type() == lua.LTString {
		*value = lua.LVAsString(val)
	}
}

func (TBaseBysLua) LuaTableValueInt(table *lua.LTable, key string, value *int) {
	val := table.RawGetString(key)
	if val.Type() == lua.LTNumber {
		*value = int(lua.LVAsNumber(val))
	}
}

func (TBaseBysLua) LuaTableValueInt64(table *lua.LTable, key string, value *int64) {
	val := table.RawGetString(key)
	if val.Type() == lua.LTNumber {
		*value = int64(lua.LVAsNumber(val))
	}
}

func (TBaseBysLua) LuaTableValueFloat32(table *lua.LTable, key string, value *float32) {
	val := table.RawGetString(key)
	if val.Type() == lua.LTNumber {
		*value = float32(lua.LVAsNumber(val))
	}
}

func (TBaseBysLua) LuaTableValueFloat64(table *lua.LTable, key string, value *float64) {
	val := table.RawGetString(key)
	if val.Type() == lua.LTNumber {
		*value = float64(lua.LVAsNumber(val))
	}
}

func (TBaseBysLua) LuaTableValueBoolean(table *lua.LTable, key string, value *bool) {
	val := table.RawGetString(key)
	if val.Type() == lua.LTBool {
		*value = lua.LVAsBool(val)
	}
}

func (TBaseBysLua) Output() {

}

/*--------------------------------------------------------------------------------------------------------------------*/

type TBysLua struct {
	TBaseBysLua
}

func (TBysLua) PathExists(L *lua.LState) int {
	p := L.ToString(1)
	b := u.PathExist(p)
	L.Push(lua.LBool(b))
	return 1
}

func (TBysLua) MD5String(L *lua.LState) int {
	val := L.ToString(1)
	L.Push(lua.LString(u.MD5(val)))
	return 1
}

func (TBysLua) MD5File(L *lua.LState) int {
	p := L.ToString(1)
	if u.PathExist(p) {
		L.Push(lua.LString(u.MD5File(p)))
	} else {
		L.Push(lua.LString(``))
	}
	return 1
}

func (TBysLua) ExtractFileName(L *lua.LState) int {
	p := L.ToString(1)
	L.Push(lua.LString(u.ExtractFileName(p)))
	return 1
}

func (TBysLua) RandString(L *lua.LState) int {
	size := L.ToInt(1)
	if size > 0 {
		L.Push(lua.LString(u.RandString(size)))
	} else {
		L.Push(lua.LString(""))
	}
	return 1
}

func (TBysLua) MakeDir(L *lua.LState) int {
	p := L.ToString(1)
	err := u.MakeDir(p)
	L.Push(lua.LBool(err == nil))
	return 1
}

func (TBysLua) WorkPath(L *lua.LState) int {
	p := L.ToString(1)
	L.Push(lua.LString(u.GetWorkPath(p)))
	return 1
}

func (TBysLua) ReadFileSize(L *lua.LState) int {
	var Result int64

	p := L.ToString(1)

	if u.PathExist(p) {
		Result = u.GetFileSize(p)
	}

	L.Push(lua.LNumber(Result))

	return 1
}

func (TBysLua) Replace(L *lua.LState) int {
	result := L.ToString(1)
	o := L.ToString(2)
	n := L.ToString(3)

	if result != "" && o != "" {
		result = strings.ReplaceAll(result, o, n)
	}

	L.Push(lua.LString(result))
	return 1
}

func (TBysLua) Sha1String(L *lua.LState) int {
	val := L.ToString(1)
	val = u.SHA1([]byte(val))
	L.Push(lua.LString(val))
	return 1
}

func (TBysLua) Output(L *lua.LState) int {
	val := L.Get(1)
	switch val.Type() {
	case lua.LTBool:
		if lua.LVAsBool(val) == true {
			fmt.Print(`true`)
		} else {
			fmt.Print(`false`)
		}
	case lua.LTNumber:
		fmt.Print(lua.LVAsNumber(val))
	case lua.LTString:
		fmt.Print(val)
	}

	return 0
}

func (TBysLua) Input(L *lua.LState) int {
	line := ``
	fmt.Print(`请输入: `)
	fmt.Scan(&line)

	L.Push(lua.LString(line))
	return 1
}

func (TBysLua) BaiduTranslate(L *lua.LState) int {
	table := L.ToTable(1)
	if table != nil {
		s := u.BaiduTranslate(table.RawGetString("appid").String(),
			table.RawGetString("appkey").String(),
			table.RawGetString("fr").String(),
			table.RawGetString("to").String(),
			table.RawGetString("query").String(),
		)
		L.Push(lua.LString(s))

		return 1
	}

	return 0
}

func (TBysLua) PTagContent(L *lua.LState) int {
	lines := strings.Split(L.ToString(1), "\n")
	result := ``
	ok := false
	if len(lines) > 0 {
		OnLine := L.ToFunction(2)
		OnCheck := L.ToFunction(3)
		Contents := strings.Builder{}
		for _, line := range lines {
			if OnCheck != nil { // 检查退出用的
				if err := L.CallByParam(lua.P{Fn: OnCheck,
					NRet:    1,
					Protect: true,
				}, lua.LString(line)); err == nil {
					res := L.Get(-1)
					L.Pop(1)

					if res.Type() == lua.LTBool {
						if lua.LVAsBool(res) {
							break
						}
					}
				}
			}

			if OnLine != nil {
				if err := L.CallByParam(lua.P{Fn: OnLine,
					NRet:    1,
					Protect: true,
				}, lua.LString(line)); err == nil {
					res := L.Get(-1)
					L.Pop(1)

					if res.Type() == lua.LTString {
						if res.String() != "" {
							Contents.WriteString(fmt.Sprintf(`<p>%s</p>`, res.String())) ////用来分割  不要用真正的换行符 因为百度翻译只会翻译第一行。。。
						}
					}
				}
			} else {
				Contents.WriteString(fmt.Sprintf(`<p>%s</p>`, line))
			}
		}

		result = Contents.String()
		ok = true
	}

	L.Push(lua.LString(result))
	L.Push(lua.LBool(ok))
	return 2
}

func (TBysLua) HasSubString(L *lua.LState) int {
	has := false
	text := L.ToString(1)
	table := L.ToTable(2)
	if text != "" && table != nil {
		text = strings.ToLower(text)
		table.ForEach(func(key lua.LValue, value lua.LValue) {
			if !has {
				if value.Type() == lua.LTString {
					if strings.Contains(text, strings.ToLower(value.String())) {
						has = true
					}
				}
			}
		})
	}

	L.Push(lua.LBool(has))
	return 1
}

func (TBysLua) ReadTextFile(L *lua.LState) int {
	text := ``

	if u.PathExist(L.ToString(1)) {
		if f, err := os.ReadFile(L.ToString(1)); err == nil {
			text = string(f)
		}
	}

	L.Push(lua.LString(text))
	return 1
}

func (TBysLua) Sleep(L *lua.LState) int {
	n := L.ToInt(1)
	if n <= 0 {
		n = 1000
	}

	time.Sleep(time.Duration(n) * time.Second)
	return 0
}

func (TBysLua) RandInt(L *lua.LState) int {
	n := L.ToInt(1)
	if n <= 0 {
		n = 10
	}

	result, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))

	L.Push(lua.LNumber(result.Int64()))
	return 1
}

func (TBysLua) RandSleep(L *lua.LState) int {
	n := L.ToInt(1)
	if n <= 0 {
		n = 10
	}
	result, _ := rand.Int(rand.Reader, big.NewInt(int64(n)))

	time.Sleep(time.Duration(result.Int64()) * time.Second)

	L.Push(lua.LNumber(n))
	return 1
}

/*--------------------------------------------------------------------------------------------------------------------*/

func BaseFunctionLoad(L *lua.LState) {
	L.Register(`pathExists`, TBysLua{}.PathExists)
	L.Register(`md5String`, TBysLua{}.MD5String)
	L.Register(`md5File`, TBysLua{}.MD5File)
	L.Register(`extractFilename`, TBysLua{}.ExtractFileName)
	L.Register(`replace`, TBysLua{}.Replace)
	L.Register(`makeDir`, TBysLua{}.MakeDir)
	L.Register(`workPath`, TBysLua{}.WorkPath)
	L.Register(`readFileSize`, TBysLua{}.ReadFileSize)
	L.Register(`replaceString`, TBysLua{}.Replace)
	L.Register(`sha1String`, TBysLua{}.Sha1String)
	L.Register(`output`, TBysLua{}.Output)
	L.Register(`input`, TBysLua{}.Input)
	L.Register(`baiduTranslate`, TBysLua{}.BaiduTranslate)
	L.Register(`ptagContent`, TBysLua{}.PTagContent)
	L.Register(`hasSubString`, TBysLua{}.HasSubString)
	L.Register(`randInt`, TBysLua{}.RandInt)
	L.Register(`sleep`, TBysLua{}.Sleep)
	L.Register(`readTextFile`, TBysLua{}.ReadTextFile)
	L.Register(`randSleep`, TBysLua{}.RandSleep)
}
