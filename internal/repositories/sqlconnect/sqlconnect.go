package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// _ infont of import for indirect use indirect

func ConnectDB(dbName string) (*sql.DB, error) {

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("HOST")
	db_port := os.Getenv("DB_PORT")
							//  (USER):(YOUR_PASSWORD)@tcp((HOST):(DB PORT))/` + dbName
	connectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",user,pass,host,db_port,dbName)

	db, err := sql.Open("mysql", connectString)
	if err != nil {
		return nil, err
	}

	fmt.Println("connected to MariaDB")

	return db, nil

}
