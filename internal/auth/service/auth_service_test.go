package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rafaeldepontes/fauthless-go/internal/auth"
	"github.com/rafaeldepontes/fauthless-go/internal/cache"
	"github.com/rafaeldepontes/fauthless-go/internal/domain"
	"github.com/rafaeldepontes/fauthless-go/internal/errorhandler"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	usernameMockTest   = "test"
	usernameMockBob    = "bob"
	hashedPasswordMock = "12345678"
	ageMock            = 18
)

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }
func ptrInt64(i int64) *int64    { return &i }

type userRepoMock struct {
	users            map[string]*domain.User
	lastSessionToken string
	lastCSRFToken    string
}

func newMockUserRepo() *userRepoMock {
	return &userRepoMock{users: make(map[string]*domain.User)}
}

func (mock *userRepoMock) FindUserByUsername(username string) (*domain.User, error) {
	user, ok := mock.users[username]
	if !ok {
		return nil, errorhandler.ErrUserNotFound
	}

	return user, nil
}

func (mock *userRepoMock) RegisterUser(u *domain.User) error {
	if _, ok := mock.users[*u.Username]; ok {
		return errorhandler.ErrUserAlreadyExists
	}
	mock.users[*u.Username] = u
	return nil
}

func (mock *userRepoMock) SetUserToken(token, csrfToken string, userId int64) error {
	mock.lastSessionToken = token
	mock.lastCSRFToken = csrfToken
	return nil
}

func (mock *userRepoMock) UpdateUserDetails(user *domain.User) error { return nil }

func (mock *userRepoMock) DeleteAccount(username string) error { return nil }

func (mock *userRepoMock) FindAllUsersCursor(cursor int64, size int) ([]domain.User, int64, error) {
	return nil, 0, nil
}

func (mock *userRepoMock) FindAllUsers(size, page int) ([]domain.User, int, error) {
	return nil, 0, nil
}

func (mock *userRepoMock) FindUserById(id int64) (*domain.User, error) { return nil, nil }

type mockSessionRepo struct {
	sessions map[string]*domain.Session
}

func newMockSessionRepo() *mockSessionRepo {
	return &mockSessionRepo{sessions: map[string]*domain.Session{}}
}

func (mock *mockSessionRepo) CreateSession(s *domain.Session) (string, error) {
	id := s.Id
	mock.sessions[id] = s
	return id, nil
}

func (mock *mockSessionRepo) FindSessionById(id string) (*domain.Session, error) {
	if s, ok := mock.sessions[id]; ok {
		return s, nil
	}
	return nil, errorhandler.ErrSessionNotFound
}

func (mock *mockSessionRepo) RevokeSession(id string) error {
	if s, ok := mock.sessions[id]; ok {
		s.IsRevoked = true
		return nil
	}
	return errorhandler.ErrSessionNotFound
}

func (mock *mockSessionRepo) DeleteSession(id string) error { return nil }

func prepareMocks() (auth.Service, *userRepoMock, *mockSessionRepo, *cache.Caches) {
	// given
	userRepo := newMockUserRepo()
	sessRepo := newMockSessionRepo()
	logg := logrus.New()
	secretKey := "secret-key"
	cache := cache.NewCacheStorage()

	return NewAuthService(userRepo, sessRepo, logg, secretKey, cache), userRepo, sessRepo, cache
}

func loginFlowMock(userRepo *userRepoMock) (*httptest.ResponseRecorder, *http.Request) {
	newUser := domain.UserLogin{
		Username: usernameMockTest,
		Password: hashedPasswordMock,
	}
	jsonReq, _ := json.Marshal(newUser)
	var r *http.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonReq))
	var w *httptest.ResponseRecorder = httptest.NewRecorder()

	hashed, _ := bcrypt.GenerateFromPassword([]byte(hashedPasswordMock), Cost)
	u := &domain.User{
		Id:             ptrInt64(1),
		Username:       ptrString(usernameMockTest),
		HashedPassword: ptrString(string(hashed)),
		Age:            ptrInt(28),
	}
	userRepo.RegisterUser(u)

	return w, r
}

// TestNewAuthService verifies the NewAuthService and mocks each service.
func TestNewAuthService(t *testing.T) {
	// when
	auth, _, _, _ := prepareMocks()

	// then
	if auth == nil {
		t.Fatal("auth.Service is nil after NewAuthService(...)")
	}
}

// TestIsValidUser verifies isValidUser and tests all cases with a
// test driven table.
func TestIsValidUser(t *testing.T) {
	// given
	logg := logrus.New()
	cache := &cache.Caches{}
	auth := &authService{
		userRepository: &userRepoMock{
			users: map[string]*domain.User{
				usernameMockTest: {
					Username: ptrString(usernameMockTest),
				},
			},
		},
		Logger: logg,
		Cache:  cache,
	}

	tests := []struct {
		name    string
		input   *domain.User
		wantOk  bool
		wantErr error
	}{
		{
			name: "empty username",
			input: &domain.User{
				Username:       ptrString(""),
				HashedPassword: ptrString(hashedPasswordMock),
				Age:            ptrInt(ageMock),
			},
			wantOk:  false,
			wantErr: errorhandler.ErrUsernameIsRequired,
		},
		{
			name: "empty password",
			input: &domain.User{
				Username:       ptrString(usernameMockTest),
				HashedPassword: ptrString(""),
				Age:            ptrInt(ageMock),
			},
			wantOk:  false,
			wantErr: errorhandler.ErrPasswordIsRequired,
		},
		{
			name: "missing age",
			input: &domain.User{
				Username:       ptrString(usernameMockTest),
				HashedPassword: ptrString(hashedPasswordMock),
				Age:            ptrInt(0),
			},
			wantOk:  false,
			wantErr: errorhandler.ErrAgeIsRequired,
		},
		{
			name: "user already exists",
			input: &domain.User{
				Username:       ptrString(usernameMockTest),
				HashedPassword: ptrString(hashedPasswordMock),
				Age:            ptrInt(ageMock),
			},
			wantOk:  false,
			wantErr: errorhandler.ErrUserAlreadyExists,
		},
		{
			name: "valid user",
			input: &domain.User{
				Username:       ptrString(usernameMockBob),
				HashedPassword: ptrString(hashedPasswordMock),
				Age:            ptrInt(ageMock),
			},
			wantOk:  true,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := isValidUser(tt.input, auth)
			if ok != tt.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOk)
			}
			if (err == nil) != (tt.wantErr == nil) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestRegister_Success verifies the Register and expects a success
// and an user in the mocked database.
func TestRegister_Success(t *testing.T) {
	// given
	auth, userRepo, _, _ := prepareMocks()
	newUser := domain.User{
		Username:       ptrString(usernameMockTest),
		HashedPassword: ptrString(hashedPasswordMock),
		Age:            ptrInt(ageMock),
	}
	jsonReq, _ := json.Marshal(newUser)
	var r *http.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(jsonReq))
	var w *httptest.ResponseRecorder = httptest.NewRecorder()

	// when
	auth.Register(w, r)

	// then
	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code 201, but it was %v", resp.StatusCode)
	}

	if _, err := userRepo.FindUserByUsername(usernameMockTest); err != nil {
		t.Fatalf("expected user to be saved, err: %v", err)
	}
}

// TestLoginCookieBased_Success verifies LoginCookieBased and expects
// a success when log in.
func TestLoginCookieBased_Success(t *testing.T) {
	// given
	auth, userRepo, _, _ := prepareMocks()
	w, r := loginFlowMock(userRepo)

	//when
	auth.LoginCookieBased(w, r)
	resp := w.Result()

	//then
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 201, but it was %v", resp.StatusCode)
	}

	if userRepo.lastSessionToken == "" {
		t.Fatalf("expected session token to be saved, got: %v", userRepo.lastCSRFToken)
	}

	if userRepo.lastCSRFToken == "" {
		t.Fatalf("expected CSRF token to be saved, got: %v", userRepo.lastCSRFToken)
	}
}

// TestLoginJwtBased_Success verifies LoginJwtBased and expects a
// success when log in and a token in the request body.
func TestLoginJwtBased_Success(t *testing.T) {
	// given
	auth, userRepo, _, _ := prepareMocks()
	w, r := loginFlowMock(userRepo)

	// when
	auth.LoginJwtBased(w, r)
	resp := w.Result()

	var tr domain.TokenResponse

	// then
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		t.Fatalf("failed decoding token response: %v", err)
	}

	if tr.Token == "" {
		t.Fatalf("expected token in response")
	}
}

// TestLoginJwtBased_Success verifies LoginJwtRefreshBased and expects a
// success when log in and a access token, a refresh token, the time both
// will expire and the session id in the request body.
func TestLoginJwtRefreshBased_Success(t *testing.T) {
	// given
	auth, userRepo, sessionRepo, _ := prepareMocks()
	w, r := loginFlowMock(userRepo)

	// when
	auth.LoginJwtRefreshBased(w, r)
	resp := w.Result()

	var tr domain.TokenRefreshResponse

	// then
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		t.Fatalf("failed decoding token response: %v", err)
	}

	if tr.AccessToken == "" {
		t.Fatalf("expected access token in response")
	}

	if tr.RefreshToken == "" {
		t.Fatalf("expected refresh token in response")
	}

	if _, err := sessionRepo.FindSessionById(tr.SessionId); err != nil {
		t.Fatalf("Got an error trying to get the session by its id: %v", err)
	}
}

// TestRenewAccessToken_Success verifies RenewAccessToken and expects
// success when trying to renew your access token. In the request
// it should contain the refresh token to validation and returns
// a new access token.
func TestRenewAccessToken_Success(t *testing.T) {
	// given
	auth, userRepo, _, _ := prepareMocks()

	w, r := loginFlowMock(userRepo)
	auth.LoginJwtRefreshBased(w, r)
	resp := w.Result()

	var tr domain.TokenRefreshResponse

	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		t.Fatalf("failed decoding token response: %v", err)
	}

	newAccessTokenReq := domain.RenewAccessTokenRequest{
		RefreshToken: tr.RefreshToken,
	}
	jsonReq, _ := json.Marshal(newAccessTokenReq)

	var req *http.Request = httptest.NewRequest(http.MethodPost, "/renew", bytes.NewReader(jsonReq))
	var rr *httptest.ResponseRecorder = httptest.NewRecorder()

	// when
	auth.RenewAccessToken(rr, req)
	resp = rr.Result()

	var ratr domain.RenewAccessTokenResponse

	// then
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&ratr); err != nil {
		t.Fatalf("failed decoding token response: %v", err)
	}

	if ratr.AccessToken == "" {
		t.Fatal("expected new access token in response")
	}

	if ratr.AccessTokenExpiresAt.Equal(time.Time{}) {
		t.Fatal("expected new expire date for the access token in response")
	}
}
