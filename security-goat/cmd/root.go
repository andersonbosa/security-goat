/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/andersonbosa/security-goat/pkg/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	repoURL = "https://github.com/andersonbosa/security-goat"
)

const (
	cfgFileName = ".security-goat"
)

var (
	cfgFile  string
	verbose  bool
	insecure bool
	dryrun   bool

	githubToken string
	githubOwner string
	githubRepo  string

	criticalSeverity int
	highSeverity     int
	mediumSeverity   int
	lowSeverity      int

	successExitCode int
	errorExitCode   int
	fatalExitCode   int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "security-goat",
	Short: "Analyze and manage security alerts for GitHub repositories",
	Long: `Security Gote is a CLI application that allows you to analyze security alerts
and manage security configurations for GitHub repositories. It provides insights
into vulnerabilities and helps you maintain a secure software ecosystem.

For more details and options, you can run 'security-goat --help'.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			glog.DefaultLogger.Level = glog.LevelDebug
			glog.Debugln("Verbose enabled")
		}
	},
	Run: GateRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(args []string) {
	// if args[0] == "version" {
	// versionCmd.Run(versionCmd, args)
	// os.Exit(successExitCode)
	// return
	// }

	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		qwe(fatalExitCode, err, "failed to execute root command")
	}
}

// Here you will define your flags and configuration settings.
// Cobra supports persistent flags, which, if defined here,
// will be global for your application.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.security-goat.yaml)")

	rootCmd.PersistentFlags().StringVarP(&githubToken, "token", "t", "", "GitHub token")
	rootCmd.PersistentFlags().StringVarP(&githubOwner, "owner", "o", "", "GitHub owner")
	rootCmd.PersistentFlags().StringVarP(&githubRepo, "repo", "r", "", "GitHub repository")

	rootCmd.PersistentFlags().IntVarP(&criticalSeverity, "critical", "C", 0, "Critical severity limit")
	rootCmd.PersistentFlags().IntVarP(&highSeverity, "high", "H", 1, "High severity limit")
	rootCmd.PersistentFlags().IntVarP(&mediumSeverity, "medium", "M", 2, "Medium severity limit")
	rootCmd.PersistentFlags().IntVarP(&lowSeverity, "low", "L", 5, "Low severity limit")

	rootCmd.PersistentFlags().IntVarP(&successExitCode, "success", "", 0, "Customize exit code to success")
	rootCmd.PersistentFlags().IntVarP(&errorExitCode, "error", "", 1, "Customize exit code to error")
	rootCmd.PersistentFlags().IntVarP(&fatalExitCode, "fatal", "", 128, "Customize exit code to fatal")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "more logs")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "skip TLS verification and use insecure http client")
	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Don't perform security gate")

	// Bind flags to Viper
	_ = viper.BindPFlag("github.token", rootCmd.PersistentFlags().Lookup("token"))
	_ = viper.BindPFlag("github.owner", rootCmd.PersistentFlags().Lookup("owner"))
	_ = viper.BindPFlag("github.repo", rootCmd.PersistentFlags().Lookup("repo"))

	_ = viper.BindPFlag("severity_limits.critical", rootCmd.PersistentFlags().Lookup("critical"))
	_ = viper.BindPFlag("severity_limits.high", rootCmd.PersistentFlags().Lookup("high"))
	_ = viper.BindPFlag("severity_limits.medium", rootCmd.PersistentFlags().Lookup("medium"))
	_ = viper.BindPFlag("severity_limits.low", rootCmd.PersistentFlags().Lookup("low"))

	_ = viper.BindPFlag("exit_codes.success", rootCmd.PersistentFlags().Lookup("success"))
	_ = viper.BindPFlag("exit_codes.error", rootCmd.PersistentFlags().Lookup("error"))
	_ = viper.BindPFlag("exit_codes.fatal", rootCmd.PersistentFlags().Lookup("fatal"))

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("insecure", rootCmd.PersistentFlags().Lookup("insecure"))
	_ = viper.BindPFlag("dryrun", rootCmd.PersistentFlags().Lookup("dryrun"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Use config file from the flag.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Look for the config file in home and project directory
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")
		viper.SetConfigName(cfgFileName)
		viper.SetEnvPrefix("goat")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		glog.Debugln("Using config file:", viper.ConfigFileUsed())
	}

	if viper.GetString("github.token") == "" {
		qwm(fatalExitCode, fmt.Sprintf("GitHub API Token configuration is required. Provide them via a config file, environment variables or command line arguments. For more information on configuration, see README on GitHub repository. %s\n", repoURL))
	}
	if viper.GetString("github.owner") == "" {
		qwm(fatalExitCode, fmt.Sprintf("GitHub owner configuration is required. Provide them via a config file, environment variables or command line arguments. For more information on configuration, see README on GitHub repository. %s\n", repoURL))
	}
	if viper.GetString("github.repo") == "" {
		qwm(fatalExitCode, fmt.Sprintf("GitHub repository name Token configuration is required. Provide them via a config file, environment variables or command line arguments. For more information on configuration, see README on GitHub repository. %s\n", repoURL))
	}
}
