package cmd

import (
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import photos into configured workspace",
	Long: `Import photos into workspace configured in config file`,
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

type Config struct {
	ImportPath 			string 		`yaml:"import_path"`
	Workspace  			string 		`yaml:"workspace"`
	CollectionDirs	[]string	`yaml:"collection_dirs"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func main() {
	data, err := ioutil.ReadFile(viper.ConfigFileUsed())
	check(err)

	var config Config

	if err := config.ParseConfig(data); err != nil {
		log.Fatal(err)
	}
}

func (config *Config) ParseConfig(data []byte) error {
	return yaml.Unmarshal(data, config)
}
