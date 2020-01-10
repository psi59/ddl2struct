package parser

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/pingcap/errors"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/types"
)

type DDLParser struct {
	TableName string
	Columns   Columns
	err       error
	p         *parser.Parser
}

func (parser *DDLParser) Parse(sql string) error {
	nodes, _, err := parser.p.Parse(sql, "", "")
	if err != nil {
		return errors.Wrap(err, "sql parsing error")
	}

	for _, node := range nodes {
		node.Accept(parser)
		if parser.err != nil {
			return errors.Wrap(err, "sql parsing error")
		}
	}

	return nil
}

func (parser DDLParser) ToStruct(withTag bool) ([]byte, error) {
	s := fmt.Sprintf("type %s struct { %s }", parser.TableName, parser.Columns.ToStructFields(withTag))
	return format.Source([]byte(s))
}

func (parser *DDLParser) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	switch n := n.(type) {
	case *ast.CreateTableStmt:
		parser.err = parser.parseCreateTableStmt(n)
	}
	return n, true
}

func (parser *DDLParser) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func (parser *DDLParser) parseCreateTableStmt(stmt *ast.CreateTableStmt) error {
	parser.TableName = stmt.Table.Name.String()
	for _, col := range stmt.Cols {
		parser.Columns = append(parser.Columns, Column{
			Name: col.Name.Name.String(),
			Type: parser.getColumnType(col.Tp.EvalType()),
		})
	}

	return nil
}

func (parser *DDLParser) getColumnType(typ types.EvalType) string {
	switch typ {
	case types.ETInt:
		return "int"
	case types.ETReal, types.ETDecimal:
		return "float64"
	case types.ETDatetime, types.ETTimestamp:
		return "time.Time"
	default:
		return "string"
	}
}

func New() *DDLParser {
	return &DDLParser{
		p: parser.New(),
	}
}

type Columns []Column

func (columns Columns) ToStructFields(withTag bool) string {
	fields := make([]string, 0)
	for _, column := range columns {
		fields = append(fields, column.ToStructField(withTag))
	}
	return strings.Join(fields, "\n")
}

type Column struct {
	Name string
	Type string
}

func (column Column) ToStructField(withTag bool) string {
	var tag string
	if withTag {
		tag = fmt.Sprintf("`json:\"%s\" gorm:\"Column:%s\"`", strcase.ToSnake(column.Name), column.Name)
	}
	return fmt.Sprintf("%s %s", strcase.ToCamel(column.Name), column.Type) + tag
}
