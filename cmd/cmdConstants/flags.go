package cmdConstants

const (
	CliFlagMySQLHost     = "mysql-host"
	CliFlagMySQLPort     = "mysql-port"
	CliFlagMySQLDatabase = "mysql-database"
	CliFlagMySQLUser     = "mysql-user"
	CliFlagMySQLPassword = "mysql-password"
	CliFlagListen        = "listen"
	CliFlagPort          = "port"
)

// CliFlags Return a list of all known CLI flags used by any and all commands.
func CliFlags() []string {
	return []string{
		CliFlagMySQLHost, CliFlagMySQLPort, CliFlagMySQLDatabase, CliFlagMySQLUser, CliFlagMySQLPassword,
		CliFlagListen, CliFlagPort,
	}
}
