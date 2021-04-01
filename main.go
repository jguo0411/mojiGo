package main

import (
	"flag"
	"os"
)

var (
	cfgPath string
)

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func main() {
	//var errFile []string
	flag.StringVar(&cfgPath, "p", "cfg/config.ini", "System Config file path.")
	flag.Parse()

	// TODO log
}
