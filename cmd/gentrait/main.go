package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
)

var (
	genType, genFile string
	verbose          bool
)

var log = logrus.New()

func init() {
	flag.StringVar(&genType, "t", "", "gen node or message type")
	flag.StringVar(&genFile, "o", "", "output file")
	flag.BoolVar(&verbose, "v", false, "verbose log")
}

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)
	if verbose {
		log.SetLevel(logrus.DebugLevel)
	}

	cwd, _ := os.Getwd()
	currentPackageName = os.Getenv("GOPACKAGE")
	log.Printf("==> gentrait in package %v, cwd is %v", currentPackageName, cwd)
	log.Printf("==> gentrait type: %v, output file: %v", genType, genFile)

	if len(genFile) == 0 {
		panic("no genFile specified")
	}

	initPackage()

	switch genType {
	case "message":
		generateMessageTrait()
	case "node":
		generateNodeTrait()
	default:
		panic("unknown genType")
	}

}

func inspectFile(pkg *packages.Package, fileName string, file *ast.File) {
	log.Printf("====> inspecting file %v\n", fileName)
	//for name, obj := range file.Scope.Objects {
	//	log.Printf("name %v  --- %v\n", name, obj.Kind)
	//}

	//messageInf := pkg.Types.Scope().Lookup("Message").Type().Underlying().(*types.Interface)
	//log.Printf("messageInf is %v\n", messageInf.String())
	for _, decl := range file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
			for _, spec := range gen.Specs {
				if ts, ok1 := spec.(*ast.TypeSpec); ok1 {
					if typ, exist := pkg.TypesInfo.Defs[ts.Name]; !exist {
						panic(fmt.Errorf("type %v not exist", ts.Name))
					} else {
						tt := typ.Type()
						ptr := types.NewPointer(tt)
						//if types.Identical(ptr, messageInf) {
						//	log.Printf("found identical: %v", tt)
						//	continue
						//}
						if types.Implements(ptr, messageInterfaceType) {
							log.Printf("%v implements message interface\n", tt)
						} else {
							//log.Printf("%v \n", tt)
						}
					}
				}
			}
		}
	}
}
