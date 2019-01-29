package lib

import (
	"database/sql"
	"fmt"
	"time"

	log "test/logging"
	"test/parser"

	//load mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DatabaseConnection : interface for database action
type DatabaseConnection interface {
	InitDb(param ...string) DbConnection
	Open() (*sql.DB, error)
	Close()
	GetRows(rows *sql.Rows) (map[int]map[string]string, error)

	GetFirstRow() (string, error)
	Query(sqlStringName string, args ...interface{}) (*sql.Rows, error)
	Exec(sqlStringName string, args ...interface{}) (int64, error)
	Queryf(sqlStringName string, args ...interface{}) (*sql.Rows, error)
	Execf(sqlStringName string, args ...interface{}) (int64, error)

	// Added by Budianto
	GetRowsbyIndex(rows *sql.Rows) (map[int]map[int]string, error)
	GetFirstData(sqlStringName string, args ...interface{}) (string, error)
	GetFirstRowByQuery(sqlStringName string, key string, args ...interface{}) (string, error)
}

// DbConnection is global config for database connection
type DbConnection struct {
	//DBType string            `yaml:"DBType"`
	//DBURL  string            `yaml:"DBURL"`
	SQL     map[string]string `yaml:"SQLCommand"`
	Db      *sql.DB
	dbTypes string
	dbName string
}

// New is use to create db connection
func New(fn string) (*DbConnection, error) {
	var c DbConnection
	if err := parser.LoadYAML(&fn, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Open function prepares dbConnection for future connection to database
func (c DbConnection) Open(optionalParam ...string) (*sql.DB, string, string, error) {
	c.Close()
	groupName := "default"
	if isset(optionalParam, 0) {
		groupName = optionalParam[0]
	}
	dbDriver, dbName, connString := prepareOpenConnection(groupName)
	// Open database connection
	var err error

	dbConn, err := sql.Open(dbDriver, connString)
	if err != nil {
		return nil, "", "", err
	}

	dbConn.SetMaxOpenConns(100)
	dbConn.SetConnMaxLifetime(time.Minute * 1)
	dbConn.SetMaxIdleConns(100)

	err = dbConn.Ping()
	if err != nil {
		return nil, "", "", err
	}
	log.Logf("Initiating database connection %s", connString)
	return dbConn, dbDriver, dbName, nil
}

// Close function closes existing dbConnection
//
func (c DbConnection) Close() {
	if c.Db != nil {
		log.Debug("Closing previous database connection.")
		c.Db.Close()
		c.Db = nil
	}
}

//GetRows parses recordset into map
func (c DbConnection) GetRows(rows *sql.Rows) (map[int]map[string]string, error) {
	var results map[int]map[string]string
	results = make(map[int]map[string]string)

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	counter := 1
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// initialize the second layer
		results[counter] = make(map[string]string)

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			results[counter][columns[i]] = value
		}
		counter++
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

//GetFirstRow parse and gets column value in first record
func (c DbConnection) GetFirstRow(rows *sql.Rows, key string) (string, error) {
	results, err := c.GetRows(rows)
	if err != nil {
		return "", err
	}
	return results[1][key], nil
}

// Query sends SELECT command to database
func (c DbConnection) Query(sqlStringName string, args ...interface{}) (*sql.Rows, error) {
	// if no dbConnection, return
	//
	if c.Db == nil {
		return nil, fmt.Errorf("database needs to be initiated first")
	}
	check, errCheck := c.CheckDB(true)
	if !check {
		return nil, errCheck
	}

	var strSQL string
	var found bool

	//if strSQL, found = sqlCommandMap[sqlStringName]; !found {
	if strSQL, found = c.SQL[sqlStringName]; !found {
		strSQL = sqlStringName
	}

	rows, err := c.Db.Query(strSQL, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

//Exec executes UPDATE/INSERT/DELETE statements and returns rows affected
func (c DbConnection) Exec(sqlStringName string, args ...interface{}) (int64, error) {
	// if no dbConnection, return
	//

	if c.Db == nil {
		return 0, fmt.Errorf("Please OpenConnection prior Query")
	}
	check, errCheck := c.CheckDB(true)
	if !check {
		return 0, errCheck
	}

	var strSQL string
	var found bool

	//if strSQL, found = sqlCommandMap[sqlStringName]; !found {
	if strSQL, found = c.SQL[sqlStringName]; !found {
		strSQL = sqlStringName
	}

	// Execute the query
	res, err := c.Db.Exec(strSQL, args...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}
	return rows, nil
}

// InsertGetLastID is use for ...
func (c DbConnection) InsertGetLastID(sqlStringName string, args ...interface{}) (int64, error) {
	// if no dbConnection, return
	//
	if c.Db == nil {
		return 0, fmt.Errorf("Please OpenConnection prior Query")
	}
	check, errCheck := c.CheckDB(true)
	if !check {
		return 0, errCheck
	}

	var strSQL string
	var found bool

	//if strSQL, found = sqlCommandMap[sqlStringName]; !found {
	if strSQL, found = c.SQL[sqlStringName]; !found {
		strSQL = sqlStringName
	}

	// Execute the query
	res, err := c.Db.Exec(strSQL, args...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	rows, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

// Queryf is use for ...
func (c DbConnection) Queryf(sql string, a ...interface{}) (*sql.Rows, error) {
	return c.Query(fmt.Sprintf(sql, a...))
}

// Execf is use for ...
func (c DbConnection) Execf(sql string, a ...interface{}) (int64, error) {
	return c.Exec(fmt.Sprintf(sql, a...))
}

//Added by Budianto

//GetRowsbyIndex parses Get row using index of row and column
func (c DbConnection) GetRowsbyIndex(rows *sql.Rows) (map[int]map[int]string, int, error) {
	var results map[int]map[int]string
	results = make(map[int]map[int]string)

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}

	//Get Column name
	values := make([]sql.RawBytes, len(columns))

	//Define dynamic variables base on column name
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	counter := 1 //row count
	for rows.Next() {
		// get RawBytes from data
		// Assign value to dynamic variables
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, 0, err
		}

		// initialize the second layer
		results[counter] = make(map[int]string)

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			results[counter][i] = value
		}
		counter++
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, counter, nil
}

//GetCustomRowColumn parses Get any number of  row or column where 0 = get all
//row start from 1, column index start from 0
//maxRow,maxColumn = row/column count, start from
func (c DbConnection) GetCustomRowColumn(rows *sql.Rows, maxRow int, maxColumn int) (map[int]map[int]string, error) {
	var results map[int]map[int]string
	results = make(map[int]map[int]string)

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	//Get Column name
	values := make([]sql.RawBytes, len(columns))

	//Define dynamic variables base on column name
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	rowCounter := 1
	for rows.Next() {
		// get RawBytes from data
		// Assign value to dynamic variables
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// initialize the second layer
		results[rowCounter] = make(map[int]string)

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for colCounter, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			results[rowCounter][colCounter] = value
			if (maxColumn > 0) && (colCounter+1 >= maxColumn) {
				break
			}
		}
		if (maxRow > 0) && (rowCounter >= maxRow) {
			break
		}
		rowCounter++
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

//GetFirstRowByQuery parse and gets column value in first record by query
func (c DbConnection) GetFirstRowByQuery(sqlStringName string, args ...interface{}) (map[int]map[int]string, error) {
	var rowret *sql.Rows
	var err error
	if c.Db == nil {
		return nil, fmt.Errorf("Please OpenConnection prior Query")
	}
	check, errCheck := c.CheckDB(true)
	if !check {
		return nil, errCheck
	}

	rowret, err = c.Query(sqlStringName, args...)
	if err != nil {
		return nil, err
	}

	results, err := c.GetCustomRowColumn(rowret, 1, 0)
	if err != nil {
		return nil, err
	}
	return results, nil
}

//GetFirstData get result from sql command that return only 1 row and 1 column ony
func (c DbConnection) GetFirstData(sqlStringName string, args ...interface{}) (string, error) {

	if c.Db == nil {
		return "", fmt.Errorf("Please OpenConnection prior Query")
	}
	check, errCheck := c.CheckDB(true)
	if !check {
		return "", errCheck
	}

	rowret, err := c.Query(sqlStringName, args...)
	if err != nil {
		return "", err
	}
	results, err := c.GetCustomRowColumn(rowret, 1, 1)
	if err != nil {
		return "", err
	}
	return results[1][0], nil
}

//SelectQuery parse and gets column value in first record by query
func (c DbConnection) SelectQuery(sqlStringName string, args ...interface{}) (map[int]map[int]string, int, error) {
	var rowret *sql.Rows
	var err error
	if c.Db == nil {
		return nil, 0, fmt.Errorf("Please OpenConnection prior Query")
	}
	check, errCheck := c.CheckDB(true)
	if !check {
		return nil, 0, errCheck
	}

	rowret, err = c.Query(sqlStringName, args...)
	if err != nil {
		return nil, 0, err
	}

	results, rowCount, err := c.GetRowsbyIndex(rowret)
	if err != nil {
		return nil, 0, err
	}
	return results, rowCount, nil
}

// CheckDB  function prepares dbConnection for future connection to database
func (c DbConnection) CheckDB(Reconnect bool) (bool, error) {
	var err error
	err = c.Db.Ping()
	if err != nil {
		log.Errorf("Database connection %s", "Failed")
		return false, fmt.Errorf("Database connection failed")
	}

	return true, err
}
