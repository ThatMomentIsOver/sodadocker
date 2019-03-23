package api

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	layersName []string
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

	if _, err := os.Stat("imagesTemp"); os.IsNotExist(err) {
		err := os.Mkdir("imagesTemp", os.ModePerm)
		errorPanic(err)
	}
	imageFile, err := os.Create("imagesTemp" + "/" + imageID + ".tar")
	errorPanic(err)
	defer imageFile.Close()

	data := string(sendHTTPReq(URI, "GET"))
	reader := strings.NewReader(data)
	_, err = io.Copy(imageFile, reader)
	errorPanic(err)
}

func DecompressLayer(imageID string) error {
	filePath := "imagesTemp" + "/" + imageID + ".tar"
	fPtr, err := os.Open(filePath)
	errorPanic(err)
	tarPtr := tar.NewReader(fPtr)
	defer destructorAll()
	for {
		header, err := tarPtr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
		path := "imagesTemp" + "/" + header.Name
		layersName = append(layersName, path)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(path); os.IsNotExist(err) {
				err := os.Mkdir(path, os.ModePerm)
				if err != nil {
					return err
				}
			}
		case tar.TypeReg:
			if _, err := os.Stat(path); os.IsNotExist(err) {
				createFilePtr, err := os.Create(path)
				if err != nil {
					return err
				}
				defer createFilePtr.Close()
				_, err = io.Copy(createFilePtr, tarPtr)
				if err != nil {
					return err
				}
			}
			log.Println("Decompressing:", path)
		}
	}
	return nil
}

func ScanImage() {
	for i, lenLayersList := 0, len(layersName); i < lenLayersList; i++ {
		path := layersName[i]
		if strings.Contains(path, "manifest") {
			filePtr, err := ioutil.ReadFile(path)
			errorPanic(err)

			data := strings.Replace(string(filePtr), "\n", "", -1)
			data = strings.Trim(strings.Trim(data, "["), "]")

			var manifest_data manifest
			err = json.Unmarshal([]uint8(data), &manifest_data)
			errorPanic(err)
			fmt.Printf("%+v", manifest_data)

		}
	}
}

func destructorAll() {
	//	os.RemoveAll("imagesTemp")
}
