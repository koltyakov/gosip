package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	ent := flag.String("ent", "", "Entity struct name")
	conf := flag.Bool("conf", false, "Has Conf() method")
	mods := flag.String("mods", "", "Modifiers comma separated list")
	flag.Parse()

	if *ent == "" {
		fmt.Printf("can't generate %+v as no entity is provided, skipping...\n", os.Args)
		return
	}

	generate(*ent, *conf, strings.Split(*mods, ","))
}

func generate(entity string, conf bool, mods []string) error {
	pkgPath, _ := filepath.Abs("./")
	pkg := filepath.Base(pkgPath)

	ent := ""
	for i, l := range entity {
		pos := string(l)
		if i == 0 {
			pos = strings.ToLower(pos)
		}
		ent += pos
	}
	genFileName := fmt.Sprintf("%s_gen.go", ent)

	code := fmt.Sprintf("// Package api :: This is auto generated file, do not edit manually\n")
	code += fmt.Sprintf("package %s\n", pkg)

	if conf {
		code += `
			// Conf receives custom request config definition, e.g. custom headers, custom OData mod
			func (` + ent + ` *` + entity + `) Conf(config *RequestConfig) *` + entity + ` {
				` + ent + `.config = config
				return ` + ent + `
			}
		`
	}

	for _, mod := range mods {
		switch mod {
		case "Select":
			code += `
				// Select adds $select OData modifier
				func (` + ent + ` *` + entity + `) Select(oDataSelect string) *` + entity + ` {
					` + ent + `.modifiers.AddSelect(oDataSelect)
					return ` + ent + `
				}
			`
		case "Expand":
			code += `
				// Expand adds $expand OData modifier
				func (` + ent + ` *` + entity + `) Expand(oDataExpand string) *` + entity + ` {
					` + ent + `.modifiers.AddExpand(oDataExpand)
					return ` + ent + `
				}
			`
		case "Filter":
			code += `
				// Filter adds $filter OData modifier
				func (` + ent + ` *` + entity + `) Filter(oDataFilter string) *` + entity + ` {
					` + ent + `.modifiers.AddFilter(oDataFilter)
					return ` + ent + `
				}
			`
		case "Top":
			code += `
				// Top adds $top OData modifier
				func (` + ent + ` *` + entity + `) Top(oDataTop int) *` + entity + ` {
					` + ent + `.modifiers.AddTop(oDataTop)
					return ` + ent + `
				}
			`
		case "Skip":
			code += `
				// Skip adds $skiptoken OData modifier
				func (` + ent + ` *` + entity + `) Skip(skipToken string) *` + entity + ` {
					` + ent + `.modifiers.AddSkip(skipToken)
					return ` + ent + `
				}
			`
		case "OrderBy":
			code += `
				// OrderBy adds $orderby OData modifier
				func (` + ent + ` *` + entity + `) OrderBy(oDataOrderBy string, ascending bool) *` + entity + ` {
					` + ent + `.modifiers.AddOrderBy(oDataOrderBy, ascending)
					return ` + ent + `
				}
			`
		}
	}

	err := ioutil.WriteFile(filepath.Join("./", genFileName), []byte(code), 0644)
	return err
}
