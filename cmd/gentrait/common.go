package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
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
	nodeTraitTagInterfaceType    *types.Interface
	messageTraitTagInterfaceType *types.Interface
	messageInterfaceType         *types.Interface
	sessionAwareInterfaceType    *types.Interface
	mainPackage                  *packages.Package
	userPackage                  []*packages.Package

	currentPackageName string
)

type templateName struct {
	Name string
}

func initPackage() {
	loadConfig := new(packages.Config)
	loadConfig.Mode = loadMode
	loadConfig.Fset = token.NewFileSet()

	// delete previously generated file so that package loader wouldn't load it
	os.Remove(genFile)

	pkgs, err := packages.Load(loadConfig, rootPackageName, ".")
	//log.Printf("pkgs len is %v\n", len(pkgs))
	if err != nil {
		panic(err)
	} else {
		for _, p := range pkgs {
			if p.PkgPath == rootPackageName {
				mainPackage = p
				scope := p.Types.Scope()
				nodeTraitTagInterfaceType = scope.Lookup("NodeTraitTag").Type().Underlying().(*types.Interface)
				messageTraitTagInterfaceType = scope.Lookup("MessageTraitTag").Type().Underlying().(*types.Interface)
				messageInterfaceType = scope.Lookup("Message").Type().Underlying().(*types.Interface)
				sessionAwareInterfaceType = scope.Lookup("SessionAware").Type().Underlying().(*types.Interface)
			} else {
				userPackage = append(userPackage, p)
			}
		}
	}

	if mainPackage == nil {
		panic("cannot find necessary interface types in root package")
	}

}

func isGenForUser() bool {
	return len(userPackage) > 0
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

func findClassImplements(pkg *packages.Package, implemented *types.Interface, f func(object types.Object, ts *ast.TypeSpec)) {
	//scope := pkg.Types.Scope()
	for _, fileSyntax := range pkg.Syntax {
		//fileName := pkg.GoFiles[index]
		//log.Infof("[message] inspecting %v", fileName)
		for _, decl := range fileSyntax.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
				// check whether the ptr of this type implements this interface
				for _, spec := range gen.Specs {
					if ts, ok1 := spec.(*ast.TypeSpec); ok1 {
						var isStruct bool
						structType, _ := pkg.TypesInfo.Defs[ts.Name]
						if _, isStruct = structType.Type().Underlying().(*types.Struct); !isStruct {
							// not a struct type
							continue
						}
						ptr := types.NewPointer(structType.Type())
						if types.Implements(ptr, implemented) {
							f(structType, ts)
						}
					}
				}
			}
		}
	}
}

func findDeclaredTypeOfType[T any](pkg *packages.Package, f func(objectType types.Object, t T)) {
	for _, fileSyntax := range pkg.Syntax {
		for _, decl := range fileSyntax.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
				for _, spec := range gen.Specs {
					if ts, ok1 := spec.(*ast.TypeSpec); ok1 {
						objectType, _ := pkg.TypesInfo.Defs[ts.Name]
						if t, ok1 := objectType.Type().Underlying().(T); !ok1 {
							continue
						} else {
							f(objectType, t)
						}

					}
				}
			}
		}
	}
}
