package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Represents the JSON structure of the EA's data output. Customise the contents of this struct.
type DataOutput struct {
	Value interface{} `json:"value"`
}

// Represents the JSON structure of the overall response. Should NOT need customisation.
type ExternalAdapterOutput struct {
	Data  *DataOutput `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// Represents tabular structure of JSON-version of Google Sheet
type SheetResponse struct {
	Values [][]string
}

// generates URL to access Google Sheet
func sheetURL() string {
	sheetId := os.Getenv("SHEETS_ID")
	apiKey := os.Getenv("SHEETS_API_KEY")
	tabName := os.Getenv("SHEETS_TAB_NAME")
	return fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?alt=json&key=%s", sheetId, tabName, apiKey)
}

// handles GET version of Cell query
func getCellHandler(c *gin.Context) {
	defer requestsProcessed.Inc() // increases the metrics counter at the end of the request
	cellQuery := CellQuery{}
	err := c.ShouldBindUri(&cellQuery)
	_cellHandler(c, cellQuery, err)
}

// handles POST version of Cell query
func postCellHandler(c *gin.Context) {
	defer requestsProcessed.Inc() // increases the metrics counter at the end of the request
	cellQuery := CellQuery{}
	err := c.ShouldBind(&cellQuery)
	_cellHandler(c, cellQuery, err)
}

// main guts of Cell query
func _cellHandler(c *gin.Context, cellQuery CellQuery, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, ExternalAdapterOutput{Error: errors.Wrap(err, "Unable to parse rqeust").Error()})
	} else {
		// Fetch the data from the source URL
		res, err := http.Get(sheetURL())
		if err != nil {
			c.JSON(http.StatusBadGateway, ExternalAdapterOutput{Error: errors.Wrap(err, "Unable to fetch data from source").Error()})
		} else {
			defer res.Body.Close() // tidies up the open input stream at the end of the request

			source := SheetResponse{}
			err := json.NewDecoder(res.Body).Decode(&source)
			if err != nil {
				c.JSON(500, ExternalAdapterOutput{Error: errors.Wrap(err, "Unable to parse data received from source").Error()})
			} else {
				valueStr := ""
				for i, rowOfCells := range source.Values {
					if uint(i+1) == cellQuery.Row {
						for j, cell := range rowOfCells {
							if j == cellQuery.ColumnIndex() {
								valueStr = cell
							}
						}
					}
				}

				value := cellQuery.returnValue(valueStr)
				returnValue := ExternalAdapterOutput{}
				returnValue.Data = &DataOutput{Value: value}
				c.JSON(http.StatusOK, returnValue)
			}
		}
	}
}
