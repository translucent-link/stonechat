package main

type ColumnQuery struct {
	Column     string `json:"column" uri:"column" `
	StartRow   uint   `json:"startRow" uri:"startRow" `
	EndRow     uint   `json:"endRow" uri:"endRow" `
	ReturnType string `json:"returnType" uri:"returnType" `
}

func (q ColumnQuery) ColumnIndex() int {
	return columnToIndex(q.Column)
}

func (q ColumnQuery) returnValue(valueStr string) interface{} {
	return cellReturnValue(valueStr, q.ReturnType)
}
