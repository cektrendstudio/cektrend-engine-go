package utpath

import (
	"os"
	"path/filepath"
	"runtime"
)

// CurrentScriptDirectory function
//
// **@Returns:** [ `$1`: cwd ]
func CurrentScriptDirectory() string {
	_, cwd, _, _ := runtime.Caller(2)
	return filepath.Dir(cwd)
}

// IsExists function
//
// **@Params:** [ `p`: path ]
//
// **@Returns:** [ `$1`: exists flag ]
func IsExists(p string) bool {
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return true
	}
	return false
}

// IsDirectory function
//
// **@Params:** [ `p`: path ]
//
// **@Returns:** [ `$1`: dir flag ]
func IsDirectory(p string) bool {
	fi, err := os.Stat(p)
	if err != nil {
		return false
	}

	mode := fi.Mode()
	return mode.IsDir()
}
