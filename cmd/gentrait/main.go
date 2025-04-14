package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	genType, genFile string
	verbose          bool
	genForRoot       bool
)

var log = logrus.New()

func init() {
	flag.StringVar(&genType, "t", "", "gen node or message type")
	flag.StringVar(&genFile, "o", "", "output file")
	flag.BoolVar(&verbose, "v", false, "verbose log")
	flag.BoolVar(&genForRoot, "gen-root", false, "generate for root package")
}

func main() {
	flag.Parse()
	log.SetOutput(os.Stdout)
	if verbose {
		log.SetLevel(logrus.DebugLevel)
	}

	cwd, _ := os.Getwd()
	workingPackageName = os.Getenv("GOPACKAGE")
	log.Printf("==> gentrait in package [%v],type: %v, cwd is %v", workingPackageName, genType, cwd)

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

func generateNodeTrait() {
	if len(userPackage) > 0 {
		for _, p := range userPackage {
			currentGeneratingPackage = p
			inspectPackageForNode(p)
			nodeEmitOnePackage()
		}
	} else {
		currentGeneratingPackage = rootPackage
		inspectPackageForNode(rootPackage)
		nodeEmitOnePackage()
	}

}

func generateMessageTrait() {
	if len(userPackage) > 0 {
		for _, p := range userPackage {
			currentGeneratingPackage = p
			inspectPackageForMessage(p)
			msgEmitOnePackage()
		}
	} else {
		currentGeneratingPackage = rootPackage
		inspectPackageForMessage(rootPackage)
	}

}
