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

package sqlbase

import (
	"github.com/readystock/noah/db/sql/pgwire/pgerror"
)

const (
	txnAbortedMsg = "current transaction is aborted, commands ignored " +
		"until end of transaction block"
	txnCommittedMsg = "current transaction is committed, commands ignored " +
		"until end of transaction block"
	txnRetryMsgPrefix = "restart transaction"
)

func NewRetryError(cause error) error {
	return pgerror.NewErrorf(
		pgerror.CodeSerializationFailureError, "%s: %s", txnRetryMsgPrefix, cause)
}

func NewError(cause error) error {
	return pgerror.NewErrorf(
		pgerror.CodeSerializationFailureError, "%s", cause)
}
