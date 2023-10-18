package Node

import (
	"IOT_Storage/src/Database"
	"IOT_Storage/src/User"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func ParseDateLocal(date string) string {
	date = strings.ReplaceAll(date, "T", " ")
	date = date + ":00"
	return date
}

func NodeIndexPageForUser(rg *gin.RouterGroup) {
	router := rg.Group("/index")
	router.Static("/assets", "./resources/webapp/assets")
	router.GET("/", func(context *gin.Context) {
		context.Cookie("identity")
		context.HTML(http.StatusOK, "DoctorIndex.html", gin.H{})
	})
}

func NodeUploadPageForUser(rg *gin.RouterGroup) {
	router := rg.Group("/upload")
	router.Static("/assets", "./resources/webapp/assets")
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "UploadIndex.html", gin.H{})
	})
}

func NodeSearchPageForUser(rg *gin.RouterGroup) {
	router := rg.Group("/search")
	router.Static("/assets", "./resources/webapp/assets")
	router.GET("/", func(context *gin.Context) {
		variety := context.Query("variety")
		if variety == "Identity" {
			context.HTML(http.StatusOK, "DoctorSearchByIdentity.html", gin.H{})
		} else {
			context.HTML(http.StatusOK, "DoctorSearchByKeywords.html", gin.H{})
		}

	})
}

func NodeLoginPage(rg *gin.RouterGroup, db *sql.DB) {
	router := rg.Group("/login")
	router.Static("/assets", "./resources/webapp/assets")
	router.GET("/", func(context *gin.Context) {
		status, err := context.Cookie("status")
		if err != nil {
			log.Println(err)
			log.Println(status)
		}
		alert := ""
		if status == "false" {
			alert = "用户名或密码错误 :("
		}
		context.HTML(http.StatusOK, "login.html", gin.H{
			"response": alert,
		})
	})
	router.POST("/", func(context *gin.Context) {
		username := context.PostForm("signin-email")
		password := context.PostForm("signin-password")
		identity := context.PostForm("identity")
		log.Println(identity)

		user := User.Doctor{
			Name:     username,
			PassWord: password,
		}
		isPasswordRight := Database.VerifyPassword(db, &user)
		if isPasswordRight == false {
			context.SetCookie("status", "false", 10, "/", context.Request.URL.Hostname(), false, true)
			location := url.URL{Path: "/login"}
			context.Redirect(http.StatusFound, location.RequestURI())
		} else {
			context.SetCookie("identity", identity, 10, "/", context.Request.URL.Hostname(), false, true)
			location := url.URL{Path: "/index"}
			context.Redirect(http.StatusFound, location.RequestURI())
		}
	})
}

func NodeSearchServerForUser(rg *gin.RouterGroup) {
	router := rg.Group("/DoctorSearch")
	router.Static("/assets", "./resources/webapp/assets")
	router.POST("/", func(context *gin.Context) {
		variety := context.Query("variety")
		log.Println(variety)
		if variety == "identity" {
			//idnumber := context.PostForm("idnumber")
			startTime := context.PostForm("starttime")
			startTime = ParseDateLocal(startTime)
			endTime := context.PostForm("endtime")
			endTime = ParseDateLocal(endTime)
			portForSendSlice := 9000
			nodeToQuery := "http://192.168.42.129:8000"
			patient := User.QueryData(nodeToQuery, startTime, endTime, portForSendSlice)
			context.JSONP(200, patient)
		} else {
			faculty := context.PostForm("faculties")
			features := context.PostFormArray("features")
			log.Println(faculty)
			log.Println(features)
			context.String(200, "OK")
		}
	})
}
