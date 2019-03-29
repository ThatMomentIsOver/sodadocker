package api

import (
	"archive/tar"
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	DockerRemoteAddress string
	DockerRemotePort    string
	DockerID            string
	TopLayerID          string
	layersName          map[string][]string
	AllDpkg             map[string]dpkgInfo
)

func errorPanic(e error) {
	if e != nil {
		panic(e)
	}
}

func errorFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
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
		log.Println("Decompressing:", path)
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
				_, err = io.Copy(createFilePtr, tarPtr)
				if err != nil {
					return err
				}
				createFilePtr.Close()
			}
		}
	}
	fPtr.Close()
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

func GetImageDpkg() {
	TopLayerID = strings.Trim(getManifestJsonData().Layers[0], "/layer.tar")
	path := "imagesTemp" + "/" + TopLayerID + "/" + "layer" +
		"/var/lib/dpkg/status"
	errd := DecompressLayer(TopLayerID)
	if errd != nil {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			//TODO: [too many open files] error causes no "dpkg status" file to be obtained
			panic(errd)
		}
	}
	file, err := os.Open(path)
	errorPanic(err)
	defer file.Close()

	dpkgList := make(map[string]dpkgInfo)
	scanner := bufio.NewScanner(file)
	var package_name, source_name string
	var validVersion string
	//package_flag := false
	//source_flag := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Package") {
			//package_flag = true
			package_name = strings.TrimPrefix(line, "Package: ")
			source_name = ""
		} else if strings.HasPrefix(line, "Source") {
			//source_flag = true
			source_name = strings.TrimPrefix(line, "Source: ")
		} else if strings.HasPrefix(line, "Version") {
			version := strings.TrimPrefix(line, "Version: ")
			less := strings.Index(version, "-")
			plus := strings.Index(version, "+")
			if less != -1 {
				validVersion = string(version[:less])
			}
			if plus != -1 {
				validVersion = string(version[:plus])
			}

			if source_name != "" {
				dpkgList[package_name] = dpkgInfo{Package: package_name, Source: source_name, Version: version, ValidVersion: validVersion}
			} else {
				dpkgList[package_name] = dpkgInfo{Package: package_name, Version: version, ValidVersion: validVersion}
			}
			/*
				if package_flag && source_flag {
					dpkgList[package_name] = dpkgInfo{Package: package_name, Source: source_name, Version: version}
				} else if package_flag {
					dpkgList[package_name] = dpkgInfo{Package: package_name, Version: version}
				}
				package_flag = false
				source_flag = false
			*/
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded dpkg done...")
	AllDpkg = dpkgList

}

func getImageFullID(short_imageID string) string {
	imageLayout := InspectImageLayers(short_imageID)
	imageID := strings.TrimPrefix(imageLayout.Id, "sha256:")
	return imageID
}

func destructorAll() {
	os.RemoveAll("imagesTemp")
	os.RemoveAll("CVEDB")
}
