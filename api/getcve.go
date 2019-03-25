package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	nvdURI = "https://nvd.nist.gov/feeds/json/cve/1.0/nvdcve-1.0-%d.json.gz"
)

func pullNvdCVEDB() {

	if _, err := os.Stat("CVEDB"); os.IsNotExist(err) {
		err := os.Mkdir("CVEDB", os.ModePerm)
		errorPanic(err)
	}

	for year := 2002; year <= 2019; year++ {
		URI := fmt.Sprintf(nvdURI, year)
		log.Println("Get CVE DB:", URI)
		f, err := os.Create("CVEDB" + "/" + strconv.Itoa(year) + ".gz")
		errorPanic(err)

		data := string(sendHTTPReq(URI, "GET"))
		reader := strings.NewReader(data)
		_, err = io.Copy(f, reader)
		errorPanic(err)
		f.Close()
	}
}

func decompressGz(filePath, descPath string) {
	f, err := os.Open(filePath)
	errorPanic(err)
	defer f.Close()
	reader, err := gzip.NewReader(f)
	errorPanic(err)
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	errorPanic(err)
	if _, err := os.Stat(descPath); os.IsNotExist(err) {
		createFilePtr, err := os.Create(descPath)
		errorPanic(err)
		err = ioutil.WriteFile(descPath, data, os.ModePerm)
		errorPanic(err)
		createFilePtr.Close()
	} else {
		log.Println(descPath + "exist, skip it...")
	}
	errorPanic(err)

}

func DecompressCVEDB() {
	src := "CVEDB/%d.gz"
	desc := "CVEDB/%d.json"
	var wg sync.WaitGroup
	for year := 2002; year <= 2019; year++ {
		srcPath := fmt.Sprintf(src, year)
		descPath := fmt.Sprintf(desc, year)
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				log.Println("decompress" + srcPath + " done")
			}()
			decompressGz(srcPath, descPath)
		}()
	}
	wg.Wait()
}
