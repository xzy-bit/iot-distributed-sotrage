package Web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       // 解析参数，默认是不会解析的
	fmt.Println(r.Form) // 这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello World!") // 这个写入到 w 的是输出到客户端的
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method", r.Method)
	if r.Method == "GET" {
		t, err := template.ParseFiles("src/Web/resources/login.html")
		if err != nil {
			log.Fatal(err)
		}
		log.Println(t.Execute(w, nil))
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Fatal("ParseForm ", err)
		}
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
