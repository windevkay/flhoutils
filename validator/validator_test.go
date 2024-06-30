package validator

import (
	"testing"

	"github.com/windevkay/flhoutils/assert"
)

func TestPermittedValue(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want bool
	}{
		{name: "Validation passes on permitted value", arg: 1, want: true},
		{name: "Validation fails on non permitted value", arg: 2, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.Check(PermittedValue(tc.arg, 1, 3, 4), "", "")

			assert.Equal(t, v.Valid(), tc.want)
		})
	}
}

func TestUniqueSlice(t *testing.T) {
	tests := []struct {
		name string
		arg  []int
		want bool
	}{
		{name: "Validation passes on unique slice", arg: []int{1, 2, 3}, want: true},
		{name: "Validation fails on non unique slice", arg: []int{1, 1, 2}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.Check(Unique(tc.arg), "", "")

			assert.Equal(t, v.Valid(), tc.want)
		})
	}
}

func TestMatches(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{name: "Validation passes on valid email", arg: "test@testemail.com", want: true},
		{name: "Validation fails on invalid email", arg: "test", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.Check(Matches(tc.arg, EmailRX), "", "")

			assert.Equal(t, v.Valid(), tc.want)
		})
	}
}
