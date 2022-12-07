package main

import (
	"bufio"
	"bytes"
	"go/format"
	"io"
	"os"
	"text/template"
)

var (
	nodeTempMethodStart = template.Must(template.New("method_start").Parse(
		`func (n *{{ .Name }}) {{.FuncName}}() {{.ReturnType}}{`))

	nodeTempConfigHandlerValue = template.Must(template.New("config_handler_value").Parse(
		`n.SetMessageHandler({{.EnumName}},func(_ {{.MessageHandler}}) {{.MessageHandler}} { return n.{{.HandlerName}} })
`))
	nodeTempAcceptStart = template.Must(template.New("accept_start").Parse(
		`    return {{.Name}}{
`))
	nodeTempAcceptValue = template.Must(template.New("accept_value").Parse(
		`        {{.Name}},
`))
	nodeTempMessageHandler = template.Must(template.New("message_handler").Parse(
		`func (n *{{ .Name }}) {{.StubHandlerName}}(evt *event.Event) {
    if msg,ok := {{.ToMessage}}[*{{.MessageType}}](evt); ok {
        n.{{.HandlerName}}(msg)   
    }
}

`))
	nodeTempNewFunc = template.Must(template.New("new_node").Parse(
		`func new{{.NodeType}}() {{.SessionAware}} {
var exist bool
node := &{{.NodeType}}{}
node.Self = node
if node.Trait,exist = {{.TraitOf}}("{{.NodeSnakeName}}"); !exist {
  panic("node type {{.NodeType}} not exist")
}
{{.ConfigHandler}}
return node
}

`))
	nodeTempTraitVar = template.Must(template.New("node_trait_var").Parse(
		`{{.NT}}[{{.TypeName}}]("{{.NodeSnakeName}}",{{.FuncName}}),
`))
	nodeTempInitTraitFunc = template.Must(template.New("node_trait_init_func").Parse(
		`func initNodeTraits () {
{{.Name}} (
`))
)

var emitNtis []*nodeTypeInfo

func nodeEmitAll() {
	enableFormat := true
	if isGenForUser() {
		emitNtis = userNodeTypeInfo
	} else {
		emitNtis = concreteNodeTypeInfo
	}
	if len(emitNtis) == 0 {
		panic("no concrete node to generate")
	}

	file, err := os.Create(genFile)
	if err != nil {
		panic(err)
	}
	defer func() { file.Close() }()

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	defer func() {
		w.Flush()
		if enableFormat {
			if formatted, err1 := format.Source(b.Bytes()); err1 != nil {
				panic(err1)
			} else {
				if enableFormat {
					file.Write(formatted)
				}
			}
		} else {
			file.Write(b.Bytes())
		}

	}()

	nodeEmitPackage(w)

	nodeEmitTraitDefs(w)
	nodeEmitMessageHandler(w)
	nodeEmitFactoryFunc(w)
	nodeEmitInitFunc(w)
}

func nodeEmitPackage(w io.Writer) {
	packageName := workingPackageName
	if !isGenForUser() {
		// generate for media project itself
		packageName = "comp"
	}
	tempPackage.Execute(w, templateName{packageName})
	if packageName != "comp" {
		// user package needs import the root package
		tempImport.Execute(w, struct{ Import string }{rootPackageName})
	}
	tempImport.Execute(w, struct{ Import string }{"github.com/appcrash/media/server/event"})
}

func nodeEmitInitFunc(w *bufio.Writer) {
	w.Write([]byte(`
func initNode() {
initNodeTraits()
}`))
}

func nodeEmitTraitDefs(w *bufio.Writer) {
	nodeTempInitTraitFunc.Execute(w, templateName{_V("RegisterNodeTrait")})
	for _, ni := range emitNtis {
		nodeTempTraitVar.Execute(w, struct{ NT, TypeName, NodeSnakeName, FuncName string }{
			_V("NT"), ni.typeName(),
			ni.snakeTypeName(), "new" + ni.typeName(),
		})
	}
	w.Write([]byte(")}\n\n")) // close func block
}

func nodeEmitFactoryFunc(w *bufio.Writer) {
	w.Write([]byte("// Node Factory Method Begin\n\n"))
	for _, n := range emitNtis {
		var configHandler string
		if len(n.acceptMessageTypes) > 0 {
			configHandler = "node.configHandler()"
		}
		nodeTempNewFunc.Execute(w, struct{ NodeType, NodeSnakeName, SessionAware, TraitOf, ConfigHandler string }{
			n.typeName(), n.snakeTypeName(), _V("SessionAware"),
			_V("NodeTraitOfType"), configHandler,
		})
	}
	w.Write([]byte("// Node Factory Method End\n\n"))
}

func nodeEmitMessageHandler(w *bufio.Writer) {
	for _, n := range emitNtis {
		if len(n.acceptMessageTypes) == 0 {
			log.Warnf("==> node %v has not any message handler", n.typeName())
			continue
		}
		nodeTempMethodStart.Execute(w, struct{ Name, FuncName, ReturnType string }{n.typeName(), "configHandler", ""})
		for i, mi := range n.acceptMessageTypes {
			nodeTempConfigHandlerValue.Execute(w, struct{ EnumName, MessageHandler, HandlerName string }{
				mi.enumName(), _V("MessageHandler"), n.stubHandlerName(i),
			})
		}
		w.Write([]byte("}\n\n")) // finish function block

		for i, _ := range n.acceptMessageTypes {
			nodeTempMessageHandler.Execute(w, struct{ Name, StubHandlerName, HandlerName, ToMessage, MessageType string }{
				n.typeName(), n.stubHandlerName(i), n.handlerName(i),
				_T(eventToMessageFunc), n.messageFullTypeName(i),
			})
		}

		if n.acceptOverridden {
			continue
		}
		// Accept() is not explicitly overridden, generate for it

		if isGenForUser() {
			nodeTempMethodStart.Execute(w, struct{ Name, FuncName, ReturnType string }{
				n.typeName(), "Accept", "[]comp.MessageType"})
			nodeTempAcceptStart.Execute(w, templateName{"[]comp.MessageType"})
		} else {
			nodeTempMethodStart.Execute(w, struct{ Name, FuncName, ReturnType string }{
				n.typeName(), "Accept", "[]MessageType"})
			nodeTempAcceptStart.Execute(w, templateName{"[]MessageType"})
		}

		for _, mi := range n.acceptMessageTypes {
			nodeTempAcceptValue.Execute(w, templateName{mi.enumName()})
		}
		w.Write([]byte("    }\n}\n\n")) // finish function block

	}
}
