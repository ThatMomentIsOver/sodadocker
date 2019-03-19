package main

import (
	"api"
)

var (
	defaultDockerRemoteAddress = "http://127.0.0.1"
	defaultDockerRemotePort    = "2356"
	defaultDockerID            = "0b1edfbffd27"
)

func main() {
	api.ExportImage(defaultDockerRemoteAddress, defaultDockerRemotePort, defaultDockerID)
}
