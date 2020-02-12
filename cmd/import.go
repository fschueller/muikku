// TODO:
// - Copy files from import to workbench
// - Copy them into correct directories (JPG - CR3 (bonus: detect common RAW formats?))
// - Refactor
// - Better error handling?

package cmd

import (
	"fmt"
	"log"
	"os"
	"io"
	"io/ioutil"
	"path/filepath"
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
	ImportPhotos(config.ImportPath, workbench)
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
}

func ImportPhotos(importPath string, workbench string) error {
	info, err := os.Lstat(importPath)
	check(err)
	return Copy(importPath, workbench, info)
}

// Recursively copying the files in src to dest

func Copy(src string, dest string, info os.FileInfo) error {
	if info.IsDir() {
		return CopyDirectory(src, dest, info)
	}
	return CopyFile(src, dest, info)
}

func CopyDirectory(src string, dest string, info os.FileInfo) error {

	c, err := ioutil.ReadDir(src)
	check(err)

	for _, p := range c {
		ps, pd := filepath.Join(src, p.Name()), filepath.Join(dest, p.Name())
		if err := Copy(ps, pd, p); err != nil {
			return err
		}
	}
	return nil
}

func CopyFile(src string, dest string, info os.FileInfo) error {
	file, err := os.Create(dest)
	check(err)
	defer file.Close()

	s, err := os.Open(src)
	check(err)
	defer s.Close()

	_, err = io.Copy(file, s)
	return err
}
