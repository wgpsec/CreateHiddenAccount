// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package common

import (
	"fmt"
	registry "golang.org/x/sys/windows/registry"
	"io/ioutil"
	"os"
	"os/exec"
)

func DetermineRegistryPermissions(path string) {
	if !DetermineSAMRegistryPermissions(path) {
		fmt.Printf("[!] Access to %s registration denied\n", path)
		fmt.Println("[!] Adding registry permissions.")
		RegSetIni := "HKEY_LOCAL_MACHINE\\" + path + " [1 17]"
		ioutil.WriteFile(".N2kvMLEQiiHHNWXFpEg7uaNmcu9ic95j8.ini", []byte(RegSetIni), 0644)
		cmd := exec.Command("regini", ".N2kvMLEQiiHHNWXFpEg7uaNmcu9ic95j8.ini")
		err := cmd.Run()
		if err != nil {
			os.Remove(".N2kvMLEQiiHHNWXFpEg7uaNmcu9ic95j8.ini")
		}
		os.Remove(".N2kvMLEQiiHHNWXFpEg7uaNmcu9ic95j8.ini")
		if DetermineSAMRegistryPermissions(path) {
			fmt.Println("[+] Added registry permissions successfully.")
		} else {
			fmt.Println("[-] Failed to add registry permissions, Please confirm whether you have administrator privileges.")
			os.Exit(3)
		}
	}
}

func DetermineSAMRegistryPermissions(path string) bool {
	_, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		if err.Error() == "Access is denied." {
			return false
		} else {
			fmt.Println("[-] Registry Path: " + path)
			fmt.Printf("[-] DetermineSAMRegistryPermissions ERR: %s\n", err)
			return false
		}
	} else {
		return true
	}
}
