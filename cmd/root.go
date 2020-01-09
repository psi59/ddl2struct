package cmd

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/atotto/clipboard"
	"github.com/iancoleman/strcase"
	"github.com/xwb1989/sqlparser"
)

var cfgFile string
var (
	inputPath       string
	outputPath      string
	tableName       string
	ddl             *sqlparser.DDL
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
	data, err := ioutil.ReadFile(inputPath)
	stmt, err := sqlparser.Parse(string(data))
	if err != nil {
		panic(err)
	}

	parsedDDL, ok := stmt.(*sqlparser.DDL)
	if !ok {
		panic("not create table")
	}

	ddl = parsedDDL
	tableName = ddl.NewName.Name.String()

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

func getColumnToType(column sqlparser.ColumnType) string {
	switch column.Type {
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

func createFieldsByColumns() string {
	columns := ddl.TableSpec.Columns
	fields := make([]string, 0)
	for _, column := range columns {
		fields = append(fields, fmt.Sprintf(
			"%s %s `json:\"%s\" gorm:\"Column:%s\"`",
			strcase.ToCamel(column.Name.String()),
			getColumnToType(column.Type),
			strcase.ToSnake(column.Name.String()),
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
