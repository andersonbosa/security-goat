package cmd

import (
	"github.com/andersonbosa/security-goat/internal/pkg/usecases/UseSecurityGoat"
	"github.com/andersonbosa/security-goat/pkg/SecurityGateSDK"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GateRun(cmd *cobra.Command, args []string) {
	// Get values from flags
	severityLimits := SecurityGateSDK.SeverityLimits{
		Critical: viper.GetInt("severity_limits.critical"),
		High:     viper.GetInt("severity_limits.high"),
		Medium:   viper.GetInt("severity_limits.medium"),
		Low:      viper.GetInt("severity_limits.low"),
	}

	exitRules := SecurityGateSDK.ExitRules{
		Success: viper.GetInt("exit_codes.success"),
		Error:   viper.GetInt("exit_codes.error"),
		Fatal:   viper.GetInt("exit_codes.fatal"),
	}

	githubConfig := SecurityGateSDK.GitHubConfig{
		Token: viper.GetString("github.token"),
		Owner: viper.GetString("github.owner"),
		Repo:  viper.GetString("github.repo"),
	}

	if viper.GetBool("dryrun") {
		// glog.Println("severityLimits", severityLimits)
		// glog.Println("exitRules", exitRules)
		// glog.Println("githubConfig", githubConfig)
		return
	}

	UseSecurityGoat.Run(githubConfig, severityLimits, exitRules)
}
