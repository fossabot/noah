package datums

type Table struct {
	TableName string
	IsGlobal bool
	IsTenantTable bool
	IdentityColumn string
}
