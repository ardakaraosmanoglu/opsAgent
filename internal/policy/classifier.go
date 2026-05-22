package policy

import (
    "strings"

    "github.com/opsagent/opsagent/internal/domain"
)

type Result struct {
    Type             domain.CommandType
    RiskLevel        domain.RiskLevel
    RequiresApproval bool
    Allowed          bool
    Reason           string
}

var blockedPatterns = []string{
    "rm -rf /",
    "mkfs",
    "dd if=",
    "chmod -R 777 /",
    "chown -R /",
    "userdel",
    "passwd",
    "iptables -F",
    "ufw disable",
    "shutdown",
    "reboot",
    "curl ",
    "wget ",
    "| bash",
}

var writePrefixes = []string{
    "apt update",
    "apt upgrade",
    "apt install",
    "systemctl restart",
    "systemctl reload",
    "journalctl --vacuum-time",
    "docker restart",
}

var readPrefixes = []string{
    "df ",
    "free ",
    "uptime",
    "ps ",
    "top ",
    "systemctl status",
    "journalctl ",
    "ss ",
    "cat /etc/os-release",
    "apt list",
    "du ",
    "docker ps",
    "docker logs",
}

func Classify(command string) Result {
    normalized := strings.TrimSpace(command)
    lower := strings.ToLower(normalized)

    for _, pattern := range blockedPatterns {
        if strings.Contains(lower, strings.ToLower(pattern)) {
            return Result{
                Type:             domain.CommandTypeBlocked,
                RiskLevel:        domain.RiskCritical,
                RequiresApproval: false,
                Allowed:          false,
                Reason:           "Command matched blocked safety policy",
            }
        }
    }

    for _, prefix := range writePrefixes {
        if strings.HasPrefix(lower, strings.ToLower(prefix)) {
            return Result{
                Type:             domain.CommandTypeWrite,
                RiskLevel:        domain.RiskMedium,
                RequiresApproval: true,
                Allowed:          true,
                Reason:           "Write operation requires approval",
            }
        }
    }

    for _, prefix := range readPrefixes {
        if strings.HasPrefix(lower, strings.ToLower(prefix)) {
            return Result{
                Type:             domain.CommandTypeRead,
                RiskLevel:        domain.RiskLow,
                RequiresApproval: false,
                Allowed:          true,
                Reason:           "Read-only command",
            }
        }
    }

    return Result{
        Type:             domain.CommandTypeDangerous,
        RiskLevel:        domain.RiskHigh,
        RequiresApproval: true,
        Allowed:          false,
        Reason:           "Command is not whitelisted",
    }
}
