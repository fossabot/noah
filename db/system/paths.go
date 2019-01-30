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

package system

import (
	"fmt"
	"strings"
)

const (
	accountsPath         = "/noah/accounts/meta/"
	accountsNodesPath    = "/noah/accounts/nodes/"
	nodesPath            = "/noah/nodes/meta/"
	shardsPath           = "/noah/shards/meta/"
	shardNodesPath       = "/noah/shards/nodes/"
	nodeShardsPath       = "/noah/nodes/shards/"
	tablesPath           = "/noah/schema/tables/"
	viewsPath            = "/noah/schema/views/"
	settingsExternalPath = "/noah/settings/public/"
	settingsInternalPath = "/noah/settings/internal/"
)

func getAccountsPath() []byte {
	return []byte(accountsPath)
}

func getAccountPath(id uint64) []byte {
	return []byte(fmt.Sprintf("%s%d", accountsPath, id))
}

func getAccountsNodesPath() []byte {
	return []byte(accountsNodesPath)
}

func getAccountsNodesAccountPath(accountId uint64) []byte {
	return []byte(fmt.Sprintf("%s%d/", accountsNodesPath, accountId))
}

func getAccountsNodesAccountNodePath(accountId, nodeId uint64) []byte {
	return []byte(fmt.Sprintf("%s%d/%d/", accountsNodesPath, accountId, nodeId))
}

func getTablesPath() []byte {
	return []byte(tablesPath)
}

func getTablePath(tableName string) []byte {
	return []byte(fmt.Sprintf("%s%s", tablesPath, strings.ToLower(tableName)))
}
