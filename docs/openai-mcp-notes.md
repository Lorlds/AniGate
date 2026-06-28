# OpenAI MCP Notes

Checked on 2026-06-26 against official OpenAI developer docs.

- ChatGPT Apps use MCP servers to expose app capabilities to ChatGPT.
- Data-only apps can skip UI resources and expose tools only.
- ChatGPT can connect to a remote MCP server endpoint.
- Secure MCP Tunnel is the official private-server path when AniGate should not
  expose inbound public ports.

Useful official pages:

- `https://developers.openai.com/apps-sdk/build/mcp-server`
- `https://developers.openai.com/apps-sdk/deploy/connect-chatgpt`
- `https://developers.openai.com/api/docs/mcp`
- `https://developers.openai.com/api/docs/guides/secure-mcp-tunnels`
