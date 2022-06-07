package main

import (
	"encoding/json"
	"net/http"
	"strconv"

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

type SheetResponse struct {
	Values [][]string
}

// /cell request handler of this external adapter. Customise to your heart's content.
func cellHandler(c *gin.Context) {
	defer requestsProcessed.Inc() // increases the metrics counter at the end of the request

	cellQuery := CellQuery{}

	err := c.ShouldBind(&cellQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, ExternalAdapterOutput{Error: errors.Wrap(err, "Unable to parse rqeust").Error()})
	} else {
		// Fetch the data from the source URL
		res, err := http.Get(cellQuery.Url)
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
							if j == cellQuery.CellIndex() {
								valueStr = cell
							}
						}
					}
				}

				var returnable interface{}
				if cellQuery.ReturnType == "number" {
					returnable, _ = strconv.Atoi(valueStr)
				} else {
					returnable = valueStr
				}

				returnValue := ExternalAdapterOutput{}
				// assign the price to the output data struct. Good place to make any modification.
				returnValue.Data = &DataOutput{Value: returnable}
				c.JSON(http.StatusOK, returnValue)
			}
		}
	}

}
