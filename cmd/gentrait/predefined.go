package main

import (
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
)

const rootPackageName = "github.com/appcrash/media/server/comp"

var (
	messageInterfaceType      *types.Interface
	sessionAwareInterfaceType *types.Interface
	mainPackage               *packages.Package
	userPackage               []*packages.Package
)

func initPackage() {
	loadConfig := new(packages.Config)
	loadConfig.Mode = loadMode
	loadConfig.Fset = token.NewFileSet()
	pkgs, err := packages.Load(loadConfig, rootPackageName, "github.com/appcrash/media/server")
	if err != nil {
		panic(err)
	} else {
		for _, p := range pkgs {
			if p.PkgPath == rootPackageName {
				mainPackage = p
				messageInterfaceType = p.Types.Scope().Lookup("Message").Type().Underlying().(*types.Interface)
				sessionAwareInterfaceType = p.Types.Scope().Lookup("SessionAware").Type().Underlying().(*types.Interface)
			} else {
				userPackage = append(userPackage, p)
			}
		}
	}

	if mainPackage == nil {
		panic("cannot find Message and SessionAware types")
	}

}
