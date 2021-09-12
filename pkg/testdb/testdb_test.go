package testdb

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDBPath = "db_test.json"

func ResetTestDB() {
	DeleteDBFile(testDBPath)
	CreateEmptyDBFile(testDBPath)
}
func TestDBNotFound(t *testing.T) {
	DeleteDBFile(testDBPath)
	db := TestDB{DBPath: testDBPath}
	err := db.New()
	assert.NoFileExists(t, testDBPath)
	assert.EqualError(t, err, "Database File is not found")
}
func TestDBStats(t *testing.T) {
	ResetTestDB()
	db := TestDB{DBPath: testDBPath}
	_ = db.New()
	assert.Equal(t, db.Stats(), "connected")
	db.Close()
	assert.Equal(t, db.Stats(), "disconnected")
}
func TestInvalidKeyValueFormat(t *testing.T) {
	ResetTestDB()

	// init the DB connection
	db := TestDB{DBPath: testDBPath}
	_ = db.New()

	// test invalid input, should return error
	err := db.Put([]byte("test"), []byte("invalid value"))
	assert.Error(t, err)

	// test valid input, should be no error
	validFormatBytes, _ := json.Marshal(map[string]string{"test": "123456"})
	err = db.Put([]byte("test"), validFormatBytes)
	assert.NoError(t, err)
}
func TestPutRecord(t *testing.T) {
	ResetTestDB()

	// inserting new record
	db := TestDB{DBPath: testDBPath}
	_ = db.New()
	newValueByte, _ := json.Marshal(map[string]string{"test": "123456"})
	_ = db.Put([]byte("test"), newValueByte)
	db.Close()

	// reconnect to the DB and test again
	validDB := TestDB{DBPath: testDBPath}
	_ = validDB.New()
	validData, _ := validDB.Get([]byte("test"))
	assert.Equal(t, string(newValueByte), string(validData))
}

func TestDBRecordRemove(t *testing.T) {
	ResetTestDB()

	// init database
	db := TestDB{DBPath: testDBPath}
	_ = db.New()
	newValueByte, _ := json.Marshal(map[string]string{"name": "testName"})
	_ = db.Put([]byte("test"), newValueByte)

	// validation of insertion
	validDB := TestDB{DBPath: testDBPath}
	_ = validDB.New()
	validData, err := validDB.Get([]byte("test"))
	assert.Equal(t, validData, newValueByte)

	// remove the record
	db.Delete([]byte("test"))

	// validation of deletion
	_ = validDB.New()
	_, err = validDB.Get([]byte("test"))
	assert.EqualError(t, err, "Key Not Found in DB")
}
