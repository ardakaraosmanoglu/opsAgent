package domain

type Metric struct {
    ID            int64
    Hostname      string
    OSInfo        string
    KernelVersion string
    CPUUsage      float64
    MemoryUsage   float64
    DiskUsage     float64
    LoadAverage1  float64
    LoadAverage5  float64
    LoadAverage15 float64
    UptimeSeconds int64
    CreatedAt     string
}
