package main

import (
	"fmt"
	"mojiGo/utils/sys"
)

// detect loop shared folder, looking for .txt file.
func detect() error {
	//source :=

	return nil
}


func main() {
	//cfgMap := file.ReadConfig(cfgPa)
	fmt.Println("=====main")
	for k, v := range sys.CfgMap["clients"]{
		fmt.Printf("%s:%s\n", k, v)
	}
	fmt.Println(sys.Test)


	//mount clients
	return

	// TODO log
}
