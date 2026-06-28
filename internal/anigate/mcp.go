package anigate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
)

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id,omitempty"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MCPTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

type toolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type toolResult struct {
	Content           []toolContent `json:"content"`
	StructuredContent any           `json:"structuredContent,omitempty"`
	IsError           bool          `json:"isError,omitempty"`
}

type toolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func ServeStdio(r io.Reader, w io.Writer, svc *Service, log *slog.Logger) int {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024), 10*1024*1024)
	bw := bufio.NewWriter(w)
	defer bw.Flush()
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		resp, ok := dispatchJSON(line, svc)
		if !ok {
			continue
		}
		if err := writeRPC(bw, resp); err != nil {
			log.Error("write rpc", "err", err)
			return 1
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error("read stdin", "err", err)
		return 1
	}
	return 0
}

func dispatchJSON(b []byte, svc *Service) (rpcResponse, bool) {
	var req rpcRequest
	if err := json.Unmarshal(b, &req); err != nil {
		return rpcResponse{JSONRPC: "2.0", Error: &rpcError{Code: -32700, Message: "parse error"}}, true
	}
	if req.ID == nil && isNotification(req.Method) {
		return rpcResponse{}, false
	}
	return dispatch(req, svc), true
}

func dispatch(req rpcRequest, svc *Service) rpcResponse {
	resp := rpcResponse{JSONRPC: "2.0", ID: req.ID}
	switch req.Method {
	case "initialize":
		resp.Result = map[string]any{
			"protocolVersion": "2025-06-18",
			"capabilities": map[string]any{
				"tools": map[string]any{"listChanged": false},
			},
			"serverInfo": map[string]any{
				"name":    "anigate",
				"version": Version,
			},
		}
	case "ping":
		resp.Result = map[string]any{}
	case "tools/list":
		resp.Result = map[string]any{"tools": svc.Tools()}
	case "tools/call":
		var params toolCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			resp.Error = &rpcError{Code: -32602, Message: "invalid params"}
			return resp
		}
		result, err := svc.CallTool(params.Name, params.Arguments)
		tr := encodeToolResult(result, err)
		resp.Result = tr
	case "resources/list", "prompts/list":
		resp.Result = map[string]any{}
	default:
		resp.Error = &rpcError{Code: -32601, Message: "method not found"}
	}
	return resp
}

func encodeToolResult(result any, err error) toolResult {
	if err != nil {
		return toolResult{
			IsError: true,
			Content: []toolContent{{
				Type: "text",
				Text: err.Error(),
			}},
		}
	}
	b, jsonErr := json.MarshalIndent(result, "", "  ")
	if jsonErr != nil {
		return toolResult{
			IsError: true,
			Content: []toolContent{{
				Type: "text",
				Text: jsonErr.Error(),
			}},
		}
	}
	return toolResult{
		Content: []toolContent{{
			Type: "text",
			Text: string(b),
		}},
		StructuredContent: result,
	}
}

func writeRPC(w *bufio.Writer, resp rpcResponse) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, string(b)); err != nil {
		return err
	}
	return w.Flush()
}

func isNotification(method string) bool {
	return method == "notifications/initialized" || method == "notifications/cancelled" || method == "notifications/progress"
}
