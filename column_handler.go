package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// handles GET version of Column query
func getColumnHandler(c *gin.Context) {
	defer requestsProcessed.Inc() // increases the metrics counter at the end of the request
	columnQuery := ColumnQuery{}
	err := c.ShouldBindUri(&columnQuery)
	_columnHandler(c, columnQuery, err)
}

// handles POST version of Column query
func postColumnHandler(c *gin.Context) {
	defer requestsProcessed.Inc() // increases the metrics counter at the end of the request
	columnQuery := ColumnQuery{}
	err := c.ShouldBind(&columnQuery)
	_columnHandler(c, columnQuery, err)
}

// main guts of Cell query
func _columnHandler(c *gin.Context, columnQuery ColumnQuery, err error) {
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
				var column []interface{}
				for i, rowOfCells := range source.Values {
					rowNo := uint(i + 1)
					if rowNo >= columnQuery.StartRow && rowNo <= columnQuery.EndRow {

						for j, cell := range rowOfCells {
							if j == columnQuery.ColumnIndex() {
								column = append(column, columnQuery.returnValue(cell))
								break
							}
						}
					}
				}

				returnValue := ExternalAdapterOutput{}
				// assign the price to the output data struct. Good place to make any modification.
				returnValue.Data = &DataOutput{Value: column}
				c.JSON(http.StatusOK, returnValue)
			}
		}
	}
}
