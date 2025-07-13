//go:build ignore

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/vektah/gqlparser/v2/ast"
)

func constraintFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	// Call default hook to proceed standard directives like goField and goTag.
	// You can omit it, if you don't need.
	if f, err := modelgen.DefaultFieldMutateHook(td, fd, f); err != nil {
		return f, err
	}

	gorm_tag := ""

	for _, c := range fd.Directives {
		switch c.Name {
		case "gorm":
			formatConstraint := c.Arguments.ForName("tags")

			if formatConstraint != nil {
				gorm_tag += formatConstraint.Value.Raw + ";"
			}

		case "primary":
			gorm_tag += "primary;"
		}
	}

	if gorm_tag != "" {
		f.Tag += " gorm:" + strconv.Quote(strings.TrimSuffix(gorm_tag, ";"))
	}

	return f, nil
}

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	// Attaching the mutation function onto modelgen plugin
	p := modelgen.Plugin{
		FieldHook: constraintFieldHook,
	}

	err = api.Generate(cfg, api.ReplacePlugin(&p))

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
