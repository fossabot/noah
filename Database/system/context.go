package system

import (
	"github.com/Ready-Stock/badger"
)

type SContext struct {
	badger *badger.DB
}

func (sctx *SContext) LogError(err error) {

}

