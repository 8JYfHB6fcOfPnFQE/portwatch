package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/rules"
)

func main() {
	configPath := flag.String("config", "configs/portwatch.yaml", "path to config file")
	rulesPath := flag.String("rules", "", "path to rules file (overrides config)")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	if *rulesPath != "" {
		cfg.RulesFile = *rulesPath
	}

	var matcher *rules.Matcher
	if cfg.RulesFile != "" {
		ruleList, err := rules.LoadFromFile(cfg.RulesFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading rules: %v\n", err)
			os.Exit(1)
		}
		matcher = rules.NewMatcher(ruleList)
	}

	notifier, err := alert.NewFormattedNotifier(cfg.Output.Format, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating notifier: %v\n", err)
		os.Exit(1)
	}

	mon := monitor.New(cfg, matcher, notifier)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("portwatch started (interval=%s, ports=%d-%d)\n",
		cfg.Interval, cfg.Ports.Start, cfg.Ports.End)

	go func() {
		<-sigCh
		fmt.Println("\nshutting down...")
		mon.Stop()
	}()

	mon.Run()
}
