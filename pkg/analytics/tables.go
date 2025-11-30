package analytics

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
	return "signoz_analytics"
}

const (
	RuleStateHistoryTableName = "distributed_rule_state_history_v0"
)
