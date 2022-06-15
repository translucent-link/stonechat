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

type rowtest struct {
	input    RowQuery
	expected interface{}
}

var rowtests = []rowtest{
	{
		RowQuery{
			StartColumn: "e",
			EndColumn:   "g",
			Row:         4,
			ReturnType:  "n",
		},
		"[7,8]",
	},
}

func TestGetRowHandler(t *testing.T) {
	godotenv.Load()

	router := setupRouter(false, false)

	for _, test := range rowtests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/row/%s/%s/%d/%s", test.input.StartColumn, test.input.EndColumn, test.input.Row, test.input.ReturnType), nil)
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
