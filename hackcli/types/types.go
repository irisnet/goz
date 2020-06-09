package types

import "os"

var (
	DefaultCLIHome = os.ExpandEnv("$HOME/.hackcli")
)
