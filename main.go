package main

import (
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"syscall"
	"text/template"
	"time"
)

type FindResult struct {
	Path string
	Info fs.FileInfo
}

var results []FindResult = []FindResult{}

func main() {

	var templateFile string
	if runtime.GOOS == "darwin" {
		templateFile = "fstemplate_darwin"
	} else {
		templateFile = "fstemplate_linux"
	}

	fstemplate, err := os.ReadFile(templateFile)
	if err != nil {
		panic(err)
	}

	rootpath := "testdata"
	dirFS := os.DirFS(".")

	err = fs.WalkDir(dirFS, rootpath, walker)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("fst").Funcs(template.FuncMap{
		"AsMode": AsMode, "AsStat": AsStat, "AsDate": AsDate}).Parse(string(fstemplate))
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, results)
	if err != nil {
		panic(err)
	}

}

func walker(path string, info fs.DirEntry, err error) error {

	if err != nil {
		return err
	}

	fsinfo, err := info.Info()
	if err != nil {
		return err
	}

	results = append(results, FindResult{Path: path, Info: fsinfo})
	return nil
}

func AsMode(val fs.FileMode) string {
	modeName := "0"
	mode := val & fs.ModeType
	switch mode {
	case fs.ModeDir:
		modeName = "fs.ModeDir"
	case fs.ModeSymlink:
		modeName = "fs.ModeSymlink"
	case fs.ModeNamedPipe:
		modeName = "fs.ModeNamedPipe"
	case fs.ModeSocket:
		modeName = "fs.ModeSocket"
	case fs.ModeDevice:
		modeName = "fs.ModeDevice"
	case fs.ModeCharDevice:
		modeName = "fs.ModeCharDevice"
	case fs.ModeIrregular:
		modeName = "fs.ModeIrregular"
	}

	return fmt.Sprintf("%s | %O", modeName, val&fs.ModePerm)
}

func AsDate(val time.Time) string {
	return fmt.Sprintf("time.Date(%d, %d, %d, %d, %d, %d, %d, time.UTC)", val.Year(), val.Month(), val.Day(), val.Hour(), val.Minute(), val.Second(), val.Nanosecond())
}

func AsStat(val any) *syscall.Stat_t {
	if s, ok := val.(*syscall.Stat_t); ok {
		return s
	} else {
		return nil
	}
}
