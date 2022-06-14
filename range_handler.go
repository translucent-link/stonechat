package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// handles GET version of Range query
func getRangeHandler(c *gin.Context) {
	defer requestsProcessed.Inc() // increases the metrics counter at the end of the request
	rangeQuery := RangeQuery{}
	err := c.ShouldBindUri(&rangeQuery)
	_rangeHandler(c, rangeQuery, err)
}

// handles POST version of Range query
func postRangeHandler(c *gin.Context) {
	defer requestsProcessed.Inc() // increases the metrics counter at the end of the request
	rangeQuery := RangeQuery{}
	err := c.ShouldBind(&rangeQuery)
	_rangeHandler(c, rangeQuery, err)
}

// main guts of Cell query
func _rangeHandler(c *gin.Context, rangeQuery RangeQuery, err error) {
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
				var table [][]interface{}
				for i, rowOfCells := range source.Values {
					rowNo := uint(i + 1)
					if rowNo >= rangeQuery.StartRow && rowNo <= rangeQuery.EndRow {

						var row []interface{}
						for j, cell := range rowOfCells {
							if j >= rangeQuery.StartColumnIndex() && j <= rangeQuery.EndColumnIndex() {
								row = append(row, rangeQuery.returnValue(cell))
							}
						}
						table = append(table, row)
					}
				}

				returnValue := ExternalAdapterOutput{}
				// assign the price to the output data struct. Good place to make any modification.
				returnValue.Data = &DataOutput{Value: table}
				c.JSON(http.StatusOK, returnValue)
			}
		}
	}
}
