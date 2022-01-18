// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package common

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func DetermineDC() bool {
	hostname, _ := os.Hostname()
	hostname = strings.ToUpper(hostname)
	cmd := exec.Command("netdom", "query", "pdc")
	stdout, _ := cmd.StdoutPipe()
	defer stdout.Close()
	err := cmd.Start()
	if err != nil {
		return false
	} else {
		opBytes, _ := ioutil.ReadAll(stdout)
		dcname := strings.Split(string(opBytes), "\n")[2]
		if hostname == dcname {
			return true
		} else {
			return false
		}
	}
}
