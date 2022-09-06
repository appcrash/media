package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"regexp"
)

const (
	msgEnumPrefix              = "Mt"
	msgEnumEnd                 = msgEnumPrefix + "UserMessageBegin"
	msgPostFix                 = "Message"
	msgConvertInterfacePostfix = "Convertable"
	msgConvertMethodPrefix     = "As"
)

var concreteMsgPattern = regexp.MustCompile(`^.+` + msgPostFix + `$`)
var msgTypeInfos []*messageTypeInfo
var concreteMessageTypeInfo []*messageTypeInfo
var msgTraitInfInfos []*messageTraitInterfaceInfo

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

type messageTraitInterfaceInfo struct {
	id            uint16
	interfaceType *types.Interface
}

func (i *messageTypeInfo) typeName() string {
	return i.structType.Name()
}

func (i *messageTypeInfo) enumName() string {
	return msgEnumPrefix + i.baseName()
}

func (i *messageTypeInfo) convertInterfaceName() string {
	return i.baseName() + msgConvertInterfacePostfix
}

func (i *messageTypeInfo) convertMethodName() string {
	return msgConvertMethodPrefix + i.typeName()
}

func generateMessageTrait() {
	if isGenForUser() {
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
	msgPassFindTraitInterface(pkg)
	msgPassCollectConcreteClass()
	msgPassCheckConvertable()
}

func msgPassFindImplementer(pkg *packages.Package) {
	var idGen uint16
	findClassImplements(pkg, messageInterfaceType, func(object types.Object, ts *ast.TypeSpec) {
		msgTypeInfos = append(msgTypeInfos, &messageTypeInfo{
			id:         idGen,
			structType: object,
			spec:       ts,
		})
		idGen++
	})
}

func msgPassCollectConcreteClass() {
	for _, i := range msgTypeInfos {
		if i.isConcrete() {
			concreteMessageTypeInfo = append(concreteMessageTypeInfo, i)
		}
	}
}

// find all interfaces embedding MessageTraitTag
func msgPassFindTraitInterface(pkg *packages.Package) {
	var idGen uint16
	findDeclaredTypeOfType(pkg, func(objectType types.Object, inf *types.Interface) {
		n := inf.NumEmbeddeds() - 1
		for n >= 0 {
			typ := inf.EmbeddedType(n)
			// CAVEAT: use Underlying when comparing type
			if types.Identical(messageTraitTagInterfaceType, typ.Underlying()) {
				info := &messageTraitInterfaceInfo{
					id:            idGen,
					interfaceType: inf,
				}
				msgTraitInfInfos = append(msgTraitInfInfos, info)
				idGen++
			}
			n--
		}
	})
}

// check convertibility between a message and all other messages
func msgPassCheckConvertable() {
	for _, i := range msgTypeInfos {
		if i.isConcrete() {
			log.Infof(" %v \n", i.enumName())
		}
	}
	return
}
