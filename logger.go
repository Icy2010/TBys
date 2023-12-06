package TBys

import (
	"fmt"
	"github.com/fatih/color"
)

type TLogger struct {
}

var (
	logger TLogger
)

func Logger() *TLogger {
	return &logger
}

func NewLogger() {
	//logger
}

func (TLogger) Info(val ...any) {
	fmt.Println(val...)
}

func (TLogger) Warn(val ...any) {
	color.New(color.BgYellow).Println(val...)
}

func (TLogger) Error(val ...any) {
	color.New(color.FgRed).Println(val...)
}

func (TLogger) Debug(val ...any) {
	color.New(color.FgGreen).Println(val...)
}
