package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
)

type AppointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) Create(appointment *domain.Appointment) error {
	// Получаем service_id по названию услуги
	var serviceID int
	serviceQuery := `SELECT id FROM services WHERE name = $1 AND deleted_at IS NULL LIMIT 1`
	err := r.db.QueryRow(serviceQuery, appointment.Service).Scan(&serviceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("услуга '%s' не найдена", appointment.Service)
		}
		return err
	}

	// Получаем doctor_id по имени врача
	var doctorID sql.NullInt64
	if appointment.Doctor != "" {
		doctorQuery := `SELECT id FROM doctors WHERE name = $1 AND deleted_at IS NULL LIMIT 1`
		var dID int
		err = r.db.QueryRow(doctorQuery, appointment.Doctor).Scan(&dID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == nil {
			doctorID.Int64 = int64(dID)
			doctorID.Valid = true
		}
	}

	query := `INSERT INTO appointments (patient_id, service_id, doctor_id, appointment_date, status, price, duration_minutes, notes)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(query, appointment.PatientID, serviceID, doctorID,
		appointment.Date, appointment.Status, appointment.Price, appointment.Duration, appointment.Notes).
		Scan(&appointment.ID, &appointment.CreatedAt, &appointment.UpdatedAt)

	return err
}

func (r *AppointmentRepository) GetByID(id int) (*domain.Appointment, error) {
	query := `SELECT a.id, a.patient_id, a.appointment_date, a.status, a.price, a.duration_minutes, a.notes, a.created_at, a.updated_at,
			  s.name as service_name, p.name as patient_name, d.name as doctor_name
			  FROM appointments a
			  LEFT JOIN services s ON a.service_id = s.id AND s.deleted_at IS NULL
			  LEFT JOIN patients p ON a.patient_id = p.id AND p.deleted_at IS NULL
			  LEFT JOIN doctors d ON a.doctor_id = d.id AND d.deleted_at IS NULL
			  WHERE a.id = $1 AND a.deleted_at IS NULL`

	appointment := &domain.Appointment{}
	var serviceName sql.NullString
	var patientName sql.NullString
	var doctorName sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&appointment.ID, &appointment.PatientID,
		&appointment.Date, &appointment.Status, &appointment.Price, &appointment.Duration, &appointment.Notes,
		&appointment.CreatedAt, &appointment.UpdatedAt, &serviceName, &patientName, &doctorName,
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

	if patientName.Valid {
		appointment.PatientName = patientName.String
	} else {
		appointment.PatientName = "Неизвестно"
	}

	if doctorName.Valid {
		appointment.Doctor = doctorName.String
	}

	appointment.Time = appointment.Date.Format("15:04")

	return appointment, nil
}

func (r *AppointmentRepository) GetAll() ([]*domain.Appointment, error) {
	query := `SELECT a.id, a.patient_id, a.appointment_date, a.status, a.price, a.duration_minutes, a.notes, a.created_at, a.updated_at,
			  s.name as service_name, p.name as patient_name, d.name as doctor_name
			  FROM appointments a
			  LEFT JOIN services s ON a.service_id = s.id AND s.deleted_at IS NULL
			  LEFT JOIN patients p ON a.patient_id = p.id AND p.deleted_at IS NULL
			  LEFT JOIN doctors d ON a.doctor_id = d.id AND d.deleted_at IS NULL
			  WHERE a.deleted_at IS NULL
			  ORDER BY a.appointment_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		appointment := &domain.Appointment{}
		var serviceName sql.NullString
		var patientName sql.NullString
		var doctorName sql.NullString
		err := rows.Scan(
			&appointment.ID, &appointment.PatientID,
			&appointment.Date, &appointment.Status, &appointment.Price, &appointment.Duration, &appointment.Notes,
			&appointment.CreatedAt, &appointment.UpdatedAt, &serviceName, &patientName, &doctorName,
		)
		if err != nil {
			return nil, err
		}

		if serviceName.Valid {
			appointment.Service = serviceName.String
		} else {
			appointment.Service = "Неизвестная услуга"
		}

		if patientName.Valid {
			appointment.PatientName = patientName.String
		} else {
			appointment.PatientName = "Неизвестно"
		}

		if doctorName.Valid {
			appointment.Doctor = doctorName.String
		}

		appointment.Time = appointment.Date.Format("15:04")

		appointments = append(appointments, appointment)
	}

	return appointments, rows.Err()
}

func (r *AppointmentRepository) GetByDateRange(startDate, endDate time.Time) ([]*domain.Appointment, error) {
	query := `SELECT a.id, a.patient_id, a.appointment_date, a.status, a.price, a.duration_minutes, a.notes, a.created_at, a.updated_at,
			  s.name as service_name, p.name as patient_name, d.name as doctor_name
			  FROM appointments a
			  LEFT JOIN services s ON a.service_id = s.id AND s.deleted_at IS NULL
			  LEFT JOIN patients p ON a.patient_id = p.id AND p.deleted_at IS NULL
			  LEFT JOIN doctors d ON a.doctor_id = d.id AND d.deleted_at IS NULL
			  WHERE a.deleted_at IS NULL AND a.appointment_date BETWEEN $1 AND $2
			  ORDER BY a.appointment_date`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		appointment := &domain.Appointment{}
		var serviceName sql.NullString
		var patientName sql.NullString
		var doctorName sql.NullString
		err := rows.Scan(
			&appointment.ID, &appointment.PatientID,
			&appointment.Date, &appointment.Status, &appointment.Price, &appointment.Duration, &appointment.Notes,
			&appointment.CreatedAt, &appointment.UpdatedAt, &serviceName, &patientName, &doctorName,
		)
		if err != nil {
			return nil, err
		}

		if serviceName.Valid {
			appointment.Service = serviceName.String
		} else {
			appointment.Service = "Неизвестная услуга"
		}

		if patientName.Valid {
			appointment.PatientName = patientName.String
		} else {
			appointment.PatientName = "Неизвестно"
		}

		if doctorName.Valid {
			appointment.Doctor = doctorName.String
		}

		appointment.Time = appointment.Date.Format("15:04")
		appointments = append(appointments, appointment)
	}

	return appointments, rows.Err()
}

func (r *AppointmentRepository) Update(appointment *domain.Appointment) error {
	// Lookup service_id by service name
	var serviceID int
	serviceQuery := `SELECT id FROM services WHERE name = $1 AND deleted_at IS NULL LIMIT 1`
	err := r.db.QueryRow(serviceQuery, appointment.Service).Scan(&serviceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("услуга '%s' не найдена", appointment.Service)
		}
		return err
	}

	// Lookup doctor_id by doctor name
	var doctorID sql.NullInt64
	if appointment.Doctor != "" {
		doctorQuery := `SELECT id FROM doctors WHERE name = $1 AND deleted_at IS NULL LIMIT 1`
		var dID int
		err = r.db.QueryRow(doctorQuery, appointment.Doctor).Scan(&dID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == nil {
			doctorID.Int64 = int64(dID)
			doctorID.Valid = true
		}
	}

	query := `UPDATE appointments SET patient_id = $1, service_id = $2, doctor_id = $3, appointment_date = $4,
			  status = $5, notes = $6, price = $7, duration_minutes = $8, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $9 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, appointment.PatientID, serviceID, doctorID,
		appointment.Date, appointment.Status, appointment.Notes, appointment.Price, appointment.Duration, appointment.ID)
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

func (r *AppointmentRepository) Delete(id int) error {
	query := `UPDATE appointments SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

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

func (r *AppointmentRepository) GetByPatientID(patientID int) ([]*domain.Appointment, error) {
	query := `SELECT a.id, a.patient_id, a.appointment_date, a.status, a.price, a.duration_minutes, a.notes, a.created_at, a.updated_at,
			  s.name as service_name, p.name as patient_name, d.name as doctor_name
			  FROM appointments a
			  LEFT JOIN services s ON a.service_id = s.id AND s.deleted_at IS NULL
			  LEFT JOIN patients p ON a.patient_id = p.id AND p.deleted_at IS NULL
			  LEFT JOIN doctors d ON a.doctor_id = d.id AND d.deleted_at IS NULL
			  WHERE a.deleted_at IS NULL AND a.patient_id = $1
			  ORDER BY a.appointment_date DESC`

	rows, err := r.db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment
	for rows.Next() {
		appointment := &domain.Appointment{}
		var serviceName sql.NullString
		var patientName sql.NullString
		var doctorName sql.NullString
		err := rows.Scan(
			&appointment.ID, &appointment.PatientID,
			&appointment.Date, &appointment.Status, &appointment.Price, &appointment.Duration, &appointment.Notes,
			&appointment.CreatedAt, &appointment.UpdatedAt, &serviceName, &patientName, &doctorName,
		)
		if err != nil {
			return nil, err
		}

		if serviceName.Valid {
			appointment.Service = serviceName.String
		} else {
			appointment.Service = "Неизвестная услуга"
		}

		if patientName.Valid {
			appointment.PatientName = patientName.String
		} else {
			appointment.PatientName = "Неизвестно"
		}

		if doctorName.Valid {
			appointment.Doctor = doctorName.String
		}

		appointment.Time = appointment.Date.Format("15:04")
		appointments = append(appointments, appointment)
	}

	return appointments, rows.Err()
}

func (r *AppointmentRepository) GetByDate(date time.Time) ([]*domain.Appointment, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return r.GetByDateRange(startOfDay, endOfDay)
}

func (r *AppointmentRepository) CheckTimeConflict(date time.Time, timeStr string, excludeID int) (bool, error) {
	query := `SELECT COUNT(*) FROM appointments
			  WHERE deleted_at IS NULL AND appointment_date = $1 AND id != $2`

	var count int
	err := r.db.QueryRow(query, date, excludeID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
