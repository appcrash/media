package main

import (
	"bufio"
	"os"
	"text/template"
)

var (
	nodeTempMethodStart = template.Must(template.New("method_start").Parse(
		`func (n *{{ .Name }}) {{.FuncName}}() {{.ReturnType}}{
`))
	nodeTempConfigHandlerValue = template.Must(template.New("config_handler_value").Parse(
		`    n.SetHandler({{.EnumName}},n.{{.HandlerName}})
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
)

var emitNtis []*nodeTypeInfo

func nodeEmitAll() {
	if isGenForUser() {
		emitNtis = userNodeTypeInfo
	} else {
		emitNtis = concreteNodeTypeInfo
	}
	if len(emitNtis) == 0 {
		panic("no concrete node to generate")
	}

	file, err := os.Create(genFile)
	defer func() { file.Close() }()
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(file)
	defer func() { w.Flush() }()

	emitPackage(w)

	nodeEmitMessageHandler(w)
}

func nodeEmitMessageHandler(w *bufio.Writer) {
	for _, n := range emitNtis {
		if len(n.acceptMessageTypes) == 0 {
			log.Errorf("node %v has not any message handler", n.typeName())
			continue
		}
		nodeTempMethodStart.Execute(w, struct{ Name, FuncName, ReturnType string }{n.typeName(), "ConfigHandler", ""})
		for i, mi := range n.acceptMessageTypes {
			nodeTempConfigHandlerValue.Execute(w, struct{ EnumName, HandlerName string }{
				_V(mi.enumName()), n.stubHandlerName(i),
			})
		}
		w.Write([]byte("}\n\n")) // finish function block

		for i, _ := range n.acceptMessageTypes {
			nodeTempMessageHandler.Execute(w, struct{ Name, StubHandlerName, HandlerName, ToMessage, MessageType string }{
				n.typeName(), n.stubHandlerName(i), n.handlerName(i),
				_T(toMessageFunc), n.messageFullTypeName(i),
			})
		}

		if n.acceptOverridden {
			continue
		}
		// Accept() is not explicitly overridden, generate for it

		nodeTempMethodStart.Execute(w, struct{ Name, FuncName, ReturnType string }{
			n.typeName(), "Accept", _V("[]MessageType")})
		nodeTempAcceptStart.Execute(w, templateName{_V("[]MessageType")})
		for _, mi := range n.acceptMessageTypes {
			nodeTempAcceptValue.Execute(w, templateName{_V(mi.enumName())})
		}
		w.Write([]byte("    }\n}\n\n")) // finish function block

	}
}