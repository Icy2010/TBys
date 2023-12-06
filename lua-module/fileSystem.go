package lua_module

import (
	u "github.com/Icy2010/TBys"
	lua "github.com/yuin/gopher-lua"
	"os"
)

type TLuaFileSystem struct {
}

func (this *TLuaFileSystem) getFiles(L *lua.LState) int {
	path := L.ToString(1)
	table := L.NewTable()

	if u.PathExist(path) {
		sf := u.TSearchFile{}
		if err := sf.Search(path, func(fileName string) {
			table.Append(lua.LString(fileName))
		}); err != nil {
			u.Logger().Error(err)
		}

	}

	L.Push(table)
	return 1
}

func (this *TLuaFileSystem) copyFile(L *lua.LState) int {
	input := L.ToString(1)
	output := L.ToString(2)
	success := false
	if input != "" && output != "" {
		err := u.CopyFile(input, output)

		success = err == nil
		if !success {
			u.Logger().Error(err)
		}
	}

	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaFileSystem) deleteFile(L *lua.LState) int {
	fileName := L.ToString(1)
	success := false
	err := os.Remove(fileName)
	success = err == nil
	if !success {
		u.Logger().Error(err)
	}
	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaFileSystem) renameFile(L *lua.LState) int {
	fileName := L.ToString(1)
	newName := L.ToString(2)
	success := false
	err := os.Rename(fileName, newName)
	success = err == nil
	if !success {
		u.Logger().Error(err)
	}
	L.Push(lua.LBool(success))
	return 1
}

/*--------------------------------------------------------------------------------------------------------------------*/

func FileSystemPreload(L *lua.LState) {
	L.PreloadModule(`fs`, func(L *lua.LState) int {
		fs := &TLuaFileSystem{}
		t := L.NewTable()

		L.SetFuncs(t, map[string]lua.LGFunction{
			`getFiles`:   fs.getFiles,
			`copyFile`:   fs.copyFile,
			`deleteFile`: fs.deleteFile,
			`renameFile`: fs.renameFile,
		})

		L.Push(t)
		return 1
	})
}
