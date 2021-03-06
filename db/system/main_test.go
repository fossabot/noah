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
	"github.com/readystock/golog"
	"github.com/readystock/noah/testutils"
	"os"
	"testing"
)

var (
	SystemCtx *SContext
	Nodes     = []NNode{
		{
			Address:   "127.0.0.1:0",
			Port:      5432,
			Database:  "postgres",
			User:      "postgres",
			Password:  "",
			ReplicaOf: 0,
			Region:    "",
			Zone:      "",
			IsAlive:   true,
		},
	}
)

func TestMain(m *testing.M) {
	golog.SetLevel("verbose")
	tempFolder := testutils.CreateTempFolder()
	defer testutils.DeleteTempFolder(tempFolder)

	sctx, err := NewSystemContext(tempFolder, "127.0.0.1:0", "", "")
	if err != nil {
		panic(err)
	}
	SystemCtx = sctx
	defer SystemCtx.Close()

	for _, node := range Nodes {
		if _, err := SystemCtx.Nodes.AddNode(node); err != nil {
			panic(err)
		}
	}

	retCode := m.Run()
	os.Exit(retCode)
}
