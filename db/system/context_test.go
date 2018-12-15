package system

import (
	"testing"
)

func Test_CreateContext(t *testing.T) {
	tempFolder := CreateTempFolder()
	defer DeleteTempFolder(tempFolder)

	ctx, err := NewSystemContext()
	if err != nil {
		panic(err)
	}
	ctx.Close()
}
