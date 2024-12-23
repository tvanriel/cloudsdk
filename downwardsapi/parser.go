package downwardsapi

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-pars/pars"
)

var parseKey = pars.Many(pars.Any(pars.Lower.Map(pars.ToString),
	pars.String("/"),
	pars.String("."),
	pars.String("-"),
	pars.String("_"),
)).Map(concatValues)

var (
	parseEquals = pars.Byte('=').Map(pars.ToString)
	parseValue  = pars.AsParser(pars.Line).Map(pars.ToString).Map(maybeUnquote)
)

var parseField = pars.Seq(parseKey, parseEquals, parseValue).Map(evaluateField)

var parseComment = pars.Seq(pars.Byte('#'), pars.Line).Map(pars.ToString)

var parseLine = pars.Any(parseComment, parseField, pars.Space)

var parseDocument = pars.Many(parseLine).Map(evaluateDocument)

func Parse(s string) ([]Field, error) {
	x, err := parseDocument.Parse(pars.FromString(s))
	if err != nil {
		return nil, err
	}

	lines, ok := x.Value.([]Field)
	if !ok {
		return nil, errors.New("lines is not of type []downwardsapi.Field")
	}

	return lines, nil
}

func concatValues(p *pars.Result) error {
	var sb strings.Builder
	for i := range p.Children {
		sb.WriteString(p.Children[i].Value.(string))
	}
	p.SetValue(sb.String())

	return nil
}

func evaluateField(p *pars.Result) error {
	key := p.Children[0].Value.(string)
	value := p.Children[2].Value.(string)

	p.SetValue(Field{
		Key:   key,
		Value: value,
	})

	return nil
}

func evaluateDocument(p *pars.Result) error {
	var lines []Field
	for i := range p.Children {
		if p.Children[i].Value == nil {
			continue
		}
		name := reflect.TypeOf(p.Children[i].Value).Name()
		if name == "Field" {
			lines = append(lines, p.Children[i].Value.(Field))
		}
	}
	p.SetValue(lines)

	return nil
}

func isQuoted(s string) bool {
	if len(s) < 2 {
		return false
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		return true
	}
	return false
}

func maybeUnquote(r *pars.Result) error {
	if s, ok := r.Value.(string); ok {
		if isQuoted(s) {
			s, err := strconv.Unquote(s)
			if err != nil {
				return err
			}
			r.SetValue(s)
		}
	}
	return nil
}
