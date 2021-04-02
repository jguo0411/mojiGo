package sys

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

//+-----+---+--------------------------+
//| rwx | 7 | Read, write and execute  |
//| rw- | 6 | Read, write              |
//| r-x | 5 | Read, and execute        |
//| r-- | 4 | Read,                    |
//| -wx | 3 | Write and execute        |
//| -w- | 2 | Write                    |
//| --x | 1 | Execute                  |
//| --- | 0 | no permissions           |
//+------------------------------------+
//
//+------------+------+-------+
//| Permission | Octal| Field |
//+------------+------+-------+
//| rwx------  | 0700 | User  |
//| ---rwx---  | 0070 | Group |
//| ------rwx  | 0007 | Other |
//+------------+------+-------+

// mountClient ..
func mountClient(cfgMap map[string]map[string]string) error {
	user := cfgMap["login_account"]["uname"]
	pwd := cfgMap["login_account"]["pwd"]
	mntBase := cfgMap["login_account"]["mbase"]


	fmt.Println("Start to connect client devices.")

	for ip, shareFolder := range cfgMap["clients"]{
		//mountPoint = path.Join("/mnt", ip)
		mntPoint := path.Join(mntBase, ip)
		_ = os.Mkdir(mntPoint, 0777)
		fmt.Println(shareFolder)
		fileSys := "//" + path.Join(ip, shareFolder)

		//cmd:=exec.Command("mount","-t", "cifs","-o","username=")
		mLine :=fmt.Sprintf("mount -t cifs -o username=%s,password=%s %s %s",
			user, pwd, fileSys, mntPoint)

		//fmt.Println(mLine)
		//TODO open mount
		cmd:=exec.Command(mLine)
		out, err:=cmd.CombinedOutput()
		if err != nil{
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		fmt.Printf("combined out:\n%s\n", string(out))
	}
	return nil
}

// mountClients Mount all the clients listed
func mountClients(cfgMap map[string]map[string]string) error {

	fmt.Println("from mount")
	for c, v := range cfgMap["clients"]{
		fmt.Printf("IP:%s, Point:%s\n", c, v)
	}
	_ = mountClient(CfgMap)
	Test = "ssss"
	return nil
}