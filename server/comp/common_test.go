package comp_test

import (
	"github.com/appcrash/media/server/comp"
	"github.com/sirupsen/logrus"
)

func init() {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	//logger.SetLevel(logrus.DebugLevel)
	comp.InitLogger(logger)
}
