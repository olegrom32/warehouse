package packager

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	s := NewService([]int{500, 250, 5000, 250, 2000, 1000, 500})

	assert.Equal(t, []int{5000, 2000, 1000, 500, 250}, s.sizes)
}

func TestService_Package(t *testing.T) {
	s := NewService([]int{500, 250, 5000, 2000, 1000})

	tests := []struct {
		arg  int
		want map[int]int
	}{
		{
			arg:  0,
			want: nil,
		},
		{
			arg:  1,
			want: map[int]int{250: 1},
		},
		{
			arg:  250,
			want: map[int]int{250: 1},
		},
		{
			arg:  251,
			want: map[int]int{500: 1},
		},
		{
			arg:  501,
			want: map[int]int{250: 1, 500: 1},
		},
		{
			arg:  12001,
			want: map[int]int{250: 1, 2000: 1, 5000: 2},
		},
		{
			arg:  751,
			want: map[int]int{1000: 1},
		},
		{
			arg:  999,
			want: map[int]int{1000: 1},
		},
		{
			arg:  15000,
			want: map[int]int{5000: 3},
		},
		{
			arg:  4999,
			want: map[int]int{5000: 1},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d items", tt.arg), func(t *testing.T) {
			assert.Equal(t, tt.want, s.Package(tt.arg))
		})
	}

	t.Run("when the only box size is 5000, should use this box to fulfil any request", func(t *testing.T) {
		assert.Equal(t, map[int]int{5000: 3}, (&Service{sizes: []int{5000}}).Package(12001))
	})
}
