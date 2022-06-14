package main

type RowQuery struct {
	StartColumn string `json:"startColumn" uri:"startColumn" `
	EndColumn   string `json:"endColumn" uri:"endColumn" `
	Row         uint   `json:"row" uri:"row" `
	ReturnType  string `json:"returnType" uri:"returnType" `
}

func (q RowQuery) StartColumnIndex() int {
	return columnToIndex(q.StartColumn)
}

func (q RowQuery) EndColumnIndex() int {
	return columnToIndex(q.EndColumn)
}

func (q RowQuery) returnValue(valueStr string) interface{} {
	return cellReturnValue(valueStr, q.ReturnType)
}
