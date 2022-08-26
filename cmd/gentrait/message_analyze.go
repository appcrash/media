package main

import (
	"fmt"
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
	msgTraitEnumPrefix         = "Mr"
	msgTraitEnumShiftEnd       = "UserMessageTraitEnumShiftBegin"
)

var (
	concreteMsgPattern      = regexp.MustCompile(`^.+` + msgPostFix + `$`)
	msgTypeInfos            []*messageTypeInfo // all message type besides generic or ill-named struct types
	concreteMessageTypeInfo []*messageTypeInfo // only concrete message type
	userMessageTypeInfo     []*messageTypeInfo // only concrete message type within user package (exclude root package)

	msgTraitInfInfos     []*messageTraitInterfaceInfo // all trait interface
	msgUserTraitInfInfos []*messageTraitInterfaceInfo // only trait interface within user package (exclude root package)

	msgIdGen      uint16
	msgTraitIdGen uint16
)

type messageTypeInfo struct {
	id          uint16
	structType  types.Object
	spec        *ast.TypeSpec
	convertedTo []*messageTypeInfo // destination message types this type can convert to
}

func (i *messageTypeInfo) isGeneric() bool {
	return i.spec != nil && i.spec.TypeParams.NumFields() > 0
}

func (i *messageTypeInfo) isConcrete() bool {
	return !i.isGeneric() && concreteMsgPattern.MatchString(i.structType.Name())
}

func (i *messageTypeInfo) packagePath() string {
	return i.structType.Pkg().Path()
}

func (i *messageTypeInfo) packageName() string {
	return i.structType.Pkg().Name()
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

func (i *messageTypeInfo) typeName() string {
	return i.structType.Name()
}

func (i *messageTypeInfo) fullTypeName() string {
	// name appears in other package
	name := i.structType.Name()
	if currentGeneratingPackage.PkgPath != i.packagePath() {
		return i.packageName() + "." + name
	} else {
		return name
	}
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

type messageTraitInterfaceInfo struct {
	id            uint16
	objectType    types.Object
	interfaceType *types.Interface
}

func (i *messageTraitInterfaceInfo) enumName() string {
	return msgTraitEnumPrefix + i.objectType.Name()
}

func generateMessageTrait() {
	if len(userPackage) > 0 {
		for _, p := range userPackage {
			currentGeneratingPackage = p
			inspectPackageForMessage(p)
		}
	} else {
		currentGeneratingPackage = rootPackage
		inspectPackageForMessage(rootPackage)
	}
	msgEmitAll()
}

func inspectPackageForMessage(pkg *packages.Package) {
	msgPassFindImplementer(pkg)
	msgPassFindTraitInterface(pkg)
	msgPassCollectConcreteClass()
	msgPassAnalyzeConvertable(pkg)
}

func msgPassFindImplementer(pkg *packages.Package) {
	findClassImplements(pkg, messageInterfaceType, func(object types.Object, ts *ast.TypeSpec) {
		msgTypeInfos = append(msgTypeInfos, &messageTypeInfo{
			id:         msgIdGen,
			structType: object,
			spec:       ts,
		})
		msgIdGen++
	})
}

func msgPassCollectConcreteClass() {
	for _, i := range msgTypeInfos {
		if i.isConcrete() {
			concreteMessageTypeInfo = append(concreteMessageTypeInfo, i)
			if isGenForUser() && i.structType.Pkg().Path() != rootPackageName {
				userMessageTypeInfo = append(userMessageTypeInfo, i)
			}
		}
	}
}

// find all interfaces embedding MessageTraitTag
func msgPassFindTraitInterface(pkg *packages.Package) {
	findAllDeclaredTypeOfType(pkg, func(objectType types.Object, inf *types.Interface) {
		n := inf.NumEmbeddeds() - 1
		for n >= 0 {
			typ := inf.EmbeddedType(n)
			// CAVEAT: use Underlying when comparing type
			if types.Identical(messageTraitTagInterfaceType, typ.Underlying()) {
				info := &messageTraitInterfaceInfo{
					id:            msgTraitIdGen,
					objectType:    objectType,
					interfaceType: inf,
				}
				msgTraitInfInfos = append(msgTraitInfInfos, info)
				msgTraitIdGen++
				if info.objectType.Pkg().Path() != rootPackageName {
					msgUserTraitInfInfos = append(msgUserTraitInfInfos, info)
				}
				return
			}
			n--
		}
	})
}

// check convertibility between a message and all other messages
func msgPassAnalyzeConvertable(pkg *packages.Package) {
	var from, to []*messageTypeInfo
	to = concreteMessageTypeInfo
	scope := pkg.Types.Scope()
	if isGenForUser() {
		// check classes of [user => (user+predefined)] conversion possibility
		from = userMessageTypeInfo
	} else {
		// check classes of [predefined => predefined] conversion possibility
		from = concreteMessageTypeInfo
	}
	methodPattern := regexp.MustCompile("^" + msgConvertMethodPrefix + ".+$")
	for _, f := range from {
		decls := findClassMethodsLike(pkg, f.typeName(), true, methodPattern)
		for _, d := range decls {
			// the function should NOT have any param and only one return value of type in message type list
			ft := d.Type
			if ft.TypeParams.NumFields() != 0 || ft.Params.NumFields() != 0 || ft.Results.NumFields() != 1 {
				continue
			}
			returnExpr := ft.Results.List[0].Type
			switch expr := returnExpr.(type) {
			case *ast.StarExpr:
				// return value type must be one-level ptr type
				if ident, ok := expr.X.(*ast.Ident); ok {
					returnType := scope.Lookup(ident.Name)
					if returnType == nil {
						// not found ?
						continue
					}
					// found an eligible function with correct signature, then match from registered message types
					// it is a very slow algorithm, but I'm too lazy to improve it :)
					for _, ti := range to {
						left := ti.structType.Type()
						right := returnType.Type()

						if types.Identical(left, right) {
							log.Debugf("====> covnert signature verified:"+
								"message type %v has convert function with returned type %v",
								f.typeName(), right)
							f.convertedTo = append(f.convertedTo, ti)
						}
					}
				}

			}

		}
	}
}
