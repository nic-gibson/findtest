package main

import (
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"syscall"
	"text/template"
)

//go:embed "fstemplate"
var fstemplate string

type FindResult struct {
	Path string
	Info fs.FileInfo
}

var results []FindResult = []FindResult{}

func main() {

	rootpath := "testdata"
	dirFS := os.DirFS(".")

	err := fs.WalkDir(dirFS, rootpath, walker)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("fst").Funcs(template.FuncMap{
		"AsMode": , AsMode, "AsStat": AsStat}).Parse(fstemplate)
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
    mode := val & fs.ModeType
    switch mode:wq

	return fmt.Sprintf("%0d", uint32(val))
}

func AsStat(val any) *syscall.Stat_t {
	if s, ok := val.(*syscall.Stat_t); ok {
		return s
	} else {
		return nil
	}
}
