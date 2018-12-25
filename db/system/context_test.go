package system

import (
	"testing"
)

func Test_CreateContext(t *testing.T) {
	tempFolder := CreateTempFolder()
	defer DeleteTempFolder(tempFolder)

	ctx, err := NewSystemContext(tempFolder, "127.0.0.1:0", "")
	if err != nil {
		panic(err)
	}

	ctx.Close()
}

func Test_GetEmptyUnit64(t *testing.T) {
	tempFolder := CreateTempFolder()
	defer DeleteTempFolder(tempFolder)

	ctx, err := NewSystemContext(tempFolder, "127.0.0.1:0", "")
	if err != nil {
		panic(err)
	}

	timestamp, err := ctx.Settings.GetSettingUint64(InitialSetupTimestamp)
	if err != nil {
		panic(err)
	}

	// Since this is a new store, the value should be nil. If it is not then this should fail.
	if timestamp != nil {
		panic("InitialSetupTimestamp is not nil. This value should be nil in a new store.")
	}

	ctx.Close()
}
