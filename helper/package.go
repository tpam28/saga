package helper

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

func FindPackageName(path string) string {
	path = filepath.Dir(path)
	files, err := filepath.Glob(path + "*.go")
	if len(files) == 0 || err != nil {
		list := strings.Split(path, "/")
		return list[len(list)-1]
	}

	file, err := os.Open(path+files[0])
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	line, _, err := r.ReadLine()
	if err != nil {
		panic(err)
	}
	packageName := bytes.TrimPrefix(line, []byte("package "))
	return string(packageName)
}
