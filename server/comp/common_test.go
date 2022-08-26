package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func initComp() {
	logger := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat: "15:04:05",
		},
	}
	comp.InitLogger(logger)
	comp.InitBuiltIn()
}

func TestMain(m *testing.M) {
	initComp()
	initComposer()
	os.Exit(m.Run())
}
