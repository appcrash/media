package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
)

const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

func main() {
	flag.Parse()

	cwd, _ := os.Getwd()
	packageName := os.Getenv("GOPACKAGE")
	log.Printf("gentrait in package %v, cwd is %v", packageName, cwd)

	initPackage()

	for index, syn := range mainPackage.Syntax {
		inspectFile(mainPackage, mainPackage.GoFiles[index], syn)
		//log.Printf("%v  imports:\n", p.GoFiles[index])
		//for _, i := range syn.Imports {
		//
		//	log.Printf("=> %v\n", i.Path.Value)
		//}
		//for _, dec := range syn.Decls {
		//	log.Printf("%v \n", dec.Pos())
		//}

	}

	for _, p := range userPackage {
		for index, syn := range p.Syntax {
			inspectFile(p, p.GoFiles[index], syn)
			//log.Printf("%v  imports:\n", p.GoFiles[index])
			//for _, i := range syn.Imports {
			//
			//	log.Printf("=> %v\n", i.Path.Value)
			//}
			//for _, dec := range syn.Decls {
			//	log.Printf("%v \n", dec.Pos())
			//}

		}

	}
}

func inspectFile(pkg *packages.Package, fileName string, file *ast.File) {
	log.Printf("====> inspecting file %v\n", fileName)
	//for name, obj := range file.Scope.Objects {
	//	log.Printf("name %v  --- %v\n", name, obj.Kind)
	//}

	//messageInf := pkg.Types.Scope().Lookup("Message").Type().Underlying().(*types.Interface)
	//log.Printf("messageInf is %v\n", messageInf.String())
	for _, decl := range file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
			for _, spec := range gen.Specs {
				if ts, ok1 := spec.(*ast.TypeSpec); ok1 {
					if typ, exist := pkg.TypesInfo.Defs[ts.Name]; !exist {
						panic(fmt.Errorf("type %v not exist", ts.Name))
					} else {
						tt := typ.Type()
						ptr := types.NewPointer(tt)
						//if types.Identical(ptr, messageInf) {
						//	log.Printf("found identical: %v", tt)
						//	continue
						//}
						if types.Implements(ptr, messageInterfaceType) {
							log.Printf("%v implements message interface\n", tt)
						} else {
							//log.Printf("%v \n", tt)
						}
					}
				}
			}
		}
	}
}
