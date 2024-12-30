package main

import "fmt"

type AgentRigctlConfig struct {
    Enabled     bool   `ini:"enabled"`
    Host        string `ini:"host"`
    Port        int    `ini:"port"`
}

type AgentWavelogConfig struct {
	URL         string `ini:"url"`
	Key         string `ini:"key"`
	Radio       string `ini:"radio"`
}

type AgentConfig struct {
	Rigctl          *AgentRigctlConfig
	Wavelog         *AgentWavelogConfig
}
