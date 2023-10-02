package SecurityGateSDK

import (
	"context"
	"fmt"
	"strings"

	"github.com/andersonbosa/security-goat/pkg/glog"
	"github.com/google/go-github/v55/github"
)

type GitHubConfig struct {
	Token string
	Owner string
	Repo  string
}

// ExitRules represents exit codes for different scenarios.
type ExitRules struct {
	Success int
	Error   int
	Fatal   int
}

// SeverityLimits holds severity limits for the analysis.
type SeverityLimits struct {
	Critical int
	High     int
	Medium   int
	Low      int
}

// Analyzer represents the security gate analyzer.
type Analyzer struct {
	client         *github.Client
	severityLimits SeverityLimits
	exitRules      ExitRules
}

var KnowSeveritiesLevel = []string{"critical", "high", "medium", "low"}

// NewAnalyzer creates a new instance of the Analyzer.
func NewAnalyzer(client *github.Client, limits SeverityLimits, exitRules ExitRules) *Analyzer {
	return &Analyzer{
		client:         client,
		severityLimits: limits,
		exitRules:      exitRules,
	}
}

// FetchGitHubAlerts fetches GitHub alerts for a given repository and applies severity limits.
func (a *Analyzer) FetchGitHubAlerts(owner, repository string) ([]*github.DependabotAlert, error) {
	alertOptions := &github.ListAlertsOptions{}
	ctx := context.Background()

	alerts, resp, err := a.client.Dependabot.ListRepoAlerts(ctx, owner, repository, alertOptions)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub alerts: %v", err)
	}

	// Check if the HTTP status code is in the 20x range
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}

	return alerts, nil
}

// AnalyzeGitHubAlerts analyzes GitHub alerts and applies severity limits.
func (a *Analyzer) AnalyzeGitHubAlerts(owner, repository string) int {
	// Print the GitHub URL being analyzed
	githubURL := GenerateGitHubURL(owner, repository)
	glog.Println("Analyzing GitHub Alerts from:", githubURL)

	alerts, err := a.FetchGitHubAlerts(owner, repository)
	if err != nil {
		glog.Warnf("%s\n", err)
		return a.exitRules.Fatal
	}

	glog.Printf("Total alerts: %d\n", len(alerts))

	severityCount := countSeverityAlerts(alerts, KnowSeveritiesLevel)

	for severity, count := range severityCount {
		glog.Debugf("Found: %d %s\n", count, severity)
	}

	ExitCode := applySeverityPolicies(severityCount, a.severityLimits, a.exitRules, KnowSeveritiesLevel)

	return ExitCode
}

func countSeverityAlerts(alerts []*github.DependabotAlert, severities []string) map[string]int {
	severityCount := make(map[string]int)
	for _, alert := range alerts {
		if alert.GetState() == "open" {
			severity := alert.GetSecurityVulnerability().GetSeverity()
			if isValidSeverity(severity) {
				severityCount[severity]++
			}
		}
	}

	return severityCount
}

func applySeverityPolicies(severityCount map[string]int, severityLimits SeverityLimits, exitRules ExitRules, severities []string) int {
	for _, severity := range severities {
		count := severityCount[severity]

		limit := parseSeverityLimits(severityLimits, severity)
		if count > limit {
			glog.Printf("More than %d %s security alerts found.", limit, severity)
			return exitRules.Error
		}
	}

	return exitRules.Success
}

func parseSeverityLimits(limits SeverityLimits, severity string) int {
	switch strings.ToLower(severity) {
	case "critical":
		return limits.Critical
	case "high":
		return limits.High
	case "medium":
		return limits.Medium
	case "low":
		return limits.Low
	default:
		return 0
	}
}

func isValidSeverity(severity string) bool {
	severity = strings.ToLower(severity)
	return severity == "critical" || severity == "high" || severity == "medium" || severity == "low"
}

func GenerateGitHubURL(owner, repository string) string {
	return fmt.Sprintf("https://github.com/%s/%s", owner, repository)
}
