package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/daohuei/go-exercise/pkg/testdb"
)

const (
	CREATE_EMPTY_DB = iota
	DELETE_DB
	RUN
	DB_PATH = "db.json"
)

func main() {
	option := flag.Int("option", RUN, "0: create empty db file\n1: delete db file\n2: run the main function")
	flag.Parse()
	fmt.Println(*option)
	if *option == CREATE_EMPTY_DB {
		err := testdb.CreateEmptyDBFile(DB_PATH)
		if err != nil {
			panic(err)
		}
		fmt.Println("DB file created")
		return
	} else if *option == DELETE_DB {
		_ = testdb.DeleteDBFile(DB_PATH)
		fmt.Println("DB file deleted")
		return
	}

	// RUN
	// init new DB object and connect to it
	currentDB := testdb.TestDB{DBPath: DB_PATH}
	err := currentDB.New()
	if err != nil {
		panic(err)
	}

	// check DB connection status
	fmt.Println(currentDB.Stats())

	// insert a new person
	fmt.Println("Inserting a new record")
	newValueByte, err := json.Marshal(map[string]string{"name": "new_person"})
	if err != nil {
		panic(err)
	}
	newPersonID := "_456"
	err = currentDB.Put([]byte(newPersonID), newValueByte)
	if err != nil {
		panic(err)
	}

	// get the inserted person
	fmt.Println("Getting the inserted record")
	data, err := currentDB.Get([]byte(newPersonID))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// remove the person
	fmt.Println("removing the inserted record")
	err = currentDB.Delete([]byte("_456"))
	if err != nil {
		panic(err)
	}

	// close DB
	currentDB.Close()
}
