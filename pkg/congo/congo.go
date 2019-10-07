package congo

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/format"
	"go/types"
	"io"
	"os"
	"strings"
	"text/template"
)

const capitalLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ParamDesc struct {
	Name string
	NameSnake string
	NameCamel string
	Ptype string
}

func TraverseParams(root *ast.File) ([]*ParamDesc, error) {
	var params []*ParamDesc
	if len(root.Decls) > 1 {
		return nil, errors.New("expected a single config structure called Desc")
	}

	ast.Inspect(root, func(node ast.Node) bool {
		desc, ok := node.(*ast.StructType)
		if !ok {
			return true
		}
		for _, field := range desc.Fields.List {
			fieldName := field.Names[0].Name
			snake := toSnakeCase(fieldName)
			camel := strings.ToUpper(string(fieldName[0])) + fieldName[1:]
			params = append(params, &ParamDesc{
				Name: fieldName,
				NameSnake:snake,
				NameCamel: camel,
				Ptype:types.ExprString(field.Type),
			})
		}

		return false
	})

	return params, nil
}

func RenderTemplate(tmplName, tmplContent, dstPath string, params []*ParamDesc, runFmt, append bool) error {
	openMode := os.O_WRONLY
	if append {
		openMode |= os.O_APPEND
	} else {
		openMode |= os.O_CREATE
	}

	dstFile, err := os.OpenFile(dstPath, openMode, 0644)
	if err != nil {
		return errors.Wrap(err,"create init.go file")
	}
	defer dstFile.Close()

	tmpl, err := template.New(tmplName).Parse(tmplContent)
	if err != nil {
		return errors.Wrap(err, "parse template")
	}

	var dst io.Writer
	if runFmt {
		dst = &bytes.Buffer{}
	} else{
		dst = dstFile
	}

	if err := tmpl.Execute(dst, params); err != nil {
		return errors.Wrap(err, "execute template")
	}

	if runFmt {
		fmted, err := format.Source(dst.(*bytes.Buffer).Bytes())
		if err != nil {
			return errors.Wrap(err, "format")
		}
		fmt.Fprint(dstFile, string(fmted))
	}

	return nil
}

// NOTE: origin should start and end with lowercase letters
func toSnakeCase(origin string) string {
	if strings.ContainsAny(capitalLetters, string(origin[0])+string(origin[len(origin)-1])) {
		fmt.Printf("failed to convert %q to snake case\n", origin)
	} else {
		for {
			capIdx := strings.IndexAny(origin, capitalLetters)
			if capIdx == -1 {
				break
			}
			origin = origin[:capIdx] + "_" + strings.ToLower(string(origin[capIdx])) + origin[capIdx+1:]
		}
	}
	return origin
}
