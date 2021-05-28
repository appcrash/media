package server

var portStart uint16 = 30000
var portEnd uint16 = 50000
var nextPort uint16 = portStart

func getNextPort() (port uint16) {
	port = nextPort
	nextPort += 2
	return
}

