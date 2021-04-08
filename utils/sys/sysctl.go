package sys

import (
	"bufio"
	"log"
	"os"
	"strings"
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
const (
	READY     byte = 1
	PMW_OFF   byte = 2
	REBOOT    byte = 3
	ERROR     byte = 4
	UNKNOWN   byte = 5
	ctlSignal      = "/home/dac/Desktop/ctl_signal"
)

func SignalWriter(comm uint8, s ...string) {
	file, _ := os.OpenFile(ctlSignal, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	writer := bufio.NewWriter(file)
	cmStr := ""
	if comm == UNKNOWN {
		cmStr = s[0]
	}
	if comm == READY || strings.HasPrefix(cmStr, "ready") {
		writer.WriteString("ready")
		writer.Flush()
	} else if comm == PMW_OFF || strings.HasPrefix(cmStr, "shutdown") {
		writer.WriteString("shutdown")
		writer.Flush()
		//os.Exit(0)
		// TODO end program
	} else if comm == REBOOT || strings.HasPrefix(cmStr, "reboot") {
		writer.WriteString("reboot")
		writer.Flush()
		//os.Exit(0)
		// TODO end program
	} else if comm == ERROR {
		writer.WriteString("!CHK LOG!")
		writer.Flush()
	}
}

func RunCommand(fp string) {
	file, err := os.Open(fp)
	defer file.Close()
	var buffer = make([]byte, 64)

	if err != nil {
		log.Printf("Open %s failed.", fp)
	}
	reader := bufio.NewReader(file)
	nBytes, _ := reader.Read(buffer)
	command := string(buffer[:nBytes])
	SignalWriter(UNKNOWN, command)

	// Delete commond file.
	_ = os.Remove(fp)
}
