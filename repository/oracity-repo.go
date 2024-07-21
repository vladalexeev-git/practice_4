package repository

import (
	"city-server/models"
	"context"
	"database/sql"
	"log"
)

type CityRepo interface {
	CreateCity(models.City) (cityId int, err error)
	GetCityById(id int) (models.City, error)
	GetAllCities() ([]models.City, error)
	UpdateCity(models.City) error
}

type cityRepo struct {
	DB *sql.DB
}

func NewCityRepo(db *sql.DB) CityRepo {
	return &cityRepo{DB: db}
}

func (repo *cityRepo) CreateCity(city models.City) (cityId int, err error) {
	var id int64
	ctx := context.Background()
	query := "INSERT INTO cities (name, population) VALUES (:name, :population) RETURNING id INTO :id"
	stmt, err := repo.DB.PrepareContext(ctx, query)
	if err != nil {
		log.Println("failed to prepare statement: ", err)
		return 0, err
	}
	defer stmt.Close()

	// executing query
	_, err = stmt.ExecContext(ctx,
		sql.Named("name", city.Name),
		sql.Named("population", city.Population),
		sql.Named("id", sql.Out{Dest: &id}),
	)

	if err != nil {
		log.Println("failed to execute statement: ", err)
		return 0, err
	}

	log.Println("city id: ", id)
	return int(id), nil
}

func (repo *cityRepo) GetCityById(id int) (models.City, error) {
	var city models.City
	err := repo.DB.QueryRow(
		"SELECT id, name, population FROM cities WHERE id = :1",
		id,
	).Scan(&city.ID, &city.Name, &city.Population)
	if err != nil {
		return models.City{}, err
	}
	return city, nil
}

func (repo *cityRepo) GetAllCities() ([]models.City, error) {
	rows, err := repo.DB.Query(
		"SELECT id, name, population FROM cities",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var city models.City
		if err := rows.Scan(&city.ID, &city.Name, &city.Population); err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cities, nil
}

func (repo *cityRepo) UpdateCity(city models.City) error {
	_, err := repo.DB.Exec(
		"UPDATE cities SET name = :1, population = :2 WHERE id = :3",
		city.Name, city.Population, city.ID,
	)
	return err
}
