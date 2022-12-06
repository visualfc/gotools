package gopls

import (
	"context"

	"github.com/visualfc/gotools/gopls/golang_org_x_tools_gopls/lsp/protocol"
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

type RPCServer struct {
	server protocol.Server
}

func (s *RPCServer) OpenFile(in OpenFileIn, out *OpenFileOut) error {
	return s.server.DidOpen(context.Background(), &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        protocol.DocumentURI(in.URL),
			LanguageID: "go",
			Version:    1,
			Text:       in.Text,
		},
	})
}

func (s *RPCServer) CloseFile(in CloseFileIn, out *CloseFileOut) error {
	return s.server.DidClose(context.Background(), &protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.DocumentURI(in.URL),
		},
	})
}
