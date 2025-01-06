package main

type AgentRigctlConfig struct {
	Enabled bool   `ini:"enabled"`
	Host    string `ini:"host"`
	Port    string `ini:"port"`
}

type AgentWavelogConfig struct {
	URL     string `ini:"url"`
	Key     string `ini:"key"`
	Radio   string `ini:"radio"`
	Profile string `ini:"profile"`
}

type AgentConfig struct {
	Rigctld *AgentRigctlConfig
	Wavelog *AgentWavelogConfig
}
