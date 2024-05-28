package service

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var PASSWORD = os.Getenv("MYSQL_ROOT_PASSWORD")

func Connect() *sql.DB {
	dsn := "root:password@tcp(db:3306)/sms_auth_db"
	//	dsn := "root:${PASSWORD}@tcp(db:3306)/sms_auth_db"

	fmt.Println(dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func ClearToken(userId string) {
	db := Connect()
	defer db.Close()

	sql := `
	delete 
	from sms_auth 
	where user_id = ? 
	`

	_, err := db.Exec(sql, userId)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
}

func RegisterToken(userId string, phoneNo string, token string) {
	ClearToken(userId)

	db := Connect()
	defer db.Close()

	sql := `
	insert 
	into sms_auth (
		user_id,
		phone_no,
		one_time_password,
		expired_at
	) values (
		?,
		?,
		?,
		date_add(now(), interval 20 minute)
	)
	`

	_, err := db.Exec(sql, userId, phoneNo, token)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
}

func VerifyToken(userId string, token string) bool {
	db := Connect()
	defer db.Close()

	sql := `
	select * 
	from sms_auth 
	where user_id = ? 
	and one_time_password = ?
	and expired_at >= now()
	`

	rows, err := db.Query(sql, userId, token)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

	return rows.Next()
}

func ConnectionTest() {
	db := Connect()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Connected to MySQL database!")
}
