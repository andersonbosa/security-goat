package UseSecurityGoat

import (
	"os"

	"github.com/andersonbosa/security-goat/pkg/GithubAPIClient"
	"github.com/andersonbosa/security-goat/pkg/SecurityGateSDK"
	"github.com/andersonbosa/security-goat/pkg/glog"
)

func Run(
	githubConfig SecurityGateSDK.GitHubConfig,
	severityLimits SecurityGateSDK.SeverityLimits,
	exitRules SecurityGateSDK.ExitRules,
) {
	glog.Println("🐐 Initializing Security Goat!")

	githubApiClient := GithubAPIClient.CreateGitHubClient(githubConfig.Token)

	analyzer := SecurityGateSDK.NewAnalyzer(
		githubApiClient,
		severityLimits,
		exitRules,
	)

	ExitCode := analyzer.AnalyzeGitHubAlerts(githubConfig.Owner, githubConfig.Repo)

	glog.Printf("Finazling the process with error code: %d.\n", ExitCode)

	os.Exit(ExitCode)
}
