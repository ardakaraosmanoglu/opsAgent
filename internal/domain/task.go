package domain

type Task struct {
    ID        int64
    Source    string
    Prompt    string
    Status    string
    CreatedAt string
    UpdatedAt string
}

type TaskPlan struct {
    ID               int64
    TaskID           int64
    Summary          string
    RiskLevel        string
    RequiresApproval bool
    MetadataJSON     string
    CreatedAt        string
}
