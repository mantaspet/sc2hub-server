package main

import (
	"testing"
)

func TestParsePaginationParam(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		from string
		want int
	}{
		{
			name: "Valid",
			from: "50",
			want: 50,
		},
		{
			name: "Invalid",
			from: "asd",
			want: 0,
		},
		{
			name: "Empty",
			from: "",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from := parsePaginationParam(tt.from)

			if from != tt.want {
				t.Errorf("want %d; got %d", tt.want, from)
			}
		})
	}
}
