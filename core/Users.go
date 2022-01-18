package core

import (
	"fmt"
	so "github.com/iamacarpet/go-win64api/shared"
	"os"
	"syscall"
	"time"
	"unsafe"
)

// Reference code: https://github.com/iamacarpet/go-win64api/blob/master/users.go

var (
	modNetapi32                = syscall.NewLazyDLL("netapi32.dll")
	usrNetUserEnum             = modNetapi32.NewProc("NetUserEnum")
	usrNetUserAdd              = modNetapi32.NewProc("NetUserAdd")
	usrNetUserDel              = modNetapi32.NewProc("NetUserDel")
	usrNetUserGetInfo          = modNetapi32.NewProc("NetUserGetInfo")
	usrNetUserSetInfo          = modNetapi32.NewProc("NetUserSetInfo")
	usrNetApiBufferFree        = modNetapi32.NewProc("NetApiBufferFree")
	usrNetLocalGroupAddMembers = modNetapi32.NewProc("NetLocalGroupAddMembers")
)

const (
	NET_API_STATUS_NERR_Success = 0
	USER_PRIV_USER              = 1
	USER_UF_SCRIPT              = 1
	USER_UF_ACCOUNTDISABLE      = 2
	USER_UF_NORMAL_ACCOUNT      = 512
	USER_UF_DONT_EXPIRE_PASSWD  = 65536
	USER_FILTER_NORMAL_ACCOUNT  = 0x0002
	USER_MAX_PREFERRED_LENGTH   = 0xFFFFFFFF
	USER_UF_LOCKOUT             = 16
	USER_UF_PASSWD_CANT_CHANGE  = 64
	USER_PRIV_ADMIN             = 2
)

type LOCALGROUP_MEMBERS_INFO_3 struct {
	Lgrmi3_domainandname *uint16
}

type USER_INFO_1 struct {
	Usri1_name         *uint16
	Usri1_password     *uint16
	Usri1_password_age uint32
	Usri1_priv         uint32
	Usri1_home_dir     *uint16
	Usri1_comment      *uint16
	Usri1_flags        uint32
	Usri1_script_path  *uint16
}

type USER_INFO_2 struct {
	Usri2_name           *uint16
	Usri2_password       *uint16
	Usri2_password_age   uint32
	Usri2_priv           uint32
	Usri2_home_dir       *uint16
	Usri2_comment        *uint16
	Usri2_flags          uint32
	Usri2_script_path    *uint16
	Usri2_auth_flags     uint32
	Usri2_full_name      *uint16
	Usri2_usr_comment    *uint16
	Usri2_parms          *uint16
	Usri2_workstations   *uint16
	Usri2_last_logon     uint32
	Usri2_last_logoff    uint32
	Usri2_acct_expires   uint32
	Usri2_max_storage    uint32
	Usri2_units_per_week uint32
	Usri2_logon_hours    uintptr
	Usri2_bad_pw_count   uint32
	Usri2_num_logons     uint32
	Usri2_logon_server   *uint16
	Usri2_country_code   uint32
	Usri2_code_page      uint32
}

type USER_INFO_1003 struct {
	Usri1003_password *uint16
}

type USER_INFO_1008 struct {
	Usri1008_flags uint32
}

type USER_INFO_1011 struct {
	Usri1011_full_name *uint16
}

// USER_INFO_1052 is the Go representation of the Windwos _USER_INFO_1052 struct
// used to set a user's profile directory.
//
// See: https://docs.microsoft.com/en-us/windows/desktop/api/lmaccess/ns-lmaccess-_user_info_1052
type USER_INFO_1052 struct {
	Useri1052_profile *uint16
}

type UserAddOptions struct {
	// Required
	Username string
	Password string

	// Optional
	FullName   string
	PrivLevel  uint32
	HomeDir    string
	Comment    string
	ScriptPath string
	UserGroup  string
}

type LocalUser struct {
	Username             string        `json:"username"`
	FullName             string        `json:"fullName"`
	IsEnabled            bool          `json:"isEnabled"`
	IsLocked             bool          `json:"isLocked"`
	IsAdmin              bool          `json:"isAdmin"`
	PasswordNeverExpires bool          `json:"passwordNeverExpires"`
	NoChangePassword     bool          `json:"noChangePassword"`
	PasswordAge          time.Duration `json:"passwordAge"`
	LastLogon            time.Time     `json:"lastLogon"`
	BadPasswordCount     uint32        `json:"badPasswordCount"`
	NumberOfLogons       uint32        `json:"numberOfLogons"`
}

// UserAddEx creates a new user account.
// As opposed to the simpler UserAdd, UserAddEx allows specification of full
// level 1 information while creating a user.
func UserAddEx(opts UserAddOptions) (bool, error) {
	var parmErr uint32
	var err error
	uInfo := USER_INFO_1{
		Usri1_priv:  opts.PrivLevel,
		Usri1_flags: USER_UF_SCRIPT | USER_UF_NORMAL_ACCOUNT | USER_UF_DONT_EXPIRE_PASSWD,
	}
	uInfo.Usri1_name, err = syscall.UTF16PtrFromString(opts.Username)
	if err != nil {
		return false, fmt.Errorf("[-] Unable to encode username to UTF16: %s\n", err)
	}
	uInfo.Usri1_password, err = syscall.UTF16PtrFromString(opts.Password)
	if err != nil {
		return false, fmt.Errorf("[-] Unable to encode password to UTF16: %s\n", err)
	}
	if opts.Comment != "" {
		uInfo.Usri1_comment, err = syscall.UTF16PtrFromString(opts.Comment)
		if err != nil {
			return false, fmt.Errorf("[-] Unable to encode comment to UTF16: %s\n", err)
		}
	}
	if opts.HomeDir != "" {
		uInfo.Usri1_home_dir, err = syscall.UTF16PtrFromString(opts.HomeDir)
		if err != nil {
			return false, fmt.Errorf("[-] Unable to encode home directory path to UTF16: %s\n", err)
		}
	}
	if opts.ScriptPath != "" {
		uInfo.Usri1_script_path, err = syscall.UTF16PtrFromString(opts.HomeDir)
		if err != nil {
			return false, fmt.Errorf("[-] Unable to encode script path to UTF16: %s\n", err)
		}
	}
	ret, _, _ := usrNetUserAdd.Call(
		uintptr(0),
		uintptr(uint32(1)),
		uintptr(unsafe.Pointer(&uInfo)),
		uintptr(unsafe.Pointer(&parmErr)),
	)
	if ret != NET_API_STATUS_NERR_Success {
		if ret == 2224 {
			fmt.Printf("[!] There is already %s user in the current system, please try another name.\n", opts.Username)
			os.Exit(3)
		} else if ret == 2245 {
			fmt.Printf("[!] The current %s password complexity is low, please use a stronger password.\n", opts.Password)
			os.Exit(3)
		} else if ret == 2202 {
			fmt.Println("[!] Invalid user name.")
			os.Exit(3)
		}
		return false, fmt.Errorf("[-] Unable to process: status=%d error=%d, Please confirm whether you have administrator privileges.\n", ret, parmErr)
	}

	fmt.Printf("[+] Successfully added %s user.\n", opts.Username)
	//common.DetermineRegistryPermissions("SAM\\SAM\\Domains\\Account\\Users\\Names\\" + opts.Username)
	return AddGroupMembership(opts.Username, opts.UserGroup)
}

// UserAdd creates a new user account with the given username, full name, and
// password.
// The new account will have the standard User privilege level.
func UserAdd(username string, password string) (bool, error) {
	return UserAddEx(UserAddOptions{
		Username:  username,
		Password:  password,
		FullName:  "",
		UserGroup: "Administrators",
		PrivLevel: USER_PRIV_USER,
	})
}

// AddGroupMembership adds the user as a member of the specified group.
func AddGroupMembership(username, groupname string) (bool, error) {
	//hn, _ := os.Hostname()
	//uPointer, err := syscall.UTF16PtrFromString(hn + `\` + username)
	uPointer, err := syscall.UTF16PtrFromString(username)
	if err != nil {
		return false, fmt.Errorf("[-] Unable to encode username to UTF16\n")
	}
	gPointer, err := syscall.UTF16PtrFromString(groupname)
	if err != nil {
		return false, fmt.Errorf("[-] unable to encode group name to UTF16\n")
	}
	var uArray = make([]LOCALGROUP_MEMBERS_INFO_3, 1)
	uArray[0] = LOCALGROUP_MEMBERS_INFO_3{
		Lgrmi3_domainandname: uPointer,
	}
	ret, _, _ := usrNetLocalGroupAddMembers.Call(
		uintptr(0),                          // servername
		uintptr(unsafe.Pointer(gPointer)),   // group name
		uintptr(uint32(3)),                  // level
		uintptr(unsafe.Pointer(&uArray[0])), // user array.
		uintptr(uint32(len(uArray))),
	)
	if ret == NET_API_STATUS_NERR_Success {
		fmt.Printf("[+] Successfully added %s user to administrator group.\n", username)
	} else if int(ret) == 1387 {
		return false, fmt.Errorf("[-] A member could not be added to or removed from the local group because the member does not exist.\n")
	} else {
		return false, fmt.Errorf("[-] unable to process, Error code: %d\n", ret)
	}
	return true, err
}

// Enable User
func EnableUser(username string) {
	eFlags, err := userGetFlags(username)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("[-] Error while getting existing flags, %s.\n", err)
	}
	eFlags &^= USER_UF_ACCOUNTDISABLE // clear bits we want to remove.
	userSetFlags(username, eFlags)
}

func userGetFlags(username string) (uint32, error) {
	var dataPointer uintptr
	uPointer, err := syscall.UTF16PtrFromString(username)
	if err != nil {
		return 0, fmt.Errorf("[-] unable to encode username to UTF16\n")
	}
	_, _, _ = usrNetUserGetInfo.Call(
		uintptr(0),                            // servername
		uintptr(unsafe.Pointer(uPointer)),     // username
		uintptr(uint32(1)),                    // level, request USER_INFO_1
		uintptr(unsafe.Pointer(&dataPointer)), // Pointer to struct.
	)
	defer usrNetApiBufferFree.Call(dataPointer)

	if dataPointer == uintptr(0) {
		return 0, fmt.Errorf("[-] unable to get data structure\n")
	}

	var data = (*USER_INFO_1)(unsafe.Pointer(dataPointer))
	return data.Usri1_flags, nil
}

func userSetFlags(username string, flags uint32) (bool, error) {
	var errParam uint32
	uPointer, err := syscall.UTF16PtrFromString(username)
	if err != nil {
		return false, fmt.Errorf("[-] Unable to encode username to UTF16\n")
	}
	ret, _, _ := usrNetUserSetInfo.Call(
		uintptr(0),                        // servername
		uintptr(unsafe.Pointer(uPointer)), // username
		uintptr(uint32(1008)),             // level
		uintptr(unsafe.Pointer(&USER_INFO_1008{Usri1008_flags: flags})),
		uintptr(unsafe.Pointer(&errParam)),
	)
	if ret != NET_API_STATUS_NERR_Success {
		return false, fmt.Errorf("[-] Unable to process. %d\n", ret)
	}
	return true, nil
}

func ListLocalUsers() ([]so.LocalUser, error) {
	var (
		dataPointer  uintptr
		resumeHandle uintptr
		entriesRead  uint32
		entriesTotal uint32
		sizeTest     USER_INFO_2
		retVal       = make([]so.LocalUser, 0)
	)

	ret, _, _ := usrNetUserEnum.Call(
		uintptr(0),         // servername
		uintptr(uint32(2)), // level, USER_INFO_2
		uintptr(uint32(USER_FILTER_NORMAL_ACCOUNT)), // filter, only "normal" accounts.
		uintptr(unsafe.Pointer(&dataPointer)),       // struct buffer for output data.
		uintptr(uint32(USER_MAX_PREFERRED_LENGTH)),  // allow as much memory as required.
		uintptr(unsafe.Pointer(&entriesRead)),
		uintptr(unsafe.Pointer(&entriesTotal)),
		uintptr(unsafe.Pointer(&resumeHandle)),
	)
	if ret != NET_API_STATUS_NERR_Success {
		return nil, fmt.Errorf("error fetching user entry")
	} else if dataPointer == uintptr(0) {
		return nil, fmt.Errorf("null pointer while fetching entry")
	}

	var iter = dataPointer
	for i := uint32(0); i < entriesRead; i++ {
		var data = (*USER_INFO_2)(unsafe.Pointer(iter))

		ud := so.LocalUser{
			Username:         UTF16toString(data.Usri2_name),
			FullName:         UTF16toString(data.Usri2_full_name),
			PasswordAge:      (time.Duration(data.Usri2_password_age) * time.Second),
			LastLogon:        time.Unix(int64(data.Usri2_last_logon), 0),
			BadPasswordCount: data.Usri2_bad_pw_count,
			NumberOfLogons:   data.Usri2_num_logons,
		}

		if (data.Usri2_flags & USER_UF_ACCOUNTDISABLE) != USER_UF_ACCOUNTDISABLE {
			ud.IsEnabled = true
		}
		if (data.Usri2_flags & USER_UF_LOCKOUT) == USER_UF_LOCKOUT {
			ud.IsLocked = true
		}
		if (data.Usri2_flags & USER_UF_PASSWD_CANT_CHANGE) == USER_UF_PASSWD_CANT_CHANGE {
			ud.NoChangePassword = true
		}
		if (data.Usri2_flags & USER_UF_DONT_EXPIRE_PASSWD) == USER_UF_DONT_EXPIRE_PASSWD {
			ud.PasswordNeverExpires = true
		}
		if data.Usri2_priv == USER_PRIV_ADMIN {
			ud.IsAdmin = true
		}

		retVal = append(retVal, ud)

		iter = uintptr(unsafe.Pointer(iter + unsafe.Sizeof(sizeTest)))
	}
	usrNetApiBufferFree.Call(dataPointer)
	return retVal, nil
}

// UTF16toString converts a pointer to a UTF16 string into a Go string.
func UTF16toString(p *uint16) string {
	return syscall.UTF16ToString((*[4096]uint16)(unsafe.Pointer(p))[:])
}
