package models

// MatchRequest represents '/partners/match' request.
type MatchRequest struct {
	Materials    []uint  `json:"materials"`
	Address      Address `json:"address"`
	SquareMeters uint    `json:"square_meters"`
	PhoneNumber  string  `json:"phone_number"`
}

// Partner represents a partner.
type Partner struct {
	ID         uint       `json:"id" gorm:"column:id"`
	Categories []Category `json:"categories" gorm:"foreignKey:PartnerID"`
	Materials  []Material `json:"materials" gorm:"foreignKey:PartnerID"`
	Address    Address    `json:"address" gorm:"embedded"`
	Radius     int        `json:"radius" gorm:"column:radius"`
	Rating     int        `json:"rating" gorm:"column:rating"`
}

// Category represents a partner's category.
type Category struct {
	ID          uint   `json:"id" gorm:"column:id"`
	PartnerID   uint   `json:"-" gorm:"column:partner_id"`
	Description string `json:"description" gorm:"column:description"`
}

// Material represents a partner's material.
type Material struct {
	ID          uint   `json:"id" gorm:"column:id"`
	PartnerID   uint   `json:"-" gorm:"column:partner_id"`
	Description string `json:"description" gorm:"column:description"`
}

// Address represents an address by its latitude and longitude.
type Address struct {
	Lat  float32 `json:"lat"`
	Long float32 `json:"long"`
}
