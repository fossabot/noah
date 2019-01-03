/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License,
 Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 software
 * distributed under the License is distributed on an "AS IS" BASIS,
	 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package system

var (
    postgresTables = []NTable{
        {
            TableName:   "pg_aggregate",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_am",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_amop",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_amproc",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_attrdef",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_attribute",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_authid",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_auth_members",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_cast",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_class",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_collation",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_constraint",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_conversion",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_database",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_db_role_setting",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_default_acl",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_depend",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_description",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_enum",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_event_trigger",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_extension",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_foreign_data_wrapper",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_foreign_server",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_foreign_table",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_index",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_inherits",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_init_privs",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_language",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_largeobject",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_largeobject_metadata",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_namespace",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_opclass",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_operator",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_opfamily",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_partitioned_table",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_pltemplate",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_policy",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_proc",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_publication",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_publication_rel",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_range",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_replication_origin",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_rewrite",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_seclabel",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_sequence",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_shdepend",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_shdescription",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_shseclabel",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_statistic",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_statistic_ext",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_subscription",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_subscription_rel",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_tablespace",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_transform",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_trigger",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_ts_config",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_ts_config_map",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_ts_dict",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_ts_parser",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_ts_template",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_type",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_user_mapping",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "SystemViews",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_available_extensions",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_available_extension_versions",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_config",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_cursors",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_file_settings",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_group",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_hba_file_rules",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_indexes",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_locks",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_matviews",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_policies",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_prepared_statements",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_prepared_xacts",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_publication_tables",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_replication_origin_status",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_replication_slots",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_roles",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_rules",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_seclabels",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_sequences",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_settings",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_shadow",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_stats",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_tables",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_timezone_abbrevs",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_timezone_names",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_user",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_user_mappings",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
        {
            TableName:   "pg_views",
            TableType:   NTableType_GLOBAL,
            IsCommitted: true,
        },
    }
)
