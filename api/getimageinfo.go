package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	defaultDockerRemoteAddress = "127.0.0.1"
	defaultDockerRemotePort    = "2346"
)

func errorPanic(e error) {
	if e != nil {
		panic(e)
	}
}

func sendHTTPReq(domain string, port string, path string, ReqMethod string) string {
	URI := domain + port + path

	var response *http.Response
	var err error

	switch ReqMethod {
	case "GET":
		response, err = http.Get(URI)
		errorPanic(err)
	case "POST":
		postValue, _ := url.ParseQuery(path)
		errorPanic(err)
		response, err = http.PostForm(URI, postValue)
		errorPanic(err)
	default:
		return "Missing Request Method"
	}

	defer response.Body.Close()

	responseResult, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("read error")
	}

	return string(responseResult)
}

func InspectImage(imageId string) {
	fmt.Println(imageId)
}
