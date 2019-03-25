package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type configData struct {
	DockerRemoteAddress string
	DockerRemotePort    string
	DockerID            string
	MySQLUserName       string `json:"mySQLUserName"`
	MySQLIP             string
	MySQLport           string
	MySQLPassword       string
	MySQLdbName         string
}

type imageInspectInfo struct {
	Id              string `json:"id"`
	Container       string `json:""container"`
	Comment         string
	Os              string
	Architecture    string
	Parent          string
	ContainerConfig ContainerConfig
	DockerVersion   string
	VirtualSize     int64
	Size            int64
	RootFS          RootFS
	//----- (possibly) invalid info: -----
	//Author
	//Created
	//GraphDriver
	//RepoDigests
	//RepoTags
	//Config
}

type ContainerConfig struct {
	Tty             bool
	Hostname        string
	Domainname      string
	AttachStdout    bool
	PublishService  string
	AttachStdin     bool
	OpenStdin       bool
	StdinOnce       bool
	NetworkDisabled bool
	OnBuild         []string
	Image           string
	User            string
	WorkingDir      string
	MacAddress      string
	AttachStderr    bool
	Lables          []string
	Env             []string
	Cmd             []string
}

type RootFS struct {
	Type   string
	Layers []string
}

type manifest struct {
	Config   string
	Layers   []string
	RepoTags string
}

type dpkgInfo struct {
	Package      string
	Source       string
	Version      string
	ValidVersion string
}

type circlCVEInfo struct {
}

func sendHTTPReq(URI string, ReqMethod string) []uint8 {
	var response *http.Response
	var err error

	switch ReqMethod {
	case "GET":
		response, err = http.Get(URI)
		errorPanic(err)
	case "POST":
		response, err = http.PostForm(URI, nil)
		errorPanic(err)
	default:
		return []uint8("Missing Request Method")
	}

	defer response.Body.Close()

	responseResult, err := ioutil.ReadAll(response.Body)
	errorPanic(err)
	return responseResult
}

func LoadConfig(path string) {
	data, err := ioutil.ReadFile(path)
	errorPanic(err)
	var c configData
	err = json.Unmarshal(data, &c)
	errorPanic(err)
	DockerRemoteAddress = c.DockerRemoteAddress
	DockerRemotePort = c.DockerRemotePort
	DockerID = getImageFullID(c.DockerID)

	MySQLUserName = c.MySQLUserName
	MySQLIP = c.MySQLIP
	MySQLport = c.MySQLport
	MySQLPassword = c.MySQLPassword
	MySQLdbName = c.MySQLdbName
}
