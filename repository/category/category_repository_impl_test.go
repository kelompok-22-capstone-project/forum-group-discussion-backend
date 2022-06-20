package category

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	var repo CategoryRepository = NewCategoryRepositoryImpl(db)

	now := time.Now()

	testCases := []struct {
		name               string
		expectedError      error
		expectedCategories []entity.Category
		mockBehaviour      func()
	}{
		{
			name:               "it should return repository.ErrDatabase, when database return an error",
			expectedError:      repository.ErrDatabase,
			expectedCategories: []entity.Category{},
			mockBehaviour: func() {
				dbMock.ExpectQuery(".*").WillReturnError(errors.New("something went wrong with the databases."))
			},
		},
		{
			name:          "it should return valid categories, when database return nil error",
			expectedError: nil,
			expectedCategories: []entity.Category{
				{
					ID:          "c-xyz",
					Name:        "Tech",
					Description: "This is tech description.",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			mockBehaviour: func() {
				returnedRows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"})
				returnedRows.AddRow("c-xyz", "Tech", "This is tech description.", now, now)
				dbMock.ExpectQuery(".*").WillReturnRows(returnedRows)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			gotCategories, gotErr := repo.FindAll(context.Background())
			if testCase.expectedError != nil {
				assert.ErrorIs(t, gotErr, testCase.expectedError)
			} else {
				assert.ElementsMatch(t, gotCategories, testCase.expectedCategories)
			}
		})
	}
}
