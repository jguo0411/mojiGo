package file

import (
	"mojiGo/utils/sys"
)

//file
var CfgMap, _ = ReadConfig("cfg/config.ini")
var HYP = new(Hyp)


// mountClients Mount all the clients listed
func MountClients(cfgmap map[string]map[string]string) error {
	_ = sys.MountClient(cfgmap)
	return nil
}

//init read config file and mount targets.
func init() {
	_ = MountClients(CfgMap)
	//time.Sleep(time.Duration(2) * time.Second)
}
