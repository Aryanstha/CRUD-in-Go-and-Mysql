package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type name struct {
	id    int
	First string
	Last  string
}

//database connection

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbname := "name"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// display data in index
var tmpl = template.Must(template.ParseGlob("frontend/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM name ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}
	n := name{}
	res := []name{}
	for selDB.Next() {
		var id int
		var first, last string
		err = selDB.Scan(&id, &first, &last)
		if err != nil {
			panic(err.Error())
		}
		n.id = id
		n.First = first
		n.Last = last
		res = append(res, n)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nid := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM name WHERE id=?", nid)
	if err != nil {
		panic(err.Error())
	}
	n := name{}
	for selDB.Next() {
		var id int
		var first, last string
		err = selDB.Scan(&id, &first, &last)
		if err != nil {
			panic(err.Error())
		}
		n.id = id
		n.First = first
		n.Last = last
	}
	tmpl.ExecuteTemplate(w, "Show", n)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}
func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nid := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM name WHERE id=?", nid)
	if err != nil {
		panic(err.Error())
	}
	n := name{}
	for selDB.Next() {
		var id int
		var first, last string
		err = selDB.Scan(&id, &first, &last)
		if err != nil {
			panic(err.Error())
		}
		n.id = id
		n.First = first
		n.Last = last
	}
	tmpl.ExecuteTemplate(w, "Edit", n)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		first := r.FormValue("first")
		last := r.FormValue("last")
		insForm, err := db.Prepare("INSERT INTO name(first, last) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(first, last)
		log.Println("INSERT: First: " + first + " | Last: " + last)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		first := r.FormValue("first")
		last := r.FormValue("last")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE name SET first=?, last=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(first, last, id)
		log.Println("UPDATE: First: " + first + " | Last: " + last)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nid := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM name WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(nid)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
