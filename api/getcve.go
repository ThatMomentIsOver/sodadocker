package api

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func PullNvdCVEDB() {
	nvdURI := "https://nvd.nist.gov/feeds/json/cve/1.0/nvdcve-1.0-%d.json.gz"

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
