package gopls

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/visualfc/gotools/gopls/golang_org_x_tools/fakenet"
	"github.com/visualfc/gotools/gopls/golang_org_x_tools/jsonrpc2"
	"github.com/visualfc/gotools/gopls/golang_org_x_tools_gopls/lsp/protocol"
)

type Server struct {
	cmd    *exec.Cmd
	server protocol.Server
	client protocol.Client
	cancel context.CancelFunc
}

func NewServer(client protocol.Client) *Server {
	return &Server{client: client}
}

func (s *Server) run(goplscmd string, args ...string) error {
	cmd := exec.Command(goplscmd, args...)
	cmd.Env = os.Environ()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe for gopls: %v", err)
	}
	go func() error {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Printf("gopls stderr: %v\n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("reading standard input: %v", err)
		}
		return nil
	}()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe for gopls: %v", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe for gopls: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start gopls: %v", err)
	}
	go func() (err error) {
		if err = cmd.Wait(); err != nil {
			err = fmt.Errorf("got error running gopls: %v", err)
		}
		return nil
	}()

	fakeconn := fakenet.NewConn("stdio", stdout, stdin)
	stream := jsonrpc2.NewHeaderStream(fakeconn)
	ctxt, cancel := context.WithCancel(context.Background())
	conn := jsonrpc2.NewConn(stream)
	server := protocol.ServerDispatcher(conn)

	handler := protocol.ClientHandler(s.client, jsonrpc2.MethodNotFound)
	handler = protocol.Handlers(handler)
	ctxt = protocol.WithClient(ctxt, s.client)

	go func() {
		conn.Go(ctxt, handler)
		<-conn.Done()
	}()
	_, err = server.Initialize(context.Background(), &protocol.ParamInitialize{})
	if err != nil {
		return err
	}
	err = server.Initialized(context.Background(), &protocol.InitializedParams{})
	if err != nil {
		return err
	}
	s.server = server
	s.cancel = cancel
	s.cmd = cmd
	return nil
}
