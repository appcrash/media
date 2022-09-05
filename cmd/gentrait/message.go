package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"regexp"
)

const (
	msgEnumPrefix              = "Mt"
	msgPostFix                 = "Message"
	msgConvertInterfacePostfix = "Convertable"
	msgConvertMethodPrefix     = "As"
)

var concreteMsgPattern = regexp.MustCompile(`^.+` + msgPostFix + `$`)
var mti []*messageTypeInfo

type messageTypeInfo struct {
	id         uint16
	structType types.Object
	spec       *ast.TypeSpec
}

func (i *messageTypeInfo) isGeneric() bool {
	return i.spec.TypeParams.NumFields() > 0
}

func (i *messageTypeInfo) isConcrete() bool {
	return !i.isGeneric() && concreteMsgPattern.MatchString(i.structType.Name())
}

func (i *messageTypeInfo) baseName() string {
	name := i.structType.Name()
	if !i.isConcrete() {
		panic(fmt.Errorf("%v has no enum const", name))
	}
	// prune postfix
	name = name[:len(name)-len(msgPostFix)]
	return name
}

func (i *messageTypeInfo) enumName() string {
	return msgEnumPrefix + i.baseName()
}

func (i *messageTypeInfo) convertInterfaceName() string {
	return i.baseName() + msgConvertInterfacePostfix
}

func (i *messageTypeInfo) convertMethodName() string {
	return msgConvertMethodPrefix + i.baseName()
}

func generateMessageTrait() {
	if len(userPackage) > 0 {
		for _, p := range userPackage {
			inspectPackageForMessage(p)
		}
	} else {
		inspectPackageForMessage(mainPackage)
	}
	msgEmitAll()
}

func inspectPackageForMessage(pkg *packages.Package) {
	msgPassFindImplementer(pkg)
	msgPassCheckGeneric()
}

func msgPassFindImplementer(pkg *packages.Package) {
	//scope := pkg.Types.Scope()
	for _, fileSyntax := range pkg.Syntax {
		//fileName := pkg.GoFiles[index]
		//log.Infof("[message] inspecting %v", fileName)
		var idGen uint16
		for _, decl := range fileSyntax.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
				// check whether the ptr of this type implements Message interface
				for _, spec := range gen.Specs {
					if ts, ok1 := spec.(*ast.TypeSpec); ok1 {
						var isStruct bool
						structType, _ := pkg.TypesInfo.Defs[ts.Name]
						if _, isStruct = structType.Type().Underlying().(*types.Struct); !isStruct {
							// not a struct type
							continue
						}
						ptr := types.NewPointer(structType.Type())
						if types.Implements(ptr, messageInterfaceType) {
							mti = append(mti, &messageTypeInfo{
								id:         idGen,
								structType: structType,
								spec:       ts,
							})
							idGen++
						}
					}
				}
			}
		}
	}
}

func msgPassCheckGeneric() {
	for _, i := range mti {
		if i.isConcrete() {
			log.Infof(" %v \n", i.enumName())
		}
	}
	return
}
