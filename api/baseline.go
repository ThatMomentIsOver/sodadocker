package api

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func CheckSSH() {
	filePath := fmt.Sprintf("imagesTemp/%s/layer/etc/rc.local", TopLayerID)
	r, err := ioutil.ReadFile(filePath)
	errorPanic(err)
	if strings.Contains(string(r), "start") {
		fmt.Println("[+]image SSH Auto Start")
	}
}
