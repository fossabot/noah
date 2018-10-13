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
 *
 * This application uses Open Source components. You can find the
 * source code of their open source projects along with license
 * information below. We acknowledge and are grateful to these
 * developers for their contributions to open source.
 *
 * Project: CockroachDB https://github.com/cockroachdb/cockroach
 * Copyright 2018 The Cockroach Authors.
 * License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE
 *
 * Project: Vitess https://github.com/vitessio/vitess
 * Copyright 2018 Google Inc.
 * License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE
 *
 * Project: Citus https://github.com/citusdata/citus
 * Copyright 2018 Citus Data, Inc.
 * License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE
 *
 * Project: pg_query_go https://github.com/lfittl/pg_query_go
 * Copyright 2018 Lukas Fittl
 * License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE
 *
 * Project: pgx https://github.com/jackc/pgx
 * Copyright 2018 Jack Christensen
 * License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE
 *
 * Project: BadgerDB https://github.com/dgraph-io/badger
 * Copyright 2018 Dgraph Labs, Inc. and Contributors
 * License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE
 *
 * Project: Sonyflake https://github.com/sony/sonyflake
 * Copyright 2018 Sony Corporation
 * License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE
 */

package system

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/go-errors"
	"strconv"
)

var (
	settingsKeys = map[NoahSetting]interface{}{
		ConnectionPoolInitConnections: int64(5),
	}
)

type SSettings baseContext

type NoahSetting string

const (
	ConnectionPoolInitConnections NoahSetting = "connection_pool_init_connections"
)

func (ctx *SSettings) SetSetting(SettingName string, SettingValue interface{}) (error) {
	if _, ok := settingsKeys[NoahSetting(SettingName)]; !ok {
		return errors.New("setting key `%s` is not valid and has not been set.").Format(SettingName)
	}
	if j, err := json.Marshal(SettingValue); err != nil {
		return err
	} else{
		return ctx.db.Set([]byte(fmt.Sprintf("%s%s", SettingsPath, SettingName)), j)
	}
}

func (ctx *SSettings) GetSettingInt64(SettingName string) (*int64, error) {
	if _, ok := settingsKeys[NoahSetting(SettingName)]; !ok {
		return nil, errors.New("setting key `%s` is not valid and cannot be returned.").Format(SettingName)
	}
	value, err := ctx.db.Get([]byte(fmt.Sprintf("%s%s", SettingsPath, SettingName)))
	if err != nil {
		return nil, err
	}
	number, err := strconv.ParseInt(string(value), 10, 64)
	if err != nil {
		return nil, err
	}
	return &number, nil
}


