package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const port = ":8081"
const version = "1.0.0"
const db_name = "projects.db"

type Response struct {
	Message string
}

type Project struct {
	gorm.Model
	Title       string
	Description string
	URL         string
}

func printRequest(method string, action string) {
	fmt.Printf("#%s %s\n", method, action)
}

func allProjects(w http.ResponseWriter, r *http.Request) {
	printRequest(r.Method, r.RequestURI)
	db, err := gorm.Open("sqlite3", db_name)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var projects []Project
	db.Find(&projects)

	json.NewEncoder(w).Encode(projects)
}

func newProject(w http.ResponseWriter, r *http.Request) {
	printRequest(r.Method, r.RequestURI)

	db, err := gorm.Open("sqlite3", db_name)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	reqBody, _ := ioutil.ReadAll(r.Body)
	var project Project
	json.Unmarshal(reqBody, &project)

	db.Create(&project)
	json.NewEncoder(w).Encode(&Response{Message: "Created new Project"})
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	printRequest(r.Method, r.RequestURI)
	db, err := gorm.Open("sqlite3", db_name)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	reqBody, _ := ioutil.ReadAll(r.Body)
	var project Project
	json.Unmarshal(reqBody, &project)

	var bm Project
	db.Where("ID = ?", project.ID).Find(&bm)
	db.Delete(&bm)

	json.NewEncoder(w).Encode(&Response{Message: "Successfully deleted Project"})
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	printRequest(r.Method, r.RequestURI)
	db, err := gorm.Open("sqlite3", db_name)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	reqBody, _ := ioutil.ReadAll(r.Body)
	var project Project
	json.Unmarshal(reqBody, &project)

	var bm Project
	db.Where("ID = ?", project.ID).Find(&bm)

	bm.Title = project.Title
	bm.URL = project.URL

	db.Save(&bm)
	json.NewEncoder(w).Encode(&Response{Message: "Successfully updated Project"})
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&Response{Message: fmt.Sprintf("API Version %s", version)})
}

func EneaHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(&Response{Message: "Congrats! You know my Name :)"})
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", db_name)
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Project{})
}

func handleRequests() {
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE", "ENEA"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:9000"})

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", EneaHandler).Methods("ENEA")
	myRouter.HandleFunc("/", RootHandler)
	myRouter.HandleFunc("/projects", allProjects).Methods("GET")
	myRouter.HandleFunc("/project", newProject).Methods("POST")
	myRouter.HandleFunc("/project", deleteProject).Methods("DELETE")
	myRouter.HandleFunc("/project", updateProject).Methods("PUT")

	log.Fatal(http.ListenAndServe(port, handlers.CORS(credentials, methods, origins)(myRouter)))
}

func main() {
	fmt.Printf("@ Started API Server http://localhost%s\n", port)

	initialMigration()
	handleRequests()
}
