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

package npgx

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Connection interface {
	QueryBytes(sql string, args ...interface{}) ([][]byte, error)
}

type ConnConfig struct {
	NodeId   uint64 // NodeId is just used by noah for managing connection pools.
	Host     string // host (e.g. localhost) or path to unix domain socket directory (e.g. /private/tmp)
	Port     uint16 // default: 5432
	Database string
	User     string
	Password string
}

// ParseConnectionString parses either a URI or a DSN connection string.
// see ParseURI and ParseDSN for details.
func ParseConnectionString(s string) (ConnConfig, error) {
	if u, err := url.Parse(s); err == nil && u.Scheme != "" {
		return ParseURI(s)
	}
	return ParseDSN(s)
}

// ParseURI parses a database URI into ConnConfig
//
// Query parameters not used by the connection process are parsed into ConnConfig.RuntimeParams.
func ParseURI(uri string) (ConnConfig, error) {
	var cp ConnConfig

	url, err := url.Parse(uri)
	if err != nil {
		return cp, err
	}

	if url.User != nil {
		cp.User = url.User.Username()
		cp.Password, _ = url.User.Password()
	}

	parts := strings.SplitN(url.Host, ":", 2)
	cp.Host = parts[0]
	if len(parts) == 2 {
		p, err := strconv.ParseUint(parts[1], 10, 16)
		if err != nil {
			return cp, err
		}
		cp.Port = uint16(p)
	}
	cp.Database = strings.TrimLeft(url.Path, "/")

	if pgtimeout := url.Query().Get("connect_timeout"); pgtimeout != "" {
		timeout, err := strconv.ParseInt(pgtimeout, 10, 64)
		if err != nil {
			return cp, err
		}
		d := defaultDialer()
		d.Timeout = time.Duration(timeout) * time.Second
		// cp.Dial = d.Dial
	}

	// tlsArgs := configTLSArgs{
	//     sslCert:     url.Query().Get("sslcert"),
	//     sslKey:      url.Query().Get("sslkey"),
	//     sslMode:     url.Query().Get("sslmode"),
	//     sslRootCert: url.Query().Get("sslrootcert"),
	// }
	// err = configTLS(tlsArgs, &cp)
	// if err != nil {
	//     return cp, err
	// }

	ignoreKeys := map[string]struct{}{
		"connect_timeout": {},
		"sslcert":         {},
		"sslkey":          {},
		"sslmode":         {},
		"sslrootcert":     {},
	}

	// cp.RuntimeParams = make(map[string]string)

	for k, v := range url.Query() {
		if _, ok := ignoreKeys[k]; ok {
			continue
		}

		if k == "host" {
			cp.Host = v[0]
			continue
		}

		// cp.RuntimeParams[k] = v[0]
	}
	// if cp.Password == "" {
	//     pgpass(&cp)
	// }
	return cp, nil
}

var dsnRegexp = regexp.MustCompile(`([a-zA-Z_]+)=((?:"[^"]+")|(?:[^ ]+))`)

// ParseDSN parses a database DSN (data source name) into a ConnConfig
//
// e.g. ParseDSN("user=username password=password host=1.2.3.4 port=5432 dbname=mydb sslmode=disable")
//
// Any options not used by the connection process are parsed into ConnConfig.RuntimeParams.
//
// e.g. ParseDSN("application_name=pgxtest search_path=admin user=username password=password host=1.2.3.4 dbname=mydb")
//
// ParseDSN tries to match libpq behavior with regard to sslmode. See comments
// for ParseEnvLibpq for more information on the security implications of
// sslmode options.
func ParseDSN(s string) (ConnConfig, error) {
	var cp ConnConfig

	m := dsnRegexp.FindAllStringSubmatch(s, -1)

	// tlsArgs := configTLSArgs{}

	// cp.RuntimeParams = make(map[string]string)

	for _, b := range m {
		switch b[1] {
		case "user":
			cp.User = b[2]
		case "password":
			cp.Password = b[2]
		case "host":
			cp.Host = b[2]
		case "port":
			p, err := strconv.ParseUint(b[2], 10, 16)
			if err != nil {
				return cp, err
			}
			cp.Port = uint16(p)
		case "dbname":
			cp.Database = b[2]
		// case "sslmode":
		//     tlsArgs.sslMode = b[2]
		// case "sslrootcert":
		//     tlsArgs.sslRootCert = b[2]
		// case "sslcert":
		//     tlsArgs.sslCert = b[2]
		// case "sslkey":
		//     tlsArgs.sslKey = b[2]
		case "connect_timeout":
			timeout, err := strconv.ParseInt(b[2], 10, 64)
			if err != nil {
				return cp, err
			}
			d := defaultDialer()
			d.Timeout = time.Duration(timeout) * time.Second
			// cp.Dial = d.Dial
		default:
			// cp.RuntimeParams[b[1]] = b[2]
		}
	}

	// err := configTLS(tlsArgs, &cp)
	// if err != nil {
	//     return cp, err
	// }
	// if cp.Password == "" {
	//     pgpass(&cp)
	// }
	return cp, nil
}

func (cc *ConnConfig) NetworkAddress() (network, address string) {
	network = "tcp"
	address = fmt.Sprintf("%s:%d", cc.Host, cc.Port)
	// Make sure that the provided address is balid.
	if _, err := os.Stat(cc.Host); err == nil {
		network = "unix"
		address = cc.Host
		if !strings.Contains(address, "/.s.PGSQL.") {
			address = filepath.Join(address, ".s.PGSQL.") + strconv.FormatInt(int64(cc.Port), 10)
		}
	}
	return network, address
}
