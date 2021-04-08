package main

import (
	"fmt"
	"mojiGo/core"
	"mojiGo/utils/file"
	"mojiGo/utils/image"
	. "mojiGo/utils/sys"
	"strings"
	"time"
)

// TODO use log instead of print
// detect loop shared folder, looking for .txt file.
func detect() error {
	source := "dropBox"
	SignalWriter(READY)
	for {
		//fmt.Println("default image path:", file.HYP.ImgPath)
		txtFiles:=file.GetTxtFile(source)
		for i:=0;i<len(txtFiles); i++{
			if strings.Contains(txtFiles[i], "command"){
				RunCommand(txtFiles[i])
				continue
			}
			_ = file.HYP.LoadCfg(txtFiles[i])
			core.PredictProgress(file.HYP)


		}
		time.Sleep(time.Duration(2)* time.Second)
		fmt.Println("keep looping")
	}
}


func main() {
	//cfgMap := file.ReadConfig(cfgPa)
	fmt.Println("===== main =====")
	for k, v := range file.CfgMap["clients"] {
		fmt.Printf("%s:%s\n", k, v)
	}

	testC := image.GoSum(34, 5)
	fmt.Println(testC)
	_ = detect()
	return
}
