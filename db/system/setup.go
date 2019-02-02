/*
 * Copyright (c) 2019 Ready Stock
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

package system

import (
	"github.com/readystock/golinq"
	"github.com/readystock/golog"
)

type SSetup baseContext

func (ctx *SSetup) DoSetup() error {
	if !ctx.db.IsLeader() {
		// If the current node is not the leader then we don't want to run this.
		return nil
	}
	schema := SSchema(*ctx)
	tables, err := schema.GetTables()
	if err != nil {
		return err
	}

	nonExistingTables := make([]NTable, 0)
	linq.From(postgresTables).WhereT(func(ptable NTable) bool {
		return !linq.From(tables).AnyWithT(func(table NTable) bool {
			return table.TableName == ptable.TableName
		})
	}).ToSlice(&nonExistingTables)

	for _, table := range nonExistingTables {
		golog.Debugf("system table [%s] does not exist, creating...", table.TableName)
		if err := schema.CreateTable(table); err != nil {
			golog.Errorf("could not create table [%s]: %s", table.TableName, err.Error())
		}
	}

	return nil
}
