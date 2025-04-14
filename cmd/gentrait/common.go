package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
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
	eventToMessageFunc           *types.Func

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

	// delete previously generated file recursively so that package loader wouldn't load it
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == genFile {
			log.Debugf("==> remove previously generated file %v", path)
			os.Remove(path)
		}
		return nil
	})

	// load packages and sub-packages from the current working directory
	pkgs, err := packages.Load(loadConfig, rootPackageName, "./...")
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
				eventToMessageFunc = scope.Lookup("EventToMessage").(*types.Func)

				// special case: if we are developing the root package itself, then we need to treat it as user package
				if genForRoot {
					userPackage = append(userPackage, p)
				}
			} else {
				userPackage = append(userPackage, p)
			}
		}
	}

	if rootPackage == nil {
		panic("cannot find necessary interface types in root package")
	}

}

func findConstInPackage[T any](p *packages.Package, constName string) (returnValue T) {
	var found bool
	iterateDeclaresInPackage(p, func(gen *ast.GenDecl) {
		if found {
			return
		}
		if gen.Tok == token.CONST {
			for _, s := range gen.Specs {
				spec := s.(*ast.ValueSpec)
				for _, n := range spec.Names {
					if n.Name == constName {
						if len(spec.Values) != 1 {
							panic("dont support multiple value assignment const")
						}
						if bl, ok := spec.Values[0].(*ast.BasicLit); !ok {
							panic("const value is not basic lit")
						} else {
							var cv interface{}
							var err error
							found = true
							strValue := bl.Value
							rv := reflect.ValueOf(returnValue)
							switch rv.Type().Kind() {
							case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
								reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
								cv, err = strconv.Atoi(strValue)
							case reflect.String:
								cv = strValue
							case reflect.Float32:
								cv, err = strconv.ParseFloat(strValue, 32)
							case reflect.Float64:
								cv, err = strconv.ParseFloat(strValue, 64)
							default:
								panic("find a not supported const type")
							}
							if err != nil {
								panic(err)
							}
							returnValue = reflect.ValueOf(cv).Convert(rv.Type()).Interface().(T)
						}
					}
				}
			}

		}
	})
	if !found {
		panic(fmt.Errorf("cannot find the const %v", constName))
	}
	return
}

func isGenForUser() bool {
	return !genForRoot && currentGeneratingPackage.PkgPath != rootPackageName
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

func iterateDeclaresInPackage[T *ast.FuncDecl | *ast.GenDecl](pkg *packages.Package, f func(obj T)) {
	for _, fileSyntax := range pkg.Syntax {
		for _, decl := range fileSyntax.Decls {
			if fd, ok := decl.(T); ok {
				f(fd)
			}
		}
	}
}

func findClassMethodsLike(pkg *packages.Package, structName string, isPtrReceiver bool, methodPattern *regexp.Regexp) (funcDecls []*ast.FuncDecl) {
	iterateDeclaresInPackage(pkg, func(fd *ast.FuncDecl) {
		if !methodPattern.MatchString(fd.Name.Name) {
			return
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
	})

	return
}

// find struct that implements provided interface type, NOTE: function receiver must be ptr not struct
func findClassImplements(pkg *packages.Package, implemented *types.Interface, f func(object types.Object, ts *ast.TypeSpec)) {
	iterateDeclaresInPackage(pkg, func(gen *ast.GenDecl) {
		if gen.Tok == token.TYPE {
			// check whether the ptr of this type implements this interface
			for _, spec := range gen.Specs {
				if ts, ok1 := spec.(*ast.TypeSpec); ok1 {
					var isStruct bool
					structType, _ := pkg.TypesInfo.Defs[ts.Name]
					//log.Debugf("checking strtuctType: ts:%v %v",ts.Name,structType)
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
	})
}

func findAllDeclaredTypeOfType[T *types.Interface | *types.Struct](pkg *packages.Package, f func(objectType types.Object, t T)) {
	iterateDeclaresInPackage(pkg, func(gen *ast.GenDecl) {
		if gen.Tok == token.TYPE {
			for _, spec := range gen.Specs {
				if ts, ok1 := spec.(*ast.TypeSpec); ok1 {
					objectType, _ := pkg.TypesInfo.Defs[ts.Name]
					if t, ok2 := objectType.Type().Underlying().(T); ok2 {
						f(objectType, t)
					}
				}
			}
		}
	})
}
