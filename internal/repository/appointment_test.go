//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/sdk17/crmstom/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppointmentRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	testDB, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer testDB.Teardown(ctx)

	patientRepo := NewPatientRepository(testDB.DB)
	serviceRepo := NewServiceRepository(testDB.DB)
	appointmentRepo := NewAppointmentRepository(testDB.DB)

	// Helper to create test patient
	createTestPatient := func(t *testing.T, name string) *domain.Patient {
		patient := &domain.Patient{
			Name:  name,
			Phone: "+7 777 000 0000",
		}
		err := patientRepo.Create(patient)
		require.NoError(t, err)
		return patient
	}

	// Helper to create test service
	createTestService := func(t *testing.T, name string) *domain.Service {
		service := &domain.Service{
			Name: name,
			Type: "Treatment",
		}
		err := serviceRepo.Create(service)
		require.NoError(t, err)
		return service
	}

	t.Run("Create", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "John Doe")
		service := createTestService(t, "Dental Cleaning")

		appointment := &domain.Appointment{
			PatientID: patient.ID,
			Service:   service.Name,
			Date:      time.Now().Add(24 * time.Hour),
			Status:    domain.StatusScheduled,
			Price:     100.00,
			Duration:  60,
			Notes:     "First appointment",
		}

		err = appointmentRepo.Create(appointment)
		require.NoError(t, err)
		assert.Greater(t, appointment.ID, 0)
		assert.False(t, appointment.CreatedAt.IsZero())
		assert.False(t, appointment.UpdatedAt.IsZero())
	})

	t.Run("Create_ServiceNotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "Jane Doe")

		appointment := &domain.Appointment{
			PatientID: patient.ID,
			Service:   "Non-existent Service",
			Date:      time.Now(),
			Status:    domain.StatusScheduled,
		}

		err = appointmentRepo.Create(appointment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найдена")
	})

	t.Run("GetByID", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "Patient GetByID")
		service := createTestService(t, "Root Canal")

		appointmentDate := time.Date(2024, 12, 15, 14, 30, 0, 0, time.UTC)
		appointment := &domain.Appointment{
			PatientID: patient.ID,
			Service:   service.Name,
			Date:      appointmentDate,
			Status:    domain.StatusScheduled,
			Price:     250.00,
			Duration:  90,
			Notes:     "Complex procedure",
		}
		err = appointmentRepo.Create(appointment)
		require.NoError(t, err)

		found, err := appointmentRepo.GetByID(appointment.ID)
		require.NoError(t, err)
		assert.Equal(t, appointment.ID, found.ID)
		assert.Equal(t, appointment.PatientID, found.PatientID)
		assert.Equal(t, service.Name, found.Service)
		assert.Equal(t, appointment.Price, found.Price)
		assert.Equal(t, appointment.Duration, found.Duration)
		// Verify patient_name and time are populated
		assert.Equal(t, "Patient GetByID", found.PatientName)
		assert.Equal(t, "14:30", found.Time)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		_, err = appointmentRepo.GetByID(9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найдена")
	})

	t.Run("GetAll", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "Patient GetAll")
		service1 := createTestService(t, "Cleaning")
		service2 := createTestService(t, "Filling")

		baseDate := time.Date(2024, 12, 15, 9, 0, 0, 0, time.UTC)
		appointments := []*domain.Appointment{
			{PatientID: patient.ID, Service: service1.Name, Date: baseDate, Status: domain.StatusScheduled},
			{PatientID: patient.ID, Service: service2.Name, Date: baseDate.Add(24 * time.Hour), Status: domain.StatusScheduled},
			{PatientID: patient.ID, Service: service1.Name, Date: baseDate.Add(48 * time.Hour), Status: domain.StatusCompleted},
		}

		for _, a := range appointments {
			err := appointmentRepo.Create(a)
			require.NoError(t, err)
		}

		all, err := appointmentRepo.GetAll()
		require.NoError(t, err)
		assert.Len(t, all, 3)

		// Verify patient_name and time are populated for all appointments
		for _, a := range all {
			assert.Equal(t, "Patient GetAll", a.PatientName)
			assert.Equal(t, "09:00", a.Time)
		}
	})

	t.Run("Update", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "Patient Update")
		service := createTestService(t, "Check-up")

		appointment := &domain.Appointment{
			PatientID: patient.ID,
			Service:   service.Name,
			Date:      time.Now().Add(24 * time.Hour),
			Status:    domain.StatusScheduled,
			Price:     100.00,
			Duration:  30,
			Notes:     "Original note",
		}
		err = appointmentRepo.Create(appointment)
		require.NoError(t, err)

		appointment.Status = domain.StatusCompleted
		appointment.Notes = "Updated note"
		appointment.Price = 200.00
		appointment.Duration = 60
		err = appointmentRepo.Update(appointment)
		require.NoError(t, err)

		found, err := appointmentRepo.GetByID(appointment.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.StatusCompleted, found.Status)
		assert.Equal(t, "Updated note", found.Notes)
		assert.Equal(t, 200.00, found.Price)
		assert.Equal(t, 60, found.Duration)
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		service := createTestService(t, "Test Service")

		appointment := &domain.Appointment{
			ID:        9999,
			PatientID: 1,
			Date:      time.Now(),
			Status:    domain.StatusScheduled,
			Service:   service.Name,
		}
		err = appointmentRepo.Update(appointment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найдена")
	})

	t.Run("Update_ServiceNotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "Patient ServiceNotFound")
		service := createTestService(t, "Original Service")

		appointment := &domain.Appointment{
			PatientID: patient.ID,
			Service:   service.Name,
			Date:      time.Now().Add(24 * time.Hour),
			Status:    domain.StatusScheduled,
		}
		err = appointmentRepo.Create(appointment)
		require.NoError(t, err)

		// Try to update with non-existent service
		appointment.Service = "Non-existent Service"
		err = appointmentRepo.Update(appointment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найдена")
	})

	t.Run("Delete", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "Patient Delete")
		service := createTestService(t, "X-Ray")

		appointment := &domain.Appointment{
			PatientID: patient.ID,
			Service:   service.Name,
			Date:      time.Now().Add(24 * time.Hour),
			Status:    domain.StatusScheduled,
		}
		err = appointmentRepo.Create(appointment)
		require.NoError(t, err)

		err = appointmentRepo.Delete(appointment.ID)
		require.NoError(t, err)

		_, err = appointmentRepo.GetByID(appointment.ID)
		assert.Error(t, err)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		err = appointmentRepo.Delete(9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найдена")
	})

	t.Run("GetByPatientID", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient1 := createTestPatient(t, "Patient1")
		patient2 := createTestPatient(t, "Patient2")
		service := createTestService(t, "Consultation")

		baseDate := time.Date(2024, 12, 15, 10, 30, 0, 0, time.UTC)

		// Create appointments for patient1
		for i := 0; i < 3; i++ {
			a := &domain.Appointment{
				PatientID: patient1.ID,
				Service:   service.Name,
				Date:      baseDate.Add(time.Duration(i*24) * time.Hour),
				Status:    domain.StatusScheduled,
			}
			err := appointmentRepo.Create(a)
			require.NoError(t, err)
		}

		// Create appointment for patient2
		a := &domain.Appointment{
			PatientID: patient2.ID,
			Service:   service.Name,
			Date:      baseDate.Add(24 * time.Hour),
			Status:    domain.StatusScheduled,
		}
		err = appointmentRepo.Create(a)
		require.NoError(t, err)

		// Get patient1's appointments
		appointments, err := appointmentRepo.GetByPatientID(patient1.ID)
		require.NoError(t, err)
		assert.Len(t, appointments, 3)
		for _, appt := range appointments {
			assert.Equal(t, "Patient1", appt.PatientName)
			assert.Equal(t, "10:30", appt.Time)
		}

		// Get patient2's appointments
		appointments, err = appointmentRepo.GetByPatientID(patient2.ID)
		require.NoError(t, err)
		assert.Len(t, appointments, 1)
		assert.Equal(t, "Patient2", appointments[0].PatientName)
		assert.Equal(t, "10:30", appointments[0].Time)
	})

	t.Run("GetByDate", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "Patient Date")
		service := createTestService(t, "Appointment Service")

		today := time.Date(2024, 12, 15, 11, 0, 0, 0, time.UTC)
		tomorrow := today.Add(24 * time.Hour)

		// Create appointments for today
		for i := 0; i < 2; i++ {
			a := &domain.Appointment{
				PatientID: patient.ID,
				Service:   service.Name,
				Date:      today.Add(time.Duration(i) * time.Hour),
				Status:    domain.StatusScheduled,
			}
			err := appointmentRepo.Create(a)
			require.NoError(t, err)
		}

		// Create appointment for tomorrow
		a := &domain.Appointment{
			PatientID: patient.ID,
			Service:   service.Name,
			Date:      tomorrow,
			Status:    domain.StatusScheduled,
		}
		err = appointmentRepo.Create(a)
		require.NoError(t, err)

		// Get today's appointments
		todayAppointments, err := appointmentRepo.GetByDate(today)
		require.NoError(t, err)
		assert.Len(t, todayAppointments, 2)
		// Verify patient_name and time are populated
		for _, appt := range todayAppointments {
			assert.Equal(t, "Patient Date", appt.PatientName)
			assert.NotEmpty(t, appt.Time)
		}

		// Get tomorrow's appointments
		tomorrowAppointments, err := appointmentRepo.GetByDate(tomorrow)
		require.NoError(t, err)
		assert.Len(t, tomorrowAppointments, 1)
		assert.Equal(t, "Patient Date", tomorrowAppointments[0].PatientName)
		assert.Equal(t, "11:00", tomorrowAppointments[0].Time)
	})

	t.Run("CheckTimeConflict", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "Patient Conflict")
		service := createTestService(t, "Conflict Service")

		appointmentDate := time.Date(2024, 12, 15, 10, 0, 0, 0, time.UTC)

		appointment := &domain.Appointment{
			PatientID: patient.ID,
			Service:   service.Name,
			Date:      appointmentDate,
			Status:    domain.StatusScheduled,
		}
		err = appointmentRepo.Create(appointment)
		require.NoError(t, err)

		// Check for conflict at the same time (excluding none)
		hasConflict, err := appointmentRepo.CheckTimeConflict(appointmentDate, "10:00", 0)
		require.NoError(t, err)
		assert.True(t, hasConflict)

		// Check for conflict excluding the same appointment
		hasConflict, err = appointmentRepo.CheckTimeConflict(appointmentDate, "10:00", appointment.ID)
		require.NoError(t, err)
		assert.False(t, hasConflict)

		// Check for no conflict at a different time
		differentDate := time.Date(2024, 12, 16, 10, 0, 0, 0, time.UTC)
		hasConflict, err = appointmentRepo.CheckTimeConflict(differentDate, "10:00", 0)
		require.NoError(t, err)
		assert.False(t, hasConflict)
	})
}
