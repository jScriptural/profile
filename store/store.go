package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	_ "modernc.org/sqlite"
	"profile/internal/models"
	"strings"
)

type DBHandle struct {
	DB *sql.DB
}

func NewDBHandle(dbPath string) *DBHandle {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := createSchema(db); err != nil {
		log.Fatal(err)
	}

	return &DBHandle{DB: db}
}

func createSchema(db *sql.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS profiles (
		id BLOB PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		gender TEXT,
		gender_probability REAL,
		sample_size INTEGER,
		age INTEGER,
		age_group TEXT,
		country_id TEXT,
		country_probability REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(schema)
	return err
}

func (d *DBHandle) SaveProfile(ctx context.Context, p *models.Profile) error {
	query := `INSERT INTO profiles 
	(id,name,gender,gender_probability,sample_size,age,age_group,country_id,country_probability,created_at)
	VALUES (?,?,?,?,?,?,?,?,?,?);`

	_, err := d.DB.ExecContext(
		ctx,
		query,
		p.ID,
		p.Name,
		p.Gender,
		p.GenderProbability,
		p.SampleSize,
		p.Age,
		p.AgeGroup,
		p.CountryID,
		p.CountryProbability,
		p.CreatedAt)

	if err != nil {
		return fmt.Errorf("Failed to save profile: %w", err)
	}

	return nil
}

func (d *DBHandle) GetProfileByName(ctx context.Context, name string) (*models.Profile, error) {

	query := `SELECT id,name,gender,gender_probability,sample_size,age,age_group,country_id,country_probability,created_at
	FROM profiles
	WHERE name = ? LIMIT 1;`

	p := models.Profile{}
	var idStr string
	err := d.DB.QueryRowContext(ctx, query, name).Scan(
		&idStr,
		&p.Name,
		&p.Gender,
		&p.GenderProbability,
		&p.SampleSize,
		&p.Age,
		&p.AgeGroup,
		&p.CountryID,
		&p.CountryProbability,
		&p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}

		return nil, fmt.Errorf("GetProfileByName(name: %v): %w", name, err)
	}

	p.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("GetProfileByName(name: %v): %w", name, err)
	}
	return &p, nil
}

func (d *DBHandle) GetProfileByID(ctx context.Context, id string) (*models.Profile, error) {

	query := `SELECT id,name,gender,gender_probability,sample_size,age,age_group,country_id,country_probability,created_at
	FROM profiles
	WHERE id = ? LIMIT 1;`

	var (
		p     models.Profile
		idStr string
	)

	err := d.DB.QueryRowContext(ctx, query, id).Scan(
		&idStr,
		&p.Name,
		&p.Gender,
		&p.GenderProbability,
		&p.SampleSize,
		&p.Age,
		&p.AgeGroup,
		&p.CountryID,
		&p.CountryProbability,
		&p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("GetProfileByID(id: %v): %w", id, models.ErrNoRows)
		}

		return nil, fmt.Errorf("GetProfileByID(id: %v): %w", id, err)
	}

	p.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("GetProfileByID(id: %v): %w", id, err)
	}

	return &p, nil
}

func (d *DBHandle) GetProfiles(ctx context.Context, gender, countryID, ageGroup string) ([]*models.Profile, error) {
	var qb strings.Builder
	qb.WriteString(`SELECT id,name,gender,age,age_group,country_id
	FROM profiles
	WHERE 1=1 `)

	var queryParam []any

	if gender != "" {
		qb.WriteString(`AND gender = ? `)
		queryParam = append(queryParam, gender)
	}

	if countryID != "" {
		qb.WriteString(`AND country_id = ? `)
		queryParam = append(queryParam, countryID)
	}

	if ageGroup != "" {
		qb.WriteString(`AND age_group = ? `)
		queryParam = append(queryParam, ageGroup)
	}
	qb.WriteString(";")

	rows, err := d.DB.QueryContext(ctx, qb.String(), queryParam...)

	sp := []*models.Profile{}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sp, nil
		}

		return nil, fmt.Errorf("GetProfiles: %w", err)
	}

	var idStr string
	for rows.Next() {
		p := models.Profile{}
		err := rows.Scan(
			&idStr,
			&p.Name,
			&p.Gender,
			&p.Age,
			&p.AgeGroup,
			&p.CountryID,
		)
		if err != nil {
			return nil, fmt.Errorf("GetProfiles: %w", err)
		}

		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("GetProfiles: %w", err)
		}

		sp = append(sp, &p)
	}

	return sp, nil
}

func (d *DBHandle) DeleteProfileByID(ctx context.Context, id string) error {
	query := `DELETE 
	FROM profiles 
	WHERE id = ?;`

	res, err := d.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeleteProfileByID(id: %v): %w", id, err)
	}

	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteProfileByID(id: %v): %w", id, err)
	}

	if rowsDeleted == 0 {
		return fmt.Errorf("DeleteProfileByID: %w", models.ErrNoRows)
	}

	return nil
}
