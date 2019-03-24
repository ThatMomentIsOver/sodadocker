package api

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	DockerRemoteAddress string
	DockerRemotePort    string
	DockerID            string
	TopLayerID          string
	layersName          map[string][]string
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

func LoadConfig(path string) {
	data, err := ioutil.ReadFile(path)
	errorPanic(err)
	var c configData
	err = json.Unmarshal(data, &c)
	errorPanic(err)
	DockerRemoteAddress = c.DockerRemoteAddress
	DockerRemotePort = c.DockerRemotePort
	DockerID = getImageFullID(c.DockerID)
}

func InspectImageLayers(imageID string) *imageInspectInfo {
	if imageID == "" {
		imageID = DockerID
	}
	URI := DockerRemoteAddress + ":" + DockerRemotePort + "/images/" + imageID + "/json"
	st := sendHTTPReq(URI, "GET")

	var imageLayout imageInspectInfo
	err := json.Unmarshal(st, &imageLayout)
	errorPanic(err)

	return &imageLayout
}

func ExportImage(imageID string) {
	if imageID == "" {
		imageID = DockerID
	}
	URI := DockerRemoteAddress + ":" + DockerRemotePort + "/images/" + imageID + "/get"
	// use long imageID for image tar file
	imageLayout := InspectImageLayers("")
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

func DecompressImage() error {
	imageID := DockerID
	if _, imageHasDecompress := layersName["imageID"]; imageHasDecompress {
		return errors.New("image decompress before")
	}
	filePath := "imagesTemp" + "/" + imageID + ".tar"
	err := DecompressTar(filePath, "imagesTemp")
	errorPanic(err)
	return nil
}

func DecompressLayer(layerID string) error {
	filePath := "imagesTemp" + "/" + layerID + "/layer.tar"
	descPath := "imagesTemp" + "/" + TopLayerID + "/" + "layer"
	if _, err := os.Stat(descPath); os.IsNotExist(err) {
		err := os.Mkdir(descPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return DecompressTar(filePath, descPath)
}

func DecompressTar(srcPath, descPath string) error {
	fPtr, err := os.Open(srcPath)
	errorPanic(err)
	tarPtr := tar.NewReader(fPtr)
	for {
		header, err := tarPtr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}

		path := descPath + "/" + header.Name
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

func getManifestJsonData() *manifest {
	filePtr, err := ioutil.ReadFile("imagesTemp" + "/" + "manifest.json")
	errorPanic(err)

	data := strings.Replace(string(filePtr), "\n", "", -1)
	data = strings.Trim(strings.Trim(data, "["), "]")

	var manifest_data manifest
	err = json.Unmarshal([]uint8(data), &manifest_data)
	errorPanic(err)

	return &manifest_data
}

func ScanImage() {
	TopLayerID = strings.Trim(getManifestJsonData().Layers[0], "/layer.tar")
	DecompressLayer(TopLayerID)

}

func getImageFullID(short_imageID string) string {
	imageLayout := InspectImageLayers(short_imageID)
	imageID := strings.TrimPrefix(imageLayout.Id, "sha256:")
	return imageID
}

func destructorAll() {
	//	os.RemoveAll("imagesTemp")
}
