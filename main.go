package main

import (
	"api"
	"flag"
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
		for k, v := range api.AllDpkg {
			fmt.Printf("%s : %+v\n", k, v)
		}
	*/
	api.PullNvdCVEDB()
}
