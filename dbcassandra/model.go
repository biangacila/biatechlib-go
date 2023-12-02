package dbcassandra

type SystemSchemaColumns struct {
	Column_name      string
	Keyspace_name    string
	Table_name       string
	Type             string
	Kind             string
	Position         int
	Clustering_order string
}
type SystemSchemaKeyspaces struct {
	Keyspace_name  string
	Replication    map[string]string
	Durable_writes bool
}
