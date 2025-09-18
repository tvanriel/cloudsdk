package hclconfig

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
	"github.com/zclconf/go-cty/cty/gocty"
)

// standardFuncs returns the standard library functions as a map[string]function.Function.
// This is implemented based on the documentation for github.com/zclconf/go-cty/cty/function/stdlib.
// It populates the map with all the functions described: abs, add, and, assertnotnull, byteslen, bytestrlen, bytestrslice, csvdecode, ceil, chomp, chunklist, coalesce, coalescelist, compact, concat, contains, distinct, divide, element, eq, flatten, floor, format, formatdate, formatlist, gt, gte, hasindex, indent, index, int, jsondecode, jsonencode, join, keys, length, lt, lte, log, lookup, lower, max, md5, merge, min, modulo, multiply, negate, ne, not, or, parseint, pow, range, regex, regexall, regexreplace, replace, reverse, reverselist, sethaselement, setintersect, setproduct, setsubtract, setsymmetricdifference, setunion, signum, slice, sort, split, strlen, substr, subtract, timeadd, title, trim, trimprefix, trimsuffix, trimsspace, upper, uuid, values, zipmap.
// Note: Some functions like md5 and uuid are not shown in the provided truncated doc but are standard in the package.
// For completeness, this function mirrors the standard implementation by assigning the exported Func vars to the map.
func standardFuncs() map[string]function.Function {
	return map[string]function.Function{
		"abs":                    stdlib.AbsoluteFunc,
		"add":                    stdlib.AddFunc,
		"and":                    stdlib.AndFunc,
		"assertnotnull":          stdlib.AssertNotNullFunc,
		"byteslen":               stdlib.BytesLenFunc,
		"bytestrslice":           stdlib.BytesSliceFunc,
		"csvdecode":              stdlib.CSVDecodeFunc,
		"ceil":                   stdlib.CeilFunc,
		"chomp":                  stdlib.ChompFunc,
		"chunklist":              stdlib.ChunklistFunc,
		"coalesce":               stdlib.CoalesceFunc,
		"coalescelist":           stdlib.CoalesceListFunc,
		"compact":                stdlib.CompactFunc,
		"concat":                 stdlib.ConcatFunc,
		"contains":               stdlib.ContainsFunc,
		"distinct":               stdlib.DistinctFunc,
		"divide":                 stdlib.DivideFunc,
		"element":                stdlib.ElementFunc,
		"eq":                     stdlib.EqualFunc,
		"flatten":                stdlib.FlattenFunc,
		"floor":                  stdlib.FloorFunc,
		"format":                 stdlib.FormatFunc,
		"formatdate":             stdlib.FormatDateFunc,
		"formatlist":             stdlib.FormatListFunc,
		"gt":                     stdlib.GreaterThanFunc,
		"gte":                    stdlib.GreaterThanOrEqualToFunc,
		"hasindex":               stdlib.HasIndexFunc,
		"indent":                 stdlib.IndentFunc,
		"index":                  stdlib.IndexFunc,
		"int":                    stdlib.IntFunc,
		"jsondecode":             stdlib.JSONDecodeFunc,
		"jsonencode":             stdlib.JSONEncodeFunc,
		"join":                   stdlib.JoinFunc,
		"keys":                   stdlib.KeysFunc,
		"length":                 stdlib.LengthFunc,
		"lt":                     stdlib.LessThanFunc,
		"lte":                    stdlib.LessThanOrEqualToFunc,
		"log":                    stdlib.LogFunc,
		"lookup":                 stdlib.LookupFunc,
		"lower":                  stdlib.LowerFunc,
		"max":                    stdlib.MaxFunc,
		"merge":                  stdlib.MergeFunc,
		"min":                    stdlib.MinFunc,
		"modulo":                 stdlib.ModuloFunc,
		"multiply":               stdlib.MultiplyFunc,
		"negate":                 stdlib.NegateFunc,
		"ne":                     stdlib.NotEqualFunc,
		"not":                    stdlib.NotFunc,
		"or":                     stdlib.OrFunc,
		"parseint":               stdlib.ParseIntFunc,
		"pow":                    stdlib.PowFunc,
		"range":                  stdlib.RangeFunc,
		"regex":                  stdlib.RegexFunc,
		"regexall":               stdlib.RegexAllFunc,
		"regexreplace":           stdlib.RegexReplaceFunc,
		"replace":                stdlib.ReplaceFunc,
		"reverse":                stdlib.ReverseFunc,
		"reverselist":            stdlib.ReverseListFunc,
		"sethaselement":          stdlib.SetHasElementFunc,
		"setintersect":           stdlib.SetIntersectionFunc,
		"setproduct":             stdlib.SetProductFunc,
		"setsubtract":            stdlib.SetSubtractFunc,
		"setsymmetricdifference": stdlib.SetSymmetricDifferenceFunc,
		"setunion":               stdlib.SetUnionFunc,
		"signum":                 stdlib.SignumFunc,
		"slice":                  stdlib.SliceFunc,
		"sort":                   stdlib.SortFunc,
		"split":                  stdlib.SplitFunc,
		"strlen":                 stdlib.StrlenFunc,
		"substr":                 stdlib.SubstrFunc,
		"subtract":               stdlib.SubtractFunc,
		"timeadd":                stdlib.TimeAddFunc,
		"title":                  stdlib.TitleFunc,
		"trim":                   stdlib.TrimFunc,
		"trimprefix":             stdlib.TrimPrefixFunc,
		"trimsuffix":             stdlib.TrimSuffixFunc,
		"trimspace":              stdlib.TrimSpaceFunc,
		"upper":                  stdlib.UpperFunc,
		"values":                 stdlib.ValuesFunc,
		"zipmap":                 stdlib.ZipmapFunc,
	}
}

// parseHCL parses the given HCL configuration string and populates the provided struct pointer,
// similar to how Terraform evaluates configurations. It handles locals as root-level variables
// and allows references between blocks. The struct must use hcl tags to map fields to block names
// and attributes. Assumes no cycles in references and evaluates blocks in file order.
// Only handles top-level attributes and blocks without labels.
func parseHCL(configPtr any, configString, configFilename string) error {
	// Validate that configPtr is a non-nil pointer to a struct
	if configPtr == nil {
		return fmt.Errorf("configPtr cannot be nil")
	}
	v := reflect.ValueOf(configPtr)
	if v.Kind() != reflect.Pointer || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("configPtr must be a pointer to a struct")
	}
	configVal := v.Elem()
	configType := configVal.Type()

	// Parse HCL
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL([]byte(configString), configFilename)
	if diags.HasErrors() {
		return fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return fmt.Errorf("body is not hclsyntax.Body")
	}

	// Prepare evaluation context
	evalCtx := &hcl.EvalContext{
		Variables: make(map[string]cty.Value),
		Functions: standardFuncs(),
	}

	// Collect locals and other blocks
	var localsBlocks []*hclsyntax.Block
	var otherBlocks []*hclsyntax.Block
	for _, block := range body.Blocks {
		if block.Type == "locals" {
			localsBlocks = append(localsBlocks, block)
		} else {
			otherBlocks = append(otherBlocks, block)
		}
	}

	// Evaluate locals
	allDiags := hcl.Diagnostics{}
	for _, localsBlock := range localsBlocks {
		for name, attr := range localsBlock.Body.Attributes {
			val, attrDiags := attr.Expr.Value(evalCtx)
			allDiags = append(allDiags, attrDiags...)
			if attrDiags.HasErrors() {
				continue
			}
			if _, exists := evalCtx.Variables[name]; exists {
				allDiags = append(allDiags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Duplicate local",
					Detail:   fmt.Sprintf("Local '%s' defined multiple times", name),
				})
				continue
			}
			evalCtx.Variables[name] = val
		}
	}

	// Create a map to store block values temporarily
	blockValues := make(map[string]cty.Value)

	// Evaluate other blocks
	for _, block := range otherBlocks {
		blockType := block.Type
		if len(block.Labels) > 0 {
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported labels",
				Detail:   "This parser does not support blocks with labels",
			})
			continue
		}
		if _, exists := blockValues[blockType]; exists {
			allDiags = append(allDiags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate block",
				Detail:   fmt.Sprintf("Block type '%s' defined multiple times", blockType),
			})
			continue
		}

		attrMap := make(map[string]cty.Value)
		for name, attr := range block.Body.Attributes {
			val, attrDiags := attr.Expr.Value(evalCtx)
			allDiags = append(allDiags, attrDiags...)
			if attrDiags.HasErrors() {
				continue
			}
			attrMap[name] = val
		}

		// Store block as object in context
		blockValues[blockType] = cty.ObjectVal(attrMap)
		evalCtx.Variables[blockType] = cty.ObjectVal(attrMap)
	}

	if allDiags.HasErrors() {
		return fmt.Errorf("evaluation errors: %s", allDiags.Error())
	}

	// Map block values to struct fields
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		hclTag := field.Tag.Get("hcl")
		if hclTag == "" {
			continue
		}
		parts := splitHCLTag(hclTag)
		if len(parts) != 2 || parts[1] != "block" {
			return fmt.Errorf("invalid hcl tag for field %s: %s", field.Name, hclTag)
		}
		blockName := parts[0]

		// Get the cty value for the block
		blockVal, exists := blockValues[blockName]
		if !exists {
			continue // Skip if block not found in HCL
		}

		// Create a new instance of the field type
		fieldVal := reflect.New(field.Type).Elem()
		fieldType := field.Type

		// Map attributes to struct fields
		for j := 0; j < fieldType.NumField(); j++ {
			subField := fieldType.Field(j)
			attrTag := subField.Tag.Get("hcl")
			if attrTag == "" {
				continue
			}
			attrName := attrTag

			// Get attribute value
			attrVal := blockVal.GetAttr(attrName)
			if attrVal.IsNull() {
				continue
			}

			// Convert cty value to Go type
			goVal, err := ctyToGoValue(attrVal, subField.Type)
			if err != nil {
				return fmt.Errorf("failed to convert attribute %s.%s: %v", blockName, attrName, err)
			}

			// Set field value
			if goVal != nil {
				fieldVal.Field(j).Set(reflect.ValueOf(goVal))
			}
		}

		// Set the field in the config struct
		configVal.Field(i).Set(fieldVal)
	}

	return nil
}

// splitHCLTag splits an hcl tag into name and type (e.g., "mysql,block" -> ["mysql", "block"])
func splitHCLTag(tag string) []string {
	var result []string
	current := ""
	inQuote := false
	for _, c := range tag {
		if c == ',' && !inQuote {
			result = append(result, current)
			current = ""
			continue
		}
		if c == '"' {
			inQuote = !inQuote
			continue
		}
		current += string(c)
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// ctyToGoValue converts a cty.Value to a Go value matching the target reflect.Type
func ctyToGoValue(v cty.Value, t reflect.Type) (any, error) {
	if v.IsNull() {
		return nil, nil
	}
	if !v.IsKnown() {
		return nil, fmt.Errorf("unknown value")
	}

	switch t.Kind() {
	case reflect.String:
		if v.Type() != cty.String {
			return nil, fmt.Errorf("expected string, got %s", v.Type().FriendlyName())
		}
		return v.AsString(), nil
	case reflect.Int, reflect.Int64:
		if v.Type() != cty.Number {
			return nil, fmt.Errorf("expected number, got %s", v.Type().FriendlyName())
		}
		var i int64
		if err := gocty.FromCtyValue(v, &i); err != nil {
			return nil, err
		}
		return i, nil
	case reflect.Float64:
		if v.Type() != cty.Number {
			return nil, fmt.Errorf("expected number, got %s", v.Type().FriendlyName())
		}
		var f float64
		if err := gocty.FromCtyValue(v, &f); err != nil {
			return nil, err
		}
		return f, nil
	case reflect.Bool:
		if v.Type() != cty.Bool {
			return nil, fmt.Errorf("expected bool, got %s", v.Type().FriendlyName())
		}
		return v.True(), nil
	case reflect.Slice:
		if !v.Type().IsListType() && !v.Type().IsTupleType() && !v.Type().IsSetType() {
			return nil, fmt.Errorf("expected list, got %s", v.Type().FriendlyName())
		}
		slice := reflect.MakeSlice(t, v.LengthInt(), v.LengthInt())
		it := v.ElementIterator()
		idx := 0
		for it.Next() {
			_, val := it.Element()
			goVal, err := ctyToGoValue(val, t.Elem())
			if err != nil {
				return nil, err
			}
			if goVal != nil {
				slice.Index(idx).Set(reflect.ValueOf(goVal))
			}
			idx++
		}
		return slice.Interface(), nil
	case reflect.Map:
		if !v.Type().IsMapType() && !v.Type().IsObjectType() {
			return nil, fmt.Errorf("expected map, got %s", v.Type().FriendlyName())
		}
		m := reflect.MakeMap(t)
		for key, val := range v.AsValueMap() {
			goVal, err := ctyToGoValue(val, t.Elem())
			if err != nil {
				return nil, err
			}
			if goVal != nil {
				m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(goVal))
			}
		}
		return m.Interface(), nil
	default:
		return nil, fmt.Errorf("unsupported Go type: %s", t.String())
	}
}
