package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var Validator = validator.New()

func DecodeAndValidate(
	r *http.Request,
	dest any,
) error {

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("failed to decode request: %w", err)
	}
	if err := Validator.Struct(dest); err != nil {
		return fmt.Errorf("failed to validate request: %w", err)
	}

	return nil
}
