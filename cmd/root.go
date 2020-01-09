package cmd

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/atotto/clipboard"
	executor "github.com/bytewatch/ddl-executor"
	"github.com/iancoleman/strcase"
)

var cfgFile string
var (
	inputPath        string
	outputPath       string
	tableName        string
	table            *executor.TableDef
	copyToClipboard  bool
	columnTypeRegexp = regexp.MustCompile(`(\w+)(\((.+)\))?`)
	tableNameRegexp  = regexp.MustCompile(`(\w+)(\((.+)\))?`)
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
	parseDDL()
	structBytes := createStruct()
	if outputPath != "" {
		generateFileFromBytes(structBytes)
	} else {
		fmt.Printf("Generated your struct!!!\n=============\n%s\n=============\n", createStruct())
	}

	if copyToClipboard {
		clipboard.WriteAll(fmt.Sprintf("%s", structBytes))
	}
}

func getColumnToType(column *executor.ColumnDef) string {
	if !columnTypeRegexp.MatchString(column.Type) {
		return "interface{}"
	}

	parsedColumnType := columnTypeRegexp.FindStringSubmatch(column.Type)
	switch strings.ToLower(parsedColumnType[1]) {
	case "varchar":
		return "string"
	case "int", "tinyint":
		if column.Unsigned {
			return "uint"
		} else {
			return "int"
		}
	case "bigint":
		if column.Unsigned {
			return "uint64"
		} else {
			return "int64"
		}
	case "datetime", "timestamp":
		return "time.Time"
	default:
		return "interface{}"
	}
}

func parseDDL() {
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		panic(err)
	}

	sqlExecutor := executor.NewExecutor(executor.NewDefaultConfig())
	ddl := fmt.Sprintf("CREATE DATABASE ddl2struct; use ddl2struct; %s", string(data))
	if err := sqlExecutor.Exec(ddl); err != nil {
		panic(err)
	}
	tableNames, err := sqlExecutor.GetTables("ddl2struct")
	if err != nil {
		panic(err)
	}
	if len(tableNames) == 0 {
		panic(fmt.Errorf("no tables"))
	}
	tableDef, err := sqlExecutor.GetTableDef("ddl2struct", tableNames[0])
	if err != nil {
		panic(err)
	}
	tableName = tableDef.Name
	table = tableDef
}

func createFieldsByColumns() string {
	fields := make([]string, 0)
	for _, column := range table.Columns {
		fields = append(fields, fmt.Sprintf(
			"%s %s `json:\"%s\" gorm:\"Column:%s\"`",
			strcase.ToCamel(column.Name),
			getColumnToType(column),
			strcase.ToSnake(column.Name),
			column.Name,
		))
	}
	return strings.Join(fields, "\n")
}

func createStruct() []byte {
	s := fmt.Sprintf("type %s struct { \n %s \n}", strcase.ToCamel(tableName), createFieldsByColumns())
	formattedStructBytes, err := format.Source([]byte(s))
	if err != nil {
		panic(err)
	}
	return formattedStructBytes
}

func generateFileFromBytes(structBytes []byte) {
	if err := ioutil.WriteFile(outputPath, structBytes, 0644); err != nil {
		panic(err)
	}
}
