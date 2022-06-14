package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// handles GET version of Row query
func getRowHandler(c *gin.Context) {
	defer requestsProcessed.Inc() // increases the metrics counter at the end of the request
	rowQuery := RowQuery{}
	err := c.ShouldBindUri(&rowQuery)
	_rowHandler(c, rowQuery, err)
}

// handles POST version of Row query
func postRowHandler(c *gin.Context) {
	defer requestsProcessed.Inc() // increases the metrics counter at the end of the request
	rowQuery := RowQuery{}
	err := c.ShouldBind(&rowQuery)
	_rowHandler(c, rowQuery, err)
}

// main guts of Cell query
func _rowHandler(c *gin.Context, rowQuery RowQuery, err error) {
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
				var returnRow []interface{}
				for i, rowOfCells := range source.Values {
					rowNo := uint(i + 1)
					if rowNo == rowQuery.Row {

						for j, cell := range rowOfCells {
							if j >= rowQuery.StartColumnIndex() && j <= rowQuery.EndColumnIndex() {
								returnRow = append(returnRow, rowQuery.returnValue(cell))
							}
						}
						break
					}
				}

				returnValue := ExternalAdapterOutput{}
				// assign the price to the output data struct. Good place to make any modification.
				returnValue.Data = &DataOutput{Value: returnRow}
				c.JSON(http.StatusOK, returnValue)
			}
		}
	}
}
