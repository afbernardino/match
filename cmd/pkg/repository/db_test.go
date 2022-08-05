package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"match/cmd/pkg/models"
	"match/cmd/pkg/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	queryGetPartnerById           = `SELECT * FROM "partners" WHERE id = $1 ORDER BY "partners"."id" LIMIT 1`
	queryGetCategoriesByPartnerId = `SELECT * FROM "categories" WHERE "categories"."partner_id" = $1`
	queryGetMaterialsByPartnerId  = `SELECT * FROM "materials" WHERE "materials"."partner_id" = $1`
	queryGetPartnersMatch         = `SELECT p2.id, p2.lat, p2.long, p2.radius, p2.rating, sub.distance FROM partners p2 JOIN materials ON materials.partner_id = p2.id AND materials.id IN ($1,$2) JOIN (SELECT p1.id, haversine(p1.lat, p1.long, $3, $4) AS distance FROM partners p1) sub ON sub.id = p2.id WHERE sub.distance < p2.radius GROUP BY p2.id, p2.rating, sub.distance HAVING COUNT(DISTINCT materials.id) = $5 ORDER BY p2.rating desc, sub.distance asc LIMIT 10`
	queryGetCategoriesMatch       = `SELECT * FROM "categories" WHERE partner_id IN ($1)`
	queryGetMaterialsMatch        = `SELECT * FROM "materials" WHERE partner_id IN ($1)`
)

func initDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating a stub for database connection: '%s'", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	var handler *gorm.DB
	handler, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening a stub database connection: '%s'", err)
	}

	return db, mock, handler
}

func TestGetMatches_NoMatchesFound(t *testing.T) {
	db, mock, handler := initDB(t)
	defer db.Close()

	repo := repository.NewDatabase(handler)

	materials := []uint{1, 2}
	lat := float32(1.1)
	long := float32(1.2)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetPartnersMatch)).
		WithArgs(materials[0], materials[1], lat, long, len(materials)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	ps, err := repo.GetMatches(context.Background(), materials, lat, long)

	if err != nil {
		t.Errorf("error mismatch: want 'nil' got '%s'", err)
	}

	if diff := cmp.Diff([]models.Partner{}, ps); diff != "" {
		t.Errorf("guest list mismatch (-want +got):\n%s", diff)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: '%s'", err)
	}
}

func TestGetMatches_Success(t *testing.T) {
	db, mock, handler := initDB(t)
	defer db.Close()

	repo := repository.NewDatabase(handler)

	pExpected := models.Partner{
		ID: 1,
		Categories: []models.Category{
			{
				ID:          2,
				PartnerID:   1,
				Description: "category 2",
			},
		},
		Materials: []models.Material{
			{
				ID:          3,
				PartnerID:   1,
				Description: "material 3",
			},
			{
				ID:          4,
				PartnerID:   1,
				Description: "material 4",
			},
		},
		Address: models.Address{
			Lat:  1.1,
			Long: 1.2,
		},
		Radius: 100,
		Rating: 5,
	}

	pRows := sqlmock.NewRows([]string{"id", "lat", "long", "radius", "rating", "distance"})
	pRows.AddRow(pExpected.ID, pExpected.Address.Lat, pExpected.Address.Long, pExpected.Radius, pExpected.Rating, 1)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetPartnersMatch)).
		WithArgs(pExpected.Materials[0].ID, pExpected.Materials[1].ID, pExpected.Address.Lat, pExpected.Address.Long, len(pExpected.Materials)).
		WillReturnRows(pRows)

	cRows := sqlmock.NewRows([]string{"id", "partner_id", "description"})
	cRows.AddRow(pExpected.Categories[0].ID, pExpected.Categories[0].PartnerID, pExpected.Categories[0].Description)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetCategoriesMatch)).
		WithArgs(pExpected.ID).
		WillReturnRows(cRows)

	mRows := sqlmock.NewRows([]string{"id", "partner_id", "description"})
	mRows.AddRow(pExpected.Materials[0].ID, pExpected.Materials[0].PartnerID, pExpected.Materials[0].Description)
	mRows.AddRow(pExpected.Materials[1].ID, pExpected.Materials[1].PartnerID, pExpected.Materials[1].Description)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetMaterialsMatch)).
		WithArgs(pExpected.ID).
		WillReturnRows(mRows)

	ps, err := repo.GetMatches(
		context.Background(),
		[]uint{pExpected.Materials[0].ID, pExpected.Materials[1].ID},
		pExpected.Address.Lat,
		pExpected.Address.Long,
	)

	if err != nil {
		t.Errorf("error mismatch: want 'nil' got '%s'", err)
	}

	psExpected := []models.Partner{pExpected}
	if diff := cmp.Diff(psExpected, ps); diff != "" {
		t.Errorf("guest list mismatch (-want +got):\n%s", diff)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: '%s'", err)
	}
}

func TestGetPartnerById_NotFoundFailure(t *testing.T) {
	db, mock, handler := initDB(t)
	defer db.Close()

	repo := repository.NewDatabase(handler)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetPartnerById)).
		WithArgs(1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.GetPartnerById(context.Background(), 1)

	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error mismatch: want '%s' got '%s'", repository.ErrNotFound, err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: '%s'", err)
	}
}

func TestGetPartnerById_NoCategoryOrMaterials(t *testing.T) {
	db, mock, handler := initDB(t)
	defer db.Close()

	repo := repository.NewDatabase(handler)

	pExpected := models.Partner{
		ID:         1,
		Categories: []models.Category{},
		Materials:  []models.Material{},
		Address: models.Address{
			Lat:  1.1,
			Long: 1.2,
		},
		Radius: 100,
		Rating: 4,
	}

	pRows := sqlmock.NewRows([]string{"id", "lat", "long", "radius", "rating"})
	pRows.AddRow(pExpected.ID, pExpected.Address.Lat, pExpected.Address.Long, pExpected.Radius, pExpected.Rating)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetPartnerById)).
		WithArgs(pExpected.ID).
		WillReturnRows(pRows)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetCategoriesByPartnerId)).
		WithArgs(pExpected.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	mock.ExpectQuery(regexp.QuoteMeta(queryGetMaterialsByPartnerId)).
		WithArgs(pExpected.ID).
		WillReturnRows(sqlmock.NewRows([]string{}))

	p, err := repo.GetPartnerById(context.Background(), pExpected.ID)

	if err != nil {
		t.Errorf("error mismatch: want 'nil' got '%s'", err)
	}

	if diff := cmp.Diff(pExpected, p); diff != "" {
		t.Errorf("guest list mismatch (-want +got):\n%s", diff)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: '%s'", err)
	}
}

func TestGetPartnerById_Success(t *testing.T) {
	db, mock, handler := initDB(t)
	defer db.Close()

	repo := repository.NewDatabase(handler)

	pExpected := models.Partner{
		ID: 1,
		Categories: []models.Category{
			{
				ID:          2,
				PartnerID:   1,
				Description: "category 2",
			},
		},
		Materials: []models.Material{
			{
				ID:          3,
				PartnerID:   1,
				Description: "material 3",
			},
		},
		Address: models.Address{
			Lat:  1.1,
			Long: 1.2,
		},
		Radius: 100,
		Rating: 4,
	}

	pRows := sqlmock.NewRows([]string{"id", "lat", "long", "radius", "rating"})
	pRows.AddRow(pExpected.ID, pExpected.Address.Lat, pExpected.Address.Long, pExpected.Radius, pExpected.Rating)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetPartnerById)).
		WithArgs(pExpected.ID).
		WillReturnRows(pRows)

	cRows := sqlmock.NewRows([]string{"id", "partner_id", "description"})
	cRows.AddRow(pExpected.Categories[0].ID, pExpected.Categories[0].PartnerID, pExpected.Categories[0].Description)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetCategoriesByPartnerId)).
		WithArgs(pExpected.ID).
		WillReturnRows(cRows)

	mRows := sqlmock.NewRows([]string{"id", "partner_id", "description"})
	mRows.AddRow(pExpected.Materials[0].ID, pExpected.Materials[0].PartnerID, pExpected.Materials[0].Description)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetMaterialsByPartnerId)).
		WithArgs(pExpected.ID).
		WillReturnRows(mRows)

	p, err := repo.GetPartnerById(context.Background(), pExpected.ID)

	if err != nil {
		t.Errorf("error mismatch: want 'nil' got '%s'", err)
	}

	if diff := cmp.Diff(pExpected, p); diff != "" {
		t.Errorf("guest list mismatch (-want +got):\n%s", diff)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations were not met: '%s'", err)
	}
}
