package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"
)

var (
	tempPackage = template.Must(template.New("package").Parse(
		`// Code generated by gentrait; DO NOT EDIT.
package {{ .Name }}

`))
	tempTraitValue = template.Must(template.New("trait_enum").Parse(
		`    {{.Name}} = uint64(1) << {{.Shift}}
`))
	tempImport = template.Must(template.New("import").Parse(`import "{{.Import}}"
`))
	msgTempEnumStart = template.Must(template.New("enum_start").Parse(
		`const (
    {{.Name}} = {{.Start}}
`))
	msgTempEnumValue = template.Must(template.New("enum").Parse(
		`    {{.Name}}
`))
	msgTempTraitShiftValue = template.Must(template.New("trait_shift").Parse(
		`    {{.Name}} = {{.Value}}
`))
	msgTempConvertableInterface = template.Must(template.New("convertable_interface").Parse(
		`type {{ .InterfaceName }} interface {
    {{ .MethodName }}() *{{ .MessageTypeName }}
}

`))
	msgTempConvertableMappingStart = template.Must(template.New("convertable_mapping_start").Parse(`func initMessageConversion() {
`))
	msgTempConvertableMappingValue = template.Must(template.New("convertable_mapping_value").Parse(
		`    m({{.From}},{{.To}})
`))
	msgTempMessageImpl = template.Must(template.New("message_impl").Parse(
		`func (m *{{.Name}}) Type() {{.MessageType}} {
    return {{.Enum}}
}

func (m *{{.Name}}) AsEvent() *event.Event {
    return event.NewEvent({{.Enum}}, m)
}

`))
	msgTempTraitFuncStart = template.Must(template.New("trait_func_start").Parse(
		`func initMessageTraits() {
    {{.Name}}(
`))
	msgTempTraitFuncValue = template.Must(template.New("trait_func_value").Parse(
		`        {{.MT}}[{{.MessageType}}]({{.MetaType}}[{{.ConvertMessageInterface}}]()),
`))
)

var emitMtis []*messageTypeInfo
var emitMriis []*messageTraitInterfaceInfo

func msgEmitAll() {
	if isGenForUser() {
		emitMtis = userMessageTypeInfo
		emitMriis = msgUserTraitInfInfos
	} else {
		emitMtis = concreteMessageTypeInfo
		emitMriis = msgTraitInfInfos
	}
	if len(emitMtis) == 0 {
		panic("no concrete message to generate")
	}

	file, err := os.Create(genFile)
	defer func() { file.Close() }()
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	defer func() {
		w.Flush()
		if formatted, err1 := format.Source(b.Bytes()); err1 != nil {
			panic(err1)
		} else {
			file.Write(formatted)
		}
	}()

	msgEmitPackage(w)

	msgEmitTraitEnum(w)

	msgEmitEnum(w)
	msgEmitConvertableInterface(w)
	msgEmitMessageImpl(w)

	msgEmitTraitFunc(w)
	msgEmitConvertableMappingFunc(w)
	msgEmitInitFunc(w)

}

func msgEmitPackage(w io.Writer) {
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

func msgEmitTraitEnum(w *bufio.Writer) {
	if len(emitMriis) == 0 {
		return
	}
	start := 0
	if isGenForUser() {
		start = findConstInPackage[int](rootPackage, msgTraitEnumShiftEnd)
	}
	w.Write([]byte("// Message Trait Enum\nconst (\n"))
	for _, traitInfo := range emitMriis {
		tempTraitValue.Execute(w, struct {
			Name  string
			Shift int
		}{
			traitInfo.enumName(), start,
		})
		start++
	}
	if !isGenForUser() {
		msgTempTraitShiftValue.Execute(w, struct {
			Name  string
			Value int
		}{
			msgTraitEnumShiftEnd, start,
		})
	}
	w.Write([]byte(")\n\n")) // finish const block
}

func msgEmitEnum(w io.Writer) {
	start := "iota"
	if isGenForUser() {
		start = "iota+" + _V(msgEnumEnd)
	}
	w.Write([]byte("// Message Type Enum\n"))
	msgTempEnumStart.Execute(w, struct{ Name, Start string }{emitMtis[0].enumName(), start})
	for _, i := range emitMtis[1:] {
		msgTempEnumValue.Execute(w, templateName{i.enumName()})
	}
	if !isGenForUser() {
		msgTempEnumValue.Execute(w, templateName{msgEnumEnd})
	}
	w.Write([]byte(")\n\n")) // finish const block
}

func msgEmitMessageImpl(w *bufio.Writer) {
	w.Write([]byte("// --------Message Implementation Begin--------\n"))
	for _, i := range emitMtis {
		msgTempMessageImpl.Execute(w, struct{ Name, Enum, MessageType string }{
			i.typeName(), i.enumName(), _V("MessageType")})
	}
	w.Write([]byte("// --------Message Implementation End--------\n\n"))
}

func msgEmitTraitFunc(w *bufio.Writer) {
	msgTempTraitFuncStart.Execute(w, templateName{_V("AddMessageTrait")})
	for _, i := range emitMtis {
		msgTempTraitFuncValue.Execute(w, struct{ MT, MetaType, MessageType, ConvertMessageInterface string }{
			_V("MT"), _V("MetaType"), i.typeName(), i.convertInterfaceName()})
	}
	w.Write([]byte(")}\n\n")) // finish function block
}

func msgEmitConvertableInterface(w io.Writer) {
	for _, i := range emitMtis {
		msgTempConvertableInterface.Execute(w, struct {
			InterfaceName, MethodName, MessageTypeName string
		}{i.convertInterfaceName(), i.convertMethodName(), i.typeName()})
	}
}

func msgEmitConvertableMappingFunc(w *bufio.Writer) {
	var nbConversion int
	msgTempConvertableMappingStart.Execute(w, nil)
	defer func() { w.Write([]byte("}\n\n")) }() // finish function block}()
	for _, i := range emitMtis {
		nbConversion += len(i.convertedTo)
	}
	if nbConversion > 0 {
		w.Write([]byte(fmt.Sprintf("m := %v", _V("SetMessageConvertable"))))
	} else {
		return
	}
	for _, i := range emitMtis {
		for _, to := range i.convertedTo {
			msgTempConvertableMappingValue.Execute(w, struct{ From, To string }{
				i.enumName(), to.enumName(),
			})
		}
	}
}

func msgEmitInitFunc(w *bufio.Writer) {
	w.Write([]byte(`
func initMessage() {
    initMessageTraits()
    initMessageConversion()
}`))
}
