package sys

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

func isMount(mntPoint string) bool {
	checkL := fmt.Sprintf("mount | grep %s", mntPoint)
	res, _ := exec.Command("sh", "-c", checkL).Output()
	if len(res) == 0 {
		return false
	}
	return true
}

// mountClient ..
func MountClient(cfgMap map[string]map[string]string) error {
	user := cfgMap["login_account"]["uname"]
	pwd := cfgMap["login_account"]["pwd"]
	mntBase := cfgMap["login_account"]["mntRoot"]

	fmt.Println("Start to connect client devices.")

	for ip, shareFolder := range cfgMap["clients"] {
		//mountPoint = path.Join("/mnt", ip)
		mntPoint := path.Join(mntBase, ip)
		_ = os.Mkdir(mntPoint, 0777)

		if isMount(mntPoint) {
			log.Printf("P: %s already mounted.\n", mntPoint)
			continue
		}
		//fmt.Println(shareFolder)
		fileSys := "//" + path.Join(ip, shareFolder)
		mlUser := fmt.Sprintf("username=%s,password=%s", user, pwd)
		cmd := exec.Command("mount", "-t", "cifs", "-o", mlUser, fileSys, mntPoint)
		//out put command res
		err := cmd.Run()
		if err != nil {
			log.Printf("Mount %s failed with %s\n", ip, err)
			log.Printf("Please try with following command:\n%s\n", cmd)
		}

		//TO check mounted
		if isMount(mntPoint) {
			log.Printf("Mounted %s\n", mntPoint)
		} else {
			log.Printf("%s is not mounted.\n", mntPoint)
		}
	}
	return nil
}
