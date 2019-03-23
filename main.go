package main

import (
	"api"
	"fmt"
	"strings"
)

var (
	defaultDockerRemoteAddress = "http://127.0.0.1"
	defaultDockerRemotePort    = "2356"
	defaultDockerID            = "0b1edfbffd27"
)

func main() {
	//api.ExportImage(defaultDockerRemoteAddress, defaultDockerRemotePort, defaultDockerID)
	imageLayout := api.InspectImageLayers(defaultDockerRemoteAddress, defaultDockerRemotePort, defaultDockerID)
	imageID := strings.TrimPrefix(imageLayout.Id, "sha256:")

	if err := api.DecompressLayer(imageID); err != nil {
		fmt.Println(err)
	}
}
