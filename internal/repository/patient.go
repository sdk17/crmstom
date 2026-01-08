package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sdk17/crmstom/internal/domain"
)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Create(patient *domain.Patient) error {
	query := `INSERT INTO patients (iin, name, phone, email, birth_date, address) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(query, patient.IIN, patient.Name, patient.Phone, patient.Email, patient.BirthDate, patient.Address).
		Scan(&patient.ID, &patient.CreatedAt, &patient.UpdatedAt)

	return err
}

func (r *PatientRepository) GetByID(id int) (*domain.Patient, error) {
	query := `SELECT id, COALESCE(iin, ''), name, phone, email, birth_date, address, created_at, updated_at
			  FROM patients WHERE id = $1 AND deleted_at IS NULL`

	patient := &domain.Patient{}
	err := r.db.QueryRow(query, id).Scan(
		&patient.ID, &patient.IIN, &patient.Name, &patient.Phone, &patient.Email,
		&patient.BirthDate, &patient.Address, &patient.CreatedAt, &patient.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пациент с ID %d не найден", id)
		}
		return nil, err
	}

	return patient, nil
}

func (r *PatientRepository) GetAll() ([]*domain.Patient, error) {
	query := `SELECT id, COALESCE(iin, ''), name, phone, email, birth_date, address, created_at, updated_at
			  FROM patients WHERE deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*domain.Patient
	for rows.Next() {
		patient := &domain.Patient{}
		err := rows.Scan(
			&patient.ID, &patient.IIN, &patient.Name, &patient.Phone, &patient.Email,
			&patient.BirthDate, &patient.Address, &patient.CreatedAt, &patient.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, nil
}

func (r *PatientRepository) Update(patient *domain.Patient) error {
	query := `UPDATE patients SET iin = $1, name = $2, phone = $3, email = $4, birth_date = $5,
			  address = $6, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $7 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, patient.IIN, patient.Name, patient.Phone, patient.Email,
		patient.BirthDate, patient.Address, patient.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("пациент с ID %d не найден", patient.ID)
	}

	return nil
}

func (r *PatientRepository) Delete(id int) error {
	query := `UPDATE patients SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("пациент с ID %d не найден", id)
	}

	return nil
}

func (r *PatientRepository) GetByPhone(phone string) (*domain.Patient, error) {
	query := `SELECT id, COALESCE(iin, ''), name, phone, email, birth_date, address, created_at, updated_at
			  FROM patients WHERE phone = $1 AND deleted_at IS NULL`

	patient := &domain.Patient{}
	err := r.db.QueryRow(query, phone).Scan(
		&patient.ID, &patient.IIN, &patient.Name, &patient.Phone, &patient.Email,
		&patient.BirthDate, &patient.Address, &patient.CreatedAt, &patient.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пациент с телефоном %s не найден", phone)
		}
		return nil, err
	}

	return patient, nil
}

func (r *PatientRepository) GetByIIN(iin string) (*domain.Patient, error) {
	query := `SELECT id, COALESCE(iin, ''), name, phone, email, birth_date, address, created_at, updated_at
			  FROM patients WHERE iin = $1 AND deleted_at IS NULL`

	patient := &domain.Patient{}
	err := r.db.QueryRow(query, iin).Scan(
		&patient.ID, &patient.IIN, &patient.Name, &patient.Phone, &patient.Email,
		&patient.BirthDate, &patient.Address, &patient.CreatedAt, &patient.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пациент с ИИН %s не найден", iin)
		}
		return nil, err
	}

	return patient, nil
}

func (r *PatientRepository) Search(query string) ([]*domain.Patient, error) {
	searchQuery := `SELECT id, COALESCE(iin, ''), name, phone, email, birth_date, address, created_at, updated_at
					FROM patients WHERE deleted_at IS NULL AND (COALESCE(iin, '') ILIKE $1 OR name ILIKE $1 OR phone ILIKE $1 OR email ILIKE $1)
					ORDER BY name`

	rows, err := r.db.Query(searchQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*domain.Patient
	for rows.Next() {
		patient := &domain.Patient{}
		err := rows.Scan(
			&patient.ID, &patient.IIN, &patient.Name, &patient.Phone, &patient.Email,
			&patient.BirthDate, &patient.Address, &patient.CreatedAt, &patient.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}

	return patients, nil
}
