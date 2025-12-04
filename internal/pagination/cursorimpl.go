package pagination

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"os"
	"strconv"

	"github.com/rafaeldepontes/auth-go/internal/domain"
	"github.com/rafaeldepontes/auth-go/internal/errorhandler"
)

type cursorService[T any] struct {
	secretKey []byte
	sigLen    int
}

func NewCursorService[T any]() *cursorService[T] {
	sigLenStr := os.Getenv("SIGNATURE_LENGTH")
	sigLen, _ := strconv.Atoi(sigLenStr)

	return &cursorService[T]{
		secretKey: []byte(os.Getenv("SECRET_CURSOR_KEY")),
		sigLen:    sigLen,
	}
}

// Encode accepts a generic T type, a slice of any data, a size of records per page and
// the next page being a pointer to the next id in the database and it will return a hash
// with all the information needed in the next request for security.
func (s *cursorService[T]) Encode(data []T, size int, nextCursor int64) (string, error) {
	rawData := domain.CursorPagination[T]{
		Data:       data,
		Size:       size,
		NextCursor: nextCursor,
	}

	sb, err := json.Marshal(rawData)
	if err != nil {
		return "", err
	}

	secretKey := os.Getenv("SECRET_CURSOR_KEY")
	var mac hash.Hash = hmac.New(sha256.New, []byte(secretKey))
	mac.Write(sb)
	signature := mac.Sum(nil)

	combined := append(sb, signature...)

	return base64.RawURLEncoding.EncodeToString(combined), nil
}

// Decode accepts a hashed source to decode, it will return the CursorPagination with the
// T type generic specified previously and an error if any.
func (s *cursorService[T]) Decode(src string) (*domain.CursorPagination[T], error) {
	combined, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return nil, err

	}
	if len(combined) < s.sigLen {
		return nil, errorhandler.ErrInvalidCursorLength
	}

	jsonBody := combined[:len(combined)-s.sigLen]
	signature := combined[len(combined)-s.sigLen:]

	var mac hash.Hash = hmac.New(sha256.New, s.secretKey)
	mac.Write(jsonBody)
	expected := mac.Sum(nil)

	if !hmac.Equal(signature, expected) {
		return nil, errorhandler.ErrInvalidCursorSignature
	}

	var cursorModel domain.CursorPagination[T]
	json.Unmarshal(jsonBody, &cursorModel)

	fmt.Println(string(jsonBody))
	fmt.Println(cursorModel)

	return &cursorModel, nil
}
