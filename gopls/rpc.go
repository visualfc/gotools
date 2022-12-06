package gopls

import (
	"context"
	"fmt"
	"net/rpc"

	"github.com/visualfc/gotools/gopls/golang_org_x_tools_gopls/lsp/protocol"
)

const (
	RPCServerName = "RPCServer"
)

type OpenFileIn struct {
	URL  string
	Text string
}

type OpenFileOut struct{}

type CloseFileIn struct {
	URL string
}

type CloseFileOut struct{}

type CompletionIn struct {
	URI       string
	Line      uint32
	Character uint32
}

type CompletionItem struct {
	Label   string
	Details string
	Kind    string
	Doc     string
}

type CompletionOut struct {
	items []*CompletionItem
}

type ExitIn struct {
	Code int
}
type ExitOut struct{}

type RPCServer struct {
	server protocol.Server
	exited chan int
}

func (s *RPCServer) OpenFile(in OpenFileIn, out *OpenFileOut) error {
	return s.server.DidOpen(context.Background(), &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        protocol.URIFromPath(in.URL),
			LanguageID: "go",
			Version:    1,
			Text:       in.Text,
		},
	})
}

func (s *RPCServer) CloseFile(in CloseFileIn, out *CloseFileOut) error {
	return s.server.DidClose(context.Background(), &protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.URIFromPath(in.URL),
		},
	})
}

func (s *RPCServer) Completion(in CompletionIn, out *CompletionOut) error {
	list, err := s.server.Completion(context.Background(), &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: protocol.URIFromPath(in.URI),
			},
			Position: protocol.Position{
				Line:      in.Line,
				Character: in.Character,
			},
		},
	})
	for _, item := range list.Items {
		out.items = append(out.items, &CompletionItem{
			Label:   item.Label,
			Details: item.Detail,
			Kind:    fmt.Sprintf("%v", item.Kind),
			Doc:     item.Documentation,
		})
	}
	return err
}

func (s *RPCServer) Exit(in ExitIn, out *ExitOut) error {
	go func() {
		s.exited <- in.Code
	}()
	return nil
}

type RPCClient struct {
	client *rpc.Client
}

func (c *RPCClient) call(method string, args interface{}, reply interface{}) error {
	return c.client.Call(RPCServerName+"."+method, args, reply)
}

func (c *RPCClient) OpenFile(in OpenFileIn, out *OpenFileOut) error {
	return c.call("OpenFile", in, out)
}

func (c *RPCClient) CloseFile(in CloseFileIn, out *CloseFileOut) error {
	return c.call("CloseFile", in, out)
}

func (c *RPCClient) Completion(in CompletionIn, out *CompletionOut) error {
	return c.call("Completion", in, out)
}

func (c *RPCClient) Exit(in ExitIn, out *ExitOut) error {
	return c.call("Exit", in, out)
}
