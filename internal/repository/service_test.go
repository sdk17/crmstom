//go:build integration

package repository

import (
	"context"
	"testing"

	"github.com/sdk17/crmstom/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	testDB, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer testDB.Teardown(ctx)

	repo := NewServiceRepository(testDB.DB)

	t.Run("Create", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		service := &domain.Service{
			Name:  "Dental Cleaning",
			Type:  "Hygiene",
			Notes: "Professional teeth cleaning",
		}

		err = repo.Create(service)
		require.NoError(t, err)
		assert.Greater(t, service.ID, 0)
		assert.False(t, service.CreatedAt.IsZero())
		assert.False(t, service.UpdatedAt.IsZero())
	})

	t.Run("GetByID", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		service := &domain.Service{
			Name:  "Root Canal",
			Type:  "Treatment",
			Notes: "Endodontic treatment",
		}
		err = repo.Create(service)
		require.NoError(t, err)

		found, err := repo.GetByID(service.ID)
		require.NoError(t, err)
		assert.Equal(t, service.ID, found.ID)
		assert.Equal(t, service.Name, found.Name)
		assert.Equal(t, service.Type, found.Type)
		assert.Equal(t, service.Notes, found.Notes)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		_, err = repo.GetByID(9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "не найдена")
	})

	t.Run("GetAll", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		services := []*domain.Service{
			{Name: "Service A", Type: "Type1"},
			{Name: "Service B", Type: "Type2"},
			{Name: "Service C", Type: "Type1"},
		}

		for _, s := range services {
			err := repo.Create(s)
			require.NoError(t, err)
		}

		all, err := repo.GetAll()
		require.NoError(t, err)
		assert.Len(t, all, 3)
	})

	t.Run("Update", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		service := &domain.Service{
			Name:  "Original Service",
			Type:  "Original Type",
			Notes: "Original notes",
		}
		err = repo.Create(service)
		require.NoError(t, err)

		service.Name = "Updated Service"
		service.Type = "Updated Type"
		service.Notes = "Updated notes"
		err = repo.Update(service)
		require.NoError(t, err)

		found, err := repo.GetByID(service.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Service", found.Name)
		assert.Equal(t, "Updated Type", found.Type)
		assert.Equal(t, "Updated notes", found.Notes)
	})

	t.Run("Delete", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		service := &domain.Service{
			Name: "To Delete",
			Type: "Temporary",
		}
		err = repo.Create(service)
		require.NoError(t, err)

		err = repo.Delete(service.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(service.ID)
		assert.Error(t, err)
	})

	t.Run("GetByCategory", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		services := []*domain.Service{
			{Name: "Cleaning 1", Type: "Hygiene"},
			{Name: "Cleaning 2", Type: "Hygiene"},
			{Name: "Extraction", Type: "Surgery"},
			{Name: "Filling", Type: "Treatment"},
		}

		for _, s := range services {
			err := repo.Create(s)
			require.NoError(t, err)
		}

		hygieneServices, err := repo.GetByCategory("Hygiene")
		require.NoError(t, err)
		assert.Len(t, hygieneServices, 2)

		surgeryServices, err := repo.GetByCategory("Surgery")
		require.NoError(t, err)
		assert.Len(t, surgeryServices, 1)
		assert.Equal(t, "Extraction", surgeryServices[0].Name)
	})

	t.Run("Search", func(t *testing.T) {
		err := testDB.TruncateTables(ctx)
		require.NoError(t, err)

		services := []*domain.Service{
			{Name: "Professional Cleaning", Type: "Hygiene", Notes: "Deep cleaning service"},
			{Name: "Basic Cleaning", Type: "Hygiene", Notes: "Standard cleaning"},
			{Name: "Root Canal", Type: "Treatment", Notes: "Complex procedure"},
		}

		for _, s := range services {
			err := repo.Create(s)
			require.NoError(t, err)
		}

		// Search by name
		results, err := repo.Search("Cleaning")
		require.NoError(t, err)
		assert.Len(t, results, 2)

		// Search by notes
		results, err = repo.Search("Complex")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Root Canal", results[0].Name)

		// Case insensitive search
		results, err = repo.Search("cleaning")
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})
}
