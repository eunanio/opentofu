package command

import (
	"github.com/opentofu/opentofu/internal/oci"
)

type OciPushModuleCommand struct {
	Meta
}

func (c *OciPushModuleCommand) Help() string {
	return `Push an image to a registry`
}

func (c *OciPushModuleCommand) Synopsis() string {
	return "tofu push module <tag> <file>"
}

func (c *OciPushModuleCommand) Run(args []string) int {
	// Write your code here
	return oci.PushModule(args)
}
