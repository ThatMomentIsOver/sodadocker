package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func CheckSSH() {
	//check from image
	filePath := fmt.Sprintf("imagesTemp/%s/layer/etc/rc.local", TopLayerID)
	fmt.Println(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println("ssh config file not exist")
	} else {
		r, err := ioutil.ReadFile(filePath)
		errorPanic(err)
		if strings.Contains(string(r), "start") {
			fmt.Println("[+]image SSH Auto Start")
		}
	}

	// check from container
	URL := fmt.Sprintf("%s:%s/containers/json", DockerRemoteAddress, DockerRemotePort)
	req := sendHTTPReq(URL, "GET")
	var containListData ContainersList
	err := json.Unmarshal(req, &containListData)
	errorPanic(err)

	iplist := make(map[string]string, 0)

	rawImageID := "sha256:" + DockerID
	for i := range containListData {
		if containListData[i].ImageID == rawImageID {
			id := containListData[i].ID
			iplist[id] = fmt.Sprintf("%s:22", containListData[i].NetworkSettings.Networks.Bridge.IPAddress)
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	//var portScanIn, portScanOut chan string
	var wg sync.WaitGroup
	SSHopenList := make(map[string]string, 0) // docker contain id : contain ip

	portDetach := func(lockWrapper func(), containID string, ip string) {
		conn, err := net.DialTimeout("tcp", ip, time.Millisecond*3000)
		defer lockWrapper()
		if err != nil {
			return
		}

		defer conn.Close()
		result, err := bufio.NewReader(conn).ReadString('\n')
		if strings.Contains(result, "SSH") {
			SSHopenList[containID] = ip
			log.Println("[+]detach ssh open on " + ip)
			return
		}
	}

	for id, ip := range iplist {
		wg.Add(1)
		go portDetach(func() {
			wg.Done()
		}, id, ip)
	}
	wg.Wait()
}
