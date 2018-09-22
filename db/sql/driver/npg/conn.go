package npg

import (
	"github.com/Ready-Stock/Noah/db/sql/driver"
	"github.com/Ready-Stock/Noah/db/sql/pgwire/pgproto"
	"github.com/Ready-Stock/Noah/db/sql/types"
	"net"
	"sync"
	"time"
)

var (
// Errors
)

var (
	// Defaults
	minimalConnInfo *types.ConnInfo
)

const (
	connStatusUnitialized = iota
	connStatusClosed
	connStatusIdle
	connStatusBusy
)

func init() {
	minimalConnInfo = types.NewConnInfo()
	minimalConnInfo.InitializeDataTypes(map[string]types.OID{
		"int4":    types.Int4OID,
		"name":    types.NameOID,
		"oid":     types.OIDOID,
		"text":    types.TextOID,
		"varchar": types.VarcharOID,
	})
}

type DialFunc func(network, addr string) (net.Conn, error)

type Conn struct {
	conn     net.Conn
	config   driver.ConnConfig
	txStatus byte
	status   byte
	mux      sync.Mutex
	frontend *pgproto.Frontend

	ConnInfo *types.ConnInfo
}

func Connect(config driver.ConnConfig) (c *Conn, err error) {
	return connect(config, minimalConnInfo)
}

func connect(config driver.ConnConfig, connInfo *types.ConnInfo) (c *Conn, err error) {
	c = new(Conn)

	c.config = config
	c.ConnInfo = connInfo

	if c.config.User == "" {
		c.config.User = "postgres"
	}

	if c.config.Port == 0 {
		c.config.Port = 5432
	}

	network, address := c.config.NetworkAddress()
	d := defaultDialer()
	err := c.connect(config, network, address, d.Dial)
}

func (c *Conn) connect(config driver.ConnConfig, network, address string, dial DialFunc) (err error) {
	c.conn, err = dial(network, address)
	if err != nil {
		return err
	}
	defer func() {
		if c != nil && err != nil {
			c.conn.Close()
			c.mux.Lock()
			c.status = connStatusClosed
			c.mux.Unlock()
		}
	}()

	c.mux.Lock()
	c.status = connStatusIdle
	c.mux.Unlock()

	c.frontend, err = pgproto.NewFrontend(c.conn, c.conn)
	if err != nil {
		return err
	}

	startupMsg := pgproto.StartupMessage{
		ProtocolVersion: pgproto.ProtocolVersionNumber,
		Parameters:      make(map[string]string),
	}
	
}

func defaultDialer() *net.Dialer {
	return &net.Dialer{KeepAlive: 5 * time.Minute}
}
