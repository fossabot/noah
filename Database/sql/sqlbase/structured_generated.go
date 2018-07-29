
package sqlbase

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import cockroach_util_hlc "github.com/Ready-Stock/Noah/Database/util/hlc"


import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ConstraintValidity int32

const (
	ConstraintValidity_Validated   ConstraintValidity = 0
	ConstraintValidity_Unvalidated ConstraintValidity = 1
)

var ConstraintValidity_name = map[int32]string{
	0: "Validated",
	1: "Unvalidated",
}
var ConstraintValidity_value = map[string]int32{
	"Validated":   0,
	"Unvalidated": 1,
}

func (x ConstraintValidity) Enum() *ConstraintValidity {
	p := new(ConstraintValidity)
	*p = x
	return p
}
func (x ConstraintValidity) String() string {
	return proto.EnumName(ConstraintValidity_name, int32(x))
}
func (x *ConstraintValidity) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ConstraintValidity_value, data, "ConstraintValidity")
	if err != nil {
		return err
	}
	*x = ConstraintValidity(value)
	return nil
}
func (ConstraintValidity) EnumDescriptor() ([]byte, []int) { return fileDescriptorStructured, []int{0} }

// These mirror the types supported by the sql/parser. See
// sql/parser/col_types.go.
type ColumnType_SemanticType int32

const (
	ColumnType_BOOL           ColumnType_SemanticType = 0
	ColumnType_INT            ColumnType_SemanticType = 1
	ColumnType_FLOAT          ColumnType_SemanticType = 2
	ColumnType_DECIMAL        ColumnType_SemanticType = 3
	ColumnType_DATE           ColumnType_SemanticType = 4
	ColumnType_TIMESTAMP      ColumnType_SemanticType = 5
	ColumnType_INTERVAL       ColumnType_SemanticType = 6
	ColumnType_STRING         ColumnType_SemanticType = 7
	ColumnType_BYTES          ColumnType_SemanticType = 8
	ColumnType_TIMESTAMPTZ    ColumnType_SemanticType = 9
	ColumnType_COLLATEDSTRING ColumnType_SemanticType = 10
	ColumnType_NAME           ColumnType_SemanticType = 11
	ColumnType_OID            ColumnType_SemanticType = 12
	// NULL is not supported as a table column type, however it can be
	// transferred through distsql streams.
	ColumnType_NULL       ColumnType_SemanticType = 13
	ColumnType_UUID       ColumnType_SemanticType = 14
	ColumnType_ARRAY      ColumnType_SemanticType = 15
	ColumnType_INET       ColumnType_SemanticType = 16
	ColumnType_TIME       ColumnType_SemanticType = 17
	ColumnType_JSON       ColumnType_SemanticType = 18
	ColumnType_TIMETZ     ColumnType_SemanticType = 19
	ColumnType_TUPLE      ColumnType_SemanticType = 20
	ColumnType_INT2VECTOR ColumnType_SemanticType = 200
	ColumnType_OIDVECTOR  ColumnType_SemanticType = 201
)

var ColumnType_SemanticType_name = map[int32]string{
	0:   "BOOL",
	1:   "INT",
	2:   "FLOAT",
	3:   "DECIMAL",
	4:   "DATE",
	5:   "TIMESTAMP",
	6:   "INTERVAL",
	7:   "STRING",
	8:   "BYTES",
	9:   "TIMESTAMPTZ",
	10:  "COLLATEDSTRING",
	11:  "NAME",
	12:  "OID",
	13:  "NULL",
	14:  "UUID",
	15:  "ARRAY",
	16:  "INET",
	17:  "TIME",
	18:  "JSON",
	19:  "TIMETZ",
	20:  "TUPLE",
	200: "INT2VECTOR",
	201: "OIDVECTOR",
}
var ColumnType_SemanticType_value = map[string]int32{
	"BOOL":           0,
	"INT":            1,
	"FLOAT":          2,
	"DECIMAL":        3,
	"DATE":           4,
	"TIMESTAMP":      5,
	"INTERVAL":       6,
	"STRING":         7,
	"BYTES":          8,
	"TIMESTAMPTZ":    9,
	"COLLATEDSTRING": 10,
	"NAME":           11,
	"OID":            12,
	"NULL":           13,
	"UUID":           14,
	"ARRAY":          15,
	"INET":           16,
	"TIME":           17,
	"JSON":           18,
	"TIMETZ":         19,
	"TUPLE":          20,
	"INT2VECTOR":     200,
	"OIDVECTOR":      201,
}

func (x ColumnType_SemanticType) Enum() *ColumnType_SemanticType {
	p := new(ColumnType_SemanticType)
	*p = x
	return p
}
func (x ColumnType_SemanticType) String() string {
	return proto.EnumName(ColumnType_SemanticType_name, int32(x))
}
func (x *ColumnType_SemanticType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ColumnType_SemanticType_value, data, "ColumnType_SemanticType")
	if err != nil {
		return err
	}
	*x = ColumnType_SemanticType(value)
	return nil
}
func (ColumnType_SemanticType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{0, 0}
}

type ColumnType_VisibleType int32

const (
	ColumnType_NONE             ColumnType_VisibleType = 0
	ColumnType_INTEGER          ColumnType_VisibleType = 1
	ColumnType_SMALLINT         ColumnType_VisibleType = 2
	ColumnType_BIGINT           ColumnType_VisibleType = 3
	ColumnType_BIT              ColumnType_VisibleType = 4
	ColumnType_REAL             ColumnType_VisibleType = 5
	ColumnType_DOUBLE_PRECISION ColumnType_VisibleType = 6
)

var ColumnType_VisibleType_name = map[int32]string{
	0: "NONE",
	1: "INTEGER",
	2: "SMALLINT",
	3: "BIGINT",
	4: "BIT",
	5: "REAL",
	6: "DOUBLE_PRECISION",
}
var ColumnType_VisibleType_value = map[string]int32{
	"NONE":             0,
	"INTEGER":          1,
	"SMALLINT":         2,
	"BIGINT":           3,
	"BIT":              4,
	"REAL":             5,
	"DOUBLE_PRECISION": 6,
}

func (x ColumnType_VisibleType) Enum() *ColumnType_VisibleType {
	p := new(ColumnType_VisibleType)
	*p = x
	return p
}
func (x ColumnType_VisibleType) String() string {
	return proto.EnumName(ColumnType_VisibleType_name, int32(x))
}
func (x *ColumnType_VisibleType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ColumnType_VisibleType_value, data, "ColumnType_VisibleType")
	if err != nil {
		return err
	}
	*x = ColumnType_VisibleType(value)
	return nil
}
func (ColumnType_VisibleType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{0, 1}
}

type ForeignKeyReference_Action int32

const (
	ForeignKeyReference_NO_ACTION   ForeignKeyReference_Action = 0
	ForeignKeyReference_RESTRICT    ForeignKeyReference_Action = 1
	ForeignKeyReference_SET_NULL    ForeignKeyReference_Action = 2
	ForeignKeyReference_SET_DEFAULT ForeignKeyReference_Action = 3
	ForeignKeyReference_CASCADE     ForeignKeyReference_Action = 4
)

var ForeignKeyReference_Action_name = map[int32]string{
	0: "NO_ACTION",
	1: "RESTRICT",
	2: "SET_NULL",
	3: "SET_DEFAULT",
	4: "CASCADE",
}
var ForeignKeyReference_Action_value = map[string]int32{
	"NO_ACTION":   0,
	"RESTRICT":    1,
	"SET_NULL":    2,
	"SET_DEFAULT": 3,
	"CASCADE":     4,
}

func (x ForeignKeyReference_Action) Enum() *ForeignKeyReference_Action {
	p := new(ForeignKeyReference_Action)
	*p = x
	return p
}
func (x ForeignKeyReference_Action) String() string {
	return proto.EnumName(ForeignKeyReference_Action_name, int32(x))
}
func (x *ForeignKeyReference_Action) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ForeignKeyReference_Action_value, data, "ForeignKeyReference_Action")
	if err != nil {
		return err
	}
	*x = ForeignKeyReference_Action(value)
	return nil
}
func (ForeignKeyReference_Action) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{1, 0}
}

// The direction of a column in the index.
type IndexDescriptor_Direction int32

const (
	IndexDescriptor_ASC  IndexDescriptor_Direction = 0
	IndexDescriptor_DESC IndexDescriptor_Direction = 1
)

var IndexDescriptor_Direction_name = map[int32]string{
	0: "ASC",
	1: "DESC",
}
var IndexDescriptor_Direction_value = map[string]int32{
	"ASC":  0,
	"DESC": 1,
}

func (x IndexDescriptor_Direction) Enum() *IndexDescriptor_Direction {
	p := new(IndexDescriptor_Direction)
	*p = x
	return p
}
func (x IndexDescriptor_Direction) String() string {
	return proto.EnumName(IndexDescriptor_Direction_name, int32(x))
}
func (x *IndexDescriptor_Direction) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(IndexDescriptor_Direction_value, data, "IndexDescriptor_Direction")
	if err != nil {
		return err
	}
	*x = IndexDescriptor_Direction(value)
	return nil
}
func (IndexDescriptor_Direction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{6, 0}
}

// The direction of a column in the index.
type IndexDescriptor_Type int32

const (
	IndexDescriptor_FORWARD  IndexDescriptor_Type = 0
	IndexDescriptor_INVERTED IndexDescriptor_Type = 1
)

var IndexDescriptor_Type_name = map[int32]string{
	0: "FORWARD",
	1: "INVERTED",
}
var IndexDescriptor_Type_value = map[string]int32{
	"FORWARD":  0,
	"INVERTED": 1,
}

func (x IndexDescriptor_Type) Enum() *IndexDescriptor_Type {
	p := new(IndexDescriptor_Type)
	*p = x
	return p
}
func (x IndexDescriptor_Type) String() string {
	return proto.EnumName(IndexDescriptor_Type_name, int32(x))
}
func (x *IndexDescriptor_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(IndexDescriptor_Type_value, data, "IndexDescriptor_Type")
	if err != nil {
		return err
	}
	*x = IndexDescriptor_Type(value)
	return nil
}
func (IndexDescriptor_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{6, 1}
}

// A descriptor within a mutation is unavailable for reads, writes
// and deletes. It is only available for implicit (internal to
// the database) writes and deletes depending on the state of the mutation.
type DescriptorMutation_State int32

const (
	// Not used.
	DescriptorMutation_UNKNOWN DescriptorMutation_State = 0
	// Operations can use this invisible descriptor to implicitly
	// delete entries.
	// Column: A descriptor in this state is invisible to
	// INSERT and UPDATE. DELETE must delete a column in this state.
	// Index: A descriptor in this state is invisible to an INSERT.
	// UPDATE must delete the old value of the index but doesn't write
	// the new value. DELETE must delete the index.
	//
	// When deleting a descriptor, all descriptor related data
	// (column or index data) can only be mass deleted once
	// all the nodes have transitioned to the DELETE_ONLY state.
	DescriptorMutation_DELETE_ONLY DescriptorMutation_State = 1
	// Operations can use this invisible descriptor to implicitly
	// write and delete entries.
	// Column: INSERT will populate this column with the default
	// value. UPDATE ignores this descriptor. DELETE must delete
	// the column.
	// Index: INSERT, UPDATE and DELETE treat this index like any
	// other index.
	//
	// When adding a descriptor, all descriptor related data
	// (column default or index data) can only be backfilled once
	// all nodes have transitioned into the DELETE_AND_WRITE_ONLY state.
	DescriptorMutation_DELETE_AND_WRITE_ONLY DescriptorMutation_State = 2
)

var DescriptorMutation_State_name = map[int32]string{
	0: "UNKNOWN",
	1: "DELETE_ONLY",
	2: "DELETE_AND_WRITE_ONLY",
}
var DescriptorMutation_State_value = map[string]int32{
	"UNKNOWN":               0,
	"DELETE_ONLY":           1,
	"DELETE_AND_WRITE_ONLY": 2,
}

func (x DescriptorMutation_State) Enum() *DescriptorMutation_State {
	p := new(DescriptorMutation_State)
	*p = x
	return p
}
func (x DescriptorMutation_State) String() string {
	return proto.EnumName(DescriptorMutation_State_name, int32(x))
}
func (x *DescriptorMutation_State) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(DescriptorMutation_State_value, data, "DescriptorMutation_State")
	if err != nil {
		return err
	}
	*x = DescriptorMutation_State(value)
	return nil
}
func (DescriptorMutation_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{7, 0}
}

// Direction of mutation.
type DescriptorMutation_Direction int32

const (
	// Not used.
	DescriptorMutation_NONE DescriptorMutation_Direction = 0
	// Descriptor is being added.
	DescriptorMutation_ADD DescriptorMutation_Direction = 1
	// Descriptor is being dropped.
	DescriptorMutation_DROP DescriptorMutation_Direction = 2
)

var DescriptorMutation_Direction_name = map[int32]string{
	0: "NONE",
	1: "ADD",
	2: "DROP",
}
var DescriptorMutation_Direction_value = map[string]int32{
	"NONE": 0,
	"ADD":  1,
	"DROP": 2,
}

func (x DescriptorMutation_Direction) Enum() *DescriptorMutation_Direction {
	p := new(DescriptorMutation_Direction)
	*p = x
	return p
}
func (x DescriptorMutation_Direction) String() string {
	return proto.EnumName(DescriptorMutation_Direction_name, int32(x))
}
func (x *DescriptorMutation_Direction) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(DescriptorMutation_Direction_value, data, "DescriptorMutation_Direction")
	if err != nil {
		return err
	}
	*x = DescriptorMutation_Direction(value)
	return nil
}
func (DescriptorMutation_Direction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{7, 1}
}

// State is set if this TableDescriptor is in the process of being added or deleted.
// A non-public table descriptor cannot be leased.
// A schema changer observing DROP set will truncate the table and delete the
// descriptor.
// It is illegal to transition DROP to any other state.
type TableDescriptor_State int32

const (
	// Not used.
	TableDescriptor_PUBLIC TableDescriptor_State = 0
	// Descriptor is being added.
	TableDescriptor_ADD TableDescriptor_State = 1
	// Descriptor is being dropped.
	TableDescriptor_DROP TableDescriptor_State = 2
)

var TableDescriptor_State_name = map[int32]string{
	0: "PUBLIC",
	1: "ADD",
	2: "DROP",
}
var TableDescriptor_State_value = map[string]int32{
	"PUBLIC": 0,
	"ADD":    1,
	"DROP":   2,
}

func (x TableDescriptor_State) Enum() *TableDescriptor_State {
	p := new(TableDescriptor_State)
	*p = x
	return p
}
func (x TableDescriptor_State) String() string {
	return proto.EnumName(TableDescriptor_State_name, int32(x))
}
func (x *TableDescriptor_State) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(TableDescriptor_State_value, data, "TableDescriptor_State")
	if err != nil {
		return err
	}
	*x = TableDescriptor_State(value)
	return nil
}
func (TableDescriptor_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{8, 0}
}

// AuditMode indicates which auditing actions to take when this table is used.
type TableDescriptor_AuditMode int32

const (
	TableDescriptor_DISABLED  TableDescriptor_AuditMode = 0
	TableDescriptor_READWRITE TableDescriptor_AuditMode = 1
)

var TableDescriptor_AuditMode_name = map[int32]string{
	0: "DISABLED",
	1: "READWRITE",
}
var TableDescriptor_AuditMode_value = map[string]int32{
	"DISABLED":  0,
	"READWRITE": 1,
}

func (x TableDescriptor_AuditMode) Enum() *TableDescriptor_AuditMode {
	p := new(TableDescriptor_AuditMode)
	*p = x
	return p
}
func (x TableDescriptor_AuditMode) String() string {
	return proto.EnumName(TableDescriptor_AuditMode_name, int32(x))
}
func (x *TableDescriptor_AuditMode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(TableDescriptor_AuditMode_value, data, "TableDescriptor_AuditMode")
	if err != nil {
		return err
	}
	*x = TableDescriptor_AuditMode(value)
	return nil
}
func (TableDescriptor_AuditMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{8, 1}
}

type ColumnType struct {
	SemanticType ColumnType_SemanticType `protobuf:"varint,1,opt,name=semantic_type,json=semanticType,enum=cockroach.sql.sqlbase.ColumnType_SemanticType" json:"semantic_type"`
	// BIT, INT, FLOAT, DECIMAL, CHAR and BINARY
	Width int32 `protobuf:"varint,2,opt,name=width" json:"width"`
	// FLOAT and DECIMAL.
	Precision int32 `protobuf:"varint,3,opt,name=precision" json:"precision"`
	// The length of each dimension in the array. A dimension of -1 means that
	// no bound was specified for that dimension.
	ArrayDimensions []int32 `protobuf:"varint,4,rep,name=array_dimensions,json=arrayDimensions" json:"array_dimensions,omitempty"`
	// Collated STRING, CHAR, and VARCHAR
	Locale *string `protobuf:"bytes,5,opt,name=locale" json:"locale,omitempty"`
	// Alias for any types where our internal representation is different than
	// the user specification. Examples are INT4, FLOAT4, etc. Mostly for Postgres
	// compatibility.
	VisibleType ColumnType_VisibleType `protobuf:"varint,6,opt,name=visible_type,json=visibleType,enum=cockroach.sql.sqlbase.ColumnType_VisibleType" json:"visible_type"`
	// Only used if the kind is ARRAY.
	ArrayContents *ColumnType_SemanticType `protobuf:"varint,7,opt,name=array_contents,json=arrayContents,enum=cockroach.sql.sqlbase.ColumnType_SemanticType" json:"array_contents,omitempty"`
	// Only used if the kind is TUPLE
	TupleContents []ColumnType `protobuf:"bytes,8,rep,name=tuple_contents,json=tupleContents" json:"tuple_contents"`
}

func (m *ColumnType) Reset()                    { *m = ColumnType{} }
func (m *ColumnType) String() string            { return proto.CompactTextString(m) }
func (*ColumnType) ProtoMessage()               {}
func (*ColumnType) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{0} }

type ForeignKeyReference struct {
	Table    ID                 `protobuf:"varint,1,opt,name=table,casttype=ID" json:"table"`
	Index    IndexID            `protobuf:"varint,2,opt,name=index,casttype=IndexID" json:"index"`
	Name     string             `protobuf:"bytes,3,opt,name=name" json:"name"`
	Validity ConstraintValidity `protobuf:"varint,4,opt,name=validity,enum=cockroach.sql.sqlbase.ConstraintValidity" json:"validity"`
	// If this FK only uses a prefix of the columns in its index, we record how
	// many to avoid spuriously counting the additional cols as used by this FK.
	SharedPrefixLen int32                      `protobuf:"varint,5,opt,name=shared_prefix_len,json=sharedPrefixLen" json:"shared_prefix_len"`
	OnDelete        ForeignKeyReference_Action `protobuf:"varint,6,opt,name=on_delete,json=onDelete,enum=cockroach.sql.sqlbase.ForeignKeyReference_Action" json:"on_delete"`
	OnUpdate        ForeignKeyReference_Action `protobuf:"varint,7,opt,name=on_update,json=onUpdate,enum=cockroach.sql.sqlbase.ForeignKeyReference_Action" json:"on_update"`
}

func (m *ForeignKeyReference) Reset()                    { *m = ForeignKeyReference{} }
func (m *ForeignKeyReference) String() string            { return proto.CompactTextString(m) }
func (*ForeignKeyReference) ProtoMessage()               {}
func (*ForeignKeyReference) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{1} }

type ColumnDescriptor struct {
	Name     string     `protobuf:"bytes,1,opt,name=name" json:"name"`
	ID       ColumnID   `protobuf:"varint,2,opt,name=id,casttype=ColumnID" json:"id"`
	Type     ColumnType `protobuf:"bytes,3,opt,name=type" json:"type"`
	Nullable bool       `protobuf:"varint,4,opt,name=nullable" json:"nullable"`
	// Default expression to use to populate the column on insert if no
	// value is provided.
	DefaultExpr *string `protobuf:"bytes,5,opt,name=default_expr,json=defaultExpr" json:"default_expr,omitempty"`
	Hidden      bool    `protobuf:"varint,6,opt,name=hidden" json:"hidden"`
	// Ids of sequences used in this column's DEFAULT expression, in calls to nextval().
	UsesSequenceIds []ID `protobuf:"varint,10,rep,name=uses_sequence_ids,json=usesSequenceIds,casttype=ID" json:"uses_sequence_ids,omitempty"`
	// Expression to use to compute the value of this column if this is a
	// computed column.
	ComputeExpr *string `protobuf:"bytes,11,opt,name=compute_expr,json=computeExpr" json:"compute_expr,omitempty"`
}

func (m *ColumnDescriptor) Reset()                    { *m = ColumnDescriptor{} }
func (m *ColumnDescriptor) String() string            { return proto.CompactTextString(m) }
func (*ColumnDescriptor) ProtoMessage()               {}
func (*ColumnDescriptor) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{2} }

// ColumnFamilyDescriptor is set of columns stored together in one kv entry.
type ColumnFamilyDescriptor struct {
	Name string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	ID   FamilyID `protobuf:"varint,2,opt,name=id,casttype=FamilyID" json:"id"`
	// A list of column names of which the family is comprised. This list
	// parallels the column_ids list. If duplicating the storage of the column
	// names here proves to be prohibitive, we could clear this field before
	// saving and reconstruct it after loading.
	ColumnNames []string `protobuf:"bytes,3,rep,name=column_names,json=columnNames" json:"column_names,omitempty"`
	// A list of column ids of which the family is comprised. This list parallels
	// the column_names list.
	ColumnIDs []ColumnID `protobuf:"varint,4,rep,name=column_ids,json=columnIds,casttype=ColumnID" json:"column_ids,omitempty"`
	// If nonzero, the column involved in the single column optimization.
	//
	// Families store columns in a ValueType_TUPLE as repeated <columnID><data>
	// entries. As a space optimization and for backward compatibility, a single
	// column is written without the column id prefix. Because more columns could
	// be added, it would be ambiguous which column was stored when read back in,
	// so this field supplies it.
	DefaultColumnID ColumnID `protobuf:"varint,5,opt,name=default_column_id,json=defaultColumnId,casttype=ColumnID" json:"default_column_id"`
}

func (m *ColumnFamilyDescriptor) Reset()                    { *m = ColumnFamilyDescriptor{} }
func (m *ColumnFamilyDescriptor) String() string            { return proto.CompactTextString(m) }
func (*ColumnFamilyDescriptor) ProtoMessage()               {}
func (*ColumnFamilyDescriptor) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{3} }

// InterleaveDescriptor represents an index (either primary or secondary) that
// is interleaved into another table's data.
//
// Example:
// Table 1 -> /a/b
// Table 2 -> /a/b/c
// Table 3 -> /a/b/c/d
//
// There are two components (table 2 is the parent and table 1 is the
// grandparent) with shared lengths 2 and 1.
type InterleaveDescriptor struct {
	// Ancestors contains the nesting of interleaves in the order they appear in
	// an encoded key. This means they are always in the far-to-near ancestor
	// order (e.g. grand-grand-parent, grand-parent, parent).
	Ancestors []InterleaveDescriptor_Ancestor `protobuf:"bytes,1,rep,name=ancestors" json:"ancestors"`
}

func (m *InterleaveDescriptor) Reset()                    { *m = InterleaveDescriptor{} }
func (m *InterleaveDescriptor) String() string            { return proto.CompactTextString(m) }
func (*InterleaveDescriptor) ProtoMessage()               {}
func (*InterleaveDescriptor) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{4} }

type InterleaveDescriptor_Ancestor struct {
	// TableID the ID of the table being interleaved into.
	TableID ID `protobuf:"varint,1,opt,name=table_id,json=tableId,casttype=ID" json:"table_id"`
	// IndexID is the ID of the parent index being interleaved into.
	IndexID IndexID `protobuf:"varint,2,opt,name=index_id,json=indexId,casttype=IndexID" json:"index_id"`
	// SharedPrefixLen is how many fields are shared between a parent and child
	// being interleaved, excluding any fields shared between parent and
	// grandparent. Thus, the sum of SharedPrefixLens in the components of an
	// InterleaveDescriptor is never more than the number of fields in the index
	// being interleaved.
	SharedPrefixLen uint32 `protobuf:"varint,3,opt,name=shared_prefix_len,json=sharedPrefixLen" json:"shared_prefix_len"`
}

func (m *InterleaveDescriptor_Ancestor) Reset()         { *m = InterleaveDescriptor_Ancestor{} }
func (m *InterleaveDescriptor_Ancestor) String() string { return proto.CompactTextString(m) }
func (*InterleaveDescriptor_Ancestor) ProtoMessage()    {}
func (*InterleaveDescriptor_Ancestor) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{4, 0}
}

// PartitioningDescriptor represents the partitioning of an index into spans
// of keys addressable by a zone config. The key encoding is unchanged. Each
// partition may optionally be itself divided into further partitions, called
// subpartitions.
type PartitioningDescriptor struct {
	// NumColumns is how large of a prefix of the columns in an index are used in
	// the function mapping column values to partitions. If this is a
	// subpartition, this is offset to start from the end of the parent
	// partition's columns. If NumColumns is 0, then there is no partitioning.
	NumColumns uint32 `protobuf:"varint,1,opt,name=num_columns,json=numColumns" json:"num_columns"`
	// Exactly one of List or Range is required to be non-empty if NumColumns is
	// non-zero.
	List  []PartitioningDescriptor_List  `protobuf:"bytes,2,rep,name=list" json:"list"`
	Range []PartitioningDescriptor_Range `protobuf:"bytes,3,rep,name=range" json:"range"`
}

func (m *PartitioningDescriptor) Reset()                    { *m = PartitioningDescriptor{} }
func (m *PartitioningDescriptor) String() string            { return proto.CompactTextString(m) }
func (*PartitioningDescriptor) ProtoMessage()               {}
func (*PartitioningDescriptor) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{5} }

// List represents a list partitioning, which maps individual tuples to
// partitions.
type PartitioningDescriptor_List struct {
	// Name is the partition name.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name"`
	// Values is an unordered set of the tuples included in this partition. Each
	// tuple is encoded with the EncDatum value encoding. DEFAULT is encoded as
	// NOT NULL followed by PartitionDefaultVal encoded as a non-sorting
	// uvarint.
	Values [][]byte `protobuf:"bytes,2,rep,name=values" json:"values,omitempty"`
	// Subpartitioning represents a further partitioning of this list partition.
	Subpartitioning PartitioningDescriptor `protobuf:"bytes,3,opt,name=subpartitioning" json:"subpartitioning"`
}

func (m *PartitioningDescriptor_List) Reset()         { *m = PartitioningDescriptor_List{} }
func (m *PartitioningDescriptor_List) String() string { return proto.CompactTextString(m) }
func (*PartitioningDescriptor_List) ProtoMessage()    {}
func (*PartitioningDescriptor_List) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{5, 0}
}

// Range represents a range partitioning, which maps ranges of tuples to
// partitions by specifying exclusive upper bounds. The range partitions in a
// PartitioningDescriptor are required be sorted by UpperBound.
type PartitioningDescriptor_Range struct {
	// Name is the partition name.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name"`
	// FromInclusive is the inclusive lower bound of this range partition. It is
	// encoded with the EncDatum value encoding. MINVALUE and MAXVALUE are
	// encoded as NOT NULL followed by a PartitionSpecialValCode encoded as a
	// non-sorting uvarint.
	FromInclusive []byte `protobuf:"bytes,3,opt,name=from_inclusive,json=fromInclusive" json:"from_inclusive,omitempty"`
	// ToExclusive is the exclusive upper bound of this range partition. It is
	// encoded in the same way as From.
	ToExclusive []byte `protobuf:"bytes,2,opt,name=to_exclusive,json=toExclusive" json:"to_exclusive,omitempty"`
}

func (m *PartitioningDescriptor_Range) Reset()         { *m = PartitioningDescriptor_Range{} }
func (m *PartitioningDescriptor_Range) String() string { return proto.CompactTextString(m) }
func (*PartitioningDescriptor_Range) ProtoMessage()    {}
func (*PartitioningDescriptor_Range) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{5, 1}
}

// IndexDescriptor describes an index (primary or secondary).
//
// Sample field values on the following table:
//
//   CREATE TABLE t (
//     k1 INT NOT NULL,   // column ID: 1
//     k2 INT NOT NULL,   // column ID: 2
//     u INT NULL,        // column ID: 3
//     v INT NULL,        // column ID: 4
//     w INT NULL,        // column ID: 5
//     CONSTRAINT "primary" PRIMARY KEY (k1, k2),
//     INDEX k1v (k1, v) STORING (w),
//     FAMILY "primary" (k1, k2, u, v, w)
//   )
//
// Primary index:
//   name:                primary
//   id:                  1
//   unique:              true
//   column_names:        k1, k2
//   column_directions:   ASC, ASC
//   column_ids:          1, 2   // k1, k2
//
// [old STORING encoding] Index k1v (k1, v) STORING (w):
//   name:                k1v
//   id:                  2
//   unique:              false
//   column_names:        k1, v
//   column_directions:   ASC, ASC
//   store_column_names:  w
//   column_ids:          1, 4   // k1, v
//   extra_column_ids:    2, 5   // k2, w
//
// [new STORING encoding] Index k1v (k1, v) STORING (w):
//   name:                k1v
//   id:                  2
//   unique:              false
//   column_names:        k1, v
//   column_directions:   ASC, ASC
//   store_column_names:  w
//   column_ids:          1, 4   // k1, v
//   extra_column_ids:    2      // k2
//   store_column_ids:    5      // w
type IndexDescriptor struct {
	Name   string  `protobuf:"bytes,1,opt,name=name" json:"name"`
	ID     IndexID `protobuf:"varint,2,opt,name=id,casttype=IndexID" json:"id"`
	Unique bool    `protobuf:"varint,3,opt,name=unique" json:"unique"`
	// An ordered list of column names of which the index is comprised; these
	// columns do not include any additional stored columns (which are in
	// stored_column_names). This list parallels the column_ids list.
	//
	// Note: if duplicating the storage of the column names here proves to be
	// prohibitive, we could clear this field before saving and reconstruct it
	// after loading.
	ColumnNames []string `protobuf:"bytes,4,rep,name=column_names,json=columnNames" json:"column_names,omitempty"`
	// The sort direction of each column in column_names.
	ColumnDirections []IndexDescriptor_Direction `protobuf:"varint,8,rep,name=column_directions,json=columnDirections,enum=cockroach.sql.sqlbase.IndexDescriptor_Direction" json:"column_directions,omitempty"`
	// An ordered list of column names which the index stores in addition to the
	// columns which are explicitly part of the index (STORING clause). Only used
	// for secondary indexes.
	StoreColumnNames []string `protobuf:"bytes,5,rep,name=store_column_names,json=storeColumnNames" json:"store_column_names,omitempty"`
	// An ordered list of column IDs of which the index is comprised. This list
	// parallels the column_names list and does not include any additional stored
	// columns.
	ColumnIDs []ColumnID `protobuf:"varint,6,rep,name=column_ids,json=columnIds,casttype=ColumnID" json:"column_ids,omitempty"`
	// An ordered list of IDs for the additional columns associated with the
	// index:
	//  - implicit columns, which are all the primary key columns that are not
	//    already part of the index (i.e. PrimaryIndex.column_ids - column_ids).
	//  - stored columns (the columns in store_column_names) if this index uses the
	//    old STORING encoding (key-encoded data).
	//
	// Only used for secondary indexes.
	// For non-unique indexes, these columns are appended to the key.
	// For unique indexes, these columns are stored in the value (unless the key
	// contains a NULL value: then the extra columns are appended to the key to
	// unique-ify it).
	// This distinction exists because we want to be able to insert an entry using
	// a single conditional put on the key.
	ExtraColumnIDs []ColumnID `protobuf:"varint,7,rep,name=extra_column_ids,json=extraColumnIds,casttype=ColumnID" json:"extra_column_ids,omitempty"`
	// An ordered list of column IDs that parallels store_column_names if this
	// index uses the new STORING encoding (value-encoded data, always in the KV
	// value).
	StoreColumnIDs []ColumnID `protobuf:"varint,14,rep,name=store_column_ids,json=storeColumnIds,casttype=ColumnID" json:"store_column_ids,omitempty"`
	// CompositeColumnIDs contains an ordered list of IDs of columns that appear
	// in the index and have a composite encoding. Includes IDs from both
	// column_ids and extra_column_ids.
	CompositeColumnIDs []ColumnID            `protobuf:"varint,13,rep,name=composite_column_ids,json=compositeColumnIds,casttype=ColumnID" json:"composite_column_ids,omitempty"`
	ForeignKey         ForeignKeyReference   `protobuf:"bytes,9,opt,name=foreign_key,json=foreignKey" json:"foreign_key"`
	ReferencedBy       []ForeignKeyReference `protobuf:"bytes,10,rep,name=referenced_by,json=referencedBy" json:"referenced_by"`
	// Interleave, if it's not the zero value, describes how this index's data is
	// interleaved into another index's data.
	Interleave InterleaveDescriptor `protobuf:"bytes,11,opt,name=interleave" json:"interleave"`
	// InterleavedBy contains a reference to every table/index that is interleaved
	// into this one.
	InterleavedBy []ForeignKeyReference `protobuf:"bytes,12,rep,name=interleaved_by,json=interleavedBy" json:"interleaved_by"`
	// Partitioning, if it's not the zero value, describes how this index's data
	// is partitioned into spans of keys each addressable by zone configs.
	Partitioning PartitioningDescriptor `protobuf:"bytes,15,opt,name=partitioning" json:"partitioning"`
	// Type is the type of index, inverted or forward.
	Type IndexDescriptor_Type `protobuf:"varint,16,opt,name=type,enum=cockroach.sql.sqlbase.IndexDescriptor_Type" json:"type"`
}

func (m *IndexDescriptor) Reset()                    { *m = IndexDescriptor{} }
func (m *IndexDescriptor) String() string            { return proto.CompactTextString(m) }
func (*IndexDescriptor) ProtoMessage()               {}
func (*IndexDescriptor) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{6} }

// A DescriptorMutation represents a column or an index that
// has either been added or dropped and hasn't yet transitioned
// into a stable state: completely backfilled and visible, or
// completely deleted. A table descriptor in the middle of a
// schema change will have a DescriptorMutation FIFO queue
// containing each column/index descriptor being added or dropped.
type DescriptorMutation struct {
	// Types that are valid to be assigned to Descriptor_:
	//	*DescriptorMutation_Column
	//	*DescriptorMutation_Index
	Descriptor_ isDescriptorMutation_Descriptor_ `protobuf_oneof:"descriptor"`
	State       DescriptorMutation_State         `protobuf:"varint,3,opt,name=state,enum=cockroach.sql.sqlbase.DescriptorMutation_State" json:"state"`
	Direction   DescriptorMutation_Direction     `protobuf:"varint,4,opt,name=direction,enum=cockroach.sql.sqlbase.DescriptorMutation_Direction" json:"direction"`
	// The mutation id used to group mutations that should be applied together.
	// This is used for situations like creating a unique column, which
	// involve adding two mutations: one for the column, and another for the
	// unique constraint index.
	MutationID MutationID `protobuf:"varint,5,opt,name=mutation_id,json=mutationId,casttype=MutationID" json:"mutation_id"`
	// Indicates that this mutation is a rollback.
	Rollback bool `protobuf:"varint,7,opt,name=rollback" json:"rollback"`
}

func (m *DescriptorMutation) Reset()                    { *m = DescriptorMutation{} }
func (m *DescriptorMutation) String() string            { return proto.CompactTextString(m) }
func (*DescriptorMutation) ProtoMessage()               {}
func (*DescriptorMutation) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{7} }

type isDescriptorMutation_Descriptor_ interface {
	isDescriptorMutation_Descriptor_()
	MarshalTo([]byte) (int, error)
	Size() int
}

type DescriptorMutation_Column struct {
	Column *ColumnDescriptor `protobuf:"bytes,1,opt,name=column,oneof"`
}
type DescriptorMutation_Index struct {
	Index *IndexDescriptor `protobuf:"bytes,2,opt,name=index,oneof"`
}

func (*DescriptorMutation_Column) isDescriptorMutation_Descriptor_() {}
func (*DescriptorMutation_Index) isDescriptorMutation_Descriptor_()  {}

func (m *DescriptorMutation) GetDescriptor_() isDescriptorMutation_Descriptor_ {
	if m != nil {
		return m.Descriptor_
	}
	return nil
}

func (m *DescriptorMutation) GetColumn() *ColumnDescriptor {
	if x, ok := m.GetDescriptor_().(*DescriptorMutation_Column); ok {
		return x.Column
	}
	return nil
}

func (m *DescriptorMutation) GetIndex() *IndexDescriptor {
	if x, ok := m.GetDescriptor_().(*DescriptorMutation_Index); ok {
		return x.Index
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*DescriptorMutation) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _DescriptorMutation_OneofMarshaler, _DescriptorMutation_OneofUnmarshaler, _DescriptorMutation_OneofSizer, []interface{}{
		(*DescriptorMutation_Column)(nil),
		(*DescriptorMutation_Index)(nil),
	}
}

func _DescriptorMutation_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*DescriptorMutation)
	// descriptor
	switch x := m.Descriptor_.(type) {
	case *DescriptorMutation_Column:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Column); err != nil {
			return err
		}
	case *DescriptorMutation_Index:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Index); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("DescriptorMutation.Descriptor_ has unexpected type %T", x)
	}
	return nil
}

func _DescriptorMutation_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*DescriptorMutation)
	switch tag {
	case 1: // descriptor.column
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ColumnDescriptor)
		err := b.DecodeMessage(msg)
		m.Descriptor_ = &DescriptorMutation_Column{msg}
		return true, err
	case 2: // descriptor.index
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(IndexDescriptor)
		err := b.DecodeMessage(msg)
		m.Descriptor_ = &DescriptorMutation_Index{msg}
		return true, err
	default:
		return false, nil
	}
}

func _DescriptorMutation_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*DescriptorMutation)
	// descriptor
	switch x := m.Descriptor_.(type) {
	case *DescriptorMutation_Column:
		s := proto.Size(x.Column)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *DescriptorMutation_Index:
		s := proto.Size(x.Index)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// A TableDescriptor represents a table or view and is stored in a
// structured metadata key. The TableDescriptor has a globally-unique ID,
// while its member {Column,Index}Descriptors have locally-unique IDs.
type TableDescriptor struct {
	// The table name. It should be normalized using NormalizeName() before
	// comparing it.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name"`
	ID   ID     `protobuf:"varint,3,opt,name=id,casttype=ID" json:"id"`
	// ID of the parent database.
	ParentID ID `protobuf:"varint,4,opt,name=parent_id,json=parentId,casttype=ID" json:"parent_id"`
	// Monotonically increasing version of the table descriptor.
	//
	// The design maintains two invariant:
	// 1. Two safe versions: A transaction at a particular timestamp is
	//    allowed to use one of two versions of a table descriptor:
	//    the one that would be read from the store at that timestamp,
	//    and the one behind it in version.
	// 2. Two leased version: There can be valid leases on at most the 2
	//    latest versions of a table in the cluster at any time. New leases
	//    are only granted on the latest version.
	//
	// The database must maintain correctness in light of there being two
	// versions of a descriptor that can be used.
	//
	// Multiple schema change mutations can be grouped together on a
	// particular version increment.
	Version DescriptorVersion `protobuf:"varint,5,opt,name=version,casttype=DescriptorVersion" json:"version"`
	// Do not use. On the path to deprecation.
	UpVersion bool `protobuf:"varint,6,opt,name=up_version,json=upVersion" json:"up_version"`
	// Last modification time of the table descriptor.
	ModificationTime cockroach_util_hlc.Timestamp `protobuf:"bytes,7,opt,name=modification_time,json=modificationTime" json:"modification_time"`
	Columns          []ColumnDescriptor           `protobuf:"bytes,8,rep,name=columns" json:"columns"`
	// next_column_id is used to ensure that deleted column ids are not reused.
	NextColumnID ColumnID                 `protobuf:"varint,9,opt,name=next_column_id,json=nextColumnId,casttype=ColumnID" json:"next_column_id"`
	Families     []ColumnFamilyDescriptor `protobuf:"bytes,22,rep,name=families" json:"families"`
	// next_family_id is used to ensure that deleted family ids are not reused.
	NextFamilyID FamilyID        `protobuf:"varint,23,opt,name=next_family_id,json=nextFamilyId,casttype=FamilyID" json:"next_family_id"`
	PrimaryIndex IndexDescriptor `protobuf:"bytes,10,opt,name=primary_index,json=primaryIndex" json:"primary_index"`
	// indexes are all the secondary indexes.
	Indexes []IndexDescriptor `protobuf:"bytes,11,rep,name=indexes" json:"indexes"`
	// next_index_id is used to ensure that deleted index ids are not reused.
	NextIndexID IndexID              `protobuf:"varint,12,opt,name=next_index_id,json=nextIndexId,casttype=IndexID" json:"next_index_id"`
	// Privileges  *PrivilegeDescriptor `protobuf:"bytes,13,opt,name=privileges" json:"privileges,omitempty"`
	// Columns or indexes being added or deleted in a FIFO order.
	Mutations []DescriptorMutation               `protobuf:"bytes,14,rep,name=mutations" json:"mutations"`
	Lease     *TableDescriptor_SchemaChangeLease `protobuf:"bytes,15,opt,name=lease" json:"lease,omitempty"`
	// An id for the next group of mutations to be applied together.
	NextMutationID MutationID `protobuf:"varint,16,opt,name=next_mutation_id,json=nextMutationId,casttype=MutationID" json:"next_mutation_id"`
	// format_version declares which sql to key:value mapping is being used to
	// represent the data in this table.
	FormatVersion FormatVersion                      `protobuf:"varint,17,opt,name=format_version,json=formatVersion,casttype=FormatVersion" json:"format_version"`
	State         TableDescriptor_State              `protobuf:"varint,19,opt,name=state,enum=cockroach.sql.sqlbase.TableDescriptor_State" json:"state"`
	Checks        []*TableDescriptor_CheckConstraint `protobuf:"bytes,20,rep,name=checks" json:"checks,omitempty"`
	// A list of draining names. The draining name entries are drained from
	// the cluster wide name caches by incrementing the version for this
	// descriptor and ensuring that there are no leases on prior
	// versions of the descriptor. This field is then cleared and the version
	// of the descriptor incremented.
	DrainingNames []TableDescriptor_NameInfo `protobuf:"bytes,21,rep,name=draining_names,json=drainingNames" json:"draining_names"`
	// The TableDescriptor is used for views in addition to tables. Views
	// use mostly the same fields as tables, but need to track the actual
	// query from the view definition as well.
	//
	// For now we only track a string representation of the query. This prevents
	// us from easily supporting things like renames of the dependencies of a
	// view. Eventually we'll want to switch to a semantic encoding of the query
	// that relies on IDs rather than names so that we can support renames of
	// fields relied on by the query, as Postgres does.
	//
	// Note: The presence of this field is used to determine whether or not
	// a TableDescriptor represents a view.
	ViewQuery string `protobuf:"bytes,24,opt,name=view_query,json=viewQuery" json:"view_query"`
	// The IDs of all relations that this depends on.
	// Only ever populated if this descriptor is for a view.
	DependsOn []ID `protobuf:"varint,25,rep,name=dependsOn,casttype=ID" json:"dependsOn,omitempty"`
	// All references to this table/view from other views in the system, tracked
	// down to the column/index so that we can restrict changes to them while
	// they're still being referred to.
	DependedOnBy []TableDescriptor_Reference `protobuf:"bytes,26,rep,name=dependedOnBy" json:"dependedOnBy"`
	// Mutation jobs queued for execution in a FIFO order. Remains synchronized
	// with the mutations list.
	MutationJobs []TableDescriptor_MutationJob `protobuf:"bytes,27,rep,name=mutationJobs" json:"mutationJobs"`
	// The presence of sequence_opts indicates that this descriptor is for a sequence.
	SequenceOpts *TableDescriptor_SequenceOpts `protobuf:"bytes,28,opt,name=sequence_opts,json=sequenceOpts" json:"sequence_opts,omitempty"`
	// The drop time is set when a table is truncated or dropped,
	// based on the current time in nanoseconds since the epoch.
	// Use this timestamp + GC TTL to start deleting the table's
	// contents.
	//
	// TODO(vivek): Replace with the ModificationTime. This has been
	// added only for migration purposes.
	DropTime int64 `protobuf:"varint,29,opt,name=drop_time,json=dropTime" json:"drop_time"`
	// ReplacementOf tracks prior IDs by which this table went -- e.g. when
	// TRUNCATE creates a replacement of a table and swaps it in for the the old
	// one, it should note on the new table the ID of the table it replaced. This
	// can be used when trying to track a table's history across truncatations.
	ReplacementOf TableDescriptor_Replacement `protobuf:"bytes,30,opt,name=replacement_of,json=replacementOf" json:"replacement_of"`
	AuditMode     TableDescriptor_AuditMode   `protobuf:"varint,31,opt,name=audit_mode,json=auditMode,enum=cockroach.sql.sqlbase.TableDescriptor_AuditMode" json:"audit_mode"`
}

func (m *TableDescriptor) Reset()                    { *m = TableDescriptor{} }
func (m *TableDescriptor) String() string            { return proto.CompactTextString(m) }
func (*TableDescriptor) ProtoMessage()               {}
func (*TableDescriptor) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{8} }

func (m *TableDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TableDescriptor) GetID() ID {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *TableDescriptor) GetParentID() ID {
	if m != nil {
		return m.ParentID
	}
	return 0
}

func (m *TableDescriptor) GetVersion() DescriptorVersion {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *TableDescriptor) GetUpVersion() bool {
	if m != nil {
		return m.UpVersion
	}
	return false
}

func (m *TableDescriptor) GetModificationTime() cockroach_util_hlc.Timestamp {
	if m != nil {
		return m.ModificationTime
	}
	return cockroach_util_hlc.Timestamp{}
}

func (m *TableDescriptor) GetColumns() []ColumnDescriptor {
	if m != nil {
		return m.Columns
	}
	return nil
}

func (m *TableDescriptor) GetNextColumnID() ColumnID {
	if m != nil {
		return m.NextColumnID
	}
	return 0
}

func (m *TableDescriptor) GetFamilies() []ColumnFamilyDescriptor {
	if m != nil {
		return m.Families
	}
	return nil
}

func (m *TableDescriptor) GetNextFamilyID() FamilyID {
	if m != nil {
		return m.NextFamilyID
	}
	return 0
}

func (m *TableDescriptor) GetPrimaryIndex() IndexDescriptor {
	if m != nil {
		return m.PrimaryIndex
	}
	return IndexDescriptor{}
}

func (m *TableDescriptor) GetIndexes() []IndexDescriptor {
	if m != nil {
		return m.Indexes
	}
	return nil
}

func (m *TableDescriptor) GetNextIndexID() IndexID {
	if m != nil {
		return m.NextIndexID
	}
	return 0
}



func (m *TableDescriptor) GetMutations() []DescriptorMutation {
	if m != nil {
		return m.Mutations
	}
	return nil
}

func (m *TableDescriptor) GetLease() *TableDescriptor_SchemaChangeLease {
	if m != nil {
		return m.Lease
	}
	return nil
}

func (m *TableDescriptor) GetNextMutationID() MutationID {
	if m != nil {
		return m.NextMutationID
	}
	return 0
}

func (m *TableDescriptor) GetFormatVersion() FormatVersion {
	if m != nil {
		return m.FormatVersion
	}
	return 0
}

func (m *TableDescriptor) GetState() TableDescriptor_State {
	if m != nil {
		return m.State
	}
	return TableDescriptor_PUBLIC
}

func (m *TableDescriptor) GetChecks() []*TableDescriptor_CheckConstraint {
	if m != nil {
		return m.Checks
	}
	return nil
}

func (m *TableDescriptor) GetDrainingNames() []TableDescriptor_NameInfo {
	if m != nil {
		return m.DrainingNames
	}
	return nil
}

func (m *TableDescriptor) GetViewQuery() string {
	if m != nil {
		return m.ViewQuery
	}
	return ""
}

func (m *TableDescriptor) GetDependsOn() []ID {
	if m != nil {
		return m.DependsOn
	}
	return nil
}

func (m *TableDescriptor) GetDependedOnBy() []TableDescriptor_Reference {
	if m != nil {
		return m.DependedOnBy
	}
	return nil
}

func (m *TableDescriptor) GetMutationJobs() []TableDescriptor_MutationJob {
	if m != nil {
		return m.MutationJobs
	}
	return nil
}

func (m *TableDescriptor) GetSequenceOpts() *TableDescriptor_SequenceOpts {
	if m != nil {
		return m.SequenceOpts
	}
	return nil
}

func (m *TableDescriptor) GetDropTime() int64 {
	if m != nil {
		return m.DropTime
	}
	return 0
}

func (m *TableDescriptor) GetReplacementOf() TableDescriptor_Replacement {
	if m != nil {
		return m.ReplacementOf
	}
	return TableDescriptor_Replacement{}
}

func (m *TableDescriptor) GetAuditMode() TableDescriptor_AuditMode {
	if m != nil {
		return m.AuditMode
	}
	return TableDescriptor_DISABLED
}

// The schema update lease. A single goroutine across a cockroach cluster
// can own it, and will execute pending schema changes for this table.
// Since the execution of a pending schema change is through transactions,
// it is legal for more than one goroutine to attempt to execute it. This
// lease reduces write contention on the schema change.
type TableDescriptor_SchemaChangeLease struct {
	//NodeID github_com_cockroachdb_cockroach_pkg_roachpb.NodeID `protobuf:"varint,1,opt,name=node_id,json=nodeId,casttype=github.com/cockroachdb/cockroach/pkg/roachpb.NodeID" json:"node_id"`
	// Nanoseconds since the Unix epoch.
	ExpirationTime int64 `protobuf:"varint,2,opt,name=expiration_time,json=expirationTime" json:"expiration_time"`
}

func (m *TableDescriptor_SchemaChangeLease) Reset()         { *m = TableDescriptor_SchemaChangeLease{} }
func (m *TableDescriptor_SchemaChangeLease) String() string { return proto.CompactTextString(m) }
func (*TableDescriptor_SchemaChangeLease) ProtoMessage()    {}
func (*TableDescriptor_SchemaChangeLease) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{8, 0}
}

type TableDescriptor_CheckConstraint struct {
	Expr     string             `protobuf:"bytes,1,opt,name=expr" json:"expr"`
	Name     string             `protobuf:"bytes,2,opt,name=name" json:"name"`
	Validity ConstraintValidity `protobuf:"varint,3,opt,name=validity,enum=cockroach.sql.sqlbase.ConstraintValidity" json:"validity"`
	// An ordered list of column IDs used by the check constraint.
	ColumnIDs []ColumnID `protobuf:"varint,5,rep,name=column_ids,json=columnIds,casttype=ColumnID" json:"column_ids,omitempty"`
}

func (m *TableDescriptor_CheckConstraint) Reset()         { *m = TableDescriptor_CheckConstraint{} }
func (m *TableDescriptor_CheckConstraint) String() string { return proto.CompactTextString(m) }
func (*TableDescriptor_CheckConstraint) ProtoMessage()    {}
func (*TableDescriptor_CheckConstraint) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{8, 1}
}

// A table descriptor is named through a name map stored in the
// system.namespace table: a map from {parent_id, table_name} -> id.
// This name map can be cached for performance on a node in the cluster
// making reassigning a name complicated. In particular, since a
// name cannot be withdrawn across a cluster in a transaction at
// timestamp T, we have to worry about the following:
//
// 1. A table is dropped at T, and the name and descriptor are still
// cached and used by transactions at timestamps >= T.
// 2. A table is renamed from foo to bar at T, and both names foo and bar
// can be used by transactions at timestamps >= T.
// 3. A name foo is reassigned from one table to another at T, and the name
// foo can reference two different tables at timestamps >= T.
//
// The system ensures that a name can be resolved only to a single
// descriptor at a timestamp thereby permitting 1 and 2, but not 3
// (the name references two tables).
//
// The transaction at T is followed by a time period when names no longer
// a part of the namespace are drained from the system. Once the old name
// is drained from the system another transaction at timestamp S is
// executed to release the name for future use. The interval from T to S
// is called the name drain interval: If the T transaction is removing
// the name foo then, at timestamps above S, foo can no longer be resolved.
//
// Consider a transaction at T in which name B is dropped, a new name C is
// created. Name C is viable as soon as the transaction commits.
// When the transaction at S commits, the name B is released for reuse.
//
// The transaction at S runs through the schema changer, with the system
// returning a response to the client initiating transaction T only after
// transaction at S is committed. So effectively the SQL transaction once
// it returns can be followed by SQL transactions that do not observe
// old name mappings.
//
// Note: an exception to this is #19925 which needs to be fixed.
//
// In order for transaction at S to act properly the system.namespace
// table entry for an old name references the descriptor who was the
// prior owner of the name requiring draining.
//
// Before T:   B -> Desc B
//
// After T and before S: B -> Desc B, C -> Desc C
//
// After S: C -> Desc C
//
// Between T and S the name B is drained and the system is unable
// to assign it to another descriptor.
//
// BEGIN;
// RENAME foo TO bar;
// CREATE foo;
//
// will fail because CREATE foo is executed at T.
//
// RENAME foo TO bar;
// CREATE foo;
//
// will succeed because the RENAME returns after S and CREATE foo is
// executed after S.
//
// The above scheme suffers from the problem that a transaction can observe
// the partial effect of a committed transaction during the drain interval.
// For instance during the drain interval a transaction can see the correct
// assignment for C, and the old assignments for B.
//
type TableDescriptor_NameInfo struct {
	// The database that the table belonged to before the rename (tables can be
	// renamed from one db to another).
	ParentID ID     `protobuf:"varint,1,opt,name=parent_id,json=parentId,casttype=ID" json:"parent_id"`
	Name     string `protobuf:"bytes,2,opt,name=name" json:"name"`
}

func (m *TableDescriptor_NameInfo) Reset()         { *m = TableDescriptor_NameInfo{} }
func (m *TableDescriptor_NameInfo) String() string { return proto.CompactTextString(m) }
func (*TableDescriptor_NameInfo) ProtoMessage()    {}
func (*TableDescriptor_NameInfo) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{8, 2}
}

type TableDescriptor_Reference struct {
	// The ID of the relation that depends on this one.
	ID ID `protobuf:"varint,1,opt,name=id,casttype=ID" json:"id"`
	// If applicable, the ID of this table's index that is referenced by the
	// dependent relation.
	IndexID IndexID `protobuf:"varint,2,opt,name=index_id,json=indexId,casttype=IndexID" json:"index_id"`
	// The IDs of this table's columns that are referenced by the dependent
	// relation.
	ColumnIDs []ColumnID `protobuf:"varint,3,rep,name=column_ids,json=columnIds,casttype=ColumnID" json:"column_ids,omitempty"`
}

func (m *TableDescriptor_Reference) Reset()         { *m = TableDescriptor_Reference{} }
func (m *TableDescriptor_Reference) String() string { return proto.CompactTextString(m) }
func (*TableDescriptor_Reference) ProtoMessage()    {}
func (*TableDescriptor_Reference) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{8, 3}
}

type TableDescriptor_MutationJob struct {
	// The mutation id of this mutation job.
	MutationID MutationID `protobuf:"varint,1,opt,name=mutation_id,json=mutationId,casttype=MutationID" json:"mutation_id"`
	// The job id for a mutation job is the id in the system.jobs table of the
	// schema change job executing the mutation referenced by mutation_id.
	JobID int64 `protobuf:"varint,2,opt,name=job_id,json=jobId" json:"job_id"`
}

func (m *TableDescriptor_MutationJob) Reset()         { *m = TableDescriptor_MutationJob{} }
func (m *TableDescriptor_MutationJob) String() string { return proto.CompactTextString(m) }
func (*TableDescriptor_MutationJob) ProtoMessage()    {}
func (*TableDescriptor_MutationJob) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{8, 4}
}

type TableDescriptor_SequenceOpts struct {
	// How much to increment the sequence by when nextval() is called.
	Increment int64 `protobuf:"varint,1,opt,name=increment" json:"increment"`
	// Minimum value of the sequence.
	MinValue int64 `protobuf:"varint,2,opt,name=min_value,json=minValue" json:"min_value"`
	// Maximum value of the sequence.
	MaxValue int64 `protobuf:"varint,3,opt,name=max_value,json=maxValue" json:"max_value"`
	// Start value of the sequence.
	Start int64 `protobuf:"varint,4,opt,name=start" json:"start"`
}

func (m *TableDescriptor_SequenceOpts) Reset()         { *m = TableDescriptor_SequenceOpts{} }
func (m *TableDescriptor_SequenceOpts) String() string { return proto.CompactTextString(m) }
func (*TableDescriptor_SequenceOpts) ProtoMessage()    {}
func (*TableDescriptor_SequenceOpts) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{8, 5}
}

type TableDescriptor_Replacement struct {
	ID   ID                           `protobuf:"varint,1,opt,name=id,casttype=ID" json:"id"`
	Time cockroach_util_hlc.Timestamp `protobuf:"bytes,2,opt,name=time" json:"time"`
}

func (m *TableDescriptor_Replacement) Reset()         { *m = TableDescriptor_Replacement{} }
func (m *TableDescriptor_Replacement) String() string { return proto.CompactTextString(m) }
func (*TableDescriptor_Replacement) ProtoMessage()    {}
func (*TableDescriptor_Replacement) Descriptor() ([]byte, []int) {
	return fileDescriptorStructured, []int{8, 6}
}

// DatabaseDescriptor represents a namespace (aka database) and is stored
// in a structured metadata key. The DatabaseDescriptor has a globally-unique
// ID shared with the TableDescriptor ID.
// Permissions are applied to all tables in the namespace.
type DatabaseDescriptor struct {
	Name       string               `protobuf:"bytes,1,opt,name=name" json:"name"`
	ID         ID                   `protobuf:"varint,2,opt,name=id,casttype=ID" json:"id"`
	//Privileges *PrivilegeDescriptor `protobuf:"bytes,3,opt,name=privileges" json:"privileges,omitempty"`
}

func (m *DatabaseDescriptor) Reset()                    { *m = DatabaseDescriptor{} }
func (m *DatabaseDescriptor) String() string            { return proto.CompactTextString(m) }
func (*DatabaseDescriptor) ProtoMessage()               {}
func (*DatabaseDescriptor) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{9} }

func (m *DatabaseDescriptor) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *DatabaseDescriptor) GetID() ID {
	if m != nil {
		return m.ID
	}
	return 0
}



// Descriptor is a union type holding either a table or database descriptor.
type Descriptor struct {
	// Types that are valid to be assigned to Union:
	//	*Descriptor_Table
	//	*Descriptor_Database
	Union isDescriptor_Union `protobuf_oneof:"union"`
}

func (m *Descriptor) Reset()                    { *m = Descriptor{} }
func (m *Descriptor) String() string            { return proto.CompactTextString(m) }
func (*Descriptor) ProtoMessage()               {}
func (*Descriptor) Descriptor() ([]byte, []int) { return fileDescriptorStructured, []int{10} }

type isDescriptor_Union interface {
	isDescriptor_Union()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Descriptor_Table struct {
	Table *TableDescriptor `protobuf:"bytes,1,opt,name=table,oneof"`
}
type Descriptor_Database struct {
	Database *DatabaseDescriptor `protobuf:"bytes,2,opt,name=database,oneof"`
}

func (*Descriptor_Table) isDescriptor_Union()    {}
func (*Descriptor_Database) isDescriptor_Union() {}

func (m *Descriptor) GetUnion() isDescriptor_Union {
	if m != nil {
		return m.Union
	}
	return nil
}

func (m *Descriptor) GetTable() *TableDescriptor {
	if x, ok := m.GetUnion().(*Descriptor_Table); ok {
		return x.Table
	}
	return nil
}

func (m *Descriptor) GetDatabase() *DatabaseDescriptor {
	if x, ok := m.GetUnion().(*Descriptor_Database); ok {
		return x.Database
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Descriptor) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Descriptor_OneofMarshaler, _Descriptor_OneofUnmarshaler, _Descriptor_OneofSizer, []interface{}{
		(*Descriptor_Table)(nil),
		(*Descriptor_Database)(nil),
	}
}

func _Descriptor_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Descriptor)
	// union
	switch x := m.Union.(type) {
	case *Descriptor_Table:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Table); err != nil {
			return err
		}
	case *Descriptor_Database:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Database); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Descriptor.Union has unexpected type %T", x)
	}
	return nil
}

func _Descriptor_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Descriptor)
	switch tag {
	case 1: // union.table
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(TableDescriptor)
		err := b.DecodeMessage(msg)
		m.Union = &Descriptor_Table{msg}
		return true, err
	case 2: // union.database
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(DatabaseDescriptor)
		err := b.DecodeMessage(msg)
		m.Union = &Descriptor_Database{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Descriptor_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Descriptor)
	// union
	switch x := m.Union.(type) {
	case *Descriptor_Table:
		s := proto.Size(x.Table)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Descriptor_Database:
		s := proto.Size(x.Database)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ColumnType)(nil), "cockroach.sql.sqlbase.ColumnType")
	proto.RegisterType((*ForeignKeyReference)(nil), "cockroach.sql.sqlbase.ForeignKeyReference")
	proto.RegisterType((*ColumnDescriptor)(nil), "cockroach.sql.sqlbase.ColumnDescriptor")
	proto.RegisterType((*ColumnFamilyDescriptor)(nil), "cockroach.sql.sqlbase.ColumnFamilyDescriptor")
	proto.RegisterType((*InterleaveDescriptor)(nil), "cockroach.sql.sqlbase.InterleaveDescriptor")
	proto.RegisterType((*InterleaveDescriptor_Ancestor)(nil), "cockroach.sql.sqlbase.InterleaveDescriptor.Ancestor")
	proto.RegisterType((*PartitioningDescriptor)(nil), "cockroach.sql.sqlbase.PartitioningDescriptor")
	proto.RegisterType((*PartitioningDescriptor_List)(nil), "cockroach.sql.sqlbase.PartitioningDescriptor.List")
	proto.RegisterType((*PartitioningDescriptor_Range)(nil), "cockroach.sql.sqlbase.PartitioningDescriptor.Range")
	proto.RegisterType((*IndexDescriptor)(nil), "cockroach.sql.sqlbase.IndexDescriptor")
	proto.RegisterType((*DescriptorMutation)(nil), "cockroach.sql.sqlbase.DescriptorMutation")
	proto.RegisterType((*TableDescriptor)(nil), "cockroach.sql.sqlbase.TableDescriptor")
	proto.RegisterType((*TableDescriptor_SchemaChangeLease)(nil), "cockroach.sql.sqlbase.TableDescriptor.SchemaChangeLease")
	proto.RegisterType((*TableDescriptor_CheckConstraint)(nil), "cockroach.sql.sqlbase.TableDescriptor.CheckConstraint")
	proto.RegisterType((*TableDescriptor_NameInfo)(nil), "cockroach.sql.sqlbase.TableDescriptor.NameInfo")
	proto.RegisterType((*TableDescriptor_Reference)(nil), "cockroach.sql.sqlbase.TableDescriptor.Reference")
	proto.RegisterType((*TableDescriptor_MutationJob)(nil), "cockroach.sql.sqlbase.TableDescriptor.MutationJob")
	proto.RegisterType((*TableDescriptor_SequenceOpts)(nil), "cockroach.sql.sqlbase.TableDescriptor.SequenceOpts")
	proto.RegisterType((*TableDescriptor_Replacement)(nil), "cockroach.sql.sqlbase.TableDescriptor.Replacement")
	proto.RegisterType((*DatabaseDescriptor)(nil), "cockroach.sql.sqlbase.DatabaseDescriptor")
	proto.RegisterType((*Descriptor)(nil), "cockroach.sql.sqlbase.Descriptor")
	proto.RegisterEnum("cockroach.sql.sqlbase.ConstraintValidity", ConstraintValidity_name, ConstraintValidity_value)
	proto.RegisterEnum("cockroach.sql.sqlbase.ColumnType_SemanticType", ColumnType_SemanticType_name, ColumnType_SemanticType_value)
	proto.RegisterEnum("cockroach.sql.sqlbase.ColumnType_VisibleType", ColumnType_VisibleType_name, ColumnType_VisibleType_value)
	proto.RegisterEnum("cockroach.sql.sqlbase.ForeignKeyReference_Action", ForeignKeyReference_Action_name, ForeignKeyReference_Action_value)
	proto.RegisterEnum("cockroach.sql.sqlbase.IndexDescriptor_Direction", IndexDescriptor_Direction_name, IndexDescriptor_Direction_value)
	proto.RegisterEnum("cockroach.sql.sqlbase.IndexDescriptor_Type", IndexDescriptor_Type_name, IndexDescriptor_Type_value)
	proto.RegisterEnum("cockroach.sql.sqlbase.DescriptorMutation_State", DescriptorMutation_State_name, DescriptorMutation_State_value)
	proto.RegisterEnum("cockroach.sql.sqlbase.DescriptorMutation_Direction", DescriptorMutation_Direction_name, DescriptorMutation_Direction_value)
	proto.RegisterEnum("cockroach.sql.sqlbase.TableDescriptor_State", TableDescriptor_State_name, TableDescriptor_State_value)
	proto.RegisterEnum("cockroach.sql.sqlbase.TableDescriptor_AuditMode", TableDescriptor_AuditMode_name, TableDescriptor_AuditMode_value)
}
func (this *ColumnType) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ColumnType)
	if !ok {
		that2, ok := that.(ColumnType)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.SemanticType != that1.SemanticType {
		return false
	}
	if this.Width != that1.Width {
		return false
	}
	if this.Precision != that1.Precision {
		return false
	}
	if len(this.ArrayDimensions) != len(that1.ArrayDimensions) {
		return false
	}
	for i := range this.ArrayDimensions {
		if this.ArrayDimensions[i] != that1.ArrayDimensions[i] {
			return false
		}
	}
	if this.Locale != nil && that1.Locale != nil {
		if *this.Locale != *that1.Locale {
			return false
		}
	} else if this.Locale != nil {
		return false
	} else if that1.Locale != nil {
		return false
	}
	if this.VisibleType != that1.VisibleType {
		return false
	}
	if this.ArrayContents != nil && that1.ArrayContents != nil {
		if *this.ArrayContents != *that1.ArrayContents {
			return false
		}
	} else if this.ArrayContents != nil {
		return false
	} else if that1.ArrayContents != nil {
		return false
	}
	if len(this.TupleContents) != len(that1.TupleContents) {
		return false
	}
	for i := range this.TupleContents {
		if !this.TupleContents[i].Equal(&that1.TupleContents[i]) {
			return false
		}
	}
	return true
}
func (m *ColumnType) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ColumnType) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.SemanticType))
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Width))
	dAtA[i] = 0x18
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Precision))
	if len(m.ArrayDimensions) > 0 {
		for _, num := range m.ArrayDimensions {
			dAtA[i] = 0x20
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	if m.Locale != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintStructured(dAtA, i, uint64(len(*m.Locale)))
		i += copy(dAtA[i:], *m.Locale)
	}
	dAtA[i] = 0x30
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.VisibleType))
	if m.ArrayContents != nil {
		dAtA[i] = 0x38
		i++
		i = encodeVarintStructured(dAtA, i, uint64(*m.ArrayContents))
	}
	if len(m.TupleContents) > 0 {
		for _, msg := range m.TupleContents {
			dAtA[i] = 0x42
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *ForeignKeyReference) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ForeignKeyReference) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Table))
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Index))
	dAtA[i] = 0x1a
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	dAtA[i] = 0x20
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Validity))
	dAtA[i] = 0x28
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.SharedPrefixLen))
	dAtA[i] = 0x30
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.OnDelete))
	dAtA[i] = 0x38
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.OnUpdate))
	return i, nil
}

func (m *ColumnDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ColumnDescriptor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ID))
	dAtA[i] = 0x1a
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Type.Size()))
	n1, err := m.Type.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	dAtA[i] = 0x20
	i++
	if m.Nullable {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i++
	if m.DefaultExpr != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintStructured(dAtA, i, uint64(len(*m.DefaultExpr)))
		i += copy(dAtA[i:], *m.DefaultExpr)
	}
	dAtA[i] = 0x30
	i++
	if m.Hidden {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i++
	if len(m.UsesSequenceIds) > 0 {
		for _, num := range m.UsesSequenceIds {
			dAtA[i] = 0x50
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	if m.ComputeExpr != nil {
		dAtA[i] = 0x5a
		i++
		i = encodeVarintStructured(dAtA, i, uint64(len(*m.ComputeExpr)))
		i += copy(dAtA[i:], *m.ComputeExpr)
	}
	return i, nil
}

func (m *ColumnFamilyDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ColumnFamilyDescriptor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ID))
	if len(m.ColumnNames) > 0 {
		for _, s := range m.ColumnNames {
			dAtA[i] = 0x1a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if len(m.ColumnIDs) > 0 {
		for _, num := range m.ColumnIDs {
			dAtA[i] = 0x20
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	dAtA[i] = 0x28
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.DefaultColumnID))
	return i, nil
}

func (m *InterleaveDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InterleaveDescriptor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Ancestors) > 0 {
		for _, msg := range m.Ancestors {
			dAtA[i] = 0xa
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *InterleaveDescriptor_Ancestor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InterleaveDescriptor_Ancestor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.TableID))
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.IndexID))
	dAtA[i] = 0x18
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.SharedPrefixLen))
	return i, nil
}

func (m *PartitioningDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PartitioningDescriptor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.NumColumns))
	if len(m.List) > 0 {
		for _, msg := range m.List {
			dAtA[i] = 0x12
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Range) > 0 {
		for _, msg := range m.Range {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *PartitioningDescriptor_List) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PartitioningDescriptor_List) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	if len(m.Values) > 0 {
		for _, b := range m.Values {
			dAtA[i] = 0x12
			i++
			i = encodeVarintStructured(dAtA, i, uint64(len(b)))
			i += copy(dAtA[i:], b)
		}
	}
	dAtA[i] = 0x1a
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Subpartitioning.Size()))
	n2, err := m.Subpartitioning.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	return i, nil
}

func (m *PartitioningDescriptor_Range) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PartitioningDescriptor_Range) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	if m.ToExclusive != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintStructured(dAtA, i, uint64(len(m.ToExclusive)))
		i += copy(dAtA[i:], m.ToExclusive)
	}
	if m.FromInclusive != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintStructured(dAtA, i, uint64(len(m.FromInclusive)))
		i += copy(dAtA[i:], m.FromInclusive)
	}
	return i, nil
}

func (m *IndexDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IndexDescriptor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ID))
	dAtA[i] = 0x18
	i++
	if m.Unique {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i++
	if len(m.ColumnNames) > 0 {
		for _, s := range m.ColumnNames {
			dAtA[i] = 0x22
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if len(m.StoreColumnNames) > 0 {
		for _, s := range m.StoreColumnNames {
			dAtA[i] = 0x2a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	if len(m.ColumnIDs) > 0 {
		for _, num := range m.ColumnIDs {
			dAtA[i] = 0x30
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	if len(m.ExtraColumnIDs) > 0 {
		for _, num := range m.ExtraColumnIDs {
			dAtA[i] = 0x38
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	if len(m.ColumnDirections) > 0 {
		for _, num := range m.ColumnDirections {
			dAtA[i] = 0x40
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	dAtA[i] = 0x4a
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ForeignKey.Size()))
	n3, err := m.ForeignKey.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	if len(m.ReferencedBy) > 0 {
		for _, msg := range m.ReferencedBy {
			dAtA[i] = 0x52
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	dAtA[i] = 0x5a
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Interleave.Size()))
	n4, err := m.Interleave.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	if len(m.InterleavedBy) > 0 {
		for _, msg := range m.InterleavedBy {
			dAtA[i] = 0x62
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.CompositeColumnIDs) > 0 {
		for _, num := range m.CompositeColumnIDs {
			dAtA[i] = 0x68
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	if len(m.StoreColumnIDs) > 0 {
		for _, num := range m.StoreColumnIDs {
			dAtA[i] = 0x70
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	dAtA[i] = 0x7a
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Partitioning.Size()))
	n5, err := m.Partitioning.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n5
	dAtA[i] = 0x80
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Type))
	return i, nil
}

func (m *DescriptorMutation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DescriptorMutation) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Descriptor_ != nil {
		nn6, err := m.Descriptor_.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn6
	}
	dAtA[i] = 0x18
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.State))
	dAtA[i] = 0x20
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Direction))
	dAtA[i] = 0x28
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.MutationID))
	dAtA[i] = 0x38
	i++
	if m.Rollback {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i++
	return i, nil
}

func (m *DescriptorMutation_Column) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Column != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintStructured(dAtA, i, uint64(m.Column.Size()))
		n7, err := m.Column.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n7
	}
	return i, nil
}
func (m *DescriptorMutation_Index) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Index != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintStructured(dAtA, i, uint64(m.Index.Size()))
		n8, err := m.Index.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n8
	}
	return i, nil
}
func (m *TableDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TableDescriptor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	dAtA[i] = 0x18
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ID))
	dAtA[i] = 0x20
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ParentID))
	dAtA[i] = 0x28
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Version))
	dAtA[i] = 0x30
	i++
	if m.UpVersion {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i++
	dAtA[i] = 0x3a
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ModificationTime.Size()))
	n9, err := m.ModificationTime.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n9
	if len(m.Columns) > 0 {
		for _, msg := range m.Columns {
			dAtA[i] = 0x42
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	dAtA[i] = 0x48
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.NextColumnID))
	dAtA[i] = 0x52
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.PrimaryIndex.Size()))
	n10, err := m.PrimaryIndex.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n10
	if len(m.Indexes) > 0 {
		for _, msg := range m.Indexes {
			dAtA[i] = 0x5a
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	dAtA[i] = 0x60
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.NextIndexID))
	// if m.Privileges != nil {
	// 	dAtA[i] = 0x6a
	// 	i++
	// 	i = encodeVarintStructured(dAtA, i, uint64(m.Privileges.Size()))
	// 	n11, err := m.Privileges.MarshalTo(dAtA[i:])
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	i += n11
	// }
	if len(m.Mutations) > 0 {
		for _, msg := range m.Mutations {
			dAtA[i] = 0x72
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.Lease != nil {
		dAtA[i] = 0x7a
		i++
		i = encodeVarintStructured(dAtA, i, uint64(m.Lease.Size()))
		n12, err := m.Lease.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n12
	}
	dAtA[i] = 0x80
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.NextMutationID))
	dAtA[i] = 0x88
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.FormatVersion))
	dAtA[i] = 0x98
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.State))
	if len(m.Checks) > 0 {
		for _, msg := range m.Checks {
			dAtA[i] = 0xa2
			i++
			dAtA[i] = 0x1
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.DrainingNames) > 0 {
		for _, msg := range m.DrainingNames {
			dAtA[i] = 0xaa
			i++
			dAtA[i] = 0x1
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Families) > 0 {
		for _, msg := range m.Families {
			dAtA[i] = 0xb2
			i++
			dAtA[i] = 0x1
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	dAtA[i] = 0xb8
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.NextFamilyID))
	dAtA[i] = 0xc2
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.ViewQuery)))
	i += copy(dAtA[i:], m.ViewQuery)
	if len(m.DependsOn) > 0 {
		for _, num := range m.DependsOn {
			dAtA[i] = 0xc8
			i++
			dAtA[i] = 0x1
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	if len(m.DependedOnBy) > 0 {
		for _, msg := range m.DependedOnBy {
			dAtA[i] = 0xd2
			i++
			dAtA[i] = 0x1
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.MutationJobs) > 0 {
		for _, msg := range m.MutationJobs {
			dAtA[i] = 0xda
			i++
			dAtA[i] = 0x1
			i++
			i = encodeVarintStructured(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.SequenceOpts != nil {
		dAtA[i] = 0xe2
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintStructured(dAtA, i, uint64(m.SequenceOpts.Size()))
		n13, err := m.SequenceOpts.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n13
	}
	dAtA[i] = 0xe8
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.DropTime))
	dAtA[i] = 0xf2
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ReplacementOf.Size()))
	n14, err := m.ReplacementOf.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n14
	dAtA[i] = 0xf8
	i++
	dAtA[i] = 0x1
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.AuditMode))
	return i, nil
}

func (m *TableDescriptor_SchemaChangeLease) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TableDescriptor_SchemaChangeLease) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(1))
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ExpirationTime))
	return i, nil
}

func (m *TableDescriptor_CheckConstraint) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TableDescriptor_CheckConstraint) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Expr)))
	i += copy(dAtA[i:], m.Expr)
	dAtA[i] = 0x12
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	dAtA[i] = 0x18
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Validity))
	if len(m.ColumnIDs) > 0 {
		for _, num := range m.ColumnIDs {
			dAtA[i] = 0x28
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	return i, nil
}

func (m *TableDescriptor_NameInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TableDescriptor_NameInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ParentID))
	dAtA[i] = 0x12
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	return i, nil
}

func (m *TableDescriptor_Reference) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TableDescriptor_Reference) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ID))
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.IndexID))
	if len(m.ColumnIDs) > 0 {
		for _, num := range m.ColumnIDs {
			dAtA[i] = 0x18
			i++
			i = encodeVarintStructured(dAtA, i, uint64(num))
		}
	}
	return i, nil
}

func (m *TableDescriptor_MutationJob) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TableDescriptor_MutationJob) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.MutationID))
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.JobID))
	return i, nil
}

func (m *TableDescriptor_SequenceOpts) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TableDescriptor_SequenceOpts) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Increment))
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.MinValue))
	dAtA[i] = 0x18
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.MaxValue))
	dAtA[i] = 0x20
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Start))
	return i, nil
}

func (m *TableDescriptor_Replacement) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TableDescriptor_Replacement) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ID))
	dAtA[i] = 0x12
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.Time.Size()))
	n15, err := m.Time.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n15
	return i, nil
}

func (m *DatabaseDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DatabaseDescriptor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintStructured(dAtA, i, uint64(len(m.Name)))
	i += copy(dAtA[i:], m.Name)
	dAtA[i] = 0x10
	i++
	i = encodeVarintStructured(dAtA, i, uint64(m.ID))
	// if m.Privileges != nil {
	// 	dAtA[i] = 0x1a
	// 	i++
	// 	i = encodeVarintStructured(dAtA, i, uint64(m.Privileges.Size()))
	// 	n16, err := m.Privileges.MarshalTo(dAtA[i:])
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	i += n16
	// }
	return i, nil
}

func (m *Descriptor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Descriptor) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Union != nil {
		nn17, err := m.Union.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn17
	}
	return i, nil
}

func (m *Descriptor_Table) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Table != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintStructured(dAtA, i, uint64(m.Table.Size()))
		n18, err := m.Table.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n18
	}
	return i, nil
}
func (m *Descriptor_Database) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Database != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintStructured(dAtA, i, uint64(m.Database.Size()))
		n19, err := m.Database.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n19
	}
	return i, nil
}
func encodeVarintStructured(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ColumnType) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(m.SemanticType))
	n += 1 + sovStructured(uint64(m.Width))
	n += 1 + sovStructured(uint64(m.Precision))
	if len(m.ArrayDimensions) > 0 {
		for _, e := range m.ArrayDimensions {
			n += 1 + sovStructured(uint64(e))
		}
	}
	if m.Locale != nil {
		l = len(*m.Locale)
		n += 1 + l + sovStructured(uint64(l))
	}
	n += 1 + sovStructured(uint64(m.VisibleType))
	if m.ArrayContents != nil {
		n += 1 + sovStructured(uint64(*m.ArrayContents))
	}
	if len(m.TupleContents) > 0 {
		for _, e := range m.TupleContents {
			l = e.Size()
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	return n
}

func (m *ForeignKeyReference) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(m.Table))
	n += 1 + sovStructured(uint64(m.Index))
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	n += 1 + sovStructured(uint64(m.Validity))
	n += 1 + sovStructured(uint64(m.SharedPrefixLen))
	n += 1 + sovStructured(uint64(m.OnDelete))
	n += 1 + sovStructured(uint64(m.OnUpdate))
	return n
}

func (m *ColumnDescriptor) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	n += 1 + sovStructured(uint64(m.ID))
	l = m.Type.Size()
	n += 1 + l + sovStructured(uint64(l))
	n += 2
	if m.DefaultExpr != nil {
		l = len(*m.DefaultExpr)
		n += 1 + l + sovStructured(uint64(l))
	}
	n += 2
	if len(m.UsesSequenceIds) > 0 {
		for _, e := range m.UsesSequenceIds {
			n += 1 + sovStructured(uint64(e))
		}
	}
	if m.ComputeExpr != nil {
		l = len(*m.ComputeExpr)
		n += 1 + l + sovStructured(uint64(l))
	}
	return n
}

func (m *ColumnFamilyDescriptor) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	n += 1 + sovStructured(uint64(m.ID))
	if len(m.ColumnNames) > 0 {
		for _, s := range m.ColumnNames {
			l = len(s)
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	if len(m.ColumnIDs) > 0 {
		for _, e := range m.ColumnIDs {
			n += 1 + sovStructured(uint64(e))
		}
	}
	n += 1 + sovStructured(uint64(m.DefaultColumnID))
	return n
}

func (m *InterleaveDescriptor) Size() (n int) {
	var l int
	_ = l
	if len(m.Ancestors) > 0 {
		for _, e := range m.Ancestors {
			l = e.Size()
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	return n
}

func (m *InterleaveDescriptor_Ancestor) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(m.TableID))
	n += 1 + sovStructured(uint64(m.IndexID))
	n += 1 + sovStructured(uint64(m.SharedPrefixLen))
	return n
}

func (m *PartitioningDescriptor) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(m.NumColumns))
	if len(m.List) > 0 {
		for _, e := range m.List {
			l = e.Size()
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	if len(m.Range) > 0 {
		for _, e := range m.Range {
			l = e.Size()
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	return n
}

func (m *PartitioningDescriptor_List) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	if len(m.Values) > 0 {
		for _, b := range m.Values {
			l = len(b)
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	l = m.Subpartitioning.Size()
	n += 1 + l + sovStructured(uint64(l))
	return n
}

func (m *PartitioningDescriptor_Range) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	if m.ToExclusive != nil {
		l = len(m.ToExclusive)
		n += 1 + l + sovStructured(uint64(l))
	}
	if m.FromInclusive != nil {
		l = len(m.FromInclusive)
		n += 1 + l + sovStructured(uint64(l))
	}
	return n
}

func (m *IndexDescriptor) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	n += 1 + sovStructured(uint64(m.ID))
	n += 2
	if len(m.ColumnNames) > 0 {
		for _, s := range m.ColumnNames {
			l = len(s)
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	if len(m.StoreColumnNames) > 0 {
		for _, s := range m.StoreColumnNames {
			l = len(s)
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	if len(m.ColumnIDs) > 0 {
		for _, e := range m.ColumnIDs {
			n += 1 + sovStructured(uint64(e))
		}
	}
	if len(m.ExtraColumnIDs) > 0 {
		for _, e := range m.ExtraColumnIDs {
			n += 1 + sovStructured(uint64(e))
		}
	}
	if len(m.ColumnDirections) > 0 {
		for _, e := range m.ColumnDirections {
			n += 1 + sovStructured(uint64(e))
		}
	}
	l = m.ForeignKey.Size()
	n += 1 + l + sovStructured(uint64(l))
	if len(m.ReferencedBy) > 0 {
		for _, e := range m.ReferencedBy {
			l = e.Size()
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	l = m.Interleave.Size()
	n += 1 + l + sovStructured(uint64(l))
	if len(m.InterleavedBy) > 0 {
		for _, e := range m.InterleavedBy {
			l = e.Size()
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	if len(m.CompositeColumnIDs) > 0 {
		for _, e := range m.CompositeColumnIDs {
			n += 1 + sovStructured(uint64(e))
		}
	}
	if len(m.StoreColumnIDs) > 0 {
		for _, e := range m.StoreColumnIDs {
			n += 1 + sovStructured(uint64(e))
		}
	}
	l = m.Partitioning.Size()
	n += 1 + l + sovStructured(uint64(l))
	n += 2 + sovStructured(uint64(m.Type))
	return n
}

func (m *DescriptorMutation) Size() (n int) {
	var l int
	_ = l
	if m.Descriptor_ != nil {
		n += m.Descriptor_.Size()
	}
	n += 1 + sovStructured(uint64(m.State))
	n += 1 + sovStructured(uint64(m.Direction))
	n += 1 + sovStructured(uint64(m.MutationID))
	n += 2
	return n
}

func (m *DescriptorMutation_Column) Size() (n int) {
	var l int
	_ = l
	if m.Column != nil {
		l = m.Column.Size()
		n += 1 + l + sovStructured(uint64(l))
	}
	return n
}
func (m *DescriptorMutation_Index) Size() (n int) {
	var l int
	_ = l
	if m.Index != nil {
		l = m.Index.Size()
		n += 1 + l + sovStructured(uint64(l))
	}
	return n
}
func (m *TableDescriptor) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	n += 1 + sovStructured(uint64(m.ID))
	n += 1 + sovStructured(uint64(m.ParentID))
	n += 1 + sovStructured(uint64(m.Version))
	n += 2
	l = m.ModificationTime.Size()
	n += 1 + l + sovStructured(uint64(l))
	if len(m.Columns) > 0 {
		for _, e := range m.Columns {
			l = e.Size()
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	n += 1 + sovStructured(uint64(m.NextColumnID))
	l = m.PrimaryIndex.Size()
	n += 1 + l + sovStructured(uint64(l))
	if len(m.Indexes) > 0 {
		for _, e := range m.Indexes {
			l = e.Size()
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	n += 1 + sovStructured(uint64(m.NextIndexID))
	// if m.Privileges != nil {
	// 	l = m.Privileges.Size()
	// 	n += 1 + l + sovStructured(uint64(l))
	// }
	if len(m.Mutations) > 0 {
		for _, e := range m.Mutations {
			l = e.Size()
			n += 1 + l + sovStructured(uint64(l))
		}
	}
	if m.Lease != nil {
		l = m.Lease.Size()
		n += 1 + l + sovStructured(uint64(l))
	}
	n += 2 + sovStructured(uint64(m.NextMutationID))
	n += 2 + sovStructured(uint64(m.FormatVersion))
	n += 2 + sovStructured(uint64(m.State))
	if len(m.Checks) > 0 {
		for _, e := range m.Checks {
			l = e.Size()
			n += 2 + l + sovStructured(uint64(l))
		}
	}
	if len(m.DrainingNames) > 0 {
		for _, e := range m.DrainingNames {
			l = e.Size()
			n += 2 + l + sovStructured(uint64(l))
		}
	}
	if len(m.Families) > 0 {
		for _, e := range m.Families {
			l = e.Size()
			n += 2 + l + sovStructured(uint64(l))
		}
	}
	n += 2 + sovStructured(uint64(m.NextFamilyID))
	l = len(m.ViewQuery)
	n += 2 + l + sovStructured(uint64(l))
	if len(m.DependsOn) > 0 {
		for _, e := range m.DependsOn {
			n += 2 + sovStructured(uint64(e))
		}
	}
	if len(m.DependedOnBy) > 0 {
		for _, e := range m.DependedOnBy {
			l = e.Size()
			n += 2 + l + sovStructured(uint64(l))
		}
	}
	if len(m.MutationJobs) > 0 {
		for _, e := range m.MutationJobs {
			l = e.Size()
			n += 2 + l + sovStructured(uint64(l))
		}
	}
	if m.SequenceOpts != nil {
		l = m.SequenceOpts.Size()
		n += 2 + l + sovStructured(uint64(l))
	}
	n += 2 + sovStructured(uint64(m.DropTime))
	l = m.ReplacementOf.Size()
	n += 2 + l + sovStructured(uint64(l))
	n += 2 + sovStructured(uint64(m.AuditMode))
	return n
}

func (m *TableDescriptor_SchemaChangeLease) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(1))
	n += 1 + sovStructured(uint64(m.ExpirationTime))
	return n
}

func (m *TableDescriptor_CheckConstraint) Size() (n int) {
	var l int
	_ = l
	l = len(m.Expr)
	n += 1 + l + sovStructured(uint64(l))
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	n += 1 + sovStructured(uint64(m.Validity))
	if len(m.ColumnIDs) > 0 {
		for _, e := range m.ColumnIDs {
			n += 1 + sovStructured(uint64(e))
		}
	}
	return n
}

func (m *TableDescriptor_NameInfo) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(m.ParentID))
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	return n
}

func (m *TableDescriptor_Reference) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(m.ID))
	n += 1 + sovStructured(uint64(m.IndexID))
	if len(m.ColumnIDs) > 0 {
		for _, e := range m.ColumnIDs {
			n += 1 + sovStructured(uint64(e))
		}
	}
	return n
}

func (m *TableDescriptor_MutationJob) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(m.MutationID))
	n += 1 + sovStructured(uint64(m.JobID))
	return n
}

func (m *TableDescriptor_SequenceOpts) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(m.Increment))
	n += 1 + sovStructured(uint64(m.MinValue))
	n += 1 + sovStructured(uint64(m.MaxValue))
	n += 1 + sovStructured(uint64(m.Start))
	return n
}

func (m *TableDescriptor_Replacement) Size() (n int) {
	var l int
	_ = l
	n += 1 + sovStructured(uint64(m.ID))
	l = m.Time.Size()
	n += 1 + l + sovStructured(uint64(l))
	return n
}

func (m *DatabaseDescriptor) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	n += 1 + l + sovStructured(uint64(l))
	n += 1 + sovStructured(uint64(m.ID))
	// if m.Privileges != nil {
	// 	l = m.Privileges.Size()
	// 	n += 1 + l + sovStructured(uint64(l))
	// }
	return n
}

func (m *Descriptor) Size() (n int) {
	var l int
	_ = l
	if m.Union != nil {
		n += m.Union.Size()
	}
	return n
}

func (m *Descriptor_Table) Size() (n int) {
	var l int
	_ = l
	if m.Table != nil {
		l = m.Table.Size()
		n += 1 + l + sovStructured(uint64(l))
	}
	return n
}
func (m *Descriptor_Database) Size() (n int) {
	var l int
	_ = l
	if m.Database != nil {
		l = m.Database.Size()
		n += 1 + l + sovStructured(uint64(l))
	}
	return n
}

func sovStructured(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozStructured(x uint64) (n int) {
	return sovStructured(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ColumnType) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ColumnType: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ColumnType: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SemanticType", wireType)
			}
			m.SemanticType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SemanticType |= (ColumnType_SemanticType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Width", wireType)
			}
			m.Width = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Width |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Precision", wireType)
			}
			m.Precision = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Precision |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType == 0 {
				var v int32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (int32(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ArrayDimensions = append(m.ArrayDimensions, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v int32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ArrayDimensions = append(m.ArrayDimensions, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ArrayDimensions", wireType)
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Locale", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Locale = &s
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VisibleType", wireType)
			}
			m.VisibleType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VisibleType |= (ColumnType_VisibleType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ArrayContents", wireType)
			}
			var v ColumnType_SemanticType
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (ColumnType_SemanticType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.ArrayContents = &v
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TupleContents", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TupleContents = append(m.TupleContents, ColumnType{})
			if err := m.TupleContents[len(m.TupleContents)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ForeignKeyReference) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ForeignKeyReference: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ForeignKeyReference: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Table", wireType)
			}
			m.Table = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Table |= (ID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= (IndexID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validity", wireType)
			}
			m.Validity = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Validity |= (ConstraintValidity(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SharedPrefixLen", wireType)
			}
			m.SharedPrefixLen = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SharedPrefixLen |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OnDelete", wireType)
			}
			m.OnDelete = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OnDelete |= (ForeignKeyReference_Action(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OnUpdate", wireType)
			}
			m.OnUpdate = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OnUpdate |= (ForeignKeyReference_Action(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ColumnDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ColumnDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ColumnDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= (ColumnID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Type.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nullable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Nullable = bool(v != 0)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultExpr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.DefaultExpr = &s
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hidden", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Hidden = bool(v != 0)
		case 10:
			if wireType == 0 {
				var v ID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.UsesSequenceIds = append(m.UsesSequenceIds, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v ID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.UsesSequenceIds = append(m.UsesSequenceIds, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field UsesSequenceIds", wireType)
			}
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ComputeExpr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.ComputeExpr = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ColumnFamilyDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ColumnFamilyDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ColumnFamilyDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= (FamilyID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ColumnNames", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ColumnNames = append(m.ColumnNames, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
			if wireType == 0 {
				var v ColumnID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ColumnID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ColumnIDs = append(m.ColumnIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v ColumnID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ColumnID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ColumnIDs = append(m.ColumnIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ColumnIDs", wireType)
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultColumnID", wireType)
			}
			m.DefaultColumnID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DefaultColumnID |= (ColumnID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *InterleaveDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: InterleaveDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InterleaveDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ancestors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ancestors = append(m.Ancestors, InterleaveDescriptor_Ancestor{})
			if err := m.Ancestors[len(m.Ancestors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *InterleaveDescriptor_Ancestor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Ancestor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Ancestor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TableID", wireType)
			}
			m.TableID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TableID |= (ID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IndexID", wireType)
			}
			m.IndexID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IndexID |= (IndexID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SharedPrefixLen", wireType)
			}
			m.SharedPrefixLen = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SharedPrefixLen |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PartitioningDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PartitioningDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PartitioningDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumColumns", wireType)
			}
			m.NumColumns = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumColumns |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field List", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.List = append(m.List, PartitioningDescriptor_List{})
			if err := m.List[len(m.List)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Range", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Range = append(m.Range, PartitioningDescriptor_Range{})
			if err := m.Range[len(m.Range)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PartitioningDescriptor_List) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: List: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: List: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Values", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Values = append(m.Values, make([]byte, postIndex-iNdEx))
			copy(m.Values[len(m.Values)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subpartitioning", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Subpartitioning.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PartitioningDescriptor_Range) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Range: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Range: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ToExclusive", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ToExclusive = append(m.ToExclusive[:0], dAtA[iNdEx:postIndex]...)
			if m.ToExclusive == nil {
				m.ToExclusive = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromInclusive", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FromInclusive = append(m.FromInclusive[:0], dAtA[iNdEx:postIndex]...)
			if m.FromInclusive == nil {
				m.FromInclusive = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *IndexDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IndexDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IndexDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= (IndexID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Unique", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Unique = bool(v != 0)
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ColumnNames", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ColumnNames = append(m.ColumnNames, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StoreColumnNames", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StoreColumnNames = append(m.StoreColumnNames, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 6:
			if wireType == 0 {
				var v ColumnID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ColumnID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ColumnIDs = append(m.ColumnIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v ColumnID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ColumnID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ColumnIDs = append(m.ColumnIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ColumnIDs", wireType)
			}
		case 7:
			if wireType == 0 {
				var v ColumnID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ColumnID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ExtraColumnIDs = append(m.ExtraColumnIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v ColumnID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ColumnID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ExtraColumnIDs = append(m.ExtraColumnIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtraColumnIDs", wireType)
			}
		case 8:
			if wireType == 0 {
				var v IndexDescriptor_Direction
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (IndexDescriptor_Direction(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ColumnDirections = append(m.ColumnDirections, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v IndexDescriptor_Direction
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (IndexDescriptor_Direction(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ColumnDirections = append(m.ColumnDirections, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ColumnDirections", wireType)
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ForeignKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ForeignKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReferencedBy", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReferencedBy = append(m.ReferencedBy, ForeignKeyReference{})
			if err := m.ReferencedBy[len(m.ReferencedBy)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Interleave", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Interleave.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterleavedBy", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InterleavedBy = append(m.InterleavedBy, ForeignKeyReference{})
			if err := m.InterleavedBy[len(m.InterleavedBy)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType == 0 {
				var v ColumnID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ColumnID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.CompositeColumnIDs = append(m.CompositeColumnIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v ColumnID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ColumnID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.CompositeColumnIDs = append(m.CompositeColumnIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field CompositeColumnIDs", wireType)
			}
		case 14:
			if wireType == 0 {
				var v ColumnID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ColumnID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.StoreColumnIDs = append(m.StoreColumnIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v ColumnID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ColumnID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.StoreColumnIDs = append(m.StoreColumnIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field StoreColumnIDs", wireType)
			}
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Partitioning", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Partitioning.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= (IndexDescriptor_Type(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DescriptorMutation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DescriptorMutation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DescriptorMutation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Column", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &ColumnDescriptor{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Descriptor_ = &DescriptorMutation_Column{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &IndexDescriptor{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Descriptor_ = &DescriptorMutation_Index{v}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= (DescriptorMutation_State(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Direction", wireType)
			}
			m.Direction = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Direction |= (DescriptorMutation_Direction(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MutationID", wireType)
			}
			m.MutationID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MutationID |= (MutationID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rollback", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Rollback = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TableDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TableDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TableDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= (ID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParentID", wireType)
			}
			m.ParentID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ParentID |= (ID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= (DescriptorVersion(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpVersion", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.UpVersion = bool(v != 0)
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ModificationTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ModificationTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Columns", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Columns = append(m.Columns, ColumnDescriptor{})
			if err := m.Columns[len(m.Columns)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextColumnID", wireType)
			}
			m.NextColumnID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextColumnID |= (ColumnID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PrimaryIndex", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.PrimaryIndex.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Indexes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Indexes = append(m.Indexes, IndexDescriptor{})
			if err := m.Indexes[len(m.Indexes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextIndexID", wireType)
			}
			m.NextIndexID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextIndexID |= (IndexID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Privileges", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			// if m.Privileges == nil {
			// 	m.Privileges = &PrivilegeDescriptor{}
			// }
			// if err := m.Privileges.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
			// 	return err
			// }
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Mutations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Mutations = append(m.Mutations, DescriptorMutation{})
			if err := m.Mutations[len(m.Mutations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Lease", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Lease == nil {
				m.Lease = &TableDescriptor_SchemaChangeLease{}
			}
			if err := m.Lease.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextMutationID", wireType)
			}
			m.NextMutationID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextMutationID |= (MutationID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 17:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FormatVersion", wireType)
			}
			m.FormatVersion = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FormatVersion |= (FormatVersion(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 19:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= (TableDescriptor_State(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 20:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Checks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Checks = append(m.Checks, &TableDescriptor_CheckConstraint{})
			if err := m.Checks[len(m.Checks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 21:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DrainingNames", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DrainingNames = append(m.DrainingNames, TableDescriptor_NameInfo{})
			if err := m.DrainingNames[len(m.DrainingNames)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 22:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Families", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Families = append(m.Families, ColumnFamilyDescriptor{})
			if err := m.Families[len(m.Families)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 23:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NextFamilyID", wireType)
			}
			m.NextFamilyID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NextFamilyID |= (FamilyID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 24:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ViewQuery", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ViewQuery = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 25:
			if wireType == 0 {
				var v ID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.DependsOn = append(m.DependsOn, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v ID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.DependsOn = append(m.DependsOn, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field DependsOn", wireType)
			}
		case 26:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DependedOnBy", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DependedOnBy = append(m.DependedOnBy, TableDescriptor_Reference{})
			if err := m.DependedOnBy[len(m.DependedOnBy)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 27:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MutationJobs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MutationJobs = append(m.MutationJobs, TableDescriptor_MutationJob{})
			if err := m.MutationJobs[len(m.MutationJobs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 28:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SequenceOpts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SequenceOpts == nil {
				m.SequenceOpts = &TableDescriptor_SequenceOpts{}
			}
			if err := m.SequenceOpts.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 29:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DropTime", wireType)
			}
			m.DropTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DropTime |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 30:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReplacementOf", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ReplacementOf.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 31:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuditMode", wireType)
			}
			m.AuditMode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AuditMode |= (TableDescriptor_AuditMode(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TableDescriptor_SchemaChangeLease) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SchemaChangeLease: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SchemaChangeLease: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeID", wireType)
			}
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				//m.NodeID |= (github_com_cockroachdb_cockroach_pkg_roachpb.NodeID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpirationTime", wireType)
			}
			m.ExpirationTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpirationTime |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TableDescriptor_CheckConstraint) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CheckConstraint: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CheckConstraint: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Expr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Expr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validity", wireType)
			}
			m.Validity = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Validity |= (ConstraintValidity(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType == 0 {
				var v ColumnID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ColumnID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ColumnIDs = append(m.ColumnIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v ColumnID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ColumnID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ColumnIDs = append(m.ColumnIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ColumnIDs", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TableDescriptor_NameInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: NameInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NameInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ParentID", wireType)
			}
			m.ParentID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ParentID |= (ID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TableDescriptor_Reference) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Reference: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Reference: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= (ID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IndexID", wireType)
			}
			m.IndexID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IndexID |= (IndexID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType == 0 {
				var v ColumnID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (ColumnID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ColumnIDs = append(m.ColumnIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthStructured
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v ColumnID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStructured
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (ColumnID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ColumnIDs = append(m.ColumnIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ColumnIDs", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TableDescriptor_MutationJob) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MutationJob: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MutationJob: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MutationID", wireType)
			}
			m.MutationID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MutationID |= (MutationID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field JobID", wireType)
			}
			m.JobID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.JobID |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TableDescriptor_SequenceOpts) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SequenceOpts: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SequenceOpts: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Increment", wireType)
			}
			m.Increment = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Increment |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinValue", wireType)
			}
			m.MinValue = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinValue |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxValue", wireType)
			}
			m.MaxValue = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxValue |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Start", wireType)
			}
			m.Start = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Start |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TableDescriptor_Replacement) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Replacement: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Replacement: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= (ID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Time.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DatabaseDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DatabaseDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DatabaseDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= (ID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Privileges", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			// if m.Privileges == nil {
			// 	m.Privileges = &PrivilegeDescriptor{}
			// }
			// if err := m.Privileges.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
			// 	return err
			// }
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Descriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Descriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Descriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Table", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &TableDescriptor{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Union = &Descriptor_Table{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Database", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStructured
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &DatabaseDescriptor{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Union = &Descriptor_Database{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStructured(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStructured
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipStructured(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStructured
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStructured
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthStructured
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowStructured
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipStructured(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthStructured = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStructured   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("sql/sqlbase/structured.proto", fileDescriptorStructured) }

var fileDescriptorStructured = []byte{
	// 2986 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x59, 0x5b, 0x6f, 0xe3, 0xc6,
	0x15, 0x36, 0x75, 0xd7, 0xd1, 0x8d, 0x9e, 0xbd, 0x44, 0xab, 0x6c, 0x6c, 0xaf, 0x92, 0x4d, 0x9d,
	0x9b, 0xbc, 0xf1, 0xa6, 0x6d, 0x90, 0x16, 0x41, 0x25, 0x91, 0xde, 0xa5, 0x57, 0x96, 0xbc, 0x94,
	0xec, 0xcd, 0x06, 0x69, 0x09, 0x5a, 0x1c, 0xdb, 0xcc, 0x52, 0xa4, 0x96, 0xa4, 0x1c, 0xeb, 0x1f,
	0xe4, 0xb1, 0x0f, 0x05, 0xfa, 0x16, 0x04, 0x41, 0xdf, 0xfa, 0x07, 0xfa, 0x13, 0xb6, 0x40, 0x1f,
	0x8a, 0x3e, 0xf5, 0xa5, 0x46, 0xeb, 0xa2, 0x40, 0x7f, 0x41, 0x1f, 0x02, 0x14, 0x2d, 0x66, 0x38,
	0x43, 0x51, 0xb6, 0xe5, 0xc8, 0xbb, 0x6f, 0xe2, 0x99, 0x73, 0xbe, 0x99, 0x73, 0xe6, 0x5c, 0x47,
	0x70, 0xdb, 0x7b, 0x6e, 0xad, 0x79, 0xcf, 0xad, 0x3d, 0xdd, 0xc3, 0x6b, 0x9e, 0xef, 0x8e, 0xfa,
	0xfe, 0xc8, 0xc5, 0x46, 0x6d, 0xe8, 0x3a, 0xbe, 0x83, 0x6e, 0xf4, 0x9d, 0xfe, 0x33, 0xd7, 0xd1,
	0xfb, 0x87, 0x35, 0xef, 0xb9, 0x55, 0x63, 0x7c, 0x95, 0xf2, 0xc8, 0x37, 0xad, 0xb5, 0x43, 0xab,
	0xbf, 0xe6, 0x9b, 0x03, 0xec, 0xf9, 0xfa, 0x60, 0x18, 0x08, 0x54, 0x5e, 0x8f, 0xc2, 0x0d, 0x5d,
	0xf3, 0xc8, 0xb4, 0xf0, 0x01, 0x66, 0x8b, 0xd7, 0x0f, 0x9c, 0x03, 0x87, 0xfe, 0x5c, 0x23, 0xbf,
	0x02, 0x6a, 0xf5, 0x7f, 0x29, 0x80, 0xa6, 0x63, 0x8d, 0x06, 0x76, 0x6f, 0x3c, 0xc4, 0xe8, 0x29,
	0x14, 0x3c, 0x3c, 0xd0, 0x6d, 0xdf, 0xec, 0x6b, 0xfe, 0x78, 0x88, 0xcb, 0xc2, 0x8a, 0xb0, 0x5a,
	0x5c, 0xaf, 0xd5, 0x2e, 0x3c, 0x4a, 0x6d, 0x22, 0x59, 0xeb, 0x32, 0x31, 0xf2, 0xd1, 0x48, 0xbc,
	0x38, 0x59, 0x5e, 0x50, 0xf3, 0x5e, 0x84, 0x86, 0x2a, 0x90, 0xfc, 0xca, 0x34, 0xfc, 0xc3, 0x72,
	0x6c, 0x45, 0x58, 0x4d, 0x32, 0x96, 0x80, 0x84, 0xaa, 0x90, 0x1d, 0xba, 0xb8, 0x6f, 0x7a, 0xa6,
	0x63, 0x97, 0xe3, 0x91, 0xf5, 0x09, 0x19, 0xbd, 0x03, 0xa2, 0xee, 0xba, 0xfa, 0x58, 0x33, 0xcc,
	0x01, 0xb6, 0x09, 0xc9, 0x2b, 0x27, 0x56, 0xe2, 0xab, 0x49, 0xb5, 0x44, 0xe9, 0x52, 0x48, 0x46,
	0x37, 0x21, 0x65, 0x39, 0x7d, 0xdd, 0xc2, 0xe5, 0xe4, 0x8a, 0xb0, 0x9a, 0x55, 0xd9, 0x17, 0xda,
	0x85, 0xfc, 0x91, 0xe9, 0x99, 0x7b, 0x16, 0x0e, 0x94, 0x4b, 0x51, 0xe5, 0x3e, 0xf8, 0x61, 0xe5,
	0x76, 0x03, 0xa9, 0x88, 0x6e, 0xb9, 0xa3, 0x09, 0x09, 0xed, 0x40, 0x31, 0x38, 0x5a, 0xdf, 0xb1,
	0x7d, 0x6c, 0xfb, 0x5e, 0x39, 0xfd, 0x32, 0x66, 0x53, 0x0b, 0x14, 0xa5, 0xc9, 0x40, 0x50, 0x1b,
	0x8a, 0xfe, 0x68, 0x68, 0xe1, 0x09, 0x6c, 0x66, 0x25, 0xbe, 0x9a, 0x5b, 0xbf, 0xf3, 0x83, 0xb0,
	0xec, 0x90, 0x05, 0x2a, 0xce, 0xf1, 0xaa, 0xbf, 0x8b, 0x41, 0x3e, 0xba, 0x1f, 0xca, 0x40, 0xa2,
	0xd1, 0xe9, 0xb4, 0xc4, 0x05, 0x94, 0x86, 0xb8, 0xd2, 0xee, 0x89, 0x02, 0xca, 0x42, 0x72, 0xa3,
	0xd5, 0xa9, 0xf7, 0xc4, 0x18, 0xca, 0x41, 0x5a, 0x92, 0x9b, 0xca, 0x56, 0xbd, 0x25, 0xc6, 0x09,
	0xab, 0x54, 0xef, 0xc9, 0x62, 0x02, 0x15, 0x20, 0xdb, 0x53, 0xb6, 0xe4, 0x6e, 0xaf, 0xbe, 0xb5,
	0x2d, 0x26, 0x51, 0x1e, 0x32, 0x4a, 0xbb, 0x27, 0xab, 0xbb, 0xf5, 0x96, 0x98, 0x42, 0x00, 0xa9,
	0x6e, 0x4f, 0x55, 0xda, 0x0f, 0xc4, 0x34, 0x81, 0x6a, 0x3c, 0xed, 0xc9, 0x5d, 0x31, 0x83, 0x4a,
	0x90, 0x0b, 0x65, 0x7a, 0x9f, 0x8b, 0x59, 0x84, 0xa0, 0xd8, 0xec, 0xb4, 0x5a, 0xf5, 0x9e, 0x2c,
	0x31, 0x7e, 0x20, 0x5b, 0xb4, 0xeb, 0x5b, 0xb2, 0x98, 0x23, 0xa7, 0xe9, 0x28, 0x92, 0x98, 0xa7,
	0xa4, 0x9d, 0x56, 0x4b, 0x2c, 0x90, 0x5f, 0x3b, 0x3b, 0x8a, 0x24, 0x16, 0x09, 0x6c, 0x5d, 0x55,
	0xeb, 0x4f, 0xc5, 0x12, 0x21, 0x2a, 0x6d, 0xb9, 0x27, 0x8a, 0xe4, 0x17, 0xd9, 0x40, 0x5c, 0x24,
	0xbf, 0x36, 0xbb, 0x9d, 0xb6, 0x88, 0xc8, 0x59, 0x08, 0xad, 0xf7, 0xb9, 0x78, 0x8d, 0x08, 0xf5,
	0x76, 0xb6, 0x5b, 0xb2, 0x78, 0x1d, 0x95, 0x00, 0x94, 0x76, 0x6f, 0x7d, 0x57, 0x6e, 0xf6, 0x3a,
	0xaa, 0xf8, 0x42, 0x40, 0x45, 0xc8, 0x76, 0x14, 0x89, 0x7d, 0xff, 0x51, 0xa8, 0x1e, 0x40, 0x2e,
	0x72, 0xdf, 0xf4, 0x0c, 0x9d, 0xb6, 0x2c, 0x2e, 0x10, 0x83, 0x10, 0x55, 0x1f, 0xc8, 0xaa, 0x28,
	0x10, 0xbd, 0xbb, 0x5b, 0xf5, 0x56, 0x8b, 0x98, 0x2d, 0x46, 0xf6, 0x6a, 0x28, 0x0f, 0xc8, 0xef,
	0x38, 0x39, 0x7d, 0x43, 0xe9, 0x89, 0x09, 0x22, 0xa9, 0xca, 0xf5, 0x96, 0x98, 0x44, 0xd7, 0x41,
	0x94, 0x3a, 0x3b, 0x8d, 0x96, 0xac, 0x6d, 0xab, 0x72, 0x53, 0xe9, 0x2a, 0x9d, 0xb6, 0x98, 0xfa,
	0x24, 0xf1, 0xef, 0x6f, 0x97, 0x85, 0xea, 0x7f, 0xe2, 0x70, 0x6d, 0xc3, 0x71, 0xb1, 0x79, 0x60,
	0x3f, 0xc2, 0x63, 0x15, 0xef, 0x63, 0x17, 0xdb, 0x7d, 0x8c, 0x56, 0x20, 0xe9, 0xeb, 0x7b, 0x56,
	0x10, 0x82, 0x85, 0x06, 0x90, 0x1b, 0xfd, 0xfe, 0x64, 0x39, 0xa6, 0x48, 0x6a, 0xb0, 0x80, 0xee,
	0x42, 0xd2, 0xb4, 0x0d, 0x7c, 0x4c, 0x23, 0xaa, 0xd0, 0x28, 0x31, 0x8e, 0xb4, 0x42, 0x88, 0x84,
	0x8d, 0xae, 0xa2, 0x32, 0x24, 0x6c, 0x7d, 0x80, 0x69, 0x5c, 0x65, 0x99, 0x67, 0x50, 0x0a, 0x7a,
	0x04, 0x99, 0x23, 0xdd, 0x32, 0x0d, 0xd3, 0x1f, 0x97, 0x13, 0xd4, 0x63, 0xdf, 0x99, 0xe9, 0x5a,
	0xb6, 0xe7, 0xbb, 0xba, 0x69, 0xfb, 0xbb, 0x4c, 0x80, 0x01, 0x85, 0x00, 0xe8, 0x1e, 0x2c, 0x7a,
	0x87, 0xba, 0x8b, 0x0d, 0x6d, 0xe8, 0xe2, 0x7d, 0xf3, 0x58, 0xb3, 0xb0, 0x4d, 0xe3, 0x8f, 0xc7,
	0x72, 0x29, 0x58, 0xde, 0xa6, 0xab, 0x2d, 0x6c, 0xa3, 0x1e, 0x64, 0x1d, 0x5b, 0x33, 0xb0, 0x85,
	0x7d, 0x1e, 0x8b, 0x1f, 0xce, 0xd8, 0xff, 0x02, 0x03, 0xd5, 0xea, 0x7d, 0xdf, 0x74, 0x6c, 0x7e,
	0x0e, 0xc7, 0x96, 0x28, 0x10, 0x43, 0x1d, 0x0d, 0x0d, 0xdd, 0xc7, 0x2c, 0x0e, 0x5f, 0x05, 0x75,
	0x87, 0x02, 0x55, 0x1f, 0x43, 0x2a, 0x58, 0x21, 0xfe, 0xdf, 0xee, 0x68, 0xf5, 0x66, 0x8f, 0x5c,
	0xe2, 0x02, 0xf1, 0x03, 0x55, 0x26, 0x3e, 0xdc, 0xec, 0x31, 0xaf, 0x90, 0x7b, 0x1a, 0x75, 0xda,
	0x18, 0x71, 0x7b, 0xf2, 0x25, 0xc9, 0x1b, 0xf5, 0x9d, 0x16, 0x71, 0x8d, 0x1c, 0xa4, 0x9b, 0xf5,
	0x6e, 0xb3, 0x2e, 0xc9, 0x62, 0xa2, 0xfa, 0xb7, 0x18, 0x88, 0x41, 0xc8, 0x4a, 0xd8, 0xeb, 0xbb,
	0xe6, 0xd0, 0x77, 0xdc, 0xf0, 0xb2, 0x84, 0x73, 0x97, 0xf5, 0x36, 0xc4, 0x4c, 0x83, 0x5d, 0xf5,
	0x4d, 0x42, 0x3f, 0xa5, 0xce, 0xf0, 0xfd, 0xc9, 0x72, 0x26, 0x40, 0x51, 0x24, 0x35, 0x66, 0x1a,
	0xe8, 0x67, 0x90, 0xa0, 0xc9, 0x8d, 0x5c, 0xf7, 0x15, 0x72, 0x05, 0x15, 0x42, 0x2b, 0x90, 0xb1,
	0x47, 0x96, 0x45, 0xfd, 0x8e, 0x78, 0x44, 0x86, 0x1b, 0x82, 0x53, 0xd1, 0x1d, 0xc8, 0x1b, 0x78,
	0x5f, 0x1f, 0x59, 0xbe, 0x86, 0x8f, 0x87, 0x2e, 0xcb, 0xb0, 0x39, 0x46, 0x93, 0x8f, 0x87, 0x2e,
	0xba, 0x0d, 0xa9, 0x43, 0xd3, 0x30, 0xb0, 0x4d, 0x2f, 0x95, 0x43, 0x30, 0x1a, 0x5a, 0x87, 0xc5,
	0x91, 0x87, 0x3d, 0xcd, 0xc3, 0xcf, 0x47, 0xc4, 0xe2, 0x9a, 0x69, 0x78, 0x65, 0x58, 0x89, 0xaf,
	0x16, 0x1a, 0x29, 0xe6, 0xdf, 0x25, 0xc2, 0xd0, 0x65, 0xeb, 0x8a, 0xe1, 0x91, 0x4d, 0xfb, 0xce,
	0x60, 0x38, 0xf2, 0x71, 0xb0, 0x69, 0x2e, 0xd8, 0x94, 0xd1, 0xc8, 0xa6, 0x9b, 0x89, 0x4c, 0x46,
	0xcc, 0x6e, 0x26, 0x32, 0x59, 0x11, 0x36, 0x13, 0x99, 0xb4, 0x98, 0xa9, 0x7e, 0x1d, 0x83, 0x9b,
	0x81, 0x9a, 0x1b, 0xfa, 0xc0, 0xb4, 0xc6, 0xaf, 0x6a, 0xe5, 0x00, 0x85, 0x59, 0x99, 0x9e, 0x88,
	0x60, 0x6b, 0x44, 0xcc, 0x2b, 0xc7, 0x57, 0xe2, 0xc1, 0x89, 0x08, 0xad, 0x4d, 0x48, 0xe8, 0x63,
	0x00, 0xc6, 0x42, 0x34, 0x4c, 0x50, 0x0d, 0x6f, 0x9d, 0x9e, 0x2c, 0x67, 0xf9, 0x75, 0x79, 0x53,
	0x77, 0x97, 0x0d, 0x98, 0x89, 0xba, 0x1d, 0x58, 0xe4, 0x36, 0x0e, 0x11, 0xa8, 0xa1, 0x0b, 0x8d,
	0x37, 0xd9, 0x99, 0x4a, 0x52, 0xc0, 0xc0, 0xc5, 0xa7, 0xa0, 0x4a, 0xc6, 0xd4, 0xa2, 0x51, 0xfd,
	0x7d, 0x0c, 0xae, 0x2b, 0xb6, 0x8f, 0x5d, 0x0b, 0xeb, 0x47, 0x38, 0x62, 0x88, 0xcf, 0x20, 0xab,
	0xdb, 0x7d, 0xec, 0xf9, 0x8e, 0xeb, 0x95, 0x05, 0x5a, 0x5d, 0x3e, 0x9a, 0xe1, 0x31, 0x17, 0xc9,
	0xd7, 0xea, 0x4c, 0x98, 0x97, 0xeb, 0x10, 0xac, 0xf2, 0x07, 0x01, 0x32, 0x7c, 0x15, 0xdd, 0x83,
	0x0c, 0x4d, 0x59, 0x44, 0x8f, 0x20, 0x9d, 0xdd, 0x60, 0x7a, 0xa4, 0x7b, 0x84, 0x4e, 0xcf, 0x4f,
	0x6e, 0x3e, 0x4d, 0xd9, 0x14, 0x03, 0xfd, 0x18, 0x32, 0x34, 0x7b, 0x69, 0xe1, 0x6d, 0x54, 0xb8,
	0x04, 0x4b, 0x6f, 0xd1, 0x4c, 0x97, 0xa6, 0xbc, 0x8a, 0x81, 0x9a, 0x17, 0x25, 0xa1, 0x38, 0x95,
	0x7f, 0x8d, 0x5b, 0xae, 0x3b, 0x9d, 0x86, 0xce, 0xe5, 0xa5, 0xea, 0xbf, 0xe2, 0x70, 0x73, 0x5b,
	0x77, 0x7d, 0x93, 0xc4, 0xbb, 0x69, 0x1f, 0x44, 0xec, 0x75, 0x17, 0x72, 0xf6, 0x68, 0xc0, 0x6e,
	0xc5, 0x63, 0xba, 0x04, 0xba, 0x83, 0x3d, 0x1a, 0x04, 0x06, 0xf7, 0x50, 0x0b, 0x12, 0x96, 0xe9,
	0xf9, 0xe5, 0x18, 0xb5, 0xe8, 0xfa, 0x0c, 0x8b, 0x5e, 0xbc, 0x47, 0xad, 0x65, 0x7a, 0x3e, 0xf7,
	0x49, 0x82, 0x82, 0x3a, 0x90, 0x74, 0x75, 0xfb, 0x00, 0x53, 0x27, 0xcb, 0xad, 0xdf, 0xbf, 0x1a,
	0x9c, 0x4a, 0x44, 0x79, 0xbb, 0x45, 0x71, 0x2a, 0xbf, 0x15, 0x20, 0x41, 0x76, 0xb9, 0x24, 0x0e,
	0x6e, 0x42, 0xea, 0x48, 0xb7, 0x46, 0xd8, 0xa3, 0x3a, 0xe4, 0x55, 0xf6, 0x85, 0x7e, 0x09, 0x25,
	0x6f, 0xb4, 0x37, 0x8c, 0x6c, 0xc5, 0x12, 0xcd, 0x07, 0x57, 0x3a, 0x55, 0x58, 0x12, 0xa6, 0xb1,
	0x2a, 0xcf, 0x20, 0x49, 0xcf, 0x7b, 0xc9, 0xc9, 0xee, 0x40, 0xde, 0x77, 0x34, 0x7c, 0xdc, 0xb7,
	0x46, 0x9e, 0x79, 0x84, 0xa9, 0x77, 0xe4, 0xd5, 0x9c, 0xef, 0xc8, 0x9c, 0x84, 0xee, 0x42, 0x71,
	0xdf, 0x75, 0x06, 0x9a, 0x69, 0x73, 0xa6, 0x38, 0x65, 0x2a, 0x10, 0xaa, 0xc2, 0x89, 0xd5, 0xff,
	0x66, 0xa0, 0x44, 0x3d, 0x68, 0xae, 0xcc, 0x70, 0x37, 0x92, 0x19, 0x6e, 0x4c, 0x65, 0x86, 0xd0,
	0x0d, 0x49, 0x62, 0xb8, 0x0d, 0xa9, 0x91, 0x6d, 0x3e, 0x1f, 0x05, 0x7b, 0x86, 0xc9, 0x2f, 0xa0,
	0x9d, 0x4b, 0x1b, 0x89, 0xf3, 0x69, 0xe3, 0x7d, 0x40, 0x24, 0x66, 0xb0, 0x36, 0xc5, 0x98, 0xa4,
	0x8c, 0x22, 0x5d, 0x69, 0xce, 0x4c, 0x32, 0xa9, 0x2b, 0x24, 0x99, 0x87, 0x20, 0xe2, 0x63, 0xdf,
	0xd5, 0xb5, 0x88, 0x7c, 0x9a, 0xca, 0x2f, 0x9d, 0x9e, 0x2c, 0x17, 0x65, 0xb2, 0x76, 0x31, 0x48,
	0x11, 0x47, 0xd6, 0x0c, 0xe2, 0x13, 0x8b, 0x0c, 0xc3, 0x30, 0x5d, 0x4c, 0xab, 0x64, 0xd0, 0xaa,
	0x16, 0xd7, 0xef, 0xcd, 0x4c, 0x26, 0x53, 0x66, 0xaf, 0x49, 0x5c, 0x50, 0x15, 0x03, 0xa8, 0x90,
	0xe0, 0xa1, 0xc7, 0x90, 0xdb, 0x0f, 0x0a, 0xb5, 0xf6, 0x0c, 0x8f, 0xcb, 0x59, 0xea, 0x6e, 0xef,
	0xce, 0x5f, 0xd2, 0x79, 0x7c, 0xee, 0x87, 0x4b, 0x68, 0x07, 0x0a, 0x2e, 0x5f, 0x36, 0xb4, 0xbd,
	0x31, 0xad, 0x3f, 0x2f, 0x03, 0x9a, 0x9f, 0xc0, 0x34, 0xc6, 0xe8, 0x31, 0x80, 0x19, 0x66, 0x49,
	0x5a, 0xa4, 0x72, 0xeb, 0xef, 0x5d, 0x21, 0x9d, 0xf2, 0x93, 0x4e, 0x40, 0xd0, 0x13, 0x28, 0x4e,
	0xbe, 0xe8, 0x51, 0xf3, 0x2f, 0x79, 0xd4, 0x42, 0x04, 0xa7, 0x31, 0x46, 0x3d, 0xb8, 0x4e, 0xca,
	0xa7, 0xe3, 0x99, 0x3e, 0x8e, 0xba, 0x40, 0x81, 0xba, 0x40, 0xf5, 0xf4, 0x64, 0x19, 0x35, 0xf9,
	0xfa, 0xc5, 0x6e, 0x80, 0xfa, 0x67, 0xd6, 0x03, 0xa7, 0x9a, 0x72, 0x5e, 0x82, 0x58, 0x9c, 0x38,
	0x55, 0x77, 0xe2, 0xbe, 0xe7, 0x9c, 0x2a, 0xe2, 0xda, 0x04, 0xe9, 0x09, 0xe4, 0xa7, 0xb2, 0x4c,
	0xe9, 0xe5, 0xb3, 0xcc, 0x14, 0x10, 0x92, 0x59, 0x7f, 0x24, 0xd2, 0xd6, 0xf0, 0xbd, 0x39, 0x1d,
	0xf4, 0x6c, 0xa7, 0x54, 0x5d, 0x82, 0x6c, 0xe8, 0xa3, 0xa4, 0xe5, 0xaf, 0x77, 0x9b, 0xe2, 0x02,
	0x1d, 0x93, 0xe4, 0x6e, 0x53, 0x14, 0xaa, 0x77, 0x20, 0x41, 0xc7, 0x87, 0x1c, 0xa4, 0x37, 0x3a,
	0xea, 0x93, 0xba, 0x2a, 0x05, 0xcd, 0xa2, 0xd2, 0xde, 0x95, 0xd5, 0x9e, 0x2c, 0x89, 0x42, 0xf5,
	0xbb, 0x04, 0xa0, 0xc9, 0x16, 0x5b, 0x23, 0x5f, 0xa7, 0x60, 0x75, 0x48, 0x05, 0xd6, 0xa3, 0x49,
	0x28, 0xb7, 0xfe, 0xa3, 0x4b, 0x5b, 0xb8, 0x09, 0xc0, 0xc3, 0x05, 0x95, 0x09, 0xa2, 0x4f, 0xa3,
	0x93, 0x41, 0x6e, 0xfd, 0xed, 0xf9, 0x94, 0x7c, 0xb8, 0xc0, 0x47, 0x86, 0x47, 0x90, 0xf4, 0x7c,
	0xd2, 0x3f, 0xc7, 0xa9, 0x91, 0xd6, 0x66, 0xc8, 0x9f, 0x3f, 0x7c, 0xad, 0x4b, 0xc4, 0x78, 0xb5,
	0xa1, 0x18, 0xe8, 0x09, 0x64, 0xc3, 0xbc, 0xc0, 0xc6, 0x8c, 0xfb, 0xf3, 0x03, 0x86, 0x46, 0xe6,
	0x2d, 0x46, 0x88, 0x85, 0xea, 0x90, 0x1b, 0x30, 0xb6, 0x49, 0x83, 0xb4, 0xc2, 0x52, 0x33, 0x70,
	0x04, 0x9a, 0xa2, 0x23, 0x5f, 0x2a, 0x70, 0x21, 0xc5, 0x20, 0xfd, 0xae, 0xeb, 0x58, 0xd6, 0x9e,
	0xde, 0x7f, 0x46, 0x67, 0x85, 0xb0, 0xdf, 0xe5, 0xd4, 0xea, 0x2f, 0x20, 0x49, 0x75, 0x22, 0x17,
	0xb9, 0xd3, 0x7e, 0xd4, 0xee, 0x3c, 0x21, 0x5d, 0x7f, 0x09, 0x72, 0x92, 0xdc, 0x92, 0x7b, 0xb2,
	0xd6, 0x69, 0xb7, 0x9e, 0x8a, 0x02, 0xba, 0x05, 0x37, 0x18, 0xa1, 0xde, 0x96, 0xb4, 0x27, 0xaa,
	0xc2, 0x97, 0x62, 0xd5, 0xd5, 0xa8, 0xa7, 0x4c, 0xa6, 0x49, 0xe2, 0x33, 0x92, 0x24, 0x0a, 0xd4,
	0x67, 0xd4, 0xce, 0xb6, 0x18, 0x6b, 0xe4, 0x01, 0x8c, 0xd0, 0x02, 0x9b, 0x89, 0x4c, 0x4a, 0x4c,
	0x57, 0xff, 0x54, 0x86, 0x12, 0xed, 0x91, 0xe6, 0x2a, 0x52, 0x2b, 0xb4, 0x48, 0x05, 0x0d, 0x8f,
	0x38, 0x55, 0xa4, 0x62, 0xac, 0x3e, 0xdd, 0x87, 0xec, 0x50, 0x77, 0xb1, 0xed, 0x13, 0x93, 0x25,
	0xa6, 0xfa, 0xdc, 0xcc, 0x36, 0x5d, 0x08, 0xd9, 0x33, 0x01, 0xa3, 0x42, 0x84, 0xd2, 0x47, 0xd8,
	0xa5, 0xaf, 0x33, 0x81, 0x95, 0x6f, 0xb1, 0x59, 0x73, 0x71, 0x72, 0xaa, 0xdd, 0x80, 0x41, 0xe5,
	0x9c, 0xe8, 0x4d, 0x80, 0xd1, 0x50, 0xe3, 0x72, 0xd1, 0x51, 0x20, 0x3b, 0x1a, 0x32, 0x6e, 0xb4,
	0x0d, 0x8b, 0x03, 0xc7, 0x30, 0xf7, 0xcd, 0x7e, 0x70, 0x8f, 0xbe, 0x39, 0x08, 0xa6, 0xb6, 0xdc,
	0xfa, 0x1b, 0x11, 0x27, 0x19, 0xf9, 0xa6, 0x55, 0x3b, 0xb4, 0xfa, 0xb5, 0x1e, 0x7f, 0xf2, 0x62,
	0x50, 0x62, 0x54, 0x9a, 0x2c, 0xa2, 0x07, 0x90, 0xe6, 0xed, 0x59, 0xf0, 0x5c, 0x32, 0x6f, 0xfc,
	0x30, 0x44, 0x2e, 0x8d, 0x36, 0xa0, 0x68, 0xe3, 0xe3, 0x68, 0x0b, 0x9e, 0x9d, 0xf2, 0xb0, 0x7c,
	0x1b, 0x1f, 0x5f, 0xdc, 0x7f, 0xe7, 0xed, 0xc9, 0x8a, 0x81, 0x1e, 0x43, 0x61, 0xe8, 0x9a, 0x03,
	0xdd, 0x1d, 0x6b, 0x41, 0x50, 0xc2, 0x55, 0x82, 0x32, 0xcc, 0x61, 0x01, 0x04, 0x5d, 0x45, 0x1b,
	0x10, 0x74, 0xbc, 0xd8, 0x2b, 0xe7, 0xa8, 0x8e, 0x57, 0x03, 0xe3, 0xc2, 0xa8, 0x01, 0x05, 0xaa,
	0x62, 0xd8, 0x6a, 0xe7, 0xa9, 0x86, 0x4b, 0x4c, 0xc3, 0x1c, 0xd1, 0xf0, 0x82, 0x76, 0x3b, 0x67,
	0x87, 0x74, 0x03, 0x6d, 0x02, 0x84, 0x4f, 0x8d, 0xa4, 0x7c, 0x5c, 0x56, 0x9d, 0xb7, 0x39, 0xe3,
	0xe4, 0x48, 0x6a, 0x44, 0x1a, 0x6d, 0x41, 0x96, 0x07, 0x67, 0x50, 0x37, 0x72, 0x33, 0x5f, 0x24,
	0xce, 0xa7, 0x0a, 0xee, 0x5c, 0x21, 0x02, 0x6a, 0x43, 0xd2, 0xc2, 0xba, 0x87, 0x59, 0xf1, 0xf8,
	0x78, 0x06, 0xd4, 0x99, 0xf0, 0xaa, 0x75, 0xfb, 0x87, 0x78, 0xa0, 0x37, 0x0f, 0x49, 0x23, 0xda,
	0x22, 0xf2, 0x6a, 0x00, 0x83, 0xda, 0x20, 0x52, 0x73, 0x45, 0xb3, 0x8e, 0x48, 0x2d, 0xf6, 0x16,
	0xb3, 0x58, 0x91, 0x58, 0x6c, 0x66, 0xe6, 0xa1, 0xfe, 0xb4, 0x35, 0xc9, 0x3e, 0x3f, 0x87, 0xe2,
	0xbe, 0xe3, 0x0e, 0x74, 0x3f, 0x8c, 0x92, 0xc5, 0x49, 0x7b, 0xf9, 0xfd, 0xc9, 0x72, 0x61, 0x83,
	0xae, 0xf2, 0xc8, 0x2a, 0xec, 0x47, 0x3f, 0xd1, 0x43, 0x9e, 0xa4, 0xaf, 0xd1, 0x9c, 0xfa, 0xfe,
	0xbc, 0xda, 0x9d, 0xcf, 0xd0, 0x6d, 0x48, 0xf5, 0x0f, 0x71, 0xff, 0x99, 0x57, 0xbe, 0x4e, 0x6d,
	0xfe, 0x93, 0x39, 0xa1, 0x9a, 0x44, 0x68, 0xf2, 0x34, 0xa4, 0x32, 0x14, 0xf4, 0x05, 0x14, 0x0d,
	0x42, 0x31, 0xed, 0x03, 0xd6, 0xbe, 0xde, 0xa0, 0xb8, 0x6b, 0x73, 0xe2, 0x92, 0xd6, 0x56, 0xb1,
	0xf7, 0x1d, 0xde, 0xb9, 0x70, 0xb0, 0xa0, 0xe5, 0xed, 0x40, 0x66, 0x9f, 0x8c, 0xe2, 0x26, 0xf6,
	0xca, 0x37, 0x29, 0xee, 0xe5, 0x2f, 0xb8, 0x67, 0xa7, 0x7f, 0x9e, 0xe2, 0x39, 0x48, 0x18, 0xe8,
	0x94, 0x30, 0x26, 0x97, 0xfa, 0xda, 0xf9, 0x40, 0xe7, 0xd3, 0xff, 0xd4, 0x4b, 0x00, 0x0d, 0x74,
	0xf6, 0x65, 0x90, 0x84, 0x77, 0x64, 0xe2, 0xaf, 0xb4, 0xe7, 0x23, 0xec, 0x8e, 0xcb, 0xe5, 0x48,
	0x72, 0xce, 0x12, 0xfa, 0x63, 0x42, 0x46, 0x1f, 0x42, 0xd6, 0xc0, 0x43, 0x6c, 0x1b, 0x5e, 0xc7,
	0x2e, 0xdf, 0xa2, 0xad, 0xd1, 0x35, 0xd2, 0xaf, 0x4b, 0x9c, 0xc8, 0x92, 0xef, 0x84, 0x0b, 0x7d,
	0x09, 0xf9, 0xe0, 0x03, 0x1b, 0x1d, 0xbb, 0x31, 0x2e, 0x57, 0xa8, 0xd2, 0xf7, 0xe6, 0x34, 0xe6,
	0xa4, 0x0f, 0xbc, 0xce, 0xf5, 0x91, 0x22, 0x68, 0xea, 0x14, 0x36, 0xfa, 0x02, 0xf2, 0xdc, 0xbb,
	0x37, 0x9d, 0x3d, 0xaf, 0xfc, 0xfa, 0xa5, 0x13, 0xec, 0xd9, 0xbd, 0xb6, 0x26, 0xa2, 0x3c, 0x6f,
	0x45, 0xd1, 0xd0, 0x67, 0x50, 0x08, 0x9f, 0x7d, 0x9c, 0xa1, 0xef, 0x95, 0x6f, 0xd3, 0xc0, 0xbc,
	0x3f, 0xaf, 0xeb, 0x32, 0xd9, 0xce, 0xd0, 0xf7, 0xd4, 0xbc, 0x17, 0xf9, 0x42, 0x77, 0x20, 0x6b,
	0xb8, 0xce, 0x30, 0xa8, 0x1f, 0x6f, 0xac, 0x08, 0xab, 0x71, 0x7e, 0xcd, 0x84, 0x4c, 0x0b, 0x83,
	0x06, 0x45, 0x17, 0x0f, 0x2d, 0xbd, 0x8f, 0x07, 0xa4, 0xfc, 0x39, 0xfb, 0xe5, 0x25, 0xba, 0xfb,
	0xfa, 0xdc, 0x86, 0x0c, 0x85, 0xb9, 0x63, 0x46, 0xf0, 0x3a, 0xfb, 0x68, 0x07, 0x40, 0x1f, 0x19,
	0xa6, 0xaf, 0x0d, 0x1c, 0x03, 0x97, 0x97, 0x69, 0x54, 0xce, 0x7b, 0x4b, 0x75, 0x22, 0xb8, 0xe5,
	0x18, 0x38, 0x7c, 0x49, 0xe1, 0x84, 0xca, 0x77, 0x02, 0x2c, 0x9e, 0x4b, 0x49, 0xe8, 0x57, 0x90,
	0xb6, 0x1d, 0x23, 0xf2, 0xa2, 0x22, 0xb3, 0xdb, 0x4d, 0xb5, 0x1d, 0x23, 0x78, 0x50, 0xb9, 0x7f,
	0x60, 0xfa, 0x87, 0xa3, 0xbd, 0x5a, 0xdf, 0x19, 0xac, 0x85, 0xa7, 0x30, 0xf6, 0x26, 0xbf, 0xd7,
	0x86, 0xcf, 0x0e, 0xd6, 0xe8, 0xaf, 0xe1, 0x5e, 0x2d, 0x10, 0x53, 0x53, 0x04, 0x55, 0x31, 0xd0,
	0x07, 0x50, 0xc2, 0xc7, 0x43, 0xd3, 0x8d, 0x94, 0xe5, 0x58, 0xc4, 0xac, 0xc5, 0xc9, 0x22, 0x31,
	0x6e, 0xe5, 0x2f, 0x02, 0x94, 0xce, 0xa4, 0x03, 0xd2, 0xa6, 0xd0, 0xd7, 0xba, 0xa9, 0x36, 0x85,
	0x50, 0xc2, 0x06, 0x26, 0x76, 0xe9, 0x93, 0x74, 0xfc, 0x55, 0x9f, 0xa4, 0xa7, 0x87, 0xe3, 0xe4,
	0xfc, 0xc3, 0xf1, 0x66, 0x22, 0x93, 0x10, 0x93, 0x95, 0xa7, 0x90, 0xe1, 0xa9, 0x68, 0xba, 0x6f,
	0x12, 0xe6, 0xec, 0x9b, 0x66, 0xea, 0x59, 0xf9, 0x46, 0x80, 0x6c, 0xf4, 0xad, 0x3f, 0x16, 0xa2,
	0x5e, 0xdc, 0xb6, 0xbd, 0xe4, 0x7b, 0xd8, 0xb4, 0x05, 0xe2, 0xf3, 0x5b, 0xa0, 0x72, 0x04, 0xb9,
	0x48, 0x34, 0x9f, 0xed, 0xb5, 0x85, 0x97, 0xe8, 0xb5, 0xdf, 0x82, 0xd4, 0x97, 0xce, 0x1e, 0x57,
	0x20, 0xde, 0x28, 0x30, 0xe9, 0xe4, 0xa6, 0xb3, 0xa7, 0x48, 0x6a, 0xf2, 0x4b, 0x67, 0x4f, 0x31,
	0x2a, 0xbf, 0x11, 0x20, 0x1f, 0x8d, 0x73, 0x54, 0x85, 0xac, 0x69, 0xf7, 0x5d, 0x1a, 0x64, 0x74,
	0x5f, 0xee, 0x82, 0x13, 0x32, 0x89, 0xfe, 0x81, 0x69, 0x6b, 0xf4, 0x8d, 0x6a, 0xca, 0x4d, 0x33,
	0x03, 0xd3, 0xde, 0x25, 0x54, 0xca, 0xa2, 0x1f, 0x33, 0x96, 0xf8, 0x14, 0x8b, 0x7e, 0x1c, 0xb0,
	0x54, 0x68, 0x41, 0x75, 0x7d, 0xda, 0x16, 0xc7, 0x23, 0x25, 0xd2, 0xf5, 0x2b, 0x87, 0x90, 0x8b,
	0xc4, 0xff, 0x1c, 0x17, 0xf6, 0x53, 0x48, 0x84, 0x41, 0x33, 0x67, 0x2f, 0x4b, 0x05, 0xaa, 0x6f,
	0xf3, 0x81, 0x03, 0x20, 0xb5, 0xbd, 0xd3, 0x68, 0x29, 0xcd, 0x0b, 0x87, 0x05, 0x32, 0x56, 0x84,
	0x49, 0x83, 0x0c, 0x96, 0x92, 0xd2, 0xad, 0x37, 0x5a, 0x32, 0x19, 0x33, 0x0b, 0x90, 0x55, 0xe5,
	0xba, 0x44, 0xa7, 0x10, 0x51, 0xf8, 0x24, 0xf1, 0xf5, 0xb7, 0xcb, 0xc2, 0x66, 0x22, 0x83, 0xc4,
	0x6b, 0xd5, 0xef, 0x04, 0x40, 0x92, 0xee, 0xeb, 0x24, 0x84, 0xae, 0x30, 0x51, 0xc4, 0x2e, 0xd1,
	0x74, 0xba, 0x01, 0x8c, 0xbf, 0x4a, 0x03, 0x18, 0x1c, 0xb5, 0xfa, 0x8d, 0x00, 0x10, 0x39, 0xdc,
	0xa7, 0xd1, 0x7f, 0xc2, 0x66, 0xf7, 0xba, 0x67, 0x52, 0x2a, 0x99, 0x66, 0x83, 0xff, 0xc9, 0x1e,
	0x40, 0xc6, 0x60, 0x2a, 0xb3, 0xeb, 0x98, 0xd9, 0x54, 0x9e, 0xb3, 0xcc, 0x43, 0x52, 0x41, 0x18,
	0xb5, 0x91, 0x86, 0xe4, 0xc8, 0x36, 0x1d, 0xfb, 0xdd, 0x8f, 0x00, 0x9d, 0x4f, 0x3f, 0xc4, 0xec,
	0xf4, 0xb7, 0xee, 0x63, 0x23, 0x98, 0x11, 0x77, 0xec, 0xa3, 0x90, 0x20, 0x34, 0xee, 0xbc, 0xf8,
	0xc7, 0xd2, 0xc2, 0x8b, 0xd3, 0x25, 0xe1, 0xcf, 0xa7, 0x4b, 0xc2, 0x5f, 0x4f, 0x97, 0x84, 0xbf,
	0x9f, 0x2e, 0x09, 0xbf, 0xfe, 0xe7, 0xd2, 0xc2, 0xe7, 0x69, 0x76, 0x80, 0xff, 0x07, 0x00, 0x00,
	0xff, 0xff, 0x0b, 0xf2, 0xfd, 0xfe, 0x11, 0x20, 0x00, 0x00,
}