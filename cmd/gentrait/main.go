package main

import (
	"flag"
	"github.com/sirupsen/logrus"
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
	workingPackageName = os.Getenv("GOPACKAGE")
	log.Printf("==> gentrait in package %v, cwd is %v", workingPackageName, cwd)
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
