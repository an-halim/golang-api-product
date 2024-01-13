package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
 gorm.Model
 ID uuid.UUID `gorm:"type:uuid;"`
 Name string    `json:"name"`
 Price int 	`json:"price"`
 Stock int 	`json:"stock"`
 Image_url string `json:"image_url"`
}

type Products struct {
 Products []Product `json:"Products"`
}
func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {

 product.ID = uuid.New()
 return
}