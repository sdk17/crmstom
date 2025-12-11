package repository

import (
	"database/sql"
	"fmt"

	"github.com/sdk17/crmstom/internal/domain"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(service *domain.Service) error {
	query := `INSERT INTO services (name, type, notes) 
			  VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, service.Name, service.Type, service.Notes).
		Scan(&service.ID, &service.CreatedAt, &service.UpdatedAt)

	return err
}

func (r *ServiceRepository) GetByID(id int) (*domain.Service, error) {
	query := `SELECT id, name, type, notes, created_at, updated_at
			  FROM services WHERE id = $1 AND deleted_at IS NULL`

	service := &domain.Service{}
	err := r.db.QueryRow(query, id).Scan(
		&service.ID, &service.Name, &service.Type, &service.Notes, &service.CreatedAt, &service.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("услуга с ID %d не найдена", id)
		}
		return nil, err
	}

	return service, nil
}

func (r *ServiceRepository) GetAll() ([]*domain.Service, error) {
	query := `SELECT id, name, type, notes, created_at, updated_at
			  FROM services WHERE deleted_at IS NULL ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*domain.Service
	for rows.Next() {
		service := &domain.Service{}
		err := rows.Scan(
			&service.ID, &service.Name, &service.Type, &service.Notes, &service.CreatedAt, &service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *ServiceRepository) Update(service *domain.Service) error {
	query := `UPDATE services SET name = $1, type = $2, notes = $3, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $4 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, service.Name, service.Type, service.Notes, service.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("услуга с ID %d не найдена", service.ID)
	}

	return nil
}

func (r *ServiceRepository) Delete(id int) error {
	query := `UPDATE services SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("услуга с ID %d не найдена", id)
	}

	return nil
}

func (r *ServiceRepository) GetByCategory(category string) ([]*domain.Service, error) {
	query := `SELECT id, name, type, notes, created_at, updated_at
			  FROM services WHERE deleted_at IS NULL AND type = $1 ORDER BY name`

	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*domain.Service
	for rows.Next() {
		service := &domain.Service{}
		err := rows.Scan(
			&service.ID, &service.Name, &service.Type, &service.Notes, &service.CreatedAt, &service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *ServiceRepository) Search(query string) ([]*domain.Service, error) {
	searchQuery := `SELECT id, name, type, notes, created_at, updated_at
					FROM services WHERE deleted_at IS NULL AND (name ILIKE $1 OR notes ILIKE $1) ORDER BY name`

	rows, err := r.db.Query(searchQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*domain.Service
	for rows.Next() {
		service := &domain.Service{}
		err := rows.Scan(
			&service.ID, &service.Name, &service.Type, &service.Notes, &service.CreatedAt, &service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}
