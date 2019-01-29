package config

/*
	For Defined your Database Connections Groups
	Created and Owned by Raditya Pratama

	Copyright 2018 - 17 Oktober 2018 test
*/

//DBGroups as constants connections groups
// PLEASE CHANGE VARIABLE NAME to DBGroups WHEN FILENAME IS CHANGED
var DBGroups1 = map[string]map[string]string{
	"default": map[string]string{
		"Driver":      "mysql", // change driver to the correct db driver
		"Host":        "mysql_host",
		"Port":        "3306", //default port for mysql
		"Protocol":    "tcp",
		"Username":    "mysql_user",
		"Password":    "mysql_pass",
		"ServiceName": "mysql_db", //as DB Name / Service
	},
	"oracleConn": map[string]string{
		"Driver":      "goracle", // change driver to the correct db driver
		"Host":        "oracle_host",
		"Port":        "1521", //default port for oracle
		"Protocol":    "tcp",
		"Username":    "oracle_user",
		"Password":    "oracle_pass",
		"ServiceName": "oracle_sn",
	},
}
