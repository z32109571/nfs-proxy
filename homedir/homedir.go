package homedir

import (
	"os"
	"runtime"
	"errors"
	"bufio"
	"syscall"
	"io"
	"strings"
	"fmt"
	"strconv"
)

var (
	ErrNoPasswdEntries = errors.New("no matching entries in passwd file")
)

type User struct {
	Name  string
	Pass  string
	Uid   int
	Gid   int
	Gecos string
	Home  string
	Shell string
}

const (
	unixPasswdPath = "/etc/passwd"
	unixGroupPath  = "/etc/group"
)

// Key returns the env var name for the user's home dir based on
// the platform being run on
func Key() string {
	if runtime.GOOS == "windows" {
		return "USERPROFILE"
	}
	return "HOME"
}

// Get returns the home directory of the current user with the help of
// environment variables depending on the target operating system.
// Returned path should be used with "path/filepath" to form new paths.
func Get() string {
	home := os.Getenv(Key())
	if home == "" && runtime.GOOS != "windows" {
		if u, err := CurrentUser(); err == nil {
			return u.Home
		}
	}
	return home
}

// GetShortcutString returns the string that is shortcut to user's home directory
// in the native shell of the platform running on.
func GetShortcutString() string {
	if runtime.GOOS == "windows" {
		return "%USERPROFILE%" // be careful while using in format functions
	}
	return "~"
}

func CurrentUser() (User, error) {
	return LookupUid(syscall.Getuid())
}

func LookupUid(uid int) (User, error) {
	return lookupUser(func(u User) bool {
		return u.Uid == uid
	})
}

func lookupUser(filter func(u User) bool) (User, error) {
	// Get operating system-specific passwd reader-closer.
	passwd, err := GetPasswd()
	if err != nil {
		return User{}, err
	}
	defer passwd.Close()

	// Get the users.
	users, err := ParsePasswdFilter(passwd, filter)
	if err != nil {
		return User{}, err
	}

	// No user entries found.
	if len(users) == 0 {
		return User{}, ErrNoPasswdEntries
	}

	// Assume the first entry is the "correct" one.
	return users[0], nil
}

func ParsePasswdFilter(r io.Reader, filter func(User) bool) ([]User, error) {
	if r == nil {
		return nil, fmt.Errorf("nil source for passwd-formatted data")
	}

	var (
		s   = bufio.NewScanner(r)
		out = []User{}
	)

	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		// see: man 5 passwd
		//  name:password:UID:GID:GECOS:directory:shell
		// Name:Pass:Uid:Gid:Gecos:Home:Shell
		//  root:x:0:0:root:/root:/bin/bash
		//  adm:x:3:4:adm:/var/adm:/bin/false
		p := User{}
		parseLine(line, &p.Name, &p.Pass, &p.Uid, &p.Gid, &p.Gecos, &p.Home, &p.Shell)

		if filter == nil || filter(p) {
			out = append(out, p)
		}
	}

	return out, nil
}

func GetPasswd() (io.ReadCloser, error) {
	return os.Open(unixPasswdPath)
}

func parseLine(line string, v ...interface{}) {
	if line == "" {
		return
	}

	parts := strings.Split(line, ":")
	for i, p := range parts {
		// Ignore cases where we don't have enough fields to populate the arguments.
		// Some configuration files like to misbehave.
		if len(v) <= i {
			break
		}

		// Use the type of the argument to figure out how to parse it, scanf() style.
		// This is legit.
		switch e := v[i].(type) {
		case *string:
			*e = p
		case *int:
			// "numbers", with conversion errors ignored because of some misbehaving configuration files.
			*e, _ = strconv.Atoi(p)
		case *[]string:
			// Comma-separated lists.
			if p != "" {
				*e = strings.Split(p, ",")
			} else {
				*e = []string{}
			}
		default:
			// Someone goof'd when writing code using this function. Scream so they can hear us.
			panic(fmt.Sprintf("parseLine only accepts {*string, *int, *[]string} as arguments! %#v is not a pointer!", e))
		}
	}
}
