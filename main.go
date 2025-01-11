package main

import (
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"strings"
	"syscall"
	"text/template"
	"time"
)

type FindResult struct {
	Data []byte
	Path string
	Info fs.FileInfo
}

var dirFS fs.ReadFileFS

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
	dirFS = os.DirFS(".").(fs.ReadFileFS)

	err = fs.WalkDir(dirFS, rootpath, walker)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("fst").Funcs(template.FuncMap{
		"AsMode": AsMode, "AsStat": AsStat, "AsTime": AsTime, "FormatData": FormatData}).Parse(string(fstemplate))
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

	data := make([]byte, 0)
	if info.Type().IsRegular() {
		data, err = dirFS.ReadFile(path)
		if err != nil {
			return err
		}
	}

	results = append(results, FindResult{Data: data, Path: path, Info: fsinfo})
	return nil
}

func AsMode(val fs.FileMode) string {
	modeName := ""
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

	if modeName == "" {
		return fmt.Sprintf("%O", val&fs.ModePerm)
	} else {
		return fmt.Sprintf("%s | %O", modeName, val&fs.ModePerm)
	}

}

func AsTime(secs int64) time.Time {
	return time.Unix(secs, 0)
}

func AsStat(val any) *syscall.Stat_t {
	if s, ok := val.(*syscall.Stat_t); ok {
		return s
	} else {
		return nil
	}
}

func FormatData(data []byte) string {
	if len(data) == 0 {
		return "[]byte{}"
	}
	chars := make([]string, len(data))
	for n, v := range data {
		chars[n] = fmt.Sprintf("0x%x", v)
	}
	return "[]byte{" + strings.Join(chars, ", ") + "}"
}
