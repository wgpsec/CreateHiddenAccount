// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package core

import (
	"fmt"
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
	return HiddenAccounts
}
