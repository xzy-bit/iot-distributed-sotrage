package Database

import (
	"IOT_Storage/src/User"
	"bytes"
	"database/sql"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tjfoc/gmsm/sm3"
	"log"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/blockchain")
	if err != nil {
		log.Println(err)
	}
	return db
}

func PasswordToHash(password string) string {
	info := bytes.Join([][]byte{
		[]byte(password),
	}, []byte{})
	h := sm3.New()
	h.Write(info)
	sum := h.Sum(nil)
	code := hex.EncodeToString(sum)
	return code
}

func AddDoctor(db *sql.DB, user *User.Doctor) {
	stmt, err := db.Prepare("INSERT into user(username,password) values (?,?)")
	if err != nil {
		log.Println(err)
		return
	}
	code := PasswordToHash(user.PassWord)
	_, err = stmt.Exec(user.Name, code)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	return
}

func VerifyPassword(db *sql.DB, user *User.Doctor) bool {
	var temp User.Doctor
	inputHash := PasswordToHash(user.PassWord)
	row := db.QueryRow("SELECT *from user where user.username=? and user.password=?", user.Name, inputHash)
	if row.Scan(&temp.Name, &temp.PassWord) != nil {
		return false
	}
	return true
}

func AddIndex(db *sql.DB, indexName string, d1 []byte, d2 []byte) {
	stmt, err := db.Prepare("INSERT into indexes(indexName,p1,p2) values (?,?,?)")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec(indexName, d1, d2)
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()
	return
}
