package controllers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	type want struct {
		res bool
	}
	tests := []struct {
		name string
		arg  int
		want want
	}{
		{
			name: "positive test #1",
			arg:  123,
			want: want{
				res: false,
			},
		},
		{
			name: "positive test #2",
			arg:  12344,
			want: want{
				res: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := valid(tt.arg)
			assert.Equal(t, tt.want.res, res, fmt.Errorf("expected result %t, got %t", tt.want.res, res))
		})
	}
}
