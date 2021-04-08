package file

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func win2mnt(winPath string, ) (string, *Err) {
	winPath = strings.Replace(strings.Replace(winPath, "\\", "/", -1), "//", "", -1)
	ip := strings.Split(winPath, "/")[0]
	replaceStr := "/" + CfgMap["clients"][ip]
	if len(replaceStr) == 1 {
		errMsg := fmt.Sprintf("IP address <%s> is not in the mount list.", ip)
		return winPath, NewError(-1, errMsg)
	}

	tempDetectPath := strings.Replace(winPath, replaceStr, "", -1)
	detectPath := filepath.Join(CfgMap["login_account"]["mntRoot"], tempDetectPath)
	if !Exist(detectPath) {
		errMsg := fmt.Sprintf("%s is not Exist.", detectPath)
		return "", NewError(-1, errMsg)
	}
	return detectPath, nil
}

//ReadConfig read .ini file. return 2D map.
func ReadConfig(path string) (map[string]map[string]string, *Err) {

	if !Exist(path){
		log.Printf("%s is missing.", path)
		os.Exit(-1)
		return nil,nil
	}
	// read file
	file, err := os.Open(path)
	if err != nil {
		errMsg :=fmt.Sprintf("read %s file failed.", path)
		return nil, NewError(-1, errMsg)
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

		if len(str) > 0 {
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
					if len(temp) != 2 {
						fmt.Printf("Wrong Format: %s\n", str)
						continue
					}
					subMap[temp[0]] = temp[1]
					cfgMap[keyList] = subMap
				}
			}
		}
	}
	return cfgMap,nil
}

//Exist return True when filename exist.
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//GetTxtFile detect watching dir path, return job txt file path.
func GetTxtFile(detectPath string) []string {
	var txtFiles []string
	_ = filepath.Walk(detectPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.Contains(path, ".txt") && !strings.Contains(path, "ERROR") {
			txtFiles = append(txtFiles, path)
		}
		return nil
	})
	return txtFiles
}

//Delete all the error files.
func DelFiles(files []string) {
	for i := 0; i < len(files); i++ {
		err := os.Remove(files[i])
		if err != nil {
			fmt.Println(err)
		}
	}
}
