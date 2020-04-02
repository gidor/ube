package apicfg

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func (d *DataConnection) open() error {
	if d.pool == nil {
		db, err := sql.Open(d.Type, d.ConnectionString)
		if err != nil {
			d.lasterror = err
		}
		d.pool = db
		d.lasterror = nil
		return nil
	}
	return nil
}

func (d *DataConnection) close() {
	if d.pool != nil {
		d.pool.Close()
	}
	d.pool = nil
}

func (d *DataConnection) Init() error {
	return d.open()
}

func (d *DataConnection) Query(stmnt string, args ...interface{}) (*sql.Rows, error) {
	if d.pool == nil {
		d.open()
	}

	if d.pool == nil {
		// log
		return nil, d.lasterror
	}
	return d.pool.Query(stmnt, args...)
}

// type Record map[string]interface{}

func (d *DataConnection) QueryRecord(stmnt string, args ...interface{}) ([]map[string]interface{}, error) {
	if d.pool == nil {
		d.open()
	}

	if d.pool == nil {
		// log
		return nil, d.lasterror
	}
	rows, err := d.pool.Query(stmnt, args...)

	if err != nil {
		return nil, err
	}
	cols, _ := rows.Columns()

	// Create a slice of interface{}'s to represent each column,
	// and a second slice to contain pointers to each item in the columns slice.
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i, _ := range columns {
		columnPointers[i] = &columns[i]
	}

	var result []map[string]interface{} // := make(RecordList)

	for rows.Next() {

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		result = append(result, m)
	}
	return result, nil
}

func (d *DataConnection) Execute(stmnt string, args ...interface{}) (sql.Result, error) {
	if d.pool == nil {
		d.open()
	}

	if d.pool == nil {
		// log
		return nil, d.lasterror
	}
	return d.pool.Exec(stmnt, args)
}

func (d *DataConnection) Ping() error {
	if d.pool == nil {
		d.open()
	}

	if d.pool == nil {
		// log
		return d.lasterror
	}
	return d.pool.Ping()
}
