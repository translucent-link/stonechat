package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type columntest struct {
	input    ColumnQuery
	expected interface{}
}

var columntests = []columntest{
	{
		ColumnQuery{
			Column:     "g",
			StartRow:   2,
			EndRow:     4,
			ReturnType: "n",
		},
		"[3,6]",
	},
}

func TestGetColumnHandler(t *testing.T) {
	godotenv.Load()

	router := setupRouter(false, false)

	for _, test := range columntests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/column/%s/%d/%d/%s", test.input.Column, test.input.StartRow, test.input.EndRow, test.input.ReturnType), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		output := ExternalAdapterOutput{}
		json.NewDecoder(w.Body).Decode(&output)

		assert.NotNil(t, output.Data)
		assert.Equal(t, output.Error, "")

		jsonActual, _ := json.Marshal(output.Data.Value)
		assert.Equal(t, test.expected, string(jsonActual))

	}

}
