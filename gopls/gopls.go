package gopls

import (
	"github.com/visualfc/gotools/pkg/command"
)

var Command = &command.Command{
	Run:       runGopls,
	UsageLine: "gopls",
	Short:     "golang gopls util",
	Long:      `golang golsp client for Go language server`,
}

func init() {
}

func runGopls(cmd *command.Command, args []string) error {
	return nil
}
