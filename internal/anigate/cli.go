package anigate

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func RunCLI(args []string, productLine ProductLine) int {
	line, err := normalizeProductLine(productLine)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	name := filepath.Base(os.Args[0])
	if len(args) == 0 {
		usage(name, line)
		return 2
	}

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	switch args[0] {
	case "version":
		fmt.Println(Version)
		return 0
	case "stdio":
		fs := flag.NewFlagSet("stdio", flag.ContinueOnError)
		configPath := fs.String("config", "", "path to config JSON")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		svc, err := loadServiceForProduct(*configPath, log, line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return ServeStdio(os.Stdin, os.Stdout, svc, log)
	case "http":
		fs := flag.NewFlagSet("http", flag.ContinueOnError)
		configPath := fs.String("config", "", "path to config JSON")
		addr := fs.String("addr", defaultHTTPAddr(name, line), "HTTP listen address")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		svc, err := loadServiceForProduct(*configPath, log, line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if err := ServeHTTP(*addr, svc, log); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	case "tools":
		fs := flag.NewFlagSet("tools", flag.ContinueOnError)
		configPath := fs.String("config", "", "path to config JSON")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		svc, err := loadServiceForProduct(*configPath, log, line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		for _, tool := range svc.Tools() {
			fmt.Printf("%s\t%s\n", tool.Name, tool.Description)
		}
		return 0
	default:
		usage(name, line)
		return 2
	}
}

func loadServiceForProduct(configPath string, log *slog.Logger, productLine ProductLine) (*Service, error) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return NewServiceWithProductLine(cfg, log, productLine)
}

func defaultHTTPAddr(commandName string, productLine ProductLine) string {
	if productLine == ProductLineMax && strings.Contains(commandName, "max") {
		return "127.0.0.1:8788"
	}
	return "127.0.0.1:8787"
}

func exampleConfig(commandName string, productLine ProductLine) string {
	if productLine == ProductLineMini {
		return "configs/anigate.mini.example.json"
	}
	if strings.Contains(commandName, "max") {
		return "configs/anigate.max.example.json"
	}
	return "configs/anigate.example.json"
}

func usage(commandName string, productLine ProductLine) {
	config := exampleConfig(commandName, productLine)
	addr := defaultHTTPAddr(commandName, productLine)
	fmt.Fprintf(os.Stderr, `usage: %s <command> [flags]

product line: %s

commands:
  version print version
  stdio   serve MCP JSON-RPC over stdin/stdout
  http    serve MCP JSON-RPC over HTTP POST /mcp
  tools   list exposed tools

examples:
  %s stdio --config %s
  %s http --addr %s --config %s
`, commandName, productLine, commandName, config, commandName, addr, config)
}
