package api

import (
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	nvdURI        = "https://nvd.nist.gov/feeds/json/cve/1.0/nvdcve-1.0-%d.json.gz"
	MySQLUserName string
	MySQLPassword string
	MySQLIP       string
	MySQLport     string
	MySQLdbName   string
	db            *sql.DB
	AllVulmap     map[string][]Vul // key = ProductName
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
		log.Println(descPath + "compressing..")
		err = ioutil.WriteFile(descPath, data, os.ModePerm)
		errorPanic(err)
		createFilePtr.Close()
	} else {
		log.Println(descPath + " exist, skip it...")
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
			}()
			decompressGz(srcPath, descPath)
		}()
	}
	wg.Wait()
}

func ConnectSQL() {
	var err error
	dbarg := "%s:%s@tcp(%s:%s/%s?charset=utf8)"
	db, err = sql.Open("mysql", fmt.Sprintf(dbarg, MySQLUserName, MySQLPassword,
		MySQLIP, MySQLport, MySQLdbName))
	errorPanic(err)
}

func nvdJsonTrans(filePath string) *nvdJson {
	f, err := ioutil.ReadFile(filePath)
	errorPanic(err)
	var nvdf nvdJson

	err = json.Unmarshal(f, &nvdf)
	errorPanic(err)

	return &nvdf
}

func UnpackNVDfile() {
	//j := nvdJsonTrans("CVEDB/2002.json")

	var CVEID string
	AllVulmap = make(map[string][]Vul)

	for y, year := 2002, 2019; y <= year; y++ {
		jsonPath := fmt.Sprintf("CVEDB/%d.json", y)
		log.Println("loading " + jsonPath)
		j := nvdJsonTrans(jsonPath)
		for i, l := 0, len(j.CVEItems); i < l; i++ {
			CVEObject := j.CVEItems[i].Cve
			CVEID = CVEObject.CVEDataMeta.ID
			vd := CVEObject.Affects.Vendor.VendorData

			//if j.CVEItems[i].Cve.CVEDataMeta.ID == "CVE-1999-0015" {
			//	vd := j.CVEItems[i].Cve.Affects.Vendor.VendorData
			for j, m := 0, len(vd); j < m; j++ {
				pro := vd[j].Product.ProductData
				for k, lenProductData := 0, len(pro); k < lenProductData; k++ {
					ProductName := pro[k].ProductName
					AllVulmap[ProductName] = append(AllVulmap[ProductName], Vul{
						CVEID:         CVEID,
						AffectVersion: pro[k].Version,
					})
				}
				//	}
			}
		}
	}
}

func CheckProductVul(productName string, version string) []string {
	var ret []string
	for k, v := range AllVulmap {
		if k == productName {
			for i, VulCount := 0, len(v); i < VulCount; i++ {
				AffectVersion := v[i].AffectVersion.VersionData
				CVEID := v[i].CVEID
				for j, versionCount := 0, len(AffectVersion); j < versionCount; j++ {
					CVEDBinfoVersion := AffectVersion[j].VersionValue
					if CVEDBinfoVersion >= version {
						ret = append(ret, CVEID)
						break
					}
				}
			}
		}
	}
	return ret
}
