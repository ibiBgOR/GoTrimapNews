/*
Copyright 2015 Oscar Ruckdeschel, Janik Schmidt, Jonathan Kuhse.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package data

import (
	"github.com/gchaincl/dotsql"
	//	"github.com/gchaincl/gotic/fs"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

var user = "root"
var passwd = ""

// drop and creates a new database
func createDatabase(user string, passwd string) {
	db, err := sql.Open("mysql", user+passwd+"@/")
	if err != nil {
		panic(err)
	}
	dot, err := dotsql.LoadFromFile("queries.sql")
	if err != nil {
		panic(err)
	}
	// TODO remove for persistent database
	_, err = dot.Exec(db, "drop-database")
	if err != nil {
		panic(err)
	}
	_, err = dot.Exec(db, "create-database")
	if err != nil {
		panic(err)
	}
	_, err = dot.Exec(db, "use-database")
	if err != nil {
		panic(err)
	}
	_, err = dot.Exec(db, "create-titles-table")
	if err != nil {
		panic(err)
	}
	_, err = dot.Exec(db, "create-trigrams-table")
	if err != nil {
		panic(err)
	}
}

// returns all news-ids of a trigram
func getIdsOfTrigram(trigram string) []int {
	db, err := sql.Open("mysql", user+passwd+"@/trigramnews")
	if err != nil {
		panic(err)
	}
	dot, err := dotsql.LoadFromFile("queries.sql")
	if err != nil {
		panic(err)
	}
	rows, err := dot.Query(db, "select-titleids-from-trigram", trigram)
	if err != nil {
		panic(err)
	}
	var ids []int
	for rows.Next() {
		var id_string string
		if err := rows.Scan(&id_string); err != nil {
			panic(err)
		}
		id, err := strconv.Atoi(id_string)
		if err != nil {
			panic(err)
		}
		ids = append(ids, id)
	}
	return ids
}

// returns news-title with this id
func getNewsTitle(id int) string {
	db, err := sql.Open("mysql", user+passwd+"@/trigramnews")
	if err != nil {
		panic(err)
	}
	dot, err := dotsql.LoadFromFile("queries.sql")
	if err != nil {
		panic(err)
	}
	rows, err := dot.Query(db, "select-title", id)
	if err != nil {
		panic(err)
	}
	var title string
	for rows.Next() {
		if err := rows.Scan(&title); err != nil {
			panic(err)
		}
	}
	return title
}

// saves trigram with this id
func putTrigram(trigram string, id int) {
	db, err := sql.Open("mysql", user+passwd+"@/trigramnews")
	if err != nil {
		panic(err)
	}
	dot, err := dotsql.LoadFromFile("queries.sql")
	if err != nil {
		panic(err)
	}
	_, err = dot.Exec(db, "insert-trigram", id, trigram)
	if err != nil {
		panic(err)
	}
}

// saves a new news title, returns id
func postNews(title string) int {
	db, err := sql.Open("mysql", user+passwd+"@/trigramnews")
	if err != nil {
		panic(err)
	}
	dot, err := dotsql.LoadFromFile("queries.sql")
	if err != nil {
		panic(err)
	}
	_, err = dot.Query(db, "insert-title", title)
	if err != nil {
		panic(err)
	}
	rows, err := dot.Query(db, "select-titleid-by-name", title)
	if err != nil {
		panic(err)
	}
	var id_string string
	for rows.Next() {
		if err := rows.Scan(&id_string); err != nil {
			panic(err)
		}
	}
	id, err := strconv.Atoi(id_string)
	if err != nil {
		panic(err)
	}
	return id
}
