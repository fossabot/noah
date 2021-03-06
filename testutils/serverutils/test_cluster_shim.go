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

package serverutils

// import (
// 	gosql "database/sql"
// 	"testing"
// )
//
// // TestClusterInterface defines TestCluster functionality used by tests.
// type TestClusterInterface interface {
// 	NumServers() int
//
// 	// Server returns the TestServerInterface corresponding to a specific node.
// 	Server(idx int) TestServerInterface
//
// 	// ServerConn returns a gosql.DB connection to a specific node.
// 	ServerConn(idx int) *gosql.DB
//
// 	// StopServer stops a single server.
// 	StopServer(idx int)
// }
//
// // TestClusterFactory encompasses the actual implementation of the shim
// // service.
// type TestClusterFactory interface {
// 	// New instantiates a test server.
// 	StartTestCluster(t testing.TB, numNodes int, args base.TestClusterArgs) TestClusterInterface
// }
//
// var clusterFactoryImpl TestClusterFactory
//
// // InitTestClusterFactory should be called once to provide the implementation
// // of the service. It will be called from a xx_test package that can import the
// // server package.
// func InitTestClusterFactory(impl TestClusterFactory) {
// 	clusterFactoryImpl = impl
// }
//
// // StartTestCluster starts up a TestCluster made up of numNodes in-memory
// // testing servers. The cluster should be stopped using Stopper().Stop().
// func StartTestCluster(t testing.TB, numNodes int, args base.TestClusterArgs) TestClusterInterface {
// 	if clusterFactoryImpl == nil {
// 		panic("TestClusterFactory not initialized. One needs to be injected " +
// 			"from the package's TestMain()")
// 	}
// 	return clusterFactoryImpl.StartTestCluster(t, numNodes, args)
// }
