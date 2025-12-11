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

func TestPatientRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	testDB, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer testDB.Teardown(ctx)

	repo := NewPatientRepository(testDB.DB)

	t.Run("Create", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := &domain.Patient{
			Name:      "John Doe",
			Phone:     "+7 777 123 4567",
			Email:     "john@example.com",
			BirthDate: time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
			Address:   "123 Main St",
		}

		err = repo.Create(patient)
		require.NoError(t, err)
		assert.Greater(t, patient.ID, 0)
		assert.False(t, patient.CreatedAt.IsZero())
		assert.False(t, patient.UpdatedAt.IsZero())
	})

	t.Run("GetByID", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := &domain.Patient{
			Name:      "Jane Doe",
			Phone:     "+7 777 234 5678",
			Email:     "jane@example.com",
			BirthDate: time.Date(1985, 3, 20, 0, 0, 0, 0, time.UTC),
			Address:   "456 Oak Ave",
		}
		err = repo.Create(patient)
		require.NoError(t, err)

		found, err := repo.GetByID(patient.ID)
		require.NoError(t, err)
		assert.Equal(t, patient.ID, found.ID)
		assert.Equal(t, patient.Name, found.Name)
		assert.Equal(t, patient.Phone, found.Phone)
		assert.Equal(t, patient.Email, found.Email)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		_, err = repo.GetByID(9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найден")
	})

	t.Run("GetAll", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patients := []*domain.Patient{
			{Name: "Patient 1", Phone: "+7 777 111 1111"},
			{Name: "Patient 2", Phone: "+7 777 222 2222"},
			{Name: "Patient 3", Phone: "+7 777 333 3333"},
		}

		for _, p := range patients {
			err := repo.Create(p)
			require.NoError(t, err)
		}

		all, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, all, 3)
	})

	t.Run("Update", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := &domain.Patient{
			Name:  "Original Name",
			Phone: "+7 777 444 4444",
		}
		err = repo.Create(patient)
		require.NoError(t, err)

		patient.Name = "Updated Name"
		patient.Email = "updated@example.com"
		err = repo.Update(patient)
		require.NoError(t, err)

		found, err := repo.GetByID(patient.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
		assert.Equal(t, "updated@example.com", found.Email)
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := &domain.Patient{
			ID:    9999,
			Name:  "Non-existent",
			Phone: "+7 777 555 5555",
		}
		err = repo.Update(patient)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найден")
	})

	t.Run("Delete", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := &domain.Patient{
			Name:  "To Delete",
			Phone: "+7 777 666 6666",
		}
		err = repo.Create(patient)
		require.NoError(t, err)

		err = repo.Delete(patient.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(patient.ID)
		assert.Error(t, err)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		err = repo.Delete(9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найден")
	})

	t.Run("GetByPhone", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patient := &domain.Patient{
			Name:  "Phone Test",
			Phone: "+7 777 777 7777",
		}
		err = repo.Create(patient)
		require.NoError(t, err)

		found, err := repo.GetByPhone("+7 777 777 7777")
		require.NoError(t, err)
		assert.Equal(t, patient.ID, found.ID)
		assert.Equal(t, patient.Name, found.Name)
	})

	t.Run("GetByPhone_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		_, err = repo.GetByPhone("+7 777 999 9999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найден")
	})

	t.Run("Search", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		patients := []*domain.Patient{
			{Name: "John Smith", Phone: "+7 777 100 0001", Email: "john@test.com"},
			{Name: "Jane Smith", Phone: "+7 777 100 0002", Email: "jane@test.com"},
			{Name: "Bob Johnson", Phone: "+7 777 100 0003", Email: "bob@other.com"},
		}

		for _, p := range patients {
			err := repo.Create(p)
			require.NoError(t, err)
		}

		// Search by name
		results, err := repo.Search("Smith")
		require.NoError(t, err)
		assert.Len(t, results, 2)

		// Search by email domain
		results, err = repo.Search("test.com")
		require.NoError(t, err)
		assert.Len(t, results, 2)

		// Search by phone
		results, err = repo.Search("100 0003")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Bob Johnson", results[0].Name)
	})
}
