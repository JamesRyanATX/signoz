package telemetryaudit

var dbName string

func Init(databaseName string) {
	if databaseName != "" {
		dbName = databaseName
	}
}

func DBName() string {
	if dbName != "" {
		return dbName
	}
	return "signoz_audit"
}

const (
	AuditLogsTableName          = "distributed_logs"
	AuditLogsLocalTableName     = "logs"
	TagAttributesTableName      = "distributed_tag_attributes"
	TagAttributesLocalTableName = "tag_attributes"
	LogAttributeKeysTblName     = "distributed_logs_attribute_keys"
	LogResourceKeysTblName      = "distributed_logs_resource_keys"
	LogsResourceTableName       = "distributed_logs_resource"
)
