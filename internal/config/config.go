package config

import (
    "fmt"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    App struct {
        Name        string `yaml:"name"`
        Environment string `yaml:"environment"`
    } `yaml:"app"`
    Server struct {
        BindAddress string `yaml:"bind_address"`
        Port        int    `yaml:"port"`
    } `yaml:"server"`
    Database struct {
        Path string `yaml:"path"`
    } `yaml:"database"`
    Logging struct {
        Level string `yaml:"level"`
        File  string `yaml:"file"`
    } `yaml:"logging"`
    Security struct {
        RequireApprovalForWriteActions bool `yaml:"require_approval_for_write_actions"`
        CommandTimeoutSeconds         int  `yaml:"command_timeout_seconds"`
        MaxOutputSizeKB               int  `yaml:"max_output_size_kb"`
    } `yaml:"security"`
    AI struct {
        Enabled  bool   `yaml:"enabled"`
        Provider string `yaml:"provider"`
        APIKey   string `yaml:"api_key"`
        Model    string `yaml:"model"`
    } `yaml:"ai"`
    Monitoring struct {
        IntervalSeconds          int     `yaml:"interval_seconds"`
        DiskWarningThreshold     float64 `yaml:"disk_warning_threshold"`
        DiskCriticalThreshold    float64 `yaml:"disk_critical_threshold"`
        MemoryWarningThreshold   float64 `yaml:"memory_warning_threshold"`
        CPUWarningThreshold      float64 `yaml:"cpu_warning_threshold"`
        CPUWarningDurationSecond int     `yaml:"cpu_warning_duration_seconds"`
    } `yaml:"monitoring"`
    Dashboard struct {
        SessionSecret string `yaml:"session_secret"`
    } `yaml:"dashboard"`
}

func Load(path string) (*Config, error) {
    raw, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config file: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(raw, &cfg); err != nil {
        return nil, fmt.Errorf("parsing config file: %w", err)
    }

    if cfg.Server.Port == 0 {
        cfg.Server.Port = 8787
    }
    if cfg.Database.Path == "" {
        cfg.Database.Path = "/var/lib/opsagent/opsagent.db"
    }

    return &cfg, nil
}
