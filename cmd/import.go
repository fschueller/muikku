// TODO:
// - Import by time range given through optional flags
// - Detect import path from platform (linux or macOS)
// - Refactor
// - Better error handling?
// - Bonus: detect common RAW formats?
// - Check for previous file/dir existence?

package cmd

import (
	"fmt"
	"time"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import photos into configured workspace",
	Long:  `Import photos into workspace configured in config file`,
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

type Config struct {
	ImportPath     string   `yaml:"import_path"`
	Workspace      string   `yaml:"workspace"`
	CollectionDirs []string `yaml:"collection_dirs"`
}

var pd string
var fromDateString string
var toDateString string
var fromDate, err = time.Parse("2006-01-02", fromDateString)
var toDate, err := time.Parse("2006-01-02", toDateString)

var jpg_ext = []string{".JPG", ".jpg", ".jpeg"}
var raw_ext = []string{".CR3", ".CR2"}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&pd, "photo directory", "p", "", "directory to import photos to. Can include subdirectories.")
	importCmd.Flags().StringVarP(&fromDateString, "from", "f", "", "timestamp of first picture to import (format: 'YYYY-MM-DD')")
	importCmd.Flags().StringVarP(&toDateString, "to", "t", "", "timestamp of last picture to import (format: 'YYYY-MM-DD')")

}

func main() {
	data, err := ioutil.ReadFile(viper.ConfigFileUsed())
	var config Config

	err = config.ParseConfig(data)


	workbench := CreateWorkbench(config.Workspace)
	SetupWorkbench(workbench, config.CollectionDirs)
	ImportPhotos(config.ImportPath, workbench)
	check(err)
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

	for _, d := range colDirs {
		err := os.MkdirAll(d, 0777)
		check(err)
	}
}

func ImportPhotos(importPath, workbench string) error {
	info, err := os.Lstat(importPath)
	check(err)
	return Copy(importPath, workbench, info)
}

// Sort while copying or sort after copy operation?
func SortbyType(dest string) string {
	ext := filepath.Ext(dest)
	dir := filepath.Dir(dest)
	filename := filepath.Base(dest)
	jpg_folder := "/JPG/"
	raw_folder := "/RAW/"

	if contains(jpg_ext, ext) {
		return constructPath(dir, filename, jpg_folder)
	}
	if contains(raw_ext, ext) {
		return constructPath(dir, filename, raw_folder)
	}
	fmt.Println("Keep " + filename + " in " + dir)
	return dest
}

func constructPath(dir, filename, folder string) string {
	new_dest := dir + folder + filename
	fmt.Println("Sorting " + filename + " into " + new_dest)
	return new_dest
}

func contains(types []string, ext string) bool {
	for _, i := range types {
		if i == ext {
			return true
		}
	}
	return false
}

// Recursively copying the files in src to dest and sorting them on the way

func Copy(src, dest string, info os.FileInfo) error {
	if info.IsDir() {
		return CopyDirectory(src, dest, info)
	}
	return CopyFile(src, dest, info)
}

func CopyDirectory(src, dest string, info os.FileInfo) error {
	c, err := ioutil.ReadDir(src)
	check(err)

	for _, i := range c {
		timestamp := i.ModTime()
		// if inTimeRange(timestamp) {
		// 	c += i
		// }
		inTimeRange(timestamp)
	}

	for _, p := range c {
		ps, pd := filepath.Join(src, p.Name()), filepath.Join(dest, p.Name())
		pd = SortbyType(pd)
		if err := Copy(ps, pd, p); err != nil {
			return err
		}
	}
	return nil
}

func inTimeRange(timestamp time.Time) bool {
	fileDate := timestamp.Format("2006-01-02")
	if fileDate.After(fromDate) {
		return fmt.Println(fileDate)
	}
	return true
}

func CopyFile(src, dest string, info os.FileInfo) error {
	file, err := os.Create(dest)
	check(err)
	defer file.Close()

	s, err := os.Open(src)
	check(err)
	defer s.Close()

	_, err = io.Copy(file, s)
	return err
}
