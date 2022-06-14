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

type rangetest struct {
	input    RangeQuery
	expected interface{}
}

var values = [][]float64{
	{1, 2, 3},
	{4, 5, 6},
	{7, 8, 9},
}

var values2 = []interface{}(
	{1, 2, 3},
)

var rangetests = []rangetest{
	{
		RangeQuery{
			StartColumn: "e",
			StartRow:    2,
			EndColumn:   "g",
			EndRow:      4,
			ReturnType:  "n",
		},
		values,
	},
}

func TestGetRangeHandler(t *testing.T) {
	godotenv.Load()

	router := setupRouter(false, false)

	for _, test := range rangetests {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/range/%s/%d/%s/%d/%s", test.input.StartColumn, test.input.StartRow, test.input.EndColumn, test.input.EndRow, test.input.ReturnType), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		output := ExternalAdapterOutput{}
		json.NewDecoder(w.Body).Decode(&output)

		assert.NotNil(t, output.Data)
		assert.Equal(t, output.Error, "")

		assert.Equal(t, test.expected, output.Data.Value)

	}

}
