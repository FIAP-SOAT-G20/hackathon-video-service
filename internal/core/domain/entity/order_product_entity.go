package entity

import (
	"time"
)

type OrderProduct struct {
	OrderID   uint64
	ProductID uint64
	Quantity  uint32
	Order     Order   // Virtual field
	Product   Product // Virtual field
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *OrderProduct) Update(quantity uint32) {
	p.Quantity = quantity
	p.UpdatedAt = time.Now()
	p.Order = Order{}
	p.Product = Product{}
}
