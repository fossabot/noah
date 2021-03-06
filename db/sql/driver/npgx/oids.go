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
	"github.com/readystock/noah/db/sql/types"
)

var (
	NameOIDs = map[string]types.OID{
		"bool":             16,
		"bytea":            17,
		"char":             18,
		"name":             19,
		"int8":             20,
		"int2":             21,
		"int2vector":       22,
		"int4":             23,
		"regproc":          24,
		"text":             25,
		"oid":              26,
		"tid":              27,
		"xid":              28,
		"cid":              29,
		"oidvector":        30,
		"json":             114,
		"xml":              142,
		"_xml":             143,
		"_json":            199,
		"pg_node_tree":     194,
		"pg_ndistinct":     3361,
		"pg_dependencies":  3402,
		"pg_ddl_command":   32,
		"smgr":             210,
		"point":            600,
		"lseg":             601,
		"path":             602,
		"box":              603,
		"polygon":          604,
		"line":             628,
		"_line":            629,
		"float4":           700,
		"float8":           701,
		"abstime":          702,
		"reltime":          703,
		"tinterval":        704,
		"unknown":          705,
		"circle":           718,
		"_circle":          719,
		"money":            790,
		"_money":           791,
		"macaddr":          829,
		"inet":             869,
		"cidr":             650,
		"macaddr8":         774,
		"_bool":            1000,
		"_bytea":           1001,
		"_char":            1002,
		"_name":            1003,
		"_int2":            1005,
		"_int2vector":      1006,
		"_int4":            1007,
		"_regproc":         1008,
		"_text":            1009,
		"_oid":             1028,
		"_tid":             1010,
		"_xid":             1011,
		"_cid":             1012,
		"_oidvector":       1013,
		"_bpchar":          1014,
		"_varchar":         1015,
		"_int8":            1016,
		"_point":           1017,
		"_lseg":            1018,
		"_path":            1019,
		"_box":             1020,
		"_float4":          1021,
		"_float8":          1022,
		"_abstime":         1023,
		"_reltime":         1024,
		"_tinterval":       1025,
		"_polygon":         1027,
		"aclitem":          1033,
		"_aclitem":         1034,
		"_macaddr":         1040,
		"_macaddr8":        775,
		"_inet":            1041,
		"_cidr":            651,
		"_cstring":         1263,
		"bpchar":           1042,
		"varchar":          1043,
		"date":             1082,
		"time":             1083,
		"timestamp":        1114,
		"_timestamp":       1115,
		"_date":            1182,
		"_time":            1183,
		"timestamptz":      1184,
		"_timestamptz":     1185,
		"interval":         1186,
		"_interval":        1187,
		"_numeric":         1231,
		"timetz":           1266,
		"_timetz":          1270,
		"bit":              1560,
		"_bit":             1561,
		"varbit":           1562,
		"_varbit":          1563,
		"numeric":          1700,
		"refcursor":        1790,
		"_refcursor":       2201,
		"regprocedure":     2202,
		"regoper":          2203,
		"regoperator":      2204,
		"regclass":         2205,
		"regtype":          2206,
		"regrole":          4096,
		"regnamespace":     4089,
		"_regprocedure":    2207,
		"_regoper":         2208,
		"_regoperator":     2209,
		"_regclass":        2210,
		"_regtype":         2211,
		"_regrole":         4097,
		"_regnamespace":    4090,
		"uuid":             2950,
		"_uuid":            2951,
		"pg_lsn":           3220,
		"_pg_lsn":          3221,
		"tsvector":         3614,
		"gtsvector":        3642,
		"tsquery":          3615,
		"regconfig":        3734,
		"regdictionary":    3769,
		"_tsvector":        3643,
		"_gtsvector":       3644,
		"_tsquery":         3645,
		"_regconfig":       3735,
		"_regdictionary":   3770,
		"jsonb":            3802,
		"_jsonb":           3807,
		"txid_snapshot":    2970,
		"_txid_snapshot":   2949,
		"int4range":        3904,
		"_int4range":       3905,
		"numrange":         3906,
		"_numrange":        3907,
		"tsrange":          3908,
		"_tsrange":         3909,
		"tstzrange":        3910,
		"_tstzrange":       3911,
		"daterange":        3912,
		"_daterange":       3913,
		"int8range":        3926,
		"_int8range":       3927,
		"record":           2249,
		"_record":          2287,
		"cstring":          2275,
		"any":              2276,
		"anyarray":         2277,
		"void":             2278,
		"trigger":          2279,
		"event_trigger":    3838,
		"language_handler": 2280,
		"internal":         2281,
		"opaque":           2282,
		"anyelement":       2283,
		"anynonarray":      2776,
		"anyenum":          3500,
		"fdw_handler":      3115,
		"index_am_handler": 325,
		"tsm_handler":      3310,
		"anyrange":         3831,
	}
)
