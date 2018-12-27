/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package sql

import (
    "github.com/juju/errors"
    "github.com/readystock/golinq"
    "github.com/readystock/golog"
    "github.com/readystock/noah/db/sql/plan"
    "github.com/readystock/noah/db/system"
    "github.com/readystock/pg_query_go/nodes"
    "strings"
)

var (
    ErrNotEnoughNodesAvailable = errors.New("not enough nodes available in cluster to create table")
    ErrTablespaceNotSpecified  = errors.New("tablespace must be specified when creating a table")
)

type CreateStatement struct {
    Statement pg_query.CreateStmt
    table     system.NTable
    IQueryStatement
}

func CreateCreateStatement(stmt pg_query.CreateStmt) *CreateStatement {
    return &CreateStatement{
        Statement: stmt,
    }
}

func (stmt *CreateStatement) Execute(ex *connExecutor, res RestrictedCommandResult) error {
    existingTable, _ := ex.SystemContext.Schema.GetTable(*stmt.Statement.Relation.Relname)
    if existingTable != nil {
        return errors.Errorf("table with name [%s] already exists in the cluster", *stmt.Statement.Relation.Relname)
    }

    targetNodes, err := stmt.getTargetNodes(ex)
    if err != nil {
        return err
    }

    plans, err := stmt.compilePlan(ex, targetNodes)
    if err != nil {
        return err
    }

    // if ex.TransactionState == NTXNoTransaction { // When creating a table try to create a new transaction if we are not in one
    //     if err := ex.execStmt(pg_query.TransactionStmt{Kind: pg_query.TRANS_STMT_BEGIN}, nil, 0); err != nil {
    //         return err
    //     }
    //     ex.TransactionStatus = NTXInProgress
    // }

    return ex.ExecutePlans(plans, res)
}

func (stmt *CreateStatement) getTargetNodes(ex *connExecutor) ([]system.NNode, error) {
    writeNodes, err := ex.SystemContext.Nodes.GetNodes()
    if err != nil {
        return nil, err
    }

    allNodes := make([]system.NNode, 0)
    linq.From(writeNodes).WhereT(func(node system.NNode) bool {
        return node.ReplicaOf == 0
    }).ToSlice(&allNodes)

    liveNodes := linq.From(allNodes).CountWithT(func(node system.NNode) bool {
        return node.IsAlive && node.ReplicaOf == 0
    })

    // Schema changes can only be made when all (non-replica) nodes are alive, if any nodes are
    // unavailable then the schema change will be rejected to ensure consistency.
    if liveNodes != len(allNodes) {
        return nil, ErrNotEnoughNodesAvailable
    }

    if liveNodes == 0 {
        return nil, errors.Errorf("no live nodes, ddl cannot be processed at this time")
    }

    return allNodes, nil
}

func (stmt *CreateStatement) compilePlan(ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error) {
    plans := make([]plan.NodeExecutionPlan, len(nodes))

    stmt.table = system.NTable{
        TableName: *stmt.Statement.Relation.Relname,
        Schema:    "default",
        Columns:   make([]*system.NColumn, 0),
    }

    // Determine the distribution of a table in the cluster
    if err := stmt.handleTableType(ex, &stmt.table); err != nil { // Handle sharding
        return nil, err
    }

    if err := stmt.handleValidation(ex, &stmt.table); err != nil {
        return nil, err
    }

    // Add handling here for custom column types.
    if err := stmt.handleColumns(ex, &stmt.table); err != nil { // Handle sharding
        return nil, err
    }

    // // Create the table in the coordinator cluster
    // if err := ex.SystemContext.Schema.CreateTable(table); err != nil {
    //     return nil, err
    // }

    compiled, err := pg_query.Deparse(stmt.Statement)
    if err != nil {
        golog.Error(err.Error())
        return nil, err
    }
    golog.Debugf("Recompiled query: %s", *compiled)
    for i := 0; i < len(plans); i++ {
        plans[i] = plan.NodeExecutionPlan{
            CompiledQuery: *compiled,
            Node:          nodes[i],
            ReadOnly:      false,
        }
    }
    return plans, nil
}

func (stmt *CreateStatement) handleValidation(ex *connExecutor, table *system.NTable) error {
    if table.TableType == system.NTableType_ACCOUNT {
        if accountsTable, err := ex.SystemContext.Schema.GetAccountsTable(); err != nil {
            return err
        } else if accountsTable != nil {
            return errors.Errorf("an accounts table named [%s] already exists in this cluster", accountsTable.TableName)
        }
    }

    return nil
}

func (stmt *CreateStatement) handleColumns(ex *connExecutor, table *system.NTable) error {
    verifyPrimaryKeyColumnType := func(column system.NColumn) error {
        switch column.ColumnTypeName {
        case "bigint", "int", "int8", "int4":
            // We only allow for these two types to be primary keys at this time.
            return nil
        default:
            // At the moment noah only supports integer column sharding.
            return errors.Errorf("column [%s] cannot be a primary key, a primary key must be an integer column", column.ColumnName)
        }
    }

    verifyForeignKeyColumn := func(table *system.NTable, column *system.NColumn, constraint pg_query.Constraint) error {
        if len(constraint.PkAttrs.Items) != 1 {
            return errors.Errorf("currently noah only supports single column foreign keys")
        }

        tableName := strings.ToLower(*constraint.Pktable.Relname)
        key := strings.ToLower(constraint.PkAttrs.Items[0].(pg_query.String).Str)

        // make sure table exists
        tbl, err := ex.SystemContext.Schema.GetTable(tableName)
        if err != nil {
            return err
        }

        if tbl == nil {
            return errors.Errorf("table with name [%s] does not exist", tableName)
        }

        if tbl.PrimaryKey == nil {
            return errors.Errorf("table [%s] does not have a primary key and cannot be used as a foreign key", tableName)
        }

        primaryKey := tbl.PrimaryKey.(*system.NTable_PKey).PKey

        if primaryKey.ColumnName != key {
            return errors.Errorf("table [%s] has primary key [%s] not [%s], foreign keys must be against a primary key", tableName, primaryKey.ColumnName, key)
        }

        foreignKey := system.NForeignKey{
            TableName:  tableName,
            ColumnName: key,
            IsShardKey: tbl.TableType == system.NTableType_ACCOUNT,
        }

        column.ForeignKey = &system.NColumn_FKey{
            FKey: &foreignKey,
        }

        if tbl.TableType == system.NTableType_ACCOUNT {
            table.ShardKey = &system.NTable_SKey{
                SKey: column,
            }
        }

        return nil
    }

    if stmt.Statement.TableElts.Items != nil && len(stmt.Statement.TableElts.Items) > 0 {
        for i, col := range stmt.Statement.TableElts.Items {
            switch tableItem := col.(type) {
            case pg_query.ColumnDef:
                noahColumn := &system.NColumn{ColumnName: *tableItem.Colname}
                table.Columns = append(table.Columns, noahColumn) // This will make it so we can assign to the index later
                // There are a few types that are handled by noah as a middle man; such as ID generation
                // because of this we want to replace serial columns with their base types since noah
                // will rewrite incoming queries to include an ID when performing inserts.
                if tableItem.TypeName != nil &&
                    tableItem.TypeName.Names.Items != nil &&
                    len(tableItem.TypeName.Names.Items) > 0 {
                    columnType := tableItem.TypeName.Names.Items[len(tableItem.TypeName.Names.Items)-1].(pg_query.String) // The last type name
                    golog.Verbosef("Processing column [%s] type [%s]", *tableItem.Colname, strings.ToLower(columnType.Str))
                    // This switch statement will handle any custom column types that we would like.
                    switch strings.ToLower(columnType.Str) {
                    case "serial": // Emulate 32 bit sequence
                        columnType.Str = "int"
                        noahColumn.IsSequence = true
                    case "snowflake": // Snowflakes are generated using twitters id system
                        noahColumn.IsSnowflake = true
                        fallthrough
                    case "bigserial": // Emulate 64 bit sequence
                        columnType.Str = "bigint"
                        noahColumn.IsSequence = true
                    default:
                        // Other column types wont be handled.
                    }
                    noahColumn.ColumnTypeName = strings.ToLower(columnType.Str)
                    tableItem.TypeName.Names.Items = []pg_query.Node{columnType}
                    stmt.Statement.TableElts.Items[i] = tableItem
                }

                // Check to see if this column is the primary key, primary keys will be used for tables
                // like account tables. If someone tries to create a table without a primary key an
                // error will be returned at this time.
                if tableItem.Constraints.Items != nil && len(tableItem.Constraints.Items) > 0 {
                    for _, c := range tableItem.Constraints.Items {
                        constraint := c.(pg_query.Constraint)
                        switch constraint.Contype {
                        case pg_query.CONSTR_PRIMARY:
                            if table.PrimaryKey != nil {
                                return errors.New("cannot define more than 1 primary key on a single table")
                            }

                            if err := verifyPrimaryKeyColumnType(*noahColumn); err != nil {
                                return err
                            }

                            table.PrimaryKey = &system.NTable_PKey{
                                PKey: noahColumn,
                            }

                            noahColumn.IsPrimaryKey = true
                        case pg_query.CONSTR_FOREIGN:
                            if err := verifyForeignKeyColumn(table, noahColumn, constraint); err != nil {
                                return err
                            }
                        }
                    }
                }

                table.Columns[i] = noahColumn
            case pg_query.Constraint:
                // Its possible for primary keys, foreign keys and identities to be defined
                // somewhere other than the column line itself, if this happens we still want to
                // handle it gracefully.
                switch tableItem.Contype {
                case pg_query.CONSTR_PRIMARY:
                    if len(tableItem.Keys.Items) != 1 {
                        return errors.Errorf("currently noah only supports single column primary keys")
                    }

                    // We want to search columns based on the column name.
                    key := tableItem.Keys.Items[0].(pg_query.String).Str

                    colIndex := linq.From(table.Columns).IndexOfT(func(column *system.NColumn) bool {
                        return strings.ToLower(column.ColumnName) == strings.ToLower(key)
                    })

                    if colIndex < 0 {
                        return errors.Errorf("could not use column [%s] as primary key, it is not defined in the create statement", key)
                    }

                    if err := verifyPrimaryKeyColumnType(*table.Columns[colIndex]); err != nil {
                        return err
                    }

                    table.Columns[colIndex].IsPrimaryKey = true
                    table.PrimaryKey = &system.NTable_PKey{
                        PKey: table.Columns[colIndex],
                    }
                case pg_query.CONSTR_FOREIGN:
                    if len(tableItem.FkAttrs.Items) != 1 {
                        return errors.Errorf("only 1 column can be used in a foreign key constraint")
                    }

                    key := tableItem.FkAttrs.Items[0].(pg_query.String).Str

                    colIndex := linq.From(table.Columns).IndexOfT(func(column *system.NColumn) bool {
                        return strings.ToLower(column.ColumnName) == strings.ToLower(key)
                    })

                    if err := verifyForeignKeyColumn(table, table.Columns[colIndex], tableItem); err != nil {
                        return err
                    }
                case pg_query.CONSTR_IDENTITY:
                }
            }
        }
    }

    if table.PrimaryKey == nil && table.TableType == system.NTableType_ACCOUNT {
        return errors.New("cannot create an account table without a primary key")
    }

    return nil
}

func (stmt *CreateStatement) handleTableType(ex *connExecutor, table *system.NTable) error {
    if stmt.Statement.Tablespacename != nil {
        switch strings.ToLower(*stmt.Statement.Tablespacename) {
        case "noah.global": // Table has the same data on all shards
            table.TableType = system.NTableType_GLOBAL
        case "noah.shard": // Table is sharded by shard column
            table.TableType = system.NTableType_SHARD
        case "noah.account": // Table contains all of the records of accounts for cluster
            table.TableType = system.NTableType_ACCOUNT
        default: // Other
            return ErrTablespaceNotSpecified
        }
    }
    return nil
}
