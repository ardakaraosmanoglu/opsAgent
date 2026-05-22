package domain

type Alert struct {
    ID           int64
    Type         string
    Severity     string
    Title        string
    Message      string
    Status       string
    MetadataJSON string
    CreatedAt    string
    UpdatedAt    string
}
