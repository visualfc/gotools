package gopls

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

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
	flagServer      bool
	flagCloseServer bool
	flagRpcAddress  string
	flagGoplsPath   string
)

func init() {
	Command.Flag.BoolVar(&flagServer, "s", false, "run gopls server")
	Command.Flag.BoolVar(&flagCloseServer, "close", false, "close gopls server")
	Command.Flag.StringVar(&flagRpcAddress, "addr", "127.0.0.1:37373", "rpc server listen address")
	Command.Flag.StringVar(&flagGoplsPath, "gopls", "gopls", "gopls filepath")
}

func runGopls(cmd *command.Command, args []string) error {
	if flagServer {
		return runServer()
	} else {
		return runClient()
	}
	return nil
}

func runClient() error {
	conn, err := net.Dial("tcp", flagRpcAddress)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := &RPCClient{client: rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))}
	if flagCloseServer {
		var out ExitOut
		client.Exit(ExitIn{0}, &out)
	}
	return nil
}

func runServer() error {
	srv := NewServer(&Client{})
	err := srv.run(flagGoplsPath, "-v", "-rpc.trace")
	if err != nil {
		return err
	}
	rcvr := &RPCServer{server: srv.server, exited: make(chan int)}
	rpc.RegisterName(RPCServerName, rcvr)

	listener, err := net.Listen("tcp", flagRpcAddress)
	if err != nil {
		return err
	}
	defer listener.Close()

	conn := make(chan net.Conn)
	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					log.Println("exit")
					break
				}
				panic(err)
			}
			conn <- c
		}
	}()
	for {
		select {
		case c := <-conn:
			go rpc.ServeCodec(jsonrpc.NewServerCodec(c))
		case code := <-rcvr.exited:
			log.Println("exit", code)
			return nil
		}
	}
	return nil
}
