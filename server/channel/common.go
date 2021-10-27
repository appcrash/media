package channel

import "github.com/sirupsen/logrus"

var logger *logrus.Entry

func InitLogger(gl *logrus.Logger) {
	logger = gl.WithFields(logrus.Fields{"module": "channel"})
}
