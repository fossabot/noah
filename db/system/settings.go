package system

import (
	"encoding/json"
	"fmt"
	"github.com/Ready-Stock/badger"
	"github.com/kataras/go-errors"
	"strconv"
)

var (
	settingsKeys = map[string]interface{}{
		"connection_pool_init_connections": int64(5),
	}
)

func (ctx *SContext) SetSetting(SettingName string, SettingValue interface{}) (error) {
	if _, ok := settingsKeys[SettingName]; !ok {
		return errors.New("setting key `%s` is not valid and has not been set.").Format(SettingName)
	}
	return ctx.Badger.Update(func(txn *badger.Txn) error {
		if j, err := json.Marshal(SettingValue); err != nil {
			return err
		} else{
			return txn.Set([]byte(fmt.Sprintf("%s%s", SettingsPath, SettingName)), j)
		}
	})
}

func (ctx *SContext) GetSettingInt64(SettingName string) (*int64, error) {
	if _, ok := settingsKeys[SettingName]; !ok {
		return nil, errors.New("setting key `%s` is not valid and cannot be returned.").Format(SettingName)
	}
	var i *int64 = nil
	return i, ctx.Badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(fmt.Sprintf("%s%s", SettingsPath, SettingName)))
		if err != nil {
			return err
		}
		if vs, err := item.Value(); err != nil {
			return err
		} else {
			x, err := strconv.ParseInt(string(vs), 10, 64)
			i = &x
			return err
		}
	})
}


