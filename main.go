package main
import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "net/http"
    "html/template"
    "log"
	"reflect"
	"strings"
)
type Product struct{
    Id int
    Model string
    Company string
    Price int
}

type User struct{
	name string
	password string
	role int
}
var database *sql.DB

func addToDB(model string, company string, price int16) {
	result, err := database.Exec("insert into productdb.Products (model, company, price) values (?, ?, ?)", 
        &model, &company, &price)
    if err != nil{
        panic(err)
    }
	fmt.Println(reflect.TypeOf(result))
}

func updateDB(rowname string, value string, id int16) {
	result, err := database.Exec("update productdb.Products set ? = ? where id = ?",rowname, value, id)
    if err != nil{
        panic(err)
    }
	fmt.Println(reflect.TypeOf(result))
}

func deleteFromDB(id int16) {
	result, err := database.Exec("delete from productdb.Products where id = ?", id)
    if err != nil{
        panic(err)
    }
	fmt.Println(reflect.TypeOf(result))
}

func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //get request method
    if r.Method == "GET" {
        t, _ := template.ParseFiles("templates/login.html")
        t.Execute(w, nil)
    } else {
        r.ParseForm()
        // logic part of log in
		//user, pass := r.Form["username"], r.Form["password"]
		user := strings.Join(r.Form["username"],"")
		pass := strings.Join(r.Form["password"],"")
        fmt.Println(user)
        fmt.Println(pass)
		rows, err := database.Query("select role from productdb.users where name=?", user)
		if err != nil {
			log.Println("There is no such user")
			//log.Println(err)
		} else {
			defer rows.Close()
			u := User{}
			for rows.Next() {
				err  := rows.Scan(&u.role)
				if err != nil {
						log.Println(err)
				}
				log.Println(u)
				log.Println("Красавчик")
				}
		}
    }
}

//table users (name varchar(50), password varchar(50), role int);
 
func indexHandler(w http.ResponseWriter, r *http.Request) {
	
    rows, err := database.Query("select * from productdb.Products")
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()
    products := []Product{}
     
    for rows.Next(){
        p := Product{}
        err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
        if err != nil{
            fmt.Println(err)
            continue
        }
        products = append(products, p)
    }
 
    tmpl, _ := template.ParseFiles("templates/base.html")
    tmpl.Execute(w, products)
}
 
func main() {

	http.HandleFunc("/login", login)
	
	db, err := sql.Open("mysql", "root:root@/productdb")
    	
    if err != nil {
        log.Println(err)
    }
    database = db
    defer db.Close()
	//addToDB("Xiaomi", "Pocophone", 19000)
	//deleteFromDB(7)
	//deleteFromDB(8)
	
    http.HandleFunc("/base", indexHandler)
 
    fmt.Println("Server is listening...")
    http.ListenAndServe(":8181", nil)
}
