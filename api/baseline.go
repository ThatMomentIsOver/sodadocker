package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func CheckSSH() {
	filePath := fmt.Sprintf("imagesTemp/%s/layer/etc/rc.local", TopLayerID)
	fmt.Println(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println("ssh config file not exist")
		return
	}
	r, err := ioutil.ReadFile(filePath)
	errorPanic(err)
	if strings.Contains(string(r), "start") {
		fmt.Println("[+]image SSH Auto Start")
	}
}
