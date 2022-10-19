package main

import (
	"github.com/appcrash/media/server/utils"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"regexp"
)

const (
	nodeSessionNodeStructName     = "SessionNode"
	nodeStubHandlerFunctionPrefix = "_convert"
	nodeHandleMessageMethodPrefix = "handle"
	nodeTraitVarNamePrefix        = "ndt"
	nodeTraitEnumPrefix           = "Nr"
)

var (
	nodeTypeInfos        []*nodeTypeInfo
	concreteNodeTypeInfo []*nodeTypeInfo // only concrete node type
	userNodeTypeInfo     []*nodeTypeInfo // only concrete node type within user package (exclude root package)

	nodeIdGen      uint16
	nodeTraitIdGen uint16
)

type nodeTypeInfo struct {
	id                 uint16
	structType         types.Object
	spec               *ast.TypeSpec
	acceptMessageTypes []*messageTypeInfo
	acceptOverridden   bool
	handlers           []string
}

func (i *nodeTypeInfo) isGeneric() bool {
	return i.spec.TypeParams.NumFields() > 0
}

func (i *nodeTypeInfo) isConcrete() bool {
	return !i.isGeneric() && i.structType.Name() != nodeSessionNodeStructName
}

func (i *nodeTypeInfo) typeName() string {
	return i.structType.Name()
}

func (i *nodeTypeInfo) snakeTypeName() string {
	return utils.CamelCaseToSnake(i.typeName())
}

func (i *nodeTypeInfo) traitVarName() string {
	return nodeTraitVarNamePrefix + i.typeName()
}

func (i *nodeTypeInfo) stubHandlerName(index int) string {
	name := i.acceptMessageTypes[index].typeName()
	return nodeStubHandlerFunctionPrefix + name
}

func (i *nodeTypeInfo) handlerName(index int) string {
	return i.handlers[index]
}

func (i *nodeTypeInfo) messageFullTypeName(index int) string {
	return i.acceptMessageTypes[index].fullTypeName()
}

type nodeTraitInterfaceInfo struct {
	id            uint16
	objectType    types.Object
	interfaceType *types.Interface
}

func (i *nodeTraitInterfaceInfo) enumName() string {
	return nodeTraitEnumPrefix + i.objectType.Name()
}

func generateNodeTrait() {
	if len(userPackage) > 0 {
		for _, p := range userPackage {
			currentGeneratingPackage = p
			inspectPackageForNode(p)
		}
	} else {
		currentGeneratingPackage = rootPackage
		inspectPackageForNode(rootPackage)
	}
	nodeEmitAll()
}

func getAnalyzeNodeInfos() []*nodeTypeInfo {
	if isGenForUser() {
		// check classes of user-defined only
		return userNodeTypeInfo
	} else {
		// check classes of predefined only
		return concreteNodeTypeInfo
	}
}

func inspectPackageForNode(p *packages.Package) {
	nodePassFindImplementer(p)
	nodePassCollectConcreteClass(p)
	nodePassFindMessageHandler(p)
	nodePassFindAcceptFunc(p)
}

func nodePassFindAcceptFunc(p *packages.Package) {
	defs := p.TypesInfo.Defs
	arrayMessageType := types.NewSlice(messageTypeInterfaceType)
	for _, info := range getAnalyzeNodeInfos() {
		decls := findClassMethodsLike(p, info.typeName(), true, regexp.MustCompile("^Accept$"))
		for _, d := range decls {
			if o, ok := defs[d.Name]; ok {
				// the check is in fact not necessary because if node define an Accept() method with another signature,
				// then it won't implement SessionAware interface as method shadowed
				sig := o.(*types.Func).Type().(*types.Signature)
				if sig.Params().Len() == 0 {
					// the Accept has no param, then see if its return type is []MessageType
					if results := sig.Results(); results.Len() == 1 {
						r := results.At(0)
						if types.Identical(r.Type(), arrayMessageType) {
							// ok, this class has overridden default Accept() method
							//log.Debugf("%v has override accept", info.typeName())
							info.acceptOverridden = true
							break
						}
					}
				}
			}
		}
	}
}

func nodePassFindMessageHandler(p *packages.Package) {
	handlerPattern := regexp.MustCompile("^" + nodeHandleMessageMethodPrefix + ".+$")
	for _, info := range getAnalyzeNodeInfos() {
		decls := findClassMethodsLike(p, info.typeName(), true, handlerPattern)
		for _, d := range decls {
			st := d.Type
			// handler function should be: no type param, no return value, only one parameter of ptr message type
			if st.TypeParams.NumFields() != 0 || st.Results.NumFields() != 0 || st.Params.NumFields() != 1 {
				continue
			}
			paramExpr := st.Params.List[0].Type
			var ident *ast.Ident
			switch expr := paramExpr.(type) {
			case *ast.StarExpr:
				// paramExpr value type must be one-level ptr type
				switch typedExpr := expr.X.(type) {
				case *ast.SelectorExpr: // in case of comp.X
					ident = typedExpr.Sel
				case *ast.Ident:
					ident = typedExpr
				}
				// as message trait is generated before node trait, all message types must be found in scope
				paramType := lookupTypeObject(ident.Name)
				if paramType != nil {
					// if the param type is resolved and its ptr type implements Message, we think the node
					// can ACCEPT this kind of message type
					ptrType := types.NewPointer(paramType.Type())
					if types.Implements(ptrType, messageInterfaceType) {
						// partially initialized msgTypeInfo, only use it to get kinds of message names
						msgTypeInfo := &messageTypeInfo{structType: paramType}
						// TODO: check duplicated handler (with same type of message)
						//log.Debugf("%v --- %v", info.typeName(), paramType)
						info.acceptMessageTypes = append(info.acceptMessageTypes, msgTypeInfo)
						info.handlers = append(info.handlers, d.Name.Name)
					}
				} else {
					log.Warnf("==> node %v method [%v]: cannot deal with parameter type %v",
						info.typeName(), d.Name.Name, ident.Name)
				}

			}
		}
	}
}

func nodePassCollectConcreteClass(p *packages.Package) {
	for _, i := range nodeTypeInfos {
		if i.isConcrete() {
			concreteNodeTypeInfo = append(concreteNodeTypeInfo, i)
			if isGenForUser() && i.structType.Pkg().Path() != rootPackageName {
				userNodeTypeInfo = append(userNodeTypeInfo, i)
			}
		}
	}
}

func nodePassFindImplementer(p *packages.Package) {
	findClassImplements(p, sessionAwareInterfaceType, func(object types.Object, ts *ast.TypeSpec) {
		nodeTypeInfos = append(nodeTypeInfos, &nodeTypeInfo{
			id:         nodeIdGen,
			structType: object,
			spec:       ts,
		})
		nodeIdGen++
	})
}
