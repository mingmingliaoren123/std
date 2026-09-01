package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"openclaw-orchestrator/orchestrator"
)

func main() {
	flags := flag.NewFlagSet("openclaw-operator", flag.ExitOnError)
	bin := flags.String("bin", "", "OpenClaw CLI path (defaults to discovery or OPENCLAW_BIN)")
	config := flags.String("config", "", "OpenClaw state config path")
	manifest := flags.String("manifest", "", "application-owned agent manifest path")
	addr := flags.String("addr", "127.0.0.1:18790", "HTTP listen address for serve")
	timeout := flags.Duration("timeout", 10*time.Minute, "command timeout")
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: openclaw-operator [flags] <status|models|plugins|agents list|agents sync|serve>")
		flags.PrintDefaults()
	}
	_ = flags.Parse(os.Args[1:])
	if flags.NArg() == 0 {
		flags.Usage()
		os.Exit(2)
	}

	service := orchestrator.Discover(*manifest)
	if *bin != "" || *config != "" {
		binaryPath := service.BinaryPath()
		if *bin != "" {
			binaryPath = *bin
		}
		configPath := *config
		if configPath == "" {
			configPath = strings.TrimSpace(os.Getenv("OPENCLAW_CONFIG_PATH"))
			if configPath == "" {
				if stateDir := strings.TrimSpace(os.Getenv("OPENCLAW_STATE_DIR")); stateDir != "" {
					configPath = filepath.Join(stateDir, "openclaw.json")
				} else if home, err := os.UserHomeDir(); err == nil {
					configPath = home + "/.openclaw/openclaw.json"
				}
			}
		}
		stateDir := strings.TrimSpace(os.Getenv("OPENCLAW_STATE_DIR"))
		if stateDir == "" && configPath != "" {
			stateDir = filepath.Clean(filepath.Dir(configPath))
		}
		service = orchestrator.New(orchestrator.Config{BinaryPath: binaryPath, ConfigPath: configPath, StateDir: stateDir, Manifest: service.ManifestPath()})
	}

	command := strings.Join(flags.Args(), " ")
	if command == "serve" {
		log.Printf("openclaw-orchestrator listening on %s", *addr)
		log.Printf("OpenClaw CLI: %s", valueOrUnavailable(service.BinaryPath()))
		log.Fatal(http.ListenAndServe(*addr, orchestrator.NewHTTPHandler(service, orchestrator.HTTPOptions{})))
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	var value any
	var err error
	switch command {
	case "status":
		value, err = service.Status(ctx)
	case "models":
		value, err = service.Models(ctx)
	case "plugins":
		var plugins []orchestrator.Plugin
		plugins, err = service.Plugins(ctx)
		value = map[string]any{"catalogVersion": orchestrator.PluginCatalogVersion, "count": len(plugins), "plugins": plugins}
	case "agents list":
		var agents []orchestrator.Agent
		agents, err = service.ListAgents(ctx)
		value = map[string]any{"count": len(agents), "agents": agents}
	case "agents sync":
		var agents []orchestrator.Agent
		agents, err = service.SyncAgents(ctx)
		value = map[string]any{"synced": err == nil, "count": len(agents), "agents": agents}
	default:
		flags.Usage()
		os.Exit(2)
	}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Fatal("operation timed out")
		}
		log.Fatal(err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		log.Fatal(err)
	}
}

func valueOrUnavailable(value string) string {
	if value == "" {
		return "unavailable"
	}
	return value
}
