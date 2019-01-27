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

package testsuite

import (
	"testing"
)

func Test_CurrentDatabase(t *testing.T) {
	DoQueryTest(t, QueryTest{
		Query: "SELECT current_database(), current_schema()",
		Expected: [][]interface{}{
			{"noah", "public"},
		},
	})
}

func Test_DataGrip_CommonQueries(t *testing.T) {
	t.Skip("cannot properly decode arrays yet.")
	DoQueryTest(t, QueryTest{
		Query: "select current_database() as a, current_schemas(false) as b",
		Expected: [][]interface{}{
			{"noah", []string{"public"}},
		},
	})
}
