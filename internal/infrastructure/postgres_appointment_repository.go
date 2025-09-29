package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
)

type PostgresAppointmentRepository struct {
	db *sql.DB
}

func NewPostgresAppointmentRepository(db *sql.DB) *PostgresAppointmentRepository {
	return &PostgresAppointmentRepository{db: db}
}

func (r *PostgresAppointmentRepository) Create(appointment *domain.Appointment) error {
	// Получаем service_id по названию услуги
	var serviceID int
	serviceQuery := `SELECT id FROM services WHERE name = $1 LIMIT 1`
	err := r.db.QueryRow(serviceQuery, appointment.Service).Scan(&serviceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("услуга '%s' не найдена", appointment.Service)
		}
		return err
	}

	query := `INSERT INTO appointments (patient_id, service_id, appointment_date, status, price, duration_minutes, notes) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(query, appointment.PatientID, serviceID,
		appointment.Date, appointment.Status, appointment.Price, appointment.Duration, appointment.Notes).
		Scan(&appointment.ID, &appointment.CreatedAt, &appointment.UpdatedAt)

	return err
}

func (r *PostgresAppointmentRepository) GetByID(id int) (*domain.Appointment, error) {
	query := `SELECT a.id, a.patient_id, a.service_id, a.appointment_date, a.status, a.price, a.duration_minutes, a.notes, a.created_at, a.updated_at,
			  s.name as service_name
			  FROM appointments a
			  LEFT JOIN services s ON a.service_id = s.id
			  WHERE a.id = $1`

	appointment := &domain.Appointment{}
	var serviceID int
	var serviceName sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&appointment.ID, &appointment.PatientID, &serviceID,
		&appointment.Date, &appointment.Status, &appointment.Price, &appointment.Duration, &appointment.Notes,
		&appointment.CreatedAt, &appointment.UpdatedAt, &serviceName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("запись с ID %d не найдена", id)
		}
		return nil, err
	}

	if serviceName.Valid {
		appointment.Service = serviceName.String
	} else {
		appointment.Service = "Неизвестная услуга"
	}

	return appointment, nil
}

func (r *PostgresAppointmentRepository) GetAll() ([]*domain.Appointment, error) {
	query := `SELECT a.id, a.patient_id, a.service_id, a.appointment_date, a.status, a.price, a.duration_minutes, a.notes, a.created_at, a.updated_at,
			  s.name as service_name
			  FROM appointments a
			  LEFT JOIN services s ON a.service_id = s.id
			  ORDER BY a.appointment_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		appointment := &domain.Appointment{}
		var serviceID int
		var serviceName sql.NullString
		err := rows.Scan(
			&appointment.ID, &appointment.PatientID, &serviceID,
			&appointment.Date, &appointment.Status, &appointment.Price, &appointment.Duration, &appointment.Notes,
			&appointment.CreatedAt, &appointment.UpdatedAt, &serviceName,
		)
		if err != nil {
			return nil, err
		}

		if serviceName.Valid {
			appointment.Service = serviceName.String
		} else {
			appointment.Service = "Неизвестная услуга"
		}

		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *PostgresAppointmentRepository) GetByDateRange(startDate, endDate time.Time) ([]*domain.Appointment, error) {
	query := `SELECT id, patient_id, service_id, appointment_date, status, notes, created_at, updated_at 
			  FROM appointments WHERE appointment_date BETWEEN $1 AND $2 
			  ORDER BY appointment_date`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		appointment := &domain.Appointment{}
		var serviceID int
		err := rows.Scan(
			&appointment.ID, &appointment.PatientID, &serviceID,
			&appointment.Date, &appointment.Status, &appointment.Notes,
			&appointment.CreatedAt, &appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// TODO: Получить название услуги по serviceID
		appointment.Service = "Услуга"
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *PostgresAppointmentRepository) Update(appointment *domain.Appointment) error {
	query := `UPDATE appointments SET patient_id = $1, service_id = $2, appointment_date = $3, 
			  status = $4, notes = $5, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = $6`

	result, err := r.db.Exec(query, appointment.PatientID, 1,
		appointment.Date, appointment.Status, appointment.Notes, appointment.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("запись с ID %d не найдена", appointment.ID)
	}

	return nil
}

func (r *PostgresAppointmentRepository) Delete(id int) error {
	query := `DELETE FROM appointments WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("запись с ID %d не найдена", id)
	}

	return nil
}

func (r *PostgresAppointmentRepository) GetByPatientID(patientID int) ([]*domain.Appointment, error) {
	query := `SELECT id, patient_id, service_id, appointment_date, status, notes, created_at, updated_at 
			  FROM appointments WHERE patient_id = $1 ORDER BY appointment_date DESC`

	rows, err := r.db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		appointment := &domain.Appointment{}
		var serviceID int
		err := rows.Scan(
			&appointment.ID, &appointment.PatientID, &serviceID,
			&appointment.Date, &appointment.Status, &appointment.Notes,
			&appointment.CreatedAt, &appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// TODO: Получить название услуги по serviceID
		appointment.Service = "Услуга"
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *PostgresAppointmentRepository) GetByDate(date time.Time) ([]*domain.Appointment, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return r.GetByDateRange(startOfDay, endOfDay)
}

func (r *PostgresAppointmentRepository) CheckTimeConflict(date time.Time, timeStr string, excludeID int) (bool, error) {
	query := `SELECT COUNT(*) FROM appointments 
			  WHERE appointment_date = $1 AND id != $2`

	var count int
	err := r.db.QueryRow(query, date, excludeID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
