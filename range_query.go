package main

type RangeQuery struct {
	StartColumn string `json:"startColumn" uri:"startColumn" `
	StartRow    uint   `json:"startRow" uri:"startRow" `
	EndColumn   string `json:"endColumn" uri:"endColumn" `
	EndRow      uint   `json:"endRow" uri:"endRow" `
	ReturnType  string `json:"returnType" uri:"returnType" `
}

func (q RangeQuery) StartColumnIndex() int {
	return columnToIndex(q.StartColumn)
}

func (q RangeQuery) EndColumnIndex() int {
	return columnToIndex(q.EndColumn)
}

func (q RangeQuery) returnValue(valueStr string) interface{} {
	return cellReturnValue(valueStr, q.ReturnType)
}
