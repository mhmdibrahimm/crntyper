package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed VERSION
var version string

var (
	commit    = "unknown"
	buildDate = "unknown"
)

func appVersion() string {
	v := strings.TrimSpace(version)
	if v == "" {
		return "dev"
	}

	return v
}

func appTitle() string {
	return "CRN Typer " + appVersion()
}

func fullVersion() string {
	details := make([]string, 0, 2)
	if commit != "" && commit != "unknown" {
		details = append(details, "commit "+commit)
	}
	if buildDate != "" && buildDate != "unknown" {
		details = append(details, "built "+buildDate)
	}

	if len(details) == 0 {
		return appTitle()
	}

	return fmt.Sprintf("%s (%s)", appTitle(), strings.Join(details, ", "))
}
