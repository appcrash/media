package main

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
)

func main() {
	cwd, _ := os.Getwd()
	packageName := os.Getenv("GOPACKAGE")
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, cwd, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		log.Fatalf("Unable to parse %s folder", cwd)
	}
	var compPkg *ast.Package
	for _, pkg := range pkgs {
		if pkg.Name == "comp" {
			compPkg = pkg
		}
	}

	conf := types.Config{Importer: importer.Default()}
	var files []*ast.File
	for _, f := range compPkg.Files {
		files = append(files, f)
	}

	if _, err := conf.Check(packageName, fset, files, nil); err != nil {
		log.Fatal(err)
	}

}

func inspectFile(file string) {

}
