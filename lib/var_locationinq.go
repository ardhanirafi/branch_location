package lib

import (
	"github.com/gomodul/envy"
)

var (
	dbport  = envy.Get("mysql_port")
	dbaddr  = envy.Get("mysql_addr")
	dbtable = envy.Get("mysql_db")
	dbusn   = envy.Get("mysql_usn")
	// Graylog_Host = os.Getenv("Graylog_Host")
	// Graylog_Port_Traffic = os.Getenv("Graylog_Port_Traffic")
	// Graylog_Port_Error = os.Getenv("Graylog_Port_Error")
)
