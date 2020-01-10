package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/realsangil/ddl2struct/pkg/parser"

	"github.com/spf13/cobra"

	"github.com/atotto/clipboard"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

var cfgFile string
var (
	inputPath       string
	outputPath      string
	tableName       string
	copyToClipboard bool
)

var rootCmd = &cobra.Command{
	Use:   "ddl2struct",
	Short: "create golang struct from ddl",
	Long:  ``,
	Run:   runCommand,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flag := rootCmd.PersistentFlags()
	flag.StringVarP(&inputPath, "input", "i", "", `sql file path`)
	flag.StringVarP(&outputPath, "output", "o", "", `output file path`)
	flag.BoolVarP(&copyToClipboard, "copy", "c", false, "copy to clipboard")
}

func runCommand(cmd *cobra.Command, args []string) {
	sql, err := ioutil.ReadFile(inputPath)
	if err != nil {
		panic(err)
	}

	parser := parser.New()
	if err := parser.Parse(string(sql)); err != nil {
		panic(err)
	}

	structBytes, err := parser.ToStruct(true)
	if err != nil {
		panic(err)
	}

	if outputPath != "" {
		if err := ioutil.WriteFile(outputPath, structBytes, 0644); err != nil {
			panic(err)
		}
	}

	if copyToClipboard {
		clipboard.WriteAll(fmt.Sprintf("%s", structBytes))
	}

	columnCount := len(parser.Columns)
	fmt.Printf(`
#   ____ ____ __      ____ __     ____ ____ ____ _  _  ___ ____ 
#  (    (    (  )    (_  _/  \   / ___(_  _(  _ / )( \/ __(_  _)
#   ) D () D / (_/\    )((  O )  \___ \ )(  )   ) \/ ( (__  )(  
#  (____(____\____/   (__)\__/   (____/(__)(__\_\____/\___)(__)

===============================================================
TableName: %s
ColumnCount: %d
Save your times: About %d seconds
===============================================================

%s

`, parser.TableName, columnCount, 2*columnCount, structBytes)
}

func generateFileFromBytes(structBytes []byte) {
	if err := ioutil.WriteFile(outputPath, structBytes, 0644); err != nil {
		panic(err)
	}
}
