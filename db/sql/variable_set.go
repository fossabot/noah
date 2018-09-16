package sql

import (
	pg_query2 "github.com/Ready-Stock/pg_query_go"
	"github.com/Ready-Stock/pg_query_go/nodes"
	"strings"
)

type VariableSetStatement struct {
	Statement pg_query.VariableSetStmt
	IQueryStatement
}

func CreateVariableSetStatement(stmt pg_query.VariableSetStmt) *VariableSetStatement {
	return &VariableSetStatement{
		Statement: stmt,
	}
}

func (stmt *VariableSetStatement) Execute(ex *connExecutor, res RestrictedCommandResult) error {
	if strings.HasPrefix(strings.ToLower(*stmt.Statement.Name), "noah") {
		setting_name := strings.Replace(strings.ToLower(*stmt.Statement.Name), "noah.", "", 1)
		setting_value, err := pg_query2.DeparseValue(stmt.Statement.Args.Items[0].(pg_query.A_Const))
		if err != nil {
			return err
		}
		return ex.SystemContext.SetSetting(setting_name, setting_value)
	}

	return nil
	target_nodes, err := stmt.getTargetNodes(ex)
	if err != nil {
		return err
	}

	plans, err := stmt.compilePlan(ex, target_nodes)
	if err != nil {
		return err
	}

	return ex.ExecutePlans(plans)
}
