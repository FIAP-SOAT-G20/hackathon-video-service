package presenter

import "encoding/json"

type CategoryJsonResponse struct {
	ID        uint64 `json:"id" example:"1"`
	Name      string `json:"name" example:"John Doe"`
	CreatedAt string `json:"created_at" example:"2024-02-09T10:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-02-09T10:00:00Z"`
}

func (r CategoryJsonResponse) String() string {
	o, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(o)
}

type CategoryJsonPaginatedResponse struct {
	JsonPagination
	Categories []CategoryJsonResponse `json:"categories"`
}

func (r CategoryJsonPaginatedResponse) String() string {
	o, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(o)
}
