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

type nvdJson struct {
	CVEDataType         string     `json:"CVE_data_type"`
	CVEDataFormat       string     `json:"CVE_data_format"`
	CVEDataVersion      string     `json:"CVE_data_version"`
	CVEDataNumberOfCVEs string     `json:"CVE_data_numberOfCVEs"`
	CVEDataTimestamp    string     `json:"CVE_data_timestamp"`
	CVEItems            []CVEItems `json:"CVE_Items"`
}

type CVEDataMeta struct {
	ID       string `json:"ID"`
	ASSIGNER string `json:"ASSIGNER"`
}

type VersionData struct {
	VersionValue    string `json:"version_value"`
	VersionAffected string `json:"version_affected"`
}

type Version struct {
	VersionData []VersionData `json:"version_data"`
}

type ProductData struct {
	ProductName string  `json:"product_name"`
	Version     Version `json:"version"`
}

type Product struct {
	ProductData []ProductData `json:"product_data"`
}

type VendorData struct {
	VendorName string  `json:"vendor_name"`
	Product    Product `json:"product"`
}

type Vendor struct {
	VendorData []VendorData `json:"vendor_data"`
}

type Affects struct {
	Vendor Vendor `json:"vendor"`
}

type DescriptionData struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type Description struct {
	DescriptionData []DescriptionData `json:"description_data"`
}

type Cve struct {
	DataType    string      `json:"data_type"`
	DataFormat  string      `json:"data_format"`
	DataVersion string      `json:"data_version"`
	CVEDataMeta CVEDataMeta `json:"CVE_data_meta"`
	Affects     Affects     `json:"affects"`
	Description Description `json:"description"`
}

type CVEItems struct {
	Cve Cve `json:"cve"`
}

/*

type ProductVul struct {
	ProductName string
	AllVul      []Vul
}
*/

type ScanResult struct {
	CVEID []string
}

type Vul struct {
	CVEID         string
	AffectVersion Version
}

type AffectVersion struct {
	LEQ     bool //  AffectVerion <= x :  LEQ = true ;  AffectVersion = x : LEQ = false
	Version float32
}

type ContainersList []struct {
	ID              string          `json:"Id"`
	Names           []string        `json:"Names"`
	Image           string          `json:"Image"`
	ImageID         string          `json:"ImageID"`
	Command         string          `json:"Command"`
	Created         int             `json:"Created"`
	Ports           []Ports         `json:"Ports"`
	Labels          Labels          `json:"Labels"`
	State           string          `json:"State"`
	Status          string          `json:"Status"`
	HostConfig      HostConfig      `json:"HostConfig"`
	NetworkSettings NetworkSettings `json:"NetworkSettings"`
	Mounts          []Mounts        `json:"Mounts"`
}
type Ports struct {
	IP          string `json:"IP"`
	PrivatePort int    `json:"PrivatePort"`
	PublicPort  int    `json:"PublicPort"`
	Type        string `json:"Type"`
}
type Labels struct {
}
type HostConfig struct {
	NetworkMode string `json:"NetworkMode"`
}
type Bridge struct {
	IPAMConfig          interface{} `json:"IPAMConfig"`
	Links               interface{} `json:"Links"`
	Aliases             interface{} `json:"Aliases"`
	NetworkID           string      `json:"NetworkID"`
	EndpointID          string      `json:"EndpointID"`
	Gateway             string      `json:"Gateway"`
	IPAddress           string      `json:"IPAddress"`
	IPPrefixLen         int         `json:"IPPrefixLen"`
	IPv6Gateway         string      `json:"IPv6Gateway"`
	GlobalIPv6Address   string      `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen int         `json:"GlobalIPv6PrefixLen"`
	MacAddress          string      `json:"MacAddress"`
}
type Networks struct {
	Bridge Bridge `json:"bridge"`
}
type NetworkSettings struct {
	Networks Networks `json:"Networks"`
}
type Mounts struct {
	Type        string `json:"Type"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
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
