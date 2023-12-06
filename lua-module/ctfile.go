package lua_module

import (
	"fmt"
	u "github.com/Icy2010/TBys"
	z "github.com/Icy2010/ZeligCTFile"
	lua "github.com/yuin/gopher-lua"
)

type TLuaCTFile struct {
	TBaseBysLua
	public bool
	z.TCTFile
}

func (this *TLuaCTFile) SetPublic(L *lua.LState) int {
	this.public = L.ToBool(1)
	return 0
}

func (this *TLuaCTFile) login(L *lua.LState) int { // 注意大写的Login 是根方法
	table := L.ToTable(1)
	success := false

	if table != nil {
		switch u.StrToInt(table.RawGetString("way").String()) {
		case 1: // 密码
			success = this.Login(table.RawGetString("email").String(), table.RawGetString("password").String()) == nil
		default:
			success = this.LoginFromToken(table.RawGetString("token").String()) == nil
		}
	}

	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaCTFile) Cloud() *z.TCTFileMethods {
	if this.public {
		return this.PublicCloud()
	}

	return this.PrivateCloud()
}

func (this *TLuaCTFile) filesFromIds(L *lua.LState) int {
	ids := L.ToTable(1)
	table := L.NewTable()
	if ids != nil {
		files, err := this.Cloud().FileIdsList(LuaTableToStrings(ids))
		if err == nil {
			if len(files) > 0 {
				for _, P := range files {
					data := L.NewTable()
					data.RawSetString("name", lua.LString(P.Name))
					data.RawSetString("icon", lua.LString(P.Icon))
					data.RawSetString("key", lua.LString(P.Key))
					data.RawSetString("size", lua.LNumber(P.Size))
					data.RawSetString("date", lua.LNumber(P.Date))
					table.Append(data)
				}
			}
		} else {
			logger(err.Error())
		}
	}

	L.Push(table)
	return 1
}

func (this *TLuaCTFile) files(L *lua.LState) int {
	params := L.ToTable(1)
	table := L.NewTable()
	if params != nil {

		files, err := this.Cloud().FileList(params.RawGetString("folderId").String(),
			int(lua.LVAsNumber(params.RawGetString("start"))),
			int(lua.LVAsNumber(params.RawGetString("reload"))),
			params.RawGetString("orderby").String(),
			params.RawGetString("filter").String(),
			params.RawGetString("keyword").String())

		if err == nil {
			if len(files) > 0 {
				for _, P := range files {
					data := L.NewTable()
					data.RawSetString("name", lua.LString(P.Name))
					data.RawSetString("icon", lua.LString(P.Icon))
					data.RawSetString("key", lua.LString(P.Key))
					data.RawSetString("size", lua.LNumber(P.Size))
					data.RawSetString("date", lua.LNumber(P.Date))
					table.Append(data)
				}
			}
		} else {
			logger(err.Error())
		}
	}

	L.Push(table)
	return 1
}

func (this *TLuaCTFile) foldersFromId(L *lua.LState) int {
	id := L.ToString(1)
	table := L.NewTable()
	if id != "" {
		folders, err := this.Cloud().FolderList(id)
		if err == nil {
			if len(folders) > 0 {
				for _, P := range folders {
					data := L.NewTable()
					data.RawSetString("name", lua.LString(P.Name))
					data.RawSetString("icon", lua.LString(P.Icon))
					data.RawSetString("key", lua.LString(P.Key))
					data.RawSetString("date", lua.LNumber(P.Date))
					table.Append(data)
				}
			}
		} else {
			logger(err.Error())
		}
	}

	L.Push(table)
	return 1
}

func (this *TLuaCTFile) folderCreate(L *lua.LState) int {
	table := L.ToTable(1)
	result := L.NewTable()
	if table != nil {
		res, err := this.Cloud().FolderCreate(table.RawGetString("folderId").String(),
			table.RawGetString("name").String(),
			table.RawGetString("description").String(),
			int(lua.LVAsNumber(table.RawGetString("isHidden"))),
		)

		if err == nil {
			result.RawSetString("folderId", lua.LString(res[`folder_id`]))
			result.RawSetString("folderPath", lua.LString(res[`folder_path`]))
		}
	}

	L.Push(result)
	return 1
}

func (this *TLuaCTFile) folderMate(L *lua.LState) int {
	id := L.ToString(1)
	result := L.NewTable()
	if id != "" {
		res, err := this.Cloud().FolderMeta(id)
		if err == nil {
			result.RawSetString("name", lua.LString(res.Name))
			result.RawSetString("path", lua.LString(res.Path))
			result.RawSetString("icon", lua.LString(res.Icon))
			result.RawSetString("key", lua.LString(res.Key))
			result.RawSetString("isHidden", lua.LNumber(res.Is_hidden))
		}
	}

	L.Push(result)
	return 1
}

func (this *TLuaCTFile) folderModify(L *lua.LState) int {
	table := L.ToTable(1)
	result := false
	if table != nil {
		_, err := this.PublicCloud().FolderModifyMeta(table.RawGetString("folderId").String(),
			table.RawGetString("name").String(),
			table.RawGetString("description").String(),
			int(lua.LVAsNumber(table.RawGetString("isHidden"))),
		)
		result = err == nil
	}

	L.Push(lua.LBool(result))
	return 1
}

func (this *TLuaCTFile) fileMeta(L *lua.LState) int {
	id := L.ToString(1)
	result := L.NewTable()
	if id != "" {
		res, err := this.Cloud().FileMeta(id)
		if err == nil {
			result.RawSetString("name", lua.LString(res.Name))
			result.RawSetString("path", lua.LString(res.Path))
			result.RawSetString("icon", lua.LString(res.Icon))
			result.RawSetString("key", lua.LString(res.Key))
			result.RawSetString("size", lua.LNumber(res.Size))
		}
	}

	L.Push(result)
	return 1
}

func (this *TLuaCTFile) fileShare(L *lua.LState) int {
	ids := L.ToTable(1)
	table := L.NewTable()
	if ids != nil {
		if res, err := this.Cloud().FileShare(LuaTableToStrings(ids)); err == nil {
			for _, p := range res {
				data := L.NewTable()
				data.RawSetString("icon", lua.LString(p.Icon))
				data.RawSetString("name", lua.LString(p.Name))
				data.RawSetString("key", lua.LString(p.Key))
				data.RawSetString("date", lua.LNumber(p.Date))
				data.RawSetString("size", lua.LNumber(p.Size))
				data.RawSetString("drLink", lua.LString(p.Drlink))
				data.RawSetString("webLink", lua.LString(p.Weblink))
				data.RawSetString("xtCode", lua.LString(p.XtCode))
				table.Append(data)
			}
		}
	}

	L.Push(table)
	return 1
}

func (this *TLuaCTFile) fileMove(L *lua.LState) int {
	id := L.ToString(1)
	ids := L.ToTable(2)
	success := false
	if id != "" && ids != nil {
		success = this.Cloud().FileMove(id, LuaTableToStrings(ids)) == nil
	}

	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaCTFile) fileDelete(L *lua.LState) int {
	ids := L.ToTable(1)
	success := false
	if ids != nil {
		success = this.Cloud().FileDelete(LuaTableToStrings(ids)) == nil
	}

	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaCTFile) fileDownload(L *lua.LState) int {
	ids := L.ToTable(1)
	table := L.NewTable()
	if ids != nil {
		if res, err := this.Cloud().FileDownload(LuaTableToStrings(ids)); err == nil {
			for _, p := range res {
				data := L.NewTable()
				data.RawSetString("icon", lua.LString(p.Icon))
				data.RawSetString("name", lua.LString(p.Name))
				data.RawSetString("key", lua.LString(p.Key))
				data.RawSetString("path", lua.LString(p.Path))
				data.RawSetString("size", lua.LNumber(p.Size))
				table.Append(data)

			}
		}
	}

	L.Push(table)
	return 1
}

func (this TLuaCTFile) fileSave(L *lua.LState) int {
	ids := L.ToTable(1)
	success := false
	if ids != nil {
		success = this.Cloud().FileSave(LuaTableToStrings(ids)) == nil
	}

	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaCTFile) fileFetchUrl(L *lua.LState) int {
	id := L.ToString(1)
	url := ``
	err := fmt.Errorf(`未知错误`)

	if id != "" {
		url, err = this.Cloud().FileFetchUrlb(id)
	}

	L.Push(lua.LString(url))
	L.Push(lua.LBool(err == nil))
	return 2
}

func (this *TLuaCTFile) fileRecycleEmpty(L *lua.LState) int {
	ids := L.ToTable(1)
	success := false
	if ids != nil {
		success = this.Cloud().FileRecycleEmpty(LuaTableToStrings(ids)) == nil
	}

	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaCTFile) fileRecycleEmptyAll(L *lua.LState) int {
	L.Push(lua.LBool(this.Cloud().FileRecycleEmptyAll() == nil))
	return 1
}

func (this *TLuaCTFile) fileRecycle(L *lua.LState) int {
	Start := L.ToInt(1)
	Reload := L.ToInt(2)
	table := L.NewTable()
	if res, err := this.Cloud().FileRecycle(Start, Reload); err == nil {
		for _, p := range res {
			data := L.NewTable()
			data.RawSetString("icon", lua.LString(p.Icon))
			data.RawSetString("name", lua.LString(p.Name))
			data.RawSetString("key", lua.LString(p.Key))
			data.RawSetString("size", lua.LNumber(p.Size))
			data.RawSetString("imgSrc", lua.LString(p.Imgsrc))
			data.RawSetString("delTime", lua.LNumber(p.Del_time))
			table.Append(data)
		}
	}

	L.Push(table)
	return 1
}

func (this *TLuaCTFile) fileUpload(L *lua.LState) int {
	id := L.ToString(1)
	fileName := L.ToString(2)
	rid := ``
	err := fmt.Errorf(`未知错误`)
	if id != "" && u.PathExist(fileName) {
		rid, err = this.Cloud().FileUpload(id, fileName)
	}

	L.Push(lua.LString(rid))
	L.Push(lua.LBool(err == nil))
	return 2
}

/*--------------------------------------------------------------------------------------------------------------------*/

func CTFilePreload(L *lua.LState) {
	L.PreloadModule("ctFile", func(L *lua.LState) int {
		ct := &TLuaCTFile{}
		ct.public = true
		t := L.NewTable()
		L.SetFuncs(t, map[string]lua.LGFunction{
			`login`:               ct.login,
			`fileFetchUrl`:        ct.fileFetchUrl,
			`files`:               ct.files,
			`filesFromIds`:        ct.filesFromIds,
			`fileDelete`:          ct.fileDelete,
			`fileDownload`:        ct.fileDownload,
			`fileUpload`:          ct.fileUpload,
			`fileShare`:           ct.fileShare,
			`fileMeta`:            ct.fileMeta,
			`fileMove`:            ct.fileMove,
			`fileRecycle`:         ct.fileRecycle,
			`fileRecycleEmpty`:    ct.fileRecycleEmpty,
			`fileRecycleEmptyAll`: ct.fileRecycleEmptyAll,
			`folderMate`:          ct.folderMate,
			`folderModify`:        ct.folderModify,
			`folderCreate`:        ct.folderCreate,
			`foldersFromId`:       ct.foldersFromId,
		})

		L.Push(t)
		return 1
	})
}
