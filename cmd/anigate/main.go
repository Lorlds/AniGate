package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"anigate/internal/anigate"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 2
	}

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	switch args[0] {
	case "version":
		fmt.Println(anigate.Version)
		return 0
	case "stdio":
		fs := flag.NewFlagSet("stdio", flag.ContinueOnError)
		configPath := fs.String("config", "", "path to config JSON")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		svc, err := loadService(*configPath, log)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return anigate.ServeStdio(os.Stdin, os.Stdout, svc, log)
	case "http":
		fs := flag.NewFlagSet("http", flag.ContinueOnError)
		configPath := fs.String("config", "", "path to config JSON")
		addr := fs.String("addr", "127.0.0.1:8787", "HTTP listen address")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		svc, err := loadService(*configPath, log)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if err := anigate.ServeHTTP(*addr, svc, log); err != nil {
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
		svc, err := loadService(*configPath, log)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		for _, tool := range svc.Tools() {
			fmt.Printf("%s\t%s\n", tool.Name, tool.Description)
		}
		return 0
	default:
		usage()
		return 2
	}
}

func loadService(configPath string, log *slog.Logger) (*anigate.Service, error) {
	cfg, err := anigate.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return anigate.NewService(cfg, log)
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: anigate <command> [flags]

commands:
  version print version
  stdio   serve MCP JSON-RPC over stdin/stdout
  http    serve MCP JSON-RPC over HTTP POST /mcp
  tools   list exposed tools

examples:
  anigate stdio --config configs/anigate.example.json
  anigate http --addr 127.0.0.1:8787 --config configs/anigate.example.json`)
}
