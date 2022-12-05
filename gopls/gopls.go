package gopls

import (
	_ "github.com/visualfc/gotools/gopls/golang_org_x_tools_gopls/lsp/protocol"
	"github.com/visualfc/gotools/pkg/command"
)

var Command = &command.Command{
	Run:       runGopls,
	UsageLine: "gopls",
	Short:     "golang gopls util",
	Long:      `golang golsp client for Go language server`,
}

var (
	flagServer bool
)

func init() {
	Command.Flag.BoolVar(&flagServer, "s", false, "run gopls server")
}

func runGopls(cmd *command.Command, args []string) error {
	if flagServer {
		srv := NewServer(&Client{})
		err := srv.run("/Users/vfc/go/bin/gopls", "-v", "-rpc.trace")
		if err != nil {
			return err
		}
		select {}
	}
	return nil
}
