package hclconfig

import (
	"fmt"
	"io"
	"os"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func parseVariables(filename string) (map[string]cty.Value, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", filename, err)
	}

	defer f.Close()

	src, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", filename, err)
	}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)

	if diags.HasErrors() {
		return nil, diags
	}

	blocks, _, diags := file.Body.PartialContent(
		&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{Type: "local"},
			},
		},
	)

	if diags.HasErrors() {
		return nil, diags
	}

	variables := map[string]cty.Value{}

	localBlocks := blocks.Blocks.OfType("local")
	for i := range localBlocks {
		attrs, diags := localBlocks[i].Body.JustAttributes()
		if diags.HasErrors() {
			return nil, fmt.Errorf("parse %q: %w", filename, diags)
		}
		for a := range attrs {
			v, diags := attrs[a].Expr.Value(&hcl.EvalContext{Variables: variables})
			if diags.HasErrors() {
				return nil, fmt.Errorf("evaluate variable %s in %s:%d+%d: %w", attrs[a].Name, attrs[a].Range.Filename, attrs[a].Range.Start.Line, attrs[a].Range.Start.Column, diags)
			}
			k := attrs[a].Name
			variables[k] = v
		}
	}

	return variables, nil
}
