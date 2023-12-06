package TBys

import (
	"io/ioutil"
)

type TSearchFile struct {
	Files []string
}

func (this *TSearchFile) Count() int {
	return len(this.Files)
}

func (this *TSearchFile) Search(path string, onRead func(fileName string)) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			fileName := path + PathSeparators() + f.Name()
			this.Files = append(this.Files, fileName)
			if onRead != nil {
				onRead(fileName)
			}
		} else {
			err = this.Search(path+PathSeparators()+f.Name(), onRead)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
