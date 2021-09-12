package testdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/*
 testdb interface
 - assume a key/value DB
*/
type BackendDatabase interface {
	New() error
	Put([]byte, []byte) error
	Get([]byte) ([]byte, error)
	Delete([]byte) error
	Close()
	Stats() string
	Flush() error
}
type TestDB struct {
	DBPath string
	dbFile *os.File
	data   map[string]interface{}
}

func (db *TestDB) New() error {
	// establish the connection with the alpha DB
	// if establish failed => return connection error
	if _, err := os.Stat(db.DBPath); os.IsNotExist(err) {
		// DB path is not exist
		return fmt.Errorf("Database File is not found")
	}
	jsonFile, err := os.Open(db.DBPath)
	if err != nil {
		return err
	}

	byteData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteData), &result)
	if err != nil {
		return err
	}

	db.dbFile = jsonFile
	db.data = result

	fmt.Println("DB Connection Established", db.data)
	return nil
}

func (db *TestDB) Put(key []byte, value []byte) error {
	if db.Stats() == "disconnected" {
		return fmt.Errorf("Database is not connected")
	}
	// if key is exist already => update it
	var result map[string]interface{}
	err := json.Unmarshal(value, &result)
	if err != nil {
		return err
	}
	db.data[string(key)] = result
	err = db.Flush()
	return err
}
func (db *TestDB) Get(key []byte) ([]byte, error) {
	// if key is not exist return error and empty data
	val, ok := db.data[string(key)]
	if !ok {
		return nil, fmt.Errorf("Key Not Found in DB")
	}
	dataByte, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}
func (db *TestDB) Delete(key []byte) error {

	_, ok := db.data[string(key)]
	if !ok {
		// if not exist, just ignore it
		return nil
	}
	delete(db.data, string(key))
	err := db.Flush()
	return err
}
func (db *TestDB) Close() {
	// close the connection
	db.dbFile.Close()
}
func (db *TestDB) Stats() string {
	_, err := db.dbFile.Stat()
	if err != nil {
		return "disconnected"
	}
	return "connected"
}
func (db *TestDB) Flush() error {
	jsonData, _ := json.MarshalIndent(db.data, "", "\t")
	_ = ioutil.WriteFile(db.DBPath, jsonData, 0644)
	fmt.Println("Update Successfully", db.data)
	return nil
}

func CreateEmptyDBFile(path string) error {
	err := ioutil.WriteFile(path, []byte("{}"), 0644)
	return err
}

func DeleteDBFile(path string) error {
	err := os.Remove(path)
	return err
}
