package domain

type CommandType string
type RiskLevel string
type CommandStatus string

const (
    CommandTypeRead      CommandType = "read"
    CommandTypeWrite     CommandType = "write"
    CommandTypeDangerous CommandType = "dangerous"
    CommandTypeBlocked   CommandType = "blocked"

    RiskLow      RiskLevel = "low"
    RiskMedium   RiskLevel = "medium"
    RiskHigh     RiskLevel = "high"
    RiskCritical RiskLevel = "critical"

    CommandPending  CommandStatus = "pending"
    CommandApproved CommandStatus = "approved"
    CommandRunning  CommandStatus = "running"
    CommandSuccess  CommandStatus = "success"
    CommandFailed   CommandStatus = "failed"
    CommandBlocked  CommandStatus = "blocked"
)

type Command struct {
    ID               int64
    TaskPlanID       int64
    Command          string
    Type             CommandType
    RiskLevel        RiskLevel
    RequiresApproval bool
    Status           CommandStatus
}
