// Copyright © 2017 Erno Rigo <erno@rigo.info>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/kardianos/osext"
	log "github.com/sirupsen/logrus"
	"github.com/heirko/go-contrib/logrusHelper"
	_ "github.com/heralight/logrus_mate/hooks/file"

	"github.com/mcree/cmdiff/session"
	"github.com/mcree/cmdiff/db"
	"github.com/mcree/cmdiff/report"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cmdiff",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmdiff.yaml and ./.cmdiff.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory
		viper.AddConfigPath(home)

		// Add current working directory to search config path
		wd, err := os.Getwd()
		if err == nil {
			viper.AddConfigPath(wd)
		}

		// Add executable directory to search config path
		ep, err := osext.ExecutableFolder()
		if err == nil {
			viper.AddConfigPath(ep)
		}

		// Set config file name to ".cmdiff" (without extension)
		viper.SetConfigName(".cmdiff")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Info("No config file:", err)
	}

	log.Debug( "Config settings: ", viper.AllSettings())

	var c = logrusHelper.UnmarshalConfiguration(viper.GetViper()) // Unmarshal configuration from Viper
	logrusHelper.SetConfig(log.StandardLogger(), c) // for e.g. apply it to logrus default instance

	log.Debug(c)

	sess := session.NewPipeline().Run()
	err := db.WriteSession(sess)
	log.Debug("Write result: ", err)
	db.DoHousekeeping()

	prev, err1 := db.PreviousSession()
	curr, err2 := db.CurrentSession()

	if err1 == nil && err2 == nil {
		rep, _ := report.NewSessionDiff(prev, curr).StringTemplate(viper.GetString("report.template"))
		fmt.Print(rep)
	}
}
