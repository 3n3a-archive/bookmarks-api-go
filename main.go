package main

import (
	"encoding/json"
    "fmt"
    "log"
    "net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/gorilla/mux"
)

// Our User Struct
type User struct {
    gorm.Model
    Name  string
    Email string
}

type Response struct {
	Message string
}

type Bookmark struct {
	gorm.Model
	Title string
	Description string
	URL string
	Category string
}

func printRequest(method string, action string) {
	fmt.Printf("#%s %s\n", method, action)
}

func allBookmarks(w http.ResponseWriter, r *http.Request) {
	printRequest(r.Method,r.RequestURI)
    db, err := gorm.Open("sqlite3", "test.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    var bookmarks []Bookmark
    db.Find(&bookmarks)

    json.NewEncoder(w).Encode(bookmarks)
}

func newBookmark(w http.ResponseWriter, r *http.Request) {
    printRequest(r.Method,r.RequestURI)

    db, err := gorm.Open("sqlite3", "test.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    vars := mux.Vars(r)
    title := vars["title"]
    desc := vars["description"]
	url := vars["url"]
	cat := vars["category"]

    db.Create(&Bookmark{Title: title, Description: desc, URL: url, Category: cat})
    json.NewEncoder(w).Encode(&Response{Message: "Created new Bookmark"})
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	printRequest(r.Method,r.RequestURI)
    db, err := gorm.Open("sqlite3", "test.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    var users []User
    db.Find(&users)

    json.NewEncoder(w).Encode(users)
}

func newUser(w http.ResponseWriter, r *http.Request) {
    printRequest(r.Method,r.RequestURI)

    db, err := gorm.Open("sqlite3", "test.db")
    if err != nil {
        panic("failed to connect database")
    }
    defer db.Close()

    vars := mux.Vars(r)
    name := vars["name"]
    email := vars["email"]

    db.Create(&User{Name: name, Email: email})
    json.NewEncoder(w).Encode(&Response{Message: "Created new User"})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	printRequest(r.Method,r.RequestURI)
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var user User
	db.Where("name = ?", name).Find(&user)
	db.Delete(&user)

	json.NewEncoder(w).Encode(&Response{Message: "Successfully deleted User"})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	printRequest(r.Method,r.RequestURI)
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	var user User
	db.Where("name = ?", name).Find(&user)

	user.Email = email

	db.Save(&user)
	json.NewEncoder(w).Encode(&Response{Message: "Successfully updated User"})
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
    if err != nil {
		fmt.Println(err.Error())
        panic("failed to connect database")
    }
    defer db.Close()
	
    // Migrate the schema
    db.AutoMigrate(&User{})
    db.AutoMigrate(&Bookmark{})
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/users", allUsers).Methods("GET")
	myRouter.HandleFunc("/user/{name}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{name}/{email}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/user/{name}/{email}", newUser).Methods("POST")
	myRouter.HandleFunc("/bms", allBookmarks).Methods("GET")
	myRouter.HandleFunc("/bm/{title}/{description}/{url}/{category}", newBookmark).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
    fmt.Println("--------------------------")
    fmt.Println("--- Started API Server ---")
    fmt.Println("--------------------------")

	initialMigration()

    handleRequests()
}