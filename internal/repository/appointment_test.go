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

		appointment := &domain.Appointment{
			PatientID: patient.ID,
			Service:   service.Name,
			Date:      time.Now().Add(48 * time.Hour),
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

		appointments := []*domain.Appointment{
			{PatientID: patient.ID, Service: service1.Name, Date: time.Now().Add(24 * time.Hour), Status: domain.StatusScheduled},
			{PatientID: patient.ID, Service: service2.Name, Date: time.Now().Add(48 * time.Hour), Status: domain.StatusScheduled},
			{PatientID: patient.ID, Service: service1.Name, Date: time.Now().Add(72 * time.Hour), Status: domain.StatusCompleted},
		}

		for _, a := range appointments {
			err := appointmentRepo.Create(a)
			require.NoError(t, err)
		}

		all, err := appointmentRepo.GetAll()
		require.NoError(t, err)
		assert.Len(t, all, 3)
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
			Notes:     "Original note",
		}
		err = appointmentRepo.Create(appointment)
		require.NoError(t, err)

		appointment.Status = domain.StatusCompleted
		appointment.Notes = "Updated note"
		err = appointmentRepo.Update(appointment)
		require.NoError(t, err)

		found, err := appointmentRepo.GetByID(appointment.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.StatusCompleted, found.Status)
		assert.Equal(t, "Updated note", found.Notes)
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		appointment := &domain.Appointment{
			ID:        9999,
			PatientID: 1,
			Date:      time.Now(),
			Status:    domain.StatusScheduled,
		}
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

		// Create appointments for patient1
		for i := 0; i < 3; i++ {
			a := &domain.Appointment{
				PatientID: patient1.ID,
				Service:   service.Name,
				Date:      time.Now().Add(time.Duration(i*24) * time.Hour),
				Status:    domain.StatusScheduled,
			}
			err := appointmentRepo.Create(a)
			require.NoError(t, err)
		}

		// Create appointment for patient2
		a := &domain.Appointment{
			PatientID: patient2.ID,
			Service:   service.Name,
			Date:      time.Now().Add(24 * time.Hour),
			Status:    domain.StatusScheduled,
		}
		err = appointmentRepo.Create(a)
		require.NoError(t, err)

		// Get patient1's appointments
		appointments, err := appointmentRepo.GetByPatientID(patient1.ID)
		require.NoError(t, err)
		assert.Len(t, appointments, 3)

		// Get patient2's appointments
		appointments, err = appointmentRepo.GetByPatientID(patient2.ID)
		require.NoError(t, err)
		assert.Len(t, appointments, 1)
	})

	t.Run("GetByDate", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := createTestPatient(t, "Patient Date")
		service := createTestService(t, "Appointment Service")

		today := time.Now()
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

		// Get tomorrow's appointments
		tomorrowAppointments, err := appointmentRepo.GetByDate(tomorrow)
		require.NoError(t, err)
		assert.Len(t, tomorrowAppointments, 1)
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
