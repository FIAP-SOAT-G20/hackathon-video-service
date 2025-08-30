package presenter

type XmlPagination struct {
	Total int64 `xml:"total" example:"100"`
	Page  int   `xml:"page" example:"1"`
	Limit int   `xml:"limit" example:"10"`
}
