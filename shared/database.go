package shared

type database struct {
	_db ticketdb.Database
}

type DBConfig struct {
	DataDir string
	Name    string
}
