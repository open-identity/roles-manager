package cmd

import (
	"github.com/open-identity/utils/cmdx"
)

func init() {
	RootCmd.AddCommand(cmdx.Version(&Version, &Commit, &Date))
}
