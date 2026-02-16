package consumer

import (
	"errors"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

type ChaosConfig struct {
	Enabled      bool
	ErrorRate    float64
	LatencyMinMs int
	LatencyMaxMs int
	Logger       *logrus.Logger
}

func InjectRandomError(config *ChaosConfig) error {
	if !config.Enabled {
		return nil
	}

	if ShouldFail(config.ErrorRate) {
		config.Logger.Warn("Chaos: Error injected")
		return errors.New("chaos engineering: simulated error")
	}

	return nil
}

func InjectLatency(config *ChaosConfig) {
	if !config.Enabled || config.LatencyMaxMs <= 0 {
		return
	}

	minMs := config.LatencyMinMs
	maxMs := config.LatencyMaxMs

	if minMs < 0 {
		minMs = 0
	}
	if maxMs < minMs {
		maxMs = minMs
	}

	latencyMs := minMs
	if maxMs > minMs {
		latencyMs = minMs + rand.Intn(maxMs-minMs)
	}

	if latencyMs > 0 {
		config.Logger.WithField("latency_ms", latencyMs).Info("Chaos: Latency injected")
		time.Sleep(time.Duration(latencyMs) * time.Millisecond)
	}
}

func ShouldFail(rate float64) bool {
	if rate <= 0 {
		return false
	}
	if rate >= 1 {
		return true
	}
	return rand.Float64() < rate
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
