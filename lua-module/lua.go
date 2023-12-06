package lua_module

import (
	lxml "github.com/ailncode/gluaxmlpath"
	lurl "github.com/cjoudrey/gluaurl"
	ljson "github.com/layeh/gopher-json"
	lua "github.com/yuin/gopher-lua"
	"sync"
)

type TLPModule int

const (
	LPCache TLPModule = iota + 1
	LPCTFile
	LPFFMpeg
	LPFileSystem
	LPHtmlParser
	LPHttpClient
	LPIni
	LPSql
	LPZipWriter
	LPJson
	LPXml
	LPUrl
)

type TLModuleFunc = func(L *lua.LState) int

type TLStatePool struct {
	m     sync.Mutex
	saved []*lua.LState

	lpModules  []TLPModule
	OnNewState func(L *lua.LState)
}

func (this *TLStatePool) preloadDefault(L *lua.LState) {
	BaseFunctionLoad(L)
	if len(this.lpModules) > 0 {
		for i := 0; i < len(this.lpModules); i++ {
			switch this.lpModules[i] {
			case LPCache:
				CacheRegisterModule(L)
			case LPCTFile:
				CTFilePreload(L)
			case LPFFMpeg:
				FFMpegPreload(L)
			case LPFileSystem:
				FileSystemPreload(L)
			case LPHtmlParser:
				HtmlParserPreload(L)
			case LPHttpClient:
				HttpClientPreload(L)
			case LPIni:
				IniPreload(L)
			case LPSql:
				SqlRegisterModule(L)
			case LPZipWriter:
				ZipWriterPreload(L)
			case LPJson:
				ljson.Preload(L)
			case LPXml:
				lxml.Preload(L)
			case LPUrl:
				L.PreloadModule(`url`, lurl.Loader)
			}
		}
	}
}

func (pl *TLStatePool) Get() *lua.LState {
	pl.m.Lock()
	defer pl.m.Unlock()
	n := len(pl.saved)
	if n == 0 {
		return pl.New()
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return x
}

func (pl *TLStatePool) New() *lua.LState {
	L := lua.NewState()
	pl.preloadDefault(L)

	if pl.OnNewState != nil {
		pl.OnNewState(L)
	}
	// setting the L up here.
	// load scripts, set global variables, share channels, etc...
	return L
}

func (pl *TLStatePool) Put(L *lua.LState) {
	pl.m.Lock()
	defer pl.m.Unlock()
	pl.saved = append(pl.saved, L)
}

func (pl *TLStatePool) Shutdown() {
	for _, L := range pl.saved {
		L.Close()
	}
}

/*--------------------------------------------------------------------------------------------------------------------*/

func NewLStatePool(modules ...TLPModule) *TLStatePool {
	return &TLStatePool{
		lpModules:  modules,
		OnNewState: nil,
	}
}
