package main

import (
	"fmt"
	"os"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
)

// runReport performs a one-shot port report and exits.
func runReport(scannerFactory func() (ports.Scanner, error), enricherFactory func() (ports.Enricher, error)) int {
	scanner, err := scannerFactory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: scanner init: %v\n", err)
		return 1
	}

	var enricher ports.Enricher
	if enricherFactory != nil {
		if e, err := enricherFactory(); err == nil {
			enricher = e
		}
	}

	h := monitor.NewReportHandler(scanner, enricher, os.Stdout)
	if err := h.Print(); err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: report: %v\n", err)
		return 1
	}
	return 0
}
