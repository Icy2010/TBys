package TBys

import (
	"io/fs"
	p "path"
	"path/filepath"
	"strings"
	"time"
)

type TFileInfo struct {
	Path    string    `json:"path,omitempty"`
	Name    string    `json:"name,omitempty"`
	Size    int64     `json:"size,omitempty"`
	Mode    uint32    `json:"mode,omitempty"`
	ModTime time.Time `json:"modTime"`
}

type TSearchFile struct {
	TBasic
	Files  []TFileInfo
	Suffix string
}

func (this *TSearchFile) Count() int {

	return len(this.Files)
}

func (this *TSearchFile) Search(DirPath string, onFile func(path string, info fs.FileInfo) bool) error {
	this.Files = make([]TFileInfo, 0)
	e := filepath.Walk(DirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			this.LogError(`TSearchFile.Search`, err.Error())
		} else {
			if !info.IsDir() {
				ok := true
				if this.Suffix != ".*" && this.Suffix != "" {
					ext := p.Ext(info.Name())
					if !strings.EqualFold(ext, this.Suffix) {
						ok = false
					}
				}

				if ok {
					if onFile != nil {
						ok = onFile(path, info)
					}
					if ok {
						this.Files = append(this.Files, TFileInfo{
							Path:    path,
							Name:    info.Name(),
							Size:    info.Size(),
							Mode:    uint32(info.Mode()),
							ModTime: info.ModTime(),
						})
					}
				}
			} else {
				if onFile != nil {
					onFile(path, info)
				}
			}
		}

		return err
	})

	return e
}
