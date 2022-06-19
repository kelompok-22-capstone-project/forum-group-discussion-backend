package user

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/config"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	var repo UserRepository = NewUserRepositoryImpl(db)

	testCases := []struct {
		name          string
		inputUser     entity.User
		expectedError error
		mockBehaviour func()
	}{
		{
			name: "it should return ErrDatabase, when database return an error",
			inputUser: entity.User{
				ID:       "u-gXyZpw",
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
				Role:     "user",
			},
			expectedError: repository.ErrDatabase,
			mockBehaviour: func() {
				mock.ExpectExec(".*").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(repository.ErrDatabase)
			},
		},
		{
			name: "it should return ErrRecordAlreadyExists, when database return an error",
			inputUser: entity.User{
				ID:       "u-gXyZpw",
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
				Role:     "user",
			},
			expectedError: repository.ErrRecordAlreadyExists,
			mockBehaviour: func() {
				mock.ExpectExec(".*").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(&pq.Error{Code: "23505"})
			},
		},
		{
			name: "it should return ErrRecordAlreadyExists, when database return an error",
			inputUser: entity.User{
				ID:       "u-gXyZpw",
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
				Role:     "user",
			},
			expectedError: repository.ErrRecordAlreadyExists,
			mockBehaviour: func() {
				mock.ExpectExec(".*").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnError(&pq.Error{Code: "23505"})
			},
		},
		{
			name: "it should return ErrDatabase, when rows affected less than 1",
			inputUser: entity.User{
				ID:       "u-gXyZpw",
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
				Role:     "user",
			},
			expectedError: repository.ErrDatabase,
			mockBehaviour: func() {
				mock.ExpectExec(".*").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(1, 0))
			},
		},
		{
			name: "it should return nil error, when successfully inserted the data to the database",
			inputUser: entity.User{
				ID:       "u-gXyZpw",
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
				Role:     "user",
			},
			expectedError: nil,
			mockBehaviour: func() {
				mock.ExpectExec(".*").WithArgs(
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			gotError := repo.Insert(context.Background(), testCase.inputUser)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, gotError)
			} else {
				assert.NoError(t, gotError)
			}
		})
	}
}

func TestFindByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	var repo UserRepository = NewUserRepositoryImpl(db)

	testCases := []struct {
		name          string
		inputUsername string
		expectedUser  entity.User
		expectedError error
		mockBehaviour func()
	}{
		{
			name:          "it should return ErrDatabase, when database return an error",
			inputUsername: "erikrios",
			expectedUser:  entity.User{},
			expectedError: repository.ErrDatabase,
			mockBehaviour: func() {
				mock.ExpectQuery(".*").WithArgs(
					sqlmock.AnyArg(),
				).WillReturnError(repository.ErrDatabase)
			},
		},
		{
			name:          "it should return ErrRecordNotFound, when database return an error",
			inputUsername: "erikrios",
			expectedUser:  entity.User{},
			expectedError: repository.ErrRecordNotFound,
			mockBehaviour: func() {
				mock.ExpectQuery(".*").WithArgs(
					sqlmock.AnyArg(),
				).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:          "it should return valid user, when database successfully return the data",
			inputUsername: "erikrios",
			expectedUser: entity.User{
				ID:       "u-gXyZpw",
				Username: "erikrios",
				Email:    "erikriosetiawan15@gmail.com",
				Name:     "Erik Rio Setiawan",
				Password: "erikriosetiawan",
				Role:     "user",
				IsActive: true,
			},
			expectedError: nil,
			mockBehaviour: func() {
				returnedRows := sqlmock.NewRows([]string{"id", "username", "email", "name", "password", "role", "is_active", "created_at", "updated_at"})
				returnedRows.AddRow("u-gXyZpw", "erikrios", "erikriosetiawan15@gmail.com", "Erik Rio Setiawan", "erikriosetiawan", "user", true, time.Now(), time.Now())
				mock.ExpectQuery(".*").WithArgs(
					sqlmock.AnyArg(),
				).WillReturnRows(returnedRows)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour()

			gotUser, gotError := repo.FindByUsername(context.Background(), testCase.inputUsername)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, gotError)
			} else {
				assert.NoError(t, gotError)
				assert.Equal(t, testCase.expectedUser.ID, gotUser.ID)
				assert.Equal(t, testCase.expectedUser.Name, gotUser.Name)
				assert.Equal(t, testCase.expectedUser.Email, gotUser.Email)
				assert.Equal(t, testCase.expectedUser.Name, gotUser.Name)
				assert.Equal(t, testCase.expectedUser.Password, gotUser.Password)
				assert.Equal(t, testCase.expectedUser.Role, gotUser.Role)
				assert.Equal(t, testCase.expectedUser.IsActive, gotUser.IsActive)
			}
		})
	}
}

func TestFindAllWithStatusAndPagination(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load("./../../.env.example")
	if err != nil {
		panic(err)
	}

	db, err := config.NewPostgreSQLDatabase()
	if err != nil {
		panic(err)
	}

	var repo UserRepository = NewUserRepositoryImpl(db)

	accessorUserID := "u-kt56R1"
	orderBy := entity.Ranking
	userStatus := entity.Active
	pageInfo := entity.PageInfo{
		Page:  1,
		Limit: 20,
	}

	if pagination, err := repo.FindAllWithStatusAndPagination(context.Background(), accessorUserID, orderBy, userStatus, pageInfo); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Pagination: %+v", pagination)
	}
}

func TestFindByUsernameWithAccessor(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load("./../../.env.example")
	if err != nil {
		panic(err)
	}

	db, err := config.NewPostgreSQLDatabase()
	if err != nil {
		panic(err)
	}

	var repo UserRepository = NewUserRepositoryImpl(db)

	accessorUserID := "u-kt56R1"
	username := "erikrios"
	if user, err := repo.FindByUsernameWithAccessor(context.Background(), accessorUserID, username); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("User: %+v", user)
	}
}

func TestBannedUser(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load("./../../.env.example")
	if err != nil {
		panic(err)
	}

	db, err := config.NewPostgreSQLDatabase()
	if err != nil {
		panic(err)
	}

	var repo UserRepository = NewUserRepositoryImpl(db)

	userID := "u-NplCrv"

	if err := repo.BannedUser(context.Background(), userID); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully banned a user with ID %s", userID)
	}
}
