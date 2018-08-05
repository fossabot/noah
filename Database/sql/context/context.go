package context

import (
	"github.com/Ready-Stock/Noah/Database/sql"
)

type Context struct {
	ClientComm sql.ClientComm
}
