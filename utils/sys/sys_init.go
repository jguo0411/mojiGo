package sys

import (
	"mojiGo/utils/file"
	"time"
)

var Test string
var CfgMap = file.ReadConfig("cfg/config.ini")


func init() {
	_ = mountClients(CfgMap)
	time.Sleep(time.Duration(2) * time.Second)
}
