//go:build integration

package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sdk17/crmstom/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoctorRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	testDB, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer testDB.Teardown(ctx)

	repo := NewDoctorRepository(testDB.DB)

	t.Run("Create", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		doctor := &domain.Doctor{
			Name:     "Dr. John Smith",
			Email:    "john.smith@clinic.com",
			Login:    "dr_john",
			Password: "secret123",
			IsAdmin:  false,
		}

		err = repo.Create(doctor)
		require.NoError(t, err)
		assert.Greater(t, doctor.ID, 0)
	})

	t.Run("Create_DuplicateLogin", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		doctor1 := &domain.Doctor{
			Name:     "Dr. First",
			Login:    "same_login",
			Password: "pass123",
		}
		err = repo.Create(doctor1)
		require.NoError(t, err)

		doctor2 := &domain.Doctor{
			Name:     "Dr. Second",
			Login:    "same_login",
			Password: "pass456",
		}
		err = repo.Create(doctor2)
		assert.Error(t, err) // Should fail due to unique constraint
	})

	t.Run("GetByID", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		doctor := &domain.Doctor{
			Name:     "Dr. Jane Doe",
			Email:    "jane.doe@clinic.com",
			Login:    "dr_jane",
			Password: "pass789",
			IsAdmin:  true,
		}
		err = repo.Create(doctor)
		require.NoError(t, err)

		found, err := repo.GetByID(doctor.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, doctor.ID, found.ID)
		assert.Equal(t, doctor.Name, found.Name)
		assert.Equal(t, doctor.Email, found.Email)
		assert.Equal(t, doctor.Login, found.Login)
		assert.Equal(t, doctor.IsAdmin, found.IsAdmin)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		found, err := repo.GetByID(9999)
		require.NoError(t, err) // Returns nil, nil for not found
		assert.Nil(t, found)
	})

	t.Run("GetAll", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		doctors := []*domain.Doctor{
			{Name: "Dr. Alpha", Login: "alpha", Password: "pass1"},
			{Name: "Dr. Beta", Login: "beta", Password: "pass2"},
			{Name: "Dr. Gamma", Login: "gamma", Password: "pass3"},
		}

		for _, d := range doctors {
			err := repo.Create(d)
			require.NoError(t, err)
		}

		all, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, all, 3)
	})

	t.Run("Update", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		doctor := &domain.Doctor{
			Name:     "Dr. Original",
			Login:    "original",
			Password: "oldpass",
			IsAdmin:  false,
		}
		err = repo.Create(doctor)
		require.NoError(t, err)

		doctor.Name = "Dr. Updated"
		doctor.Email = "updated@clinic.com"
		doctor.IsAdmin = true
		err = repo.Update(doctor)
		require.NoError(t, err)

		found, err := repo.GetByID(doctor.ID)
		require.NoError(t, err)
		assert.Equal(t, "Dr. Updated", found.Name)
		assert.Equal(t, "updated@clinic.com", found.Email)
		assert.True(t, found.IsAdmin)
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		doctor := &domain.Doctor{
			ID:       9999,
			Name:     "Non-existent",
			Login:    "nonexistent",
			Password: "pass",
		}
		err = repo.Update(doctor)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("Delete", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		doctor := &domain.Doctor{
			Name:     "Dr. ToDelete",
			Login:    "todelete",
			Password: "pass",
		}
		err = repo.Create(doctor)
		require.NoError(t, err)

		err = repo.Delete(doctor.ID)
		require.NoError(t, err)

		found, err := repo.GetByID(doctor.ID)
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		err = repo.Delete(9999)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("GetByLogin", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		doctor := &domain.Doctor{
			Name:     "Dr. LoginTest",
			Login:    "unique_login",
			Password: "testpass",
		}
		err = repo.Create(doctor)
		require.NoError(t, err)

		found, err := repo.GetByLogin("unique_login")
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, doctor.ID, found.ID)
		assert.Equal(t, doctor.Name, found.Name)
	})

	t.Run("GetByLogin_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		found, err := repo.GetByLogin("nonexistent_login")
		require.NoError(t, err) // Returns nil, nil for not found
		assert.Nil(t, found)
	})
}
