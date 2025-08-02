package sqlconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// _ infont of import for indirect use indirect

func ConnectDB(dbName string) (*sql.DB, error){

	connectString := `root:(YOUR_PASSWORD)@tcp(127.0.0.1:3306)/` + dbName

	db, err := sql.Open("mysql", connectString)
	if err != nil {
		return nil, err
	}

	fmt.Println("connected to MariaDB")

	return db, nil

}
