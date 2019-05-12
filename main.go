package main

import (
	"api"
	"flag"
	//"github.com/go-delve/delve/service/api"
)

var (
	defaultDockerRemoteAddress = "http://127.0.0.1"
	defaultDockerRemotePort    = "2356"
	defaultDockerID            = "0b1edfbffd27"
)

func main() {
	var configPath = flag.String("config", "./src/api/config.json", "sodadocker config path")
	flag.Parse()
	api.LoadConfig(*configPath)
	/*
		api.ExportImage("")

		if err := api.DecompressImage(); err != nil {
			fmt.Println(err)
		}
		api.GetImageDpkg()
		api.CheckSSH()
		api.PullNvdCVEDB()
		api.DecompressCVEDB()
		api.UnpackNVDfile()
		//	api.CheckProductVul("gcc", "0")
		api.ScanPackge()
	*/
	api.CheckSSH()
}
