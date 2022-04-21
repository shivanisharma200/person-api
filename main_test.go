package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
)

func TestMain(m *testing.M) {
	g := gofr.New()

	db := g.DB()
	if db == nil {
		return 
	}

	query := `
	CREATE TABLE IF NOT EXISTS person (
		id serial PRIMARY KEY,
		name varchar(50) NOT NULL,
		age real NOT NULL,
		address varchar(70) NOT NULL)
	`
	if g.Config.Get("DB_DIALECT") == "postgres" {
		query = `
		IF NOT EXISTS
		(SELECT [name] FROM sys.tables WHERE [name] = 'person'
		) CREATE TABLE person (id serial PRIMARY KEY,
			name varchar(50) NOT NULL,
			age real NOT NULL,
			address varchar(70) NOT NULL)
			`
	}

	if _, err := db.Exec(query); err != nil {
		g.Logger.Errorf("got error sourcing the schema: ", err)
	}

	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	go main() 
	time.Sleep(time.Second * 5)
	req, _ := request.NewMock(http.MethodGet, "http://localhost:9000/persons", nil)
	c := http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		t.Errorf("Test failed.\tHTTP request encountered Err: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Failed.\tExpected %v\tGot %v\n", http.StatusOK, resp.StatusCode)
	}

	_ = resp.Body.Close()
}
