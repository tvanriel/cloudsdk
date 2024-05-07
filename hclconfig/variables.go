package hclconfig

import (
	"io"
	"os"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func parseVariables(filename string) (map[string]cty.Value, error) {
	f, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	src, err := io.ReadAll(f)

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

	variables := make(map[string]cty.Value)

	localBlocks := blocks.Blocks.OfType("local")
	for i := range localBlocks {
		attrs, diags := localBlocks[i].Body.JustAttributes()
		if diags.HasErrors() {
			return nil, diags
		}
		for a := range attrs {
			v, diags := attrs[a].Expr.Value(nil)
			if diags.HasErrors() {
				return nil, diags
			}
			k := attrs[a].Name
			variables[k] = v
		}
	}

	return variables, nil

}
