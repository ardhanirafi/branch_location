package lib

import (
	"os"
)

var (
	dbport = os.Getenv("mysql_port")
	dbaddr = os.Getenv("mysql_addr")
	// Graylog_Host = os.Getenv("Graylog_Host")
	// Graylog_Port_Traffic = os.Getenv("Graylog_Port_Traffic")
	// Graylog_Port_Error = os.Getenv("Graylog_Port_Error")
)
