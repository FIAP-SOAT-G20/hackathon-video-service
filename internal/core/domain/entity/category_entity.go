package entity

import "time"

type Category struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Category) Update(name string) {
	p.Name = name
	p.UpdatedAt = time.Now()
}
