package fxctx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Key
// Is a type for contextual keys naming.
// Developer should make correctly scoped variables to be able to read a context w/ the usage of this Keys
type Key string

// StringerKey
// is a constraint for generics functions
type StringerKey interface {
	~string
}

// AssingCtxValue
// Function that assigns value `v` of any type to a key `k` of String constrain type
func AssingCtxValue[V any, K StringerKey](ctx context.Context, k K, v *V) context.Context {
	return context.WithValue(ctx, k, v)
}

var ErrMissingValue error = errors.New("the value is missing")

// errValueIsMissing
// helps to return human readable error while retrieving a value from the context
func errValueIsMissing[K StringerKey](k K) error {
	return fmt.Errorf("%w. provided key is: `%v`", ErrMissingValue, k)
}

// GetCtxValue
// Retrieves a value from the context
func GetCtxValue[V any, K StringerKey](ctx context.Context, k K) (*V, error) {
	u := ctx.Value(k)

	if value, ok := u.(*V); ok {
		return value, nil
	}

	return nil, errValueIsMissing(k)
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

type Validator interface {
	// Valid checks the object and returns any
	// problems. If len(problems) == 0 then
	// the object is valid.
	Valid() (problems map[string]string)
}

func DecodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}
	if problems := v.Valid(); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}
