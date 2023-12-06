package lua_module

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

type TTestLuaModule struct {
	GetText func(L *lua.LState) int
}

func (this *TTestLuaModule) GetName(L *lua.LState) int {

	fmt.Println(`Hello world`)
	return 0
}

func Test_Main(t *testing.T) {
	Pool := NewLStatePool(LPZipWriter, LPHttpClient)
	state := Pool.Get()
	if err := state.DoString(`local http = require("httpClient")
       http.gZip(true)
       local res = http.get('https://baidu.com')
       print(res)
   `); err != nil {
		t.Error(err)
	}
}
