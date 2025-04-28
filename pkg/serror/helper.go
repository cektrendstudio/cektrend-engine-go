package serror

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utstring"
)

func IsLocal() bool {
	return strings.ToLower(utstring.Env("APP_ENV", "local")) == "local"
}

func printErr(m string) {
	fmt.Fprintln(os.Stderr, m)
}

func exit() {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		os.Exit(1)
	}
}

func getPath(val string) string {
	for _, v := range rootPaths {
		if strings.HasPrefix(val, v) {
			val = utstring.Sub(val, len(v), 0)
			return val
		}
	}

	return val
}
