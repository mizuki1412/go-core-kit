package sqlconst

const (
	Postgres  = "postgres"
	Kingbase  = "kingbase" // todo
	Mysql     = "mysql"
	SqlServer = "mssql"
	Oracle    = "oracle"
	Sqlite3   = "sqlite3"
	DM        = "dm"      // 达梦
	TaosSql   = "taosSql" // tdengine
	TaosWS    = "taosWS"  // tdengine
)

func IsPostgresType(t string) bool {
	return t == Postgres || t == Kingbase
}

// IsSingleDBSchema 单数据库模式
func IsSingleDBSchema(t string) bool {
	return t == DM || t == Oracle ||
		t == Mysql || t == TaosSql || t == TaosWS
}

func IsTaos(t string) bool {
	return t == TaosSql || t == TaosWS
}
