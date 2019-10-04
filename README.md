Static config generator based on code generation.
Generates two files: 
- `init.go` (opens configuration file in the root path of project and validates parameters)
- `<PROJECT_NAME>.toml` (sample configuration file with parameters set to default values)

There are some assumptions your project should satisfy to use this tool:
1. Expected layout: `/pkg/config/config.go`
2. `config.go` must contain structure `Desc` which consists of demanded parameters

Supported parameter types:
- strings
- integers

To generate files for your project just add comment string ```//go:generate congo``` to the top of `config.go` and run `go generate ./...` from the root.