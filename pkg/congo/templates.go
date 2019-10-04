package congo

const InitTmplContent = `// Code generated by congo. DO NOT EDIT.
// source: config.go

package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

type Config interface {
	{{range .}}Get{{.NameCamel}}() {{.Ptype}}
	{{end}}
}

{{range .}}
func (d *Desc) Get{{.NameCamel}}() {{.Ptype}} {
	return d.{{.Name}}
}
{{end}}

func Init(configPath string) (Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "open config file")
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "read config file content")
	}
	c := &Desc{}
	if _, err := toml.Decode(string(content), c); err != nil {
		return nil, errors.Wrapf(err, "decode config file content")
	}
	if err := validate(c); err != nil {
		return nil, errors.Wrapf(err, "validate config")
	}
	return c, nil
}

func validate(d *Desc) (err error) {
{{range .}}if d.{{.Name}} == {{if eq .Ptype "string"}} "" {{else}} 0 {{end}} {
		return errors.Wrapf(err, "%q missing", "{{.NameSnake}}")
	}
{{end}}return nil
}
`

const CfgTmplContent = `{{range .}}{{.NameSnake}}={{if eq .Ptype "string"}}""{{else}}0{{end}}
{{end}}`