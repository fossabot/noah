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
	"encoding/json"
	"fmt"
	"github.com/kataras/go-errors"
	"github.com/readystock/noah/db/sql/types"
	"strconv"
)

var (
	settingsKeys = map[NoahSetting]interface{}{
		ConnectionPoolInitConnections: int64(5),
		QueryReplicationFactor:        int64(2),
		InitialSetupTimestamp:         nil,
	}

	settingsInfo = map[NoahSetting]struct {
		Setting     NoahSetting
		Description string
		Type        types.DataType
	}{
		ConnectionPoolInitConnections: {
			Setting:     ConnectionPoolInitConnections,
			Description: "Number of connections to initialize in a connection pool.",
		},
	}
)

type SSettings baseContext

type NoahSetting string

const (
	ConnectionPoolInitConnections NoahSetting = "connection_pool_init_connections"
	QueryReplicationFactor        NoahSetting = "query_replication_factor"
	InitialSetupTimestamp         NoahSetting = "initial_setup_timestamp"
)

func (ctx *SSettings) SetSetting(SettingName string, SettingValue interface{}) error {
	if _, ok := settingsKeys[NoahSetting(SettingName)]; !ok {
		return errors.New("setting key `%s` is not valid and has not been set.").Format(SettingName)
	}
	if j, err := json.Marshal(SettingValue); err != nil {
		return err
	} else {
		return ctx.db.Set([]byte(fmt.Sprintf("%s%s", settingsExternalPath, SettingName)), j)
	}
}

func (ctx *SSettings) SetSettingString(SettingName, SettingValue string) error {
	if _, ok := settingsKeys[NoahSetting(SettingName)]; !ok {
		return errors.New("setting key `%s` is not valid and has not been set.").Format(SettingName)
	}
	return ctx.db.Set([]byte(fmt.Sprintf("%s%s", settingsExternalPath, SettingName)), []byte(SettingValue))
}

func (ctx *SSettings) GetSetting(SettingName NoahSetting) (*string, error) {
	if _, ok := settingsKeys[NoahSetting(SettingName)]; !ok {
		return nil, errors.New("setting key `%s` is not valid and cannot be returned.").Format(SettingName)
	}
	value, err := ctx.db.Get([]byte(fmt.Sprintf("%s%s", settingsExternalPath, SettingName)))
	if err != nil {
		return nil, err
	}
	valueString := string(value)
	return &valueString, nil
}

func (ctx *SSettings) GetSettingInt64(SettingName NoahSetting) (*int64, error) {
	if _, ok := settingsKeys[NoahSetting(SettingName)]; !ok {
		return nil, errors.New("setting key `%s` is not valid and cannot be returned.").Format(SettingName)
	}
	value, err := ctx.db.Get([]byte(fmt.Sprintf("%s%s", settingsExternalPath, SettingName)))
	if err != nil {
		return nil, err
	}
	number, err := strconv.ParseInt(string(value), 10, 64)
	if err != nil {
		val := settingsKeys[SettingName].(int64)
		return &val, nil
	}
	return &number, nil
}

func (ctx *SSettings) GetSettingUint64(SettingName NoahSetting) (*uint64, error) {
	if _, ok := settingsKeys[NoahSetting(SettingName)]; !ok {
		return nil, errors.New("setting key `%s` is not valid and cannot be returned.").Format(SettingName)
	}
	value, err := ctx.db.Get([]byte(fmt.Sprintf("%s%s", settingsExternalPath, SettingName)))
	if err != nil {
		return nil, err
	}

	if len(value) == 0 {
		if def := settingsKeys[SettingName]; def == nil {
			return nil, nil
		} else {
			val := def.(uint64)
			return &val, nil
		}
	}

	number, err := strconv.ParseUint(string(value), 10, 64)
	if err != nil {
		return nil, err
	}
	return &number, nil
}
