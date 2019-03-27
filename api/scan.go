package api

import "fmt"

func ScanPackge() {
	for k, v := range AllDpkg {
		result := CheckProductVul(k, v.ValidVersion)
		if result != nil {
			fmt.Println("[+]Detected Vul in package: ", k, " CVEList :", result)
		}

	}
}
