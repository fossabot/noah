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

package pgwire

import (
	"github.com/readystock/noah/db/system"
	"github.com/readystock/noah/testutils"
	"os"
	"testing"
)

var (
	SystemCtx *system.SContext
)

func TestMain(m *testing.M) {
	tempFolder := testutils.CreateTempFolder()
	defer testutils.DeleteTempFolder(tempFolder)

	sctx, err := system.NewSystemContext(tempFolder, "127.0.0.1:0", "", "")
	if err != nil {
		panic(err)
	}
	SystemCtx = sctx
	defer SystemCtx.Close()

	retCode := m.Run()
	os.Exit(retCode)
}
