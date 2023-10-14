package Database

import (
	"IOT_Storage/src/User"
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	db := ConnectDB()
	user := User.Doctor{
		Name:     "JIA",
		PassWord: "xxoo",
	}
	AddDoctor(db, &user)
}

func TestVerifyPassword(t *testing.T) {
	db := ConnectDB()
	userTrue := User.Doctor{
		Name:     "JIA",
		PassWord: "xxoo",
	}
	userFalse := User.Doctor{
		Name:     "JIA",
		PassWord: "ooxx",
	}
	result := VerifyPassword(db, &userTrue)
	fmt.Println(result)
	result = VerifyPassword(db, &userFalse)
	fmt.Println(result)
}
