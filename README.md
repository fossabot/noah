# Noah
Noah is a database sharding tool for Postgres inspired by Vitess and Citus. 
At it's core is a gutted version of CockroachDB, rewritten to use multiple PostgreSQL databases as it's primary storage engine.
The query parsing of Cockroach has also been replaced by PostgreSQL's actual query parser in C. 
It communicates with the database servers using a bare-bone and customized version of PGX, with all non-essential features removed and two-phase transaction handling added.

