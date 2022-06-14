package main

type CellQuery struct {
	Column     string `json:"column" uri:"column" `
	Row        uint   `json:"row" uri:"row" `
	ReturnType string `json:"returnType" uri:"returnType" `
}

func (q CellQuery) ColumnIndex() int {
	return columnToIndex(q.Column)
}

func (q CellQuery) validReturnType() bool {
	return q.ReturnType == "n" || q.ReturnType == "s" || q.ReturnType == "b" || q.ReturnType == "t" || q.ReturnType == "o"
}

func (q CellQuery) Valid() bool {
	return q.validReturnType()
}

func (q CellQuery) returnValue(valueStr string) interface{} {
	return cellReturnValue(valueStr, q.ReturnType)
}
