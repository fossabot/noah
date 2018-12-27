package system

import (
    "github.com/readystock/pg_query_go/nodes"
)

type SQuery baseContext

func (ctx *SQuery) GetQueryType(stmt pg_query.Stmt) NQueryType {
    return 0
}
