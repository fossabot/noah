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

package log_test

import (
	"github.com/readystock/noah/db/security"
	"github.com/readystock/noah/db/security/securitytest"
	"github.com/readystock/noah/db/settings/cluster"
	"github.com/readystock/noah/db/util/log"
	"github.com/readystock/noah/db/util/randutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	randutil.SeedForTests()
	security.SetAssetLoader(securitytest.EmbeddedAssets)
	// serverutils.InitTestServerFactory(server.TestServerFactory)

	// MakeTestingClusterSettings initializes log.ReportingSettings to this
	// instance of setting values.
	st := cluster.MakeTestingClusterSettings()
	log.DiagnosticsReportingEnabled.Override(&st.SV, false)
	log.CrashReports.Override(&st.SV, false)

	os.Exit(m.Run())
}
