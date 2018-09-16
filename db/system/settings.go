package system

import (
	"encoding/json"
	"fmt"
	"github.com/Ready-Stock/badger"
	"github.com/kataras/go-errors"
)

var (
	settingsKeys = map[string]interface{}{
		"connection_pool_init_connections": 5,
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


