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

var pd string

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&pd, "photo directory", "p", "", "directory to import photos to. Can include subdirectories.")
}

func main() {
	data, err := ioutil.ReadFile(viper.ConfigFileUsed())
	check(err)

	var config Config

	if err := config.ParseConfig(data); err != nil {
		log.Fatal(err)
	}
	workbench := CreateWorkbench(config.Workspace)
	SetupWorkbench(workbench, config.CollectionDirs)
}

func (config *Config) ParseConfig(data []byte) error {
	return yaml.Unmarshal(data, config)
}

func CreateWorkbench(Workspace string) string {
	workbench := filepath.Join(Workspace, pd)
	err := os.MkdirAll(workbench, 0777)
	check(err)

	return workbench
}

func SetupWorkbench(workbench string, colDirs []string) {
	err := os.Chdir(workbench)
	check(err)

	for d := 0; d < len(colDirs); d++ {
		err := os.MkdirAll(colDirs[d], 0777)
		check(err)
	}

	// Debug output
	c, err := ioutil.ReadDir(workbench)
	check(err)
	fmt.Println("Contents of " + workbench)
	for _, d := range c {
		fmt.Println(" ", d.Name())
	}
	//
}
