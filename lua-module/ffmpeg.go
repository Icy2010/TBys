package lua_module

import (
	"bytes"
	"fmt"
	u "github.com/Icy2010/TBys"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	lua "github.com/yuin/gopher-lua"
	"io"
	"os"
)

type TLuaFFMpeg struct {
}

func (this *TLuaFFMpeg) fast(L *lua.LState) int {
	input := L.ToString(1)
	output := L.ToString(2)
	success := false
	if input != "" && output != "" {
		err := ffmpeg.Input(input).
			Output(output, ffmpeg.KwArgs{"vcodec": "copy", "acodec": "copy"}).
			OverWriteOutput().ErrorToStdOut().Run()
		success = err == nil
	}

	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaFFMpeg) addWatermark(L *lua.LState) int {
	table := L.ToTable(1)
	success := false
	if table != nil {
		overlay := ffmpeg.Input(lua.LVAsString(table.RawGetString(`image`))).Filter("scale", ffmpeg.Args{lua.LVAsString(table.RawGetString(`scale`))})
		err := ffmpeg.Filter(
			[]*ffmpeg.Stream{
				ffmpeg.Input(lua.LVAsString(table.RawGetString(`input`))),
				overlay,
			}, "overlay", ffmpeg.Args{lua.LVAsString(table.RawGetString(`overlay`))}, ffmpeg.KwArgs{"enable": "gte(t,1)"}).
			Output(lua.LVAsString(table.RawGetString(`output`))).OverWriteOutput().ErrorToStdOut().Run()
		if err != nil {
			u.Logger().Error(err)
		} else {
			success = true
		}
	}

	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaFFMpeg) takeFrame(L *lua.LState) int {
	table := L.ToTable(1)
	success := false
	if table != nil {
		ExampleReadFrameAsJpeg := func(inFileName string, frameNum int) (io.Reader, bool) {
			buf := bytes.NewBuffer(nil)
			err := ffmpeg.Input(inFileName).
				Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
				Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
				WithOutput(buf, os.Stdout).
				Run()
			if err != nil {
				u.Logger().Error(err)

				return nil, false
			}
			return buf, true
		}

		if reader, ok := ExampleReadFrameAsJpeg(lua.LVAsString(table.RawGetString(`input`)), int(lua.LVAsNumber(table.RawGetString(`num`)))); ok {
			img, err := imaging.Decode(reader)
			if err != nil {
				u.Logger().Error(err)
			} else {
				err = imaging.Save(img, lua.LVAsString(table.RawGetString(`output`)))
				if err != nil {
					u.Logger().Error(err)
				} else {
					success = true
				}
			}
		}
	}

	L.Push(lua.LBool(success))
	return 1
}

func (this *TLuaFFMpeg) ExtractMP3(L *lua.LState) int {
	input := L.ToString(1)
	output := L.ToString(2)
	success := false
	if input != "" && output != "" {
		err := ffmpeg.Input(input).
			Output(output, ffmpeg.KwArgs{"f": "mp3", "vn": ""}).
			OverWriteOutput().ErrorToStdOut().Run()
		success = err == nil
	}

	L.Push(lua.LBool(success))
	return 1
}

/*-------------------------------------------------------------------------------------------------------------------*/

func FFMpegPreload(L *lua.LState) {
	L.PreloadModule(`ffmpeg`, func(L *lua.LState) int {
		ff := &TLuaFFMpeg{}
		t := L.NewTable()
		L.SetFuncs(t, map[string]lua.LGFunction{
			`fast`:         ff.fast,
			`addWatermark`: ff.addWatermark,
			`takeFrame`:    ff.takeFrame,
			`extractMP3`:   ff.ExtractMP3,
		})

		L.Push(t)
		return 1
	})
}
