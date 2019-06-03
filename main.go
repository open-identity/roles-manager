package main

import (
	"github.com/open-identity/roles-manager/cmd"
	"github.com/open-identity/utils/profilex"
)

func main() {
	defer profilex.Profile().Stop()

	cmd.Execute()
}
