package entity

import (
	"time"

	"github.com/google/uuid"
)

// Package represents a subscription plan or feature package
type Package struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Features    []string  `json:"features" db:"features"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// // NewPackage creates a new package instance
// func NewPackage(name, description string, price float64, features []string) *Package {
// 	now := time.Now()
// 	return &Package{
// 		ID:          uuid.New(),
// 		Name:        name,
// 		Description: description,
// 		Price:       price,
// 		Features:    features,
// 		CreatedAt:   now,
// 		UpdatedAt:   now,
// 	}
// }

// // UpdatePackage updates package information
// func (p *Package) UpdatePackage(name, description string, price float64, features []string) {
// 	p.Name = name
// 	p.Description = description
// 	p.Price = price
// 	p.Features = features
// 	p.UpdatedAt = time.Now()
// }

// // HasFeature checks if the package includes a specific feature
// func (p *Package) HasFeature(feature string) bool {
// 	for _, f := range p.Features {
// 		if f == feature {
// 			return true
// 		}
// 	}
// 	return false
// } 