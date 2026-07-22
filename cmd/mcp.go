package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// MCP JSON-RPC message
type mcpMessage struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      json.Number     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type mcpResult struct {
	ID      json.Number `json:"id"`
	Result  any         `json:"result,omitempty"`
	Error   *mcpError   `json:"error,omitempty"`
	Jsonrpc string      `json:"jsonrpc"`
}

type mcpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP types
type mcpTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

// proto-gen tool schema
var toolProtoGen = mcpTool{
	Name:        "proto-gen",
	Description: "从 .proto 文件生成 Go/TypeScript 代码。需要 protoc 编译器已安装。",
	InputSchema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"proto_dir":   map[string]any{"type": "string", "description": "包含 .proto 文件的目录"},
			"go_out":      map[string]any{"type": "string", "description": "Go 代码输出目录"},
			"ts_out":      map[string]any{"type": "string", "description": "TypeScript 输出目录（需 protoc-gen-es）"},
			"go_module":   map[string]any{"type": "string", "description": "Go module 路径"},
			"include":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "protoc -I 附加包含路径"},
			"inject_tag":  map[string]any{"type": "boolean", "description": "编译后注入 struct tag（需 protoc-go-inject-tag）"},
			"plugin":      map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "protoc 插件，格式: name[,binary=<path>][,out=<dir>][,module=<mod>]"},
			"no_recursive": map[string]any{"type": "boolean", "description": "只扫描顶层目录，不递归子目录"},
		},
		"required": []any{"proto_dir"},
	},
}

var toolExcelGen = mcpTool{
	Name:        "excel-gen",
	Description: "从 Excel 文件生成 JSON 数据以及 Go struct / TypeScript interface",
	InputSchema: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"excel_dir":  map[string]any{"type": "string", "description": "包含 .xlsx/.xls 文件的目录"},
			"json_out":   map[string]any{"type": "string", "description": "JSON 输出目录（默认: <excel-dir>/json）"},
			"go_out":     map[string]any{"type": "string", "description": "Go struct 输出目录"},
			"go_package": map[string]any{"type": "string", "description": "Go 包名（与 go_out 配合使用）"},
			"ts_out":     map[string]any{"type": "string", "description": "TypeScript interface 输出目录"},
		},
		"required": []any{"excel_dir"},
	},
}

var mcpTools = []mcpTool{toolProtoGen, toolExcelGen}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "以 MCP stdio 模式启动，作为 AI 工具的 MCP server",
	Long: `MCP (Model Context Protocol) stdio server 模式。
AI 工具通过 stdio 与本 CLI 通信，调用 proto-gen / excel-gen 等工具。

用法（在 AI 工具的 MCP 配置中）:
{
  "mcpServers": {
    "game-dev-cli": {
      "command": "game-dev-cli",
      "args": ["mcp"]
    }
  }
}`,
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")
		return runMCPServer(verbose)
	},
}

func runMCPServer(verbose bool) error {
	scanner := bufio.NewScanner(os.Stdin)
	// max message size: 4MB
	scanner.Buffer(make([]byte, 0, 1024*1024), 4*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var msg mcpMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "invalid mcp message: %v\n", err)
			}
			continue
		}

		var result any
		var rpcErr *mcpError

		switch msg.Method {
		case "initialize":
			result = map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]any{},
				"serverInfo": map[string]any{
					"name":    "game-dev-cli",
					"version": "0.1.0",
				},
			}

		case "tools/list":
			result = map[string]any{"tools": mcpTools}

		case "tools/call":
			var params struct {
				Name      string          `json:"name"`
				Arguments json.RawMessage `json:"arguments"`
			}
			if err := json.Unmarshal(msg.Params, &params); err != nil {
				rpcErr = &mcpError{Code: -32602, Message: fmt.Sprintf("invalid params: %v", err)}
				break
			}
			out, err := callTool(params.Name, params.Arguments)
			if err != nil {
				rpcErr = &mcpError{Code: -32000, Message: err.Error()}
			} else {
				result = map[string]any{
					"content": []map[string]any{
						{"type": "text", "text": out},
					},
				}
			}

		default:
			rpcErr = &mcpError{Code: -32601, Message: fmt.Sprintf("unknown method: %s", msg.Method)}
		}

		resp := mcpResult{
			ID:      msg.ID,
			Result:  result,
			Error:   rpcErr,
			Jsonrpc: "2.0",
		}

		data, _ := json.Marshal(resp)
		fmt.Println(string(data))
	}

	return scanner.Err()
}

func callTool(name string, args json.RawMessage) (string, error) {
	switch name {
	case "proto-gen":
		return callProtoGen(args)
	case "excel-gen":
		return callExcelGen(args)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// ponytail: direct Go function call would be cleaner; currently shelling out to self for simplicity.
func callProtoGen(args json.RawMessage) (string, error) {
	var p struct {
		ProtoDir    string   `json:"proto_dir"`
		GoOut       string   `json:"go_out"`
		TsOut       string   `json:"ts_out"`
		GoModule    string   `json:"go_module"`
		Include     []string `json:"include"`
		InjectTag   bool     `json:"inject_tag"`
		Plugin      []string `json:"plugin"`
		NoRecursive bool     `json:"no_recursive"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return "", fmt.Errorf("invalid proto-gen args: %w", err)
	}

	cliArgs := []string{"proto", "gen", p.ProtoDir}
	if p.GoOut != "" {
		cliArgs = append(cliArgs, "--go-out", p.GoOut)
	}
	if p.TsOut != "" {
		cliArgs = append(cliArgs, "--ts-out", p.TsOut)
	}
	if p.GoModule != "" {
		cliArgs = append(cliArgs, "--go-module", p.GoModule)
	}
	for _, inc := range p.Include {
		cliArgs = append(cliArgs, "--include", inc)
	}
	if p.InjectTag {
		cliArgs = append(cliArgs, "--inject-tag")
	}
	for _, pl := range p.Plugin {
		cliArgs = append(cliArgs, "--plugin", pl)
	}
	if p.NoRecursive {
		cliArgs = append(cliArgs, "--no-recursive")
	}

	return runSelf(cliArgs)
}

func callExcelGen(args json.RawMessage) (string, error) {
	var p struct {
		ExcelDir  string `json:"excel_dir"`
		JsonOut   string `json:"json_out"`
		GoOut     string `json:"go_out"`
		GoPackage string `json:"go_package"`
		TsOut     string `json:"ts_out"`
	}
	if err := json.Unmarshal(args, &p); err != nil {
		return "", fmt.Errorf("invalid excel-gen args: %w", err)
	}

	cliArgs := []string{"excel", "gen", p.ExcelDir}
	if p.JsonOut != "" {
		cliArgs = append(cliArgs, "--json-out", p.JsonOut)
	}
	if p.GoOut != "" {
		cliArgs = append(cliArgs, "--go-out", p.GoOut)
	}
	if p.GoPackage != "" {
		cliArgs = append(cliArgs, "--go-package", p.GoPackage)
	}
	if p.TsOut != "" {
		cliArgs = append(cliArgs, "--ts-out", p.TsOut)
	}

	return runSelf(cliArgs)
}

func runSelf(args []string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	// resolve symlinks to get canonical path
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(exe, args...)
	cmd.Dir, _ = os.Getwd()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", strings.TrimSpace(string(out)), err)
	}
	return string(out), nil
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
