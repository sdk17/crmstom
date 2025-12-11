package repository

import (
	"database/sql"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
)

type DoctorRepository struct {
	db *sql.DB
}

func NewDoctorRepository(db *sql.DB) *DoctorRepository {
	return &DoctorRepository{db: db}
}

// Create создает нового врача
func (r *DoctorRepository) Create(doctor *domain.Doctor) error {
	query := `
		INSERT INTO doctors (name, email, login, password, is_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	now := time.Now()
	return r.db.QueryRow(
		query,
		doctor.Name,
		doctor.Email,
		doctor.Login,
		doctor.Password,
		doctor.IsAdmin,
		now,
		now,
	).Scan(&doctor.ID)
}

// GetByID получает врача по ID
func (r *DoctorRepository) GetByID(id int) (*domain.Doctor, error) {
	query := `
		SELECT id, name, email, login, password, is_admin, created_at, updated_at
		FROM doctors
		WHERE id = $1 AND deleted_at IS NULL`

	doctor := &domain.Doctor{}
	err := r.db.QueryRow(query, id).Scan(
		&doctor.ID,
		&doctor.Name,
		&doctor.Email,
		&doctor.Login,
		&doctor.Password,
		&doctor.IsAdmin,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return doctor, nil
}

// GetAll получает всех врачей
func (r *DoctorRepository) GetAll() ([]*domain.Doctor, error) {
	query := `
		SELECT id, name, email, login, password, is_admin, created_at, updated_at
		FROM doctors
		WHERE deleted_at IS NULL
		ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	doctors := make([]*domain.Doctor, 0)
	for rows.Next() {
		doctor := &domain.Doctor{}
		err := rows.Scan(
			&doctor.ID,
			&doctor.Name,
			&doctor.Email,
			&doctor.Login,
			&doctor.Password,
			&doctor.IsAdmin,
			&doctor.CreatedAt,
			&doctor.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
	}

	return doctors, rows.Err()
}

// Update обновляет врача
func (r *DoctorRepository) Update(doctor *domain.Doctor) error {
	query := `
		UPDATE doctors
		SET name = $1, email = $2, login = $3, password = $4, is_admin = $5, updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL`

	result, err := r.db.Exec(
		query,
		doctor.Name,
		doctor.Email,
		doctor.Login,
		doctor.Password,
		doctor.IsAdmin,
		time.Now(),
		doctor.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete удаляет врача (soft delete)
func (r *DoctorRepository) Delete(id int) error {
	query := `UPDATE doctors SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetByLogin получает врача по логину
func (r *DoctorRepository) GetByLogin(login string) (*domain.Doctor, error) {
	query := `
		SELECT id, name, email, login, password, is_admin, created_at, updated_at
		FROM doctors
		WHERE login = $1 AND deleted_at IS NULL`

	doctor := &domain.Doctor{}
	err := r.db.QueryRow(query, login).Scan(
		&doctor.ID,
		&doctor.Name,
		&doctor.Email,
		&doctor.Login,
		&doctor.Password,
		&doctor.IsAdmin,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return doctor, nil
}
