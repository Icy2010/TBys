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
	Pool := NewLStatePool(LPZipWriter, LPHttpClient, LPHtmlParser)
	state := Pool.Get()
	if err := state.DoString(`local hParser = require('htmlParser')
local http = require('httpClient')
 

function ReadPageContent(Src,OnContent)
   return hParser.find(Src,[[.//div[@class="w-post-elm post_content"]],OnContent)
end

--http.proxy('http://102.0.0.1:1080/')

local data =  hParser.href(http.get('https://downloadly.net/?s=delphi'),[[//h2//a/@href]])
if #data > 0 then  
   print(data[1].text,data[1].href)
end
 
   `); err != nil {
		t.Error(err)
	}
}
