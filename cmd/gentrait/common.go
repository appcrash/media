package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"regexp"
)

const rootPackageName = "github.com/appcrash/media/server/comp"

const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

var (
	messageInterfaceType      *types.Interface
	sessionAwareInterfaceType *types.Interface
	mainPackage               *packages.Package
	userPackage               []*packages.Package

	currentPackageName string
)

type templateName struct {
	Name string
}

func initPackage() {
	loadConfig := new(packages.Config)
	loadConfig.Mode = loadMode
	loadConfig.Fset = token.NewFileSet()
	pkgs, err := packages.Load(loadConfig, rootPackageName, ".")
	log.Printf("pkgs len is %v\n", len(pkgs))
	if err != nil {
		panic(err)
	} else {
		for _, p := range pkgs {
			if p.PkgPath == rootPackageName {
				mainPackage = p
				messageInterfaceType = p.Types.Scope().Lookup("Message").Type().Underlying().(*types.Interface)
				sessionAwareInterfaceType = p.Types.Scope().Lookup("SessionAware").Type().Underlying().(*types.Interface)
			} else {
				userPackage = append(userPackage, p)
			}
		}
	}

	if mainPackage == nil {
		panic("cannot find Message and SessionAware types")
	}

}

func findClassMethodsLike(pkg *packages.Package, structName string, isPtrReceiver bool, methodPattern *regexp.Regexp) (funcDecls []*ast.FuncDecl) {
	for _, fileSyntax := range pkg.Syntax {
		for _, decl := range fileSyntax.Decls {
			if fd, ok := decl.(*ast.FuncDecl); ok {
				if !methodPattern.MatchString(fd.Name.Name) {
					continue
				}
				for i := 0; i < fd.Recv.NumFields(); i++ {
					expr := fd.Recv.List[i].Type
					switch expr := expr.(type) {
					case *ast.Ident:
						if !isPtrReceiver && expr.Name == structName {
							funcDecls = append(funcDecls, fd)
						}
					case *ast.StarExpr:
						if isPtrReceiver && expr.X.(*ast.Ident).Name == structName {
							funcDecls = append(funcDecls, fd)
						}
					}
				}
			}
		}
	}
	return
}
