package congo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go/parser"
	"go/token"
	"testing"
)

func TestTraverseParams(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	for _, tcase := range []struct{
		src string
		expected []ParamDesc
	} {
		{
			src: `package config

type Desc struct {
	fooBar string
	barFoo string
}
`,
			expected: []ParamDesc{
				{
					Name: "fooBar",
					NameSnake: "foo_bar",
					NameCamel: "FooBar",
					Ptype: "string",
				},
				{
					Name: "barFoo",
					NameSnake: "bar_foo",
					NameCamel: "BarFoo",
					Ptype: "string",
				},
			},
		},
		{
			src: `package config

type Desc struct {
	fooBar int64
}
`,
			expected: []ParamDesc{
				{
					Name: "fooBar",
					NameSnake: "foo_bar",
					NameCamel: "FooBar",
					Ptype: "int64",
				},
			},
		},
	} {
		fset := token.NewFileSet()
		root, err := parser.ParseFile(fset, "", tcase.src, parser.ParseComments)
		require.NoError(err)
		params, err := TraverseParams(root)
		require.NoError(err)
		assert.Equal(tcase.expected, params)
	}
}

func TestToSnakeCase(t *testing.T) {
	assert := assert.New(t)
	for _, tcase := range []struct{
		origin string
		expected string
	} {
		{
			origin: "superParam",
			expected:"super_param",
		},
		{
			origin: "a",
			expected: "a",
		},
		{
			origin: "ParamStartsWithCapitalLetter", // unexpected origin
			expected: "ParamStartsWithCapitalLetter",
		},
		{
			origin: "paramEndsWithCapitalLetteR", // unexpected origin
			expected: "paramEndsWithCapitalLetteR",
		},
	} {
		assert.Equal(tcase.expected, toSnakeCase(tcase.origin))
	}
}
