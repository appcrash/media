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
	messageTypeInterfaceType     *types.Named
	sessionAwareInterfaceType    *types.Interface
	toMessageFunc                *types.Func

	rootPackage              *packages.Package
	userPackage              []*packages.Package
	currentGeneratingPackage *packages.Package

	workingPackageName string // the package in which go:generate is called
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
	log.Debugf("==> total loaded packages length is %v", len(pkgs))
	if err != nil {
		panic(err)
	} else {
		for _, p := range pkgs {
			if p.PkgPath == rootPackageName {
				rootPackage = p
				scope := p.Types.Scope()
				nodeTraitTagInterfaceType = scope.Lookup("NodeTraitTag").Type().Underlying().(*types.Interface)
				messageTraitTagInterfaceType = scope.Lookup("MessageTraitTag").Type().Underlying().(*types.Interface)
				messageInterfaceType = scope.Lookup("Message").Type().Underlying().(*types.Interface)
				messageTypeInterfaceType = scope.Lookup("MessageType").Type().(*types.Named)
				sessionAwareInterfaceType = scope.Lookup("SessionAware").Type().Underlying().(*types.Interface)
				toMessageFunc = scope.Lookup("ToMessage").(*types.Func)
			} else {
				userPackage = append(userPackage, p)
			}
		}
	}

	if rootPackage == nil {
		panic("cannot find necessary interface types in root package")
	}

}

func isGenForUser() bool {
	return currentGeneratingPackage.PkgPath != rootPackageName
}

func _V(name string) string {
	if isGenForUser() {
		return "comp." + name
	}
	return name
}

func _T(obj types.Object) string {
	if obj.Pkg().Path() != currentGeneratingPackage.PkgPath {
		return obj.Pkg().Name() + "." + obj.Name()
	} else {
		return obj.Name()
	}
}

// search a type object in all packages
func lookupTypeObject(name string) types.Object {
	for _, p := range append([]*packages.Package{rootPackage}, userPackage...) {
		obj := p.Types.Scope().Lookup(name)
		if obj != nil {
			return obj
		}
	}
	return nil
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

// find struct that implements provided interface type, NOTE: function receiver must be ptr not struct
func findClassImplements(pkg *packages.Package, implemented *types.Interface, f func(object types.Object, ts *ast.TypeSpec)) {
	for _, fileSyntax := range pkg.Syntax {
		//fileName := pkg.GoFiles[index]
		//log.Debugf("====> searching interface implementer in %v", fileName)
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

func findAllDeclaredTypeOfType[T *types.Interface | *types.Struct](pkg *packages.Package, f func(objectType types.Object, t T)) {
	for _, fileSyntax := range pkg.Syntax {
		for _, decl := range fileSyntax.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
				for _, spec := range gen.Specs {
					if ts, ok1 := spec.(*ast.TypeSpec); ok1 {
						objectType, _ := pkg.TypesInfo.Defs[ts.Name]
						if t, ok2 := objectType.Type().Underlying().(T); !ok2 {
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
