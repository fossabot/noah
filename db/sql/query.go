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

package sql

import (
	"context"
	"github.com/readystock/noah/db/sql/plan"
	"github.com/readystock/noah/db/system"
)

type IQueryStatement interface {
	Execute(ctx context.Context, ex *connExecutor, res RestrictedCommandResult, pinfo *plan.PlaceholderInfo) error
	compilePlan(ctx context.Context, ex *connExecutor, nodes []system.NNode) ([]plan.NodeExecutionPlan, error)
	getTargetNodes(ctx context.Context, ex *connExecutor) ([]system.NNode, error)
}
