package repository

import (
	"context"
	"errors"
	"fmt"

	"match/cmd/pkg/models"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
)

// Database can communicate with the persistent storage.
type Database struct {
	handler *gorm.DB
}

// NewDatabase creates a new instance of Database with the given SQL database handler.
func NewDatabase(handler *gorm.DB) *Database {
	return &Database{handler: handler}
}

// GetMatches returns all partners that have a radius that cover given latitude and longitude values.
func (db *Database) GetMatches(ctx context.Context, materials []uint, lat, long float32) ([]models.Partner, error) {
	var ps []models.Partner

	subQuery := db.handler.
		WithContext(ctx).
		Select("p1.id, haversine(p1.lat, p1.long, ?, ?) AS distance", lat, long).
		Table("partners p1")

	err := db.handler.
		WithContext(ctx).
		Select("p2.id, p2.lat, p2.long, p2.radius, p2.rating, sub.distance").
		Table("partners p2").
		Joins("JOIN materials ON materials.partner_id = p2.id AND materials.id IN (?)", materials).
		Joins("JOIN (?) sub ON sub.id = p2.id", subQuery).
		Where("sub.distance < p2.radius").
		Group("p2.id, p2.rating, sub.distance").
		Having("COUNT(DISTINCT materials.id) = ?", len(materials)).
		Order("p2.rating desc, sub.distance asc").
		Limit(10).
		Find(&ps).
		Error

	if err != nil {
		return nil, fmt.Errorf("error trying to retrieve the partners from the database: %w", err)
	}

	// if not matches were found just return
	if len(ps) == 0 {
		return ps, nil
	}

	var psIds []uint
	for _, p := range ps {
		psIds = append(psIds, p.ID)
	}

	var cs []models.Category
	err = db.handler.
		WithContext(ctx).
		Model(&models.Category{}).
		Where("partner_id IN (?)", psIds).
		Find(&cs).
		Error

	if err != nil {
		return nil, fmt.Errorf("error trying to retrieve the categories from the database: %w", err)
	}

	var ms []models.Material
	err = db.handler.
		WithContext(ctx).
		Model(&models.Material{}).
		Where("partner_id IN (?)", psIds).
		Find(&ms).
		Error

	if err != nil {
		return nil, fmt.Errorf("error trying to retrieve the materials from the database: %w", err)
	}

	for i := range ps {
		for j := range cs {
			if ps[i].ID == cs[j].PartnerID {
				ps[i].Categories = append(ps[i].Categories, cs[j])
			}
		}
		for j := range ms {
			if ps[i].ID == ms[j].PartnerID {
				ps[i].Materials = append(ps[i].Materials, ms[j])
			}
		}
	}

	return ps, nil
}

// GetPartnerById returns a partner by id.
func (db *Database) GetPartnerById(ctx context.Context, id uint) (models.Partner, error) {
	var p models.Partner

	err := db.handler.
		WithContext(ctx).
		Model(&models.Partner{}).
		Preload("Categories").
		Preload("Materials").
		Where("id = ?", id).
		First(&p).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Partner{}, ErrNotFound
		}
		return models.Partner{}, fmt.Errorf("error trying to retrieve the partner from the database: %w", err)
	}

	return p, nil
}
