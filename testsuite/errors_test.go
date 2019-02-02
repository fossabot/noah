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
	"fmt"
	"testing"
)

func Test_Errors_Insert_MultipleAccounts(t *testing.T) {
	DoExecTest(t, ExecTest{
		Query: `DROP TABLE IF EXISTS public.accounts CASCADE;`,
	})
	DoExecTest(t, ExecTest{
		Query: `CREATE TABLE public.accounts (account_id BIGSERIAL PRIMARY KEY, account_number BIGINT) TABLESPACE "noah.account";`,
	})
	defer DoExecTest(t, ExecTest{
		Query: `DROP TABLE public.accounts CASCADE;`,
	})

	DoExecTest(t, ExecTest{
		Query: `DROP TABLE IF EXISTS public.users CASCADE;`,
	})
	DoExecTest(t, ExecTest{
		Query: `CREATE TABLE public.users (user_id BIGSERIAL PRIMARY KEY, account_id BIGINT NOT NULL REFERENCES public.accounts (account_id), user_number BIGINT) TABLESPACE "noah.shard";`,
	})
	defer DoExecTest(t, ExecTest{
		Query: `DROP TABLE public.users CASCADE;`,
	})

	DoQueryTest(t, QueryTest{
		Query:         `INSERT INTO public.users (account_id, user_number) VALUES (1, 1), (2, 1);`,
		ExpectedError: fmt.Errorf("ERROR: cannot insert into more than 1 account ID at this time (SQLSTATE 40001)"),
	})
}
