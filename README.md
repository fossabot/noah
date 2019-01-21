# Noah 
[![Build Status](https://travis-ci.com/readystock/noah.svg?token=QvXZjJzgiir2JHLaKFrG&branch=master)](https://travis-ci.com/readystock/noah)
[![CodeFactor](https://www.codefactor.io/repository/github/readystock/noah/badge)](https://www.codefactor.io/repository/github/readystock/noah)

Noah is a database sharding tool for Postgres inspired by Vitess and Citus. 
At it's core is a gutted version of CockroachDB, rewritten to use multiple PostgreSQL databases as it's primary storage engine.
The query parsing of Cockroach has also been replaced by PostgreSQL's actual query parser in C. 
It communicates with the database servers using a bare-bone and customized version of PGX, with all non-essential features removed and two-phase transaction handling added.


## Building

Noah requires `protoc` to be built. Protoc will generate the protos needed for
some of the key value store data and wire protocol. Cockroach had some protos
that have been "hard-coded" into noah to make builds easier. These protos also
won't be modified so there should be no need to regenerate them.

Any tools that you might need can be found in [tools](./docs/Tools.md).

To build noah.
```bash
go get -d -u github.com/readystock/noah
cd $GOPATH/src/github.com/readystock/noah
make
```

This will create a `bin` folder if it does not already exist and create the noah executable within 
that folder.

#### Clean build
If you need to do a completely clean build of noah run the following:
```bash
make fresh
```
This will recompile all of noah's dependencies and may take a few minutes.