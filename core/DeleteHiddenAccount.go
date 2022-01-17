// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package core

import (
	"../common"
	"fmt"
	registry "github.com/golang/sys/windows/registry"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func DeleteHiddenAccount(username string) {
	if !DetermineUserExists(username) {
		fmt.Printf("[-] Not Found %s User.\n", username)
	} else {
		ret := DeleteWindowsApiUser(username)
		if ret != NET_API_STATUS_NERR_Success {
			common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account\\Users\\Names\\" + username)
			DeleteHiddenAccountRegistry(username)
			if !DetermineUserExists(username) {
				fmt.Printf("[+] Succeeded to delete %s user using registry.\n", username)
			} else {
				fmt.Printf("[-] Failed to delete %s user.\n", username)
			}
		} else {
			fmt.Printf("[+] Succeeded to Delete %s user using windows api.\n", username)
		}
	}
}

func DeleteWindowsApiUser(username string) uintptr {
	uPointer, err := syscall.UTF16PtrFromString(username)
	if err != nil {
		fmt.Printf("[-] Failed to Delete %s User, Unable to encode username to UTF16.\n", username)
	}
	ret, _, _ := usrNetUserDel.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(uPointer)),
	)
	return ret
}

func DeleteHiddenAccountRegistry(username string) {
	NamesUserNamePath := "SAM\\SAM\\Domains\\Account\\Users\\Names\\" + username
	_, UserNameIntRID := GetRegistryValue(NamesUserNamePath, "")
	UserNameHexRID := "00000" + strconv.FormatUint(uint64(UserNameIntRID), 16)
	UsersUserNamePath := "SAM\\SAM\\Domains\\Account\\Users\\" + strings.ToUpper(UserNameHexRID)
	err := registry.DeleteKey(registry.LOCAL_MACHINE, UsersUserNamePath)
	if err != nil {
		fmt.Printf("[-] %s\n", err)
	}
	err = registry.DeleteKey(registry.LOCAL_MACHINE, NamesUserNamePath)
}

func DetermineUserExists(username string) bool {
	HiddenAccounts := CheckHiddenAccountsSub()
	if len(HiddenAccounts) == 0 {
		return false
	} else {
		for _, i := range HiddenAccounts {
			if username == i {
				return true
			}
		}
	}
	return false
}
