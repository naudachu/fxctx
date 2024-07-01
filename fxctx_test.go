package fxctx

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestErrWrapping(t *testing.T) {
	type F string
	err := errValueIsMissing[F]("somevalue")

	t.Run("check errorsIs equalation", func(t *testing.T) {
		b := errors.Is(err, ErrMissingValue)
		if b != true {
			t.Errorf("Error wrapped incorrectly")
		}
	})
}

type Stuff struct {
	Name string
}

func TestAssignAndGet(t *testing.T) {

	stuff := Stuff{
		Name: "some_stuff",
	}

	key := Key("key")
	ctx := AssingCtxValue(context.Background(), key, &stuff)

	t.Run("check if get ctx value ok with objects", func(t *testing.T) {
		smth, err := GetCtxValue[Stuff](ctx, key)
		if err != nil {
			t.Errorf(err.Error())
		}

		if !reflect.DeepEqual(stuff, *smth) {
			t.Errorf("not equal")
		}
	})
}
