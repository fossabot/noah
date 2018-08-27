package system

import (
	"github.com/Ready-Stock/badger"
)

type SContext struct {
	Badger *badger.DB
}

func (sctx *SContext) LogError(err error) {

}

