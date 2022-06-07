package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// URLs to test for
// GET /cell/a/2/string
// GET /cell/b/2/number
// GET /range/b/2/b/4/number
// GET /range/a/1/b/4/object

type test struct {
	input    CellQuery
	expected interface{}
}

func TestCellHandler(t *testing.T) {

	router := setupRouter(false, false)

	tests := []test{
		{
			CellQuery{
				Url:        "https://sheets.googleapis.com/v4/spreadsheets/12rb-uVQsjMWk1GLv6MWXoHHPTY-erUaKVxY4JrdFdUE/values/Sheet1?alt=json&key=AIzaSyDUx0ASGdWyfIbcjT6REUfnwgTg5LbgMY0",
				Column:     "b",
				Row:        2,
				ReturnType: "number",
			},
			24.0,
		},
		{
			CellQuery{
				Url:        "https://sheets.googleapis.com/v4/spreadsheets/12rb-uVQsjMWk1GLv6MWXoHHPTY-erUaKVxY4JrdFdUE/values/Sheet1?alt=json&key=AIzaSyDUx0ASGdWyfIbcjT6REUfnwgTg5LbgMY0",
				Column:     "a",
				Row:        1,
				ReturnType: "string",
			},
			"Person",
		},
	}

	for _, test := range tests {
		qJSON, _ := json.Marshal(test.input)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/cell", bytes.NewReader(qJSON))
		req.Header.Add("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		output := ExternalAdapterOutput{}
		json.NewDecoder(w.Body).Decode(&output)

		assert.NotNil(t, output.Data)
		assert.Equal(t, output.Error, "")
		assert.Equal(t, test.expected, output.Data.Value)

	}

}

func TestCellQueryCellIndex(t *testing.T) {
	q := CellQuery{Column: "C"}
	assert.Equal(t, 2, q.CellIndex())
}
