package cmd_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/open-identity/roles-manager/cmd"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {

	for _, c := range []struct {
		args      []string
		wait      func() bool
		expectErr bool
	}{
		{args: []string{"version"}},
	} {
		cmd.RootCmd.SetArgs(c.args)

		t.Run(fmt.Sprintf("command=%v", c.args), func(t *testing.T) {
			if c.wait != nil {
				go func() {
					assert.Nil(t, cmd.RootCmd.Execute())
				}()
			}

			if c.wait != nil {
				var count = 0
				for c.wait() {
					t.Logf("Ports are not yet open, retrying attempt #%d...", count)
					count++
					if count > 15 {
						t.FailNow()
					}
					time.Sleep(time.Second)
				}
			} else {
				err := cmd.RootCmd.Execute()
				if c.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
