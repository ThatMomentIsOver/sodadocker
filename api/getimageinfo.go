package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func errorPanic(e error) {
	if e != nil {
		panic(e)
	}
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

func InspectImageLayers(domain string, port string, imageID string) *imageInspectInfo {
	URI := domain + ":" + port + "/images/" + imageID + "/json"
	st := sendHTTPReq(URI, "GET")

	var imageLayout imageInspectInfo
	err := json.Unmarshal(st, &imageLayout)
	errorPanic(err)

	return &imageLayout
}

func ExportImage(domain string, port string, imageID string) {
	URI := domain + ":" + port + "/images/" + imageID + "/get"
	// use long imageID for image tar file
	imageLayout := InspectImageLayers(domain, port, imageID)
	imageID = strings.TrimPrefix(imageLayout.Id, "sha256:")

	imageFile, err := os.Create(imageID + ".tar.gz")
	errorPanic(err)
	defer imageFile.Close()

	data := string(sendHTTPReq(URI, "GET"))
	reader := strings.NewReader(data)
	_, err = io.Copy(imageFile, reader)
	errorPanic(err)
}
