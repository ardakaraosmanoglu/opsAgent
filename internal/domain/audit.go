package domain

type AuditLog struct {
    ID            int64
    EventType     string
    Actor         string
    Message       string
    MetadataJSON  string
    CreatedAt     string
}
