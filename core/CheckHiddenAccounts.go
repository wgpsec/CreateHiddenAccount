// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package core

import (
	"../common"
	"fmt"
	registry "github.com/golang/sys/windows/registry"
	"strings"
)

func CheckHiddenAccounts() {
	HiddenAccounts := CheckHiddenAccountsSub()
	if len(HiddenAccounts) > 1 {
		fmt.Println("[+] Found Hidden Account: ")
		for _, j := range HiddenAccounts {
			fmt.Println(j)
		}
	} else if len(HiddenAccounts) == 1 {
		fmt.Println("[+] Found Hidden Account: ", HiddenAccounts[0])
	} else {
		fmt.Println("[!] Not Found Hidden Account.")
	}
}

func CheckHiddenAccountsSub() []string {
	var HiddenAccounts []string
	usernames, _ := ListLocalUsers()
	for _, i := range usernames {
		if strings.Contains(i.Username, "$") {
			HiddenAccounts = append(HiddenAccounts, i.Username)
		}
	}
	if !common.DetermineDC() {
		common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account\\Users\\Names")
		key, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SAM\\SAM\\Domains\\Account\\Users\\Names", registry.ALL_ACCESS)
		defer key.Close()
		v, _ := key.ReadSubKeyNames(0)
		for _, j := range v {
			if strings.Contains(j, "$") {
				n := 0
				for _, k := range HiddenAccounts {
					if k == j {
						n++
					}
				}
				if n == 0 {
					HiddenAccounts = append(HiddenAccounts, j)
				}
			}
		}
	}
	return HiddenAccounts
}
