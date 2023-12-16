package lua_module

import (
	"bytes"
	u "github.com/Icy2010/TBys"
	"github.com/alexmullins/zip"
	lua "github.com/yuin/gopher-lua"
	"io"
	"io/ioutil"
	"os"
)

type TZipFileName struct {
	Name string
	Path string
	Body []byte
}

type TLuaZipWriter struct {
	files    []TZipFileName
	password string
}

func (this *TLuaZipWriter) addFile(L *lua.LState) int {
	file := L.ToString(1)
	if u.PathExist(file) {
		name := L.ToString(2)
		if name == "" {
			name = u.ExtractFileName(file)
		}
		this.files = append(this.files, TZipFileName{
			Name: name,
			Path: file,
			Body: nil,
		})
	}

	return 0
}

func (this *TLuaZipWriter) addBody(L *lua.LState) int {
	name := L.ToString(1)
	body := L.ToString(2)
	if name != "" && body != "" {
		this.files = append(this.files, TZipFileName{
			Name: name,
			Path: "",
			Body: []byte(body),
		})
	}
	return 0
}

func (this *TLuaZipWriter) setPassword(L *lua.LState) int {
	this.password = L.ToString(1)
	return 0
}

func (this *TLuaZipWriter) clear(L *lua.LState) int {
	this.files = make([]TZipFileName, 0)
	return 0
}

func (this *TLuaZipWriter) saveTo(L *lua.LState) int {
	if len(this.files) > 0 {
		path := L.ToString(1)
		fzip, err := os.Create(path)
		if err != nil {
			u.Logger().Error("TLuaZipWriter", err)
			return 0
		}

		zipw := zip.NewWriter(fzip)
		defer zipw.Close()

		getBody := func(file TZipFileName) ([]byte, error) {
			if u.PathExist(file.Path) {
				if Body, e := ioutil.ReadFile(file.Path); e != nil {
					u.Logger().Error("TLuaZipWriter", e)

				} else {
					return Body, nil
				}

			}

			return file.Body, nil

		}

		getWriteZip := func(name string) (io.Writer, error) {
			if this.password != "" {
				return zipw.Encrypt(name, this.password)
			} else {
				return zipw.Create(name)
			}
		}

		doAdded := L.ToFunction(2)

		for _, file := range this.files {
			Body := make([]byte, 0)
			Body, err = getBody(file)
			if err != nil {
				u.Logger().Error("TLuaZipWriter", err)
				return 0
			}

			w, e := getWriteZip(file.Name)
			if e != nil {
				u.Logger().Error("TLuaZipWriter", err)
				return 0
			}

			_, err = io.Copy(w, bytes.NewReader(Body))
			if err != nil {
				u.Logger().Error("TLuaZipWriter", err)
				return 0
			}

			if (doAdded) != nil {
				if err := L.CallByParam(lua.P{
					Fn:      doAdded,
					NRet:    0,
					Protect: true,
				}, lua.LString(file.Name)); err != nil {
					u.Logger().Error("TluaZipWriter", err)
				}
			}
		}

		err = zipw.Flush()
		if err != nil {
			u.Logger().Error("TLuaZipWriter", err)
		}
	}
	return 0
}

func ZipWriterPreload(L *lua.LState) {
	L.PreloadModule("zipWriter", func(L *lua.LState) int {
		wzip := &TLuaZipWriter{}
		wzip.password = ""
		t := L.NewTable()
		L.SetFuncs(t, map[string]lua.LGFunction{
			`addFile`:     wzip.addFile,
			`addBody`:     wzip.addBody,
			`setPassword`: wzip.setPassword,
			`save`:        wzip.saveTo,
			`clear`:       wzip.clear,
		})
		L.Push(t)
		return 1
	})
}
