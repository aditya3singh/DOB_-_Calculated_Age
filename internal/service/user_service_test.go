package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		dob  time.Time
		want int
	}{
		{
			name: "birthday already passed this year",
			dob:  time.Date(now.Year()-30, now.Month()-1, now.Day(), 0, 0, 0, 0, time.UTC),
			want: 30,
		},
		{
			name: "birthday has not yet occurred this year",
			dob:  time.Date(now.Year()-30, now.Month()+1, now.Day(), 0, 0, 0, 0, time.UTC),
			want: 29,
		},
		{
			name: "birthday is today",
			dob:  time.Date(now.Year()-25, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			want: 25,
		},
		{
			name: "newborn (same day)",
			dob:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			want: 0,
		},
		{
			name: "century old",
			dob:  time.Date(now.Year()-100, now.Month()-1, now.Day(), 0, 0, 0, 0, time.UTC),
			want: 100,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateAge(tc.dob)
			if got != tc.want {
				t.Errorf("CalculateAge(%v) = %d, want %d", tc.dob, got, tc.want)
			}
		})
	}
}
