/*
 * Copyright (c) 2018 Ready Stock
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

package serverutils

import (
	gosql "database/sql"
	"github.com/readystock/noah/db/base"
	"github.com/readystock/noah/db/util/hlc"
	"github.com/readystock/noah/db/util/protoutil"
	"github.com/readystock/noah/db/util/stop"
	"net/http"
	"net/url"
	"testing"
)

// TestServerInterface defines test server functionality that tests need; it is
// implemented by server.TestServer.
type TestServerInterface interface {
	Stopper() *stop.Stopper

	Start(params base.TestServerArgs) error


	// ServingAddr returns the server's advertised address.
	ServingAddr() string

	// HTTPAddr returns the server's http address.
	HTTPAddr() string

	// Addr returns the server's address.
	Addr() string

	// DB returns a *client.DB instance for talking to this KV server.
	// DB() *client.DB

	// RPCContext returns the rpc context used by the test server.
	// RPCContext() *rpc.Context

	// LeaseManager() returns the *sql.LeaseManager as an interface{}.
	LeaseManager() interface{}

	// Executor() returns the *sql.Executor as an interface{}.
	Executor() interface{}

	// InternalExecutor returns a *sqlutil.InternalExecutor as an interface{}.
	InternalExecutor() interface{}

	// Gossip returns the gossip used by the TestServer.
	// Gossip() *gossip.Gossip

	// Clock returns the clock used by the TestServer.
	Clock() *hlc.Clock

	// DistSender returns the DistSender used by the TestServer.
	// DistSender() *kv.DistSender

	// DistSQLServer returns the *distsqlrun.ServerImpl as an interface{}.
	DistSQLServer() interface{}

	// JobRegistry returns the *jobs.Registry as an interface{}.
	JobRegistry() interface{}

	// SetDistSQLSpanResolver changes the SpanResolver used for DistSQL inside the
	// server's executor. The argument must be a distsqlplan.SpanResolver
	// instance.
	//
	// This method exists because we cannot pass the fake span resolver with the
	// server or cluster params: the fake span resolver needs the node IDs and
	// addresses of the servers in a cluster, which are not available before we
	// start the servers.
	//
	// It is the caller's responsibility to make sure no queries are being run
	// with DistSQL at the same time.
	SetDistSQLSpanResolver(spanResolver interface{})

	// AdminURL returns the URL for the admin UI.
	AdminURL() string
	// GetHTTPClient returns an http client configured with the client TLS
	// config required by the TestServer's configuration.
	GetHTTPClient() (http.Client, error)
	// GetAuthenticatedHTTPClient returns an http client which has been
	// authenticated to access Admin API methods (via a cookie).
	GetAuthenticatedHTTPClient() (http.Client, error)

	// MustGetSQLCounter returns the value of a counter metric from the server's
	// SQL Executor. Runs in O(# of metrics) time, which is fine for test code.
	MustGetSQLCounter(name string) int64
	// MustGetSQLNetworkCounter returns the value of a counter metric from the
	// server's SQL server. Runs in O(# of metrics) time, which is fine for test
	// code.
	MustGetSQLNetworkCounter(name string) int64
	// WriteSummaries records summaries of time-series data, which is required for
	// any tests that query server stats.
	WriteSummaries() error

	// GetFirstStoreID is a utility function returning the StoreID of the first
	// store on this node.
	// GetFirstStoreID() roachpb.StoreID

	// GetStores returns the collection of stores from this TestServer's node.
	// The return value is of type *storage.Stores.
	GetStores() interface{}

	// ClusterSettings returns the ClusterSettings shared by all components of
	// this server.
	// ClusterSettings() *cluster.Settings

	// SplitRange splits the range containing splitKey.
	// SplitRange(
	// 	splitKey roachpb.Key,
	// ) (left roachpb.RangeDescriptor, right roachpb.RangeDescriptor, err error)

	// ExpectedInitialRangeCount returns the expected number of ranges that should
	// be on the server after initial (asynchronous) splits have been completed,
	// assuming no additional information is added outside of the normal bootstrap
	// process.
	ExpectedInitialRangeCount() (int, error)
}

// TestServerFactory encompasses the actual implementation of the shim
// service.
type TestServerFactory interface {
	// New instantiates a test server.
	New(params base.TestServerArgs) interface{}
}

var srvFactoryImpl TestServerFactory

// InitTestServerFactory should be called once to provide the implementation
// of the service. It will be called from a xx_test package that can import the
// server package.
func InitTestServerFactory(impl TestServerFactory) {
	srvFactoryImpl = impl
}

// StartServer creates a test server and sets up a gosql DB connection.
// The server should be stopped by calling server.Stopper().Stop().
func StartServer(
	t testing.TB, params base.TestServerArgs,
) (TestServerInterface, *gosql.DB, *client.DB) {
	server, err := StartServerRaw(params)
	if err != nil {
		t.Fatal(err)
	}

	pgURL, cleanupGoDB := sqlutils.PGUrl(
		t, server.ServingAddr(), "StartServer", url.User(security.RootUser))
	pgURL.Path = params.UseDatabase
	if params.Insecure {
		pgURL.RawQuery = "sslmode=disable"
	}
	goDB, err := gosql.Open("postgres", pgURL.String())
	if err != nil {
		t.Fatal(err)
	}
	server.Stopper().AddCloser(
		stop.CloserFn(func() {
			_ = goDB.Close()
			cleanupGoDB()
		}))
	return server, goDB, server.DB()
}

// StartServerRaw creates and starts a TestServer.
// Generally StartServer() should be used. However this function can be used
// directly when opening a connection to the server is not desired.
func StartServerRaw(args base.TestServerArgs) (TestServerInterface, error) {
	if srvFactoryImpl == nil {
		panic("TestServerFactory not initialized. One needs to be injected " +
			"from the package's TestMain()")
	}
	server := srvFactoryImpl.New(args).(TestServerInterface)
	if err := server.Start(args); err != nil {
		return nil, err
	}
	return server, nil
}

// GetJSONProto uses the supplied client to GET the URL specified by the parameters
// and unmarshals the result into response.
func GetJSONProto(ts TestServerInterface, path string, response protoutil.Message) error {
	httpClient, err := ts.GetAuthenticatedHTTPClient()
	if err != nil {
		return err
	}
	return httputil.GetJSON(httpClient, ts.AdminURL()+path, response)
}

// PostJSONProto uses the supplied client to POST request to the URL specified by
// the parameters and unmarshals the result into response.
func PostJSONProto(ts TestServerInterface, path string, request, response protoutil.Message) error {
	httpClient, err := ts.GetAuthenticatedHTTPClient()
	if err != nil {
		return err
	}
	return httputil.PostJSON(httpClient, ts.AdminURL()+path, request, response)
}
