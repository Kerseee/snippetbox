package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Initialize test cases
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2021, 11, 18, 16, 57, 0, 0, time.UTC),
			want: "18 Nov 2021 at 16:57",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
	}

	// Run tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hd := humanDate(test.tm)
			if hd != test.want {
				t.Errorf("want %q; got %q", test.want, hd)
			}
		})
	}
}
