package telemetrytraces

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
	return "signoz_traces"
}

const (
	SpanIndexV3TableName          = "distributed_signoz_index_v3"
	SpanIndexV3LocalTableName     = "signoz_index_v3"
	TagAttributesV2TableName      = "distributed_tag_attributes_v2"
	TagAttributesV2LocalTableName = "tag_attributes_v2"
	TopLevelOperationsTableName   = "distributed_top_level_operations"
	TraceSummaryTableName         = "distributed_trace_summary"
	SpanAttributesKeysTblName     = "distributed_span_attributes_keys"
)
