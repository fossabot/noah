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
	"testing"
)

func Test_Context_InitSetup(t *testing.T) {
	timestamp, err := SystemCtx.Settings.GetSettingUint64(InitialSetupTimestamp)
	if err != nil {
		panic(err)
	}

	// Since this is a new store, the value should be nil. If it is not then this should fail.
	if timestamp != nil {
		panic("InitialSetupTimestamp is not nil. This value should be nil in a new store.")
	}
}
