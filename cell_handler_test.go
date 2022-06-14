package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type celltest struct {
	input    CellQuery
	expected interface{}
}

var tests = []celltest{
	{
		CellQuery{
			Column:     "b",
			Row:        2,
			ReturnType: "n",
		},
		27.0,
	},
	{
		CellQuery{
			Column:     "a",
			Row:        1,
			ReturnType: "s",
		},
		"Person",
	},
	{
		CellQuery{
			Column:     "c",
			Row:        2,
			ReturnType: "b",
		},
		true,
	},
	{
		CellQuery{
			Column:     "d",
			Row:        2,
			ReturnType: "t",
		},
		float64(1617231600),
	},
}

func TestPostCellHandler(t *testing.T) {
	godotenv.Load()

	router := setupRouter(false, false)

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

func TestGetCellHandler(t *testing.T) {
	godotenv.Load()

	router := setupRouter(false, false)

	for _, test := range tests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/cell/%s/%d/%s", test.input.Column, test.input.Row, test.input.ReturnType), nil)
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
	assert.Equal(t, 2, q.ColumnIndex())
}
