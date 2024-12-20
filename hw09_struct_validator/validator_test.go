package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		// Case 0: Верная структура User
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil, // Все валидно, ошибок нет
		},
		// Case 1: Неверная структура User
		{
			in: User{
				ID:     "short-id",
				Name:   "John",
				Age:    17,
				Email:  "not-an-email",
				Role:   "unknown",
				Phones: []string{"12345", "invalid"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrInvalidLength},
				{Field: "Age", Err: ErrInvalidMin},
				{Field: "Email", Err: ErrInvalidRegexp},
				{Field: "Role", Err: ErrInvalidIn},
				{Field: "Phones", Err: ErrInvalidLength},
				{Field: "Phones", Err: ErrInvalidLength},
			},
		},
		// Case 2: Неверная структура App
		{
			in: App{
				Version: "1.0.0.1",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrInvalidLength},
			},
		},
		// Case 3: Неверная структура Response
		{
			in: Response{
				Code: 123,
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrInvalidIn},
			},
		},
		// Case 4: Это вообще не структура
		{
			in:          "not a struct",
			expectedErr: ErrNotStruct,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}
