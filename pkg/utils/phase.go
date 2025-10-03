package utils

import (
	"os"
	"strings"
)

const (
	DefaultPhase    = "dev"
	ProductionPhase = "prod"
)

var phase = DefaultPhase

func init() {
	if p, ok := os.LookupEnv("PHASE"); ok {
		phase = strings.ToLower(p)
	}
}

func GetPhase() string {
	return phase
}

func IsProductionPhase() bool {
	return phase == ProductionPhase
}

func IsDefaultPhase() bool {
	return phase == DefaultPhase
}
