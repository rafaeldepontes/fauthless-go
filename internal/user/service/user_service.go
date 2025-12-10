package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rafaeldepontes/fauthless-go/internal/cache"
	"github.com/rafaeldepontes/fauthless-go/internal/domain"
	"github.com/rafaeldepontes/fauthless-go/internal/errorhandler"
	"github.com/rafaeldepontes/fauthless-go/internal/pagination"
	"github.com/rafaeldepontes/fauthless-go/internal/user"
	log "github.com/sirupsen/logrus"
)

type userService struct {
	repo   user.Repository
	Logger *log.Logger
	Cache  *cache.Caches
}

// NewUserService initialize a new UserService containing a UserRepository.
func NewUserService(userRepo user.Repository, logg *log.Logger, cache *cache.Caches) user.Service {
	return &userService{
		repo:   userRepo,
		Logger: logg,
		Cache:  cache,
	}
}

func (s *userService) FindAllUsersHashedCursorPagination(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("Listing all the users in the database using the cursor pagination with a hash...")

	cs := pagination.NewCursorService[domain.User]()

	var cursor domain.CursorResquest
	json.NewDecoder(r.Body).Decode(&cursor)

	var (
		hashSrc    string
		cursorId   int64
		size       int
		nextCursor int64
	)

	cursorId = 1
	size = 11

	if cursor.HashedCursor != "" {
		cursorBody, err := cs.Decode(cursor.HashedCursor)
		if err != nil {
			s.Logger.Errorf("An error occurred: %v", err)
			errorhandler.InternalErrorHandler(w)
			return
		}
		cursorBody.Size++

		cursorId = cursorBody.NextCursor
		size = cursorBody.Size
	}

	var users []domain.User
	users, nextCursor, err := s.repo.FindAllUsersCursor(cursorId, size)
	if err != nil {
		errorhandler.BadRequestErrorHandler(w, err, r.URL.Path)
		s.Logger.Errorf("An error occurred: %v", err)
		return
	}

	s.Logger.Infof("Found %v users, the next should be %v", len(users), nextCursor)

	hashSrc, err = cs.Encode(size, nextCursor)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	pageModel := pagination.NewCursorHashedPagination(users, hashSrc)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pageModel)
}

func (s *userService) FindAllUsersCursorPagination(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("Listing all the users in the database using the cursor pagination...")

	defaultValueCursor := "100"
	cursor, err := getQueryParam[int64](r, "cursor", defaultValueCursor)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	defaultValueSize := "25"
	size, err := getQueryParam[int](r, "size", defaultValueSize)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	if size <= 0 {
		size = 25
	}
	// Don't change this, the size should be +1 so that we can have
	// the next cursor value
	size++

	var users []domain.User
	users, nextCursor, err := s.repo.FindAllUsersCursor(cursor, size)
	if err != nil {
		errorhandler.BadRequestErrorHandler(w, err, r.URL.Path)
		s.Logger.Errorf("An error occurred: %v", err)
		return
	}

	s.Logger.Infof("Found %v users, the next should be %v", len(users), nextCursor)

	pageModel := pagination.NewCursorPagination(users, size-1, nextCursor)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pageModel)
}

// FindAllUsers list all the users without a filter and returns each
// one with pagination and a few datas missing for LGPD.
func (s *userService) FindAllUsersOffSetPagination(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("Listing all the users in the database using the offset pagination...")

	defaultValueSize := "25"
	size, err := getQueryParam[int](r, "size", defaultValueSize)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	if size <= 0 {
		size = 25
	}

	defaultValuePage := "1"
	currentPage, err := getQueryParam[int](r, "page", defaultValuePage)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	if currentPage <= 0 {
		currentPage = 1
	}

	var users []domain.User
	users, totalRecords, err := s.repo.FindAllUsers(size, currentPage)
	if err != nil {
		errorhandler.BadRequestErrorHandler(w, err, r.URL.Path)
		s.Logger.Errorf("An error occurred: %v", err)
		return
	}

	pageModel := pagination.NewOffSetPagination(users, uint(currentPage), uint(totalRecords), uint(size))

	s.Logger.Infof("Found %v users from a total of %v", len(users), totalRecords)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pageModel)
}

// FindUserById list an user by his id and returns a none
// pagination result and a few datas missing for LGPD.
func (s *userService) FindUserById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	s.Logger.Infof("Listing user by id - %v", idStr)

	if idStr == "" {
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrIdIsRequired, r.URL.Path)
		s.Logger.Errorf("An error occurred: %v", errorhandler.ErrIdIsRequired)
		return
	}

	pathId, _ := strconv.Atoi(idStr)
	id := int64(pathId)

	var user *domain.User
	user, err := s.repo.FindUserById(id)
	if err != nil {
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrUserNotFound, r.URL.Path)
		s.Logger.Errorf("An error occurred: %v", err)
		return
	}

	s.Logger.Infof("User found! username: %v\n", *user.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}

func (s *userService) FindUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	s.Logger.Infof("Listing user by username: %v", username)

	if username == "" {
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrUsernameIsRequired, r.URL.Path)
		s.Logger.Errorf("An error occurred: %v", errorhandler.ErrUsernameIsRequired)
		return
	}

	var user *domain.User
	user, err := s.repo.FindUserByUsername(username)
	if err != nil {
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrUserNotFound, r.URL.Path)
		s.Logger.Errorf("An error occurred: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}

// UpdateUserDetails changes the user age and/or name if it's the account owner.
func (s *userService) UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("Updating an user")

	var user *domain.User
	username := r.PathValue("username")
	user, err := s.repo.FindUserByUsername(username)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrUserNotFound, r.URL.Path)
		return
	}

	var newUserDetails domain.UserDetails
	if err := json.NewDecoder(r.Body).Decode(&newUserDetails); err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	if err := isValidUserDetails(user, &newUserDetails); err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.BadRequestErrorHandler(w, err, r.URL.Path)
		return
	}

	user.Age = &newUserDetails.Age

	err = s.repo.UpdateUserDetails(user)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	s.Logger.Infof("User updated successfully! username: %v\n", *user.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// DeleteAccount deletes the user from the database by his username
// if it's the account owner.
func (s *userService) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	s.Logger.Infof("Deleting an account by his username: %v\n", username)

	err := s.repo.DeleteAccount(username)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	s.Logger.Infof("Account deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("Account deleted successfully")
}

func isValidUserDetails(user *domain.User, userRequest *domain.UserDetails) error {
	if userRequest.Age <= 0 {
		return errorhandler.ErrAgeIsRequired
	}

	if userRequest.Age == *user.Age {
		return errorhandler.ErrEqualAge
	}

	return nil
}

func getQueryParam[T int | int64 | string | bool | float64](r *http.Request, key string, defaultVal string) (T, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		valStr = defaultVal
	}

	var zeroVal T

	switch any(zeroVal).(type) {
	case int:
		value, err := strconv.Atoi(valStr)
		if err != nil {
			return zeroVal, err
		}
		return any(value).(T), nil
	case int64:
		value, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return zeroVal, err
		}
		return any(value).(T), nil
	case bool:
		value, err := strconv.ParseBool(valStr)
		if err != nil {
			return zeroVal, err
		}
		return any(value).(T), nil
	case float64:
		value, err := strconv.ParseBool(valStr)
		if err != nil {
			return zeroVal, err
		}
		return any(value).(T), nil
	case string:
		return any(valStr).(T), nil
	default:
		return zeroVal, errorhandler.ErrInvalidType
	}
}
