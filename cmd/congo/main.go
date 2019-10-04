package main

import (
	"congo/pkg/congo"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path"
	"github.com/BurntSushi/toml"
)


func main() {
	fatalf := func(format string, a ...interface{}) {
		fmt.Printf(format+"\n", a...)
		os.Exit(1)
	}

	srcDirPath, err := os.Getwd()
	if err != nil {
		fatalf("find out config source path: %v", err)
	}

	// Traverse config description
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, path.Join(srcDirPath, "config.go"), nil, parser.ParseComments)
	if err != nil {
		fatalf("parse file: %v", err)
	}
	params, err := congo.TraverseParams(root)
	if err != nil {
		fatalf("traverse params: %v", err)
	}

	// Fill init.go
	const initFileName = "init.go"
	initPath := path.Join(srcDirPath, initFileName)
	if err := congo.RenderTemplate(initFileName,
		congo.InitTmplContent,
		initPath,
		params,
		true,
		false,
	); err != nil {
		fatalf("render template init.go: %v", err)
	}

	fmt.Println("generated init.go")

	// Append new parameters to config.toml
	projPath := path.Dir(path.Dir(srcDirPath))
	projName := path.Base(projPath)
	cfgSampleName := projName + ".toml"
	cfgSamplePath := path.Join(projPath, cfgSampleName)
	cfgSampleExists := false
	if _, err := os.Stat(cfgSamplePath); err == nil {
		cfgSampleExists = true
		existingParams := make(map[string]interface{})
		if _, err := toml.DecodeFile(cfgSamplePath, &existingParams); err != nil {
			fatalf("decode existing toml config: %v", err)
		}
		var absentParams []*congo.ParamDesc
		for _, param := range params {
			if _, ok := existingParams[param.NameSnake]; !ok {
				absentParams = append(absentParams, param)
			}
		}
		params = absentParams
	}

	if err := congo.RenderTemplate(cfgSampleName,
		congo.CfgTmplContent,
		path.Join(projPath, cfgSampleName),
		params,
		false,
		cfgSampleExists,
	); err != nil {
		fatalf("render template swatcher-server.toml: %v", err)
	}

	fmt.Printf("generated %s.toml\n", projName)
}
