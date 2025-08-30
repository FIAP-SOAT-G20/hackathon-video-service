package presenter

type ProductXmlResponse struct {
	ID          uint64  `xml:"id" example:"1"`
	Name        string  `xml:"name" example:"Product A"`
	Description string  `xml:"description" example:"Description of product A"`
	Price       float64 `xml:"price" example:"99.99"`
	CategoryID  uint64  `xml:"category_id" example:"1"`
	CreatedAt   string  `xml:"created_at" example:"2024-02-09T10:00:00Z"`
	UpdatedAt   string  `xml:"updated_at" example:"2024-02-09T10:00:00Z"`
}

type ProductXmlPaginatedResponse struct {
	XmlPagination
	Products []ProductXmlResponse `xml:"products"`
}
