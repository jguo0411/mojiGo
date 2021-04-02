package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//ReadConfig ...
func ReadConfig(path string) map[string]map[string]string {

	if !Exist(path) {
		return nil
	}
	// read file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("read init file failed.")
		return nil
	}
	//
	defer func() {
		err := file.Close()
		if err == nil {
			return
		}
	}()

	reader := bufio.NewReader(file)
	cfgMap := make(map[string]map[string]string)
	subMap := make(map[string]string)
	var keyList string

	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		str = strings.Replace(str, " ", "", -1)
		str = strings.Replace(str, "\n", "", -1)

		if len(str)>0 {
			switch str[0] {
			case ';', '#':
				continue
			case '[':
				keyList = strings.Replace(str[1:], "]", "", -1)
				subMap = make(map[string]string)
			default:
				if keyList != "" {
					temp := strings.Split(str, ":")
					if len(temp) != 2 {
						temp = strings.Split(str, "=")
					}
					if len(temp) !=2{
						fmt.Printf("Wrong Format: %s\n", str)
						continue
					}
					subMap[temp[0]] = temp[1]
					cfgMap[keyList] = subMap
				}
			}
		}
	}
	return cfgMap
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
