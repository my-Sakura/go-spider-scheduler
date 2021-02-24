package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	var err error
	dsn := "root:password@tcp(192.168.0.253:3307)/mysql?charset=utf8mb4"
	db, err = sql.Open("mysql", dsn)
	if err != nil {

		panic("db connect error")
	}
	CreateTableAcademician()
}

const (
	academicianTableCreate = iota
	academicianTableInsert
	academicianTableDelete
	academicianTableUpdate
	academicianTableGetAll
)

var (
	db               *sql.DB
	errInvalidInsert = errors.New("errInvalidInsert")
	errInvalidDelete = errors.New("errInvalidDelete")
	errInvalidUpDate = errors.New("errInvalidUpDate")
)

var academicianSQLString = []string{
	`CREATE TABLE IF NOT EXISTS academician(
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(10) NOT NULL,
		content VARCHAR(300) NOT NULL
		)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	`INSERT INTO academician (name, content) VALUES(?, ?)`,
	`DELETE FROM academician WHERE id=?`,
	`UPDATE academician SET content=? WHERE id=?`,
	`SELECT name, content FROM academician WHERE id=?`,
}

func CreateTableAcademician() {
	_, err := db.Exec(academicianSQLString[academicianTableCreate])
	if err != nil {
		fmt.Println(err)
		panic("academicianTable create fail")
	}
}

func InsertTableAcademician(name, content string) {
	result, err := db.Exec(academicianSQLString[academicianTableInsert], name, content)
	if err != nil {
		errInsert := fmt.Sprintf("stateDepartmentTable Insert fail\n name: %s,\ncontent: %s", name, content)
		panic(errInsert)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		panic(errInvalidInsert)
	}
}

func DeleteTableAcademician(id int) {
	result, err := db.Exec(academicianSQLString[academicianTableDelete], id)
	if err != nil {
		errDelete := fmt.Sprintf("stateDepartmentTable Delete fail\n id: %d", id)
		panic(errDelete)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		panic(errInvalidDelete)
	}
}

func UpdateTableAcademician(id int) {
	result, err := db.Exec(academicianSQLString[academicianTableUpdate], id)
	if err != nil {
		errDelete := fmt.Sprintf("stateDepartmentTable Update fail\n id: %d", id)
		panic(errDelete)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		panic(errInvalidUpDate)
	}
}

func GetTableAllAcademician(id int) {
	_, err := db.Exec(academicianSQLString[academicianTableGetAll], id)
	if err != nil {
		errGet := fmt.Sprintf("stateDepartmentTable Get fail\n id: %d", id)
		panic(errGet)
	}
}
