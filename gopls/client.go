package gopls

import (
	"context"

	"github.com/visualfc/gotools/gopls/golang_org_x_tools_gopls/lsp/protocol"
)

type Client struct {
}

var _ protocol.Client = &Client{}

func (c *Client) LogTrace(context.Context, *protocol.LogTraceParams) error {
	return nil
}
func (c *Client) Progress(context.Context, *protocol.ProgressParams) error {
	return nil
}
func (c *Client) RegisterCapability(context.Context, *protocol.RegistrationParams) error {
	return nil
}
func (c *Client) UnregisterCapability(context.Context, *protocol.UnregistrationParams) error {
	return nil
}
func (c *Client) Event(context.Context, *interface{}) error {
	return nil
}
func (c *Client) PublishDiagnostics(context.Context, *protocol.PublishDiagnosticsParams) error {
	return nil
}
func (c *Client) LogMessage(context.Context, *protocol.LogMessageParams) error {
	return nil
}
func (c *Client) ShowDocument(context.Context, *protocol.ShowDocumentParams) (*protocol.ShowDocumentResult, error) {
	return nil, nil
}
func (c *Client) ShowMessage(context.Context, *protocol.ShowMessageParams) error {
	return nil
}
func (c *Client) ShowMessageRequest(context.Context, *protocol.ShowMessageRequestParams) (*protocol.MessageActionItem, error) {
	return nil, nil
}
func (c *Client) WorkDoneProgressCreate(context.Context, *protocol.WorkDoneProgressCreateParams) error {
	return nil
}
func (c *Client) ApplyEdit(context.Context, *protocol.ApplyWorkspaceEditParams) (*protocol.ApplyWorkspaceEditResult, error) {
	return nil, nil
}
func (c *Client) CodeLensRefresh(context.Context) error {
	return nil
}
func (c *Client) Configuration(context.Context, *protocol.ParamConfiguration) ([]protocol.LSPAny, error) {
	return nil, nil
}
func (c *Client) WorkspaceFolders(context.Context) ([]protocol.WorkspaceFolder, error) {
	return nil, nil
}
