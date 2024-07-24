package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBits(t *testing.T) {
	input := []struct {
		b        byte
		expected [8]uint8
	}{
		{b: 0, expected: [8]uint8{0, 0, 0, 0, 0, 0, 0, 0}},
		{b: 1, expected: [8]uint8{0, 0, 0, 0, 0, 0, 0, 1}},
		{b: 2, expected: [8]uint8{0, 0, 0, 0, 0, 0, 1, 0}},
		{b: 4, expected: [8]uint8{0, 0, 0, 0, 0, 1, 0, 0}},
		{b: 8, expected: [8]uint8{0, 0, 0, 0, 1, 0, 0, 0}},
		{b: 16, expected: [8]uint8{0, 0, 0, 1, 0, 0, 0, 0}},
		{b: 32, expected: [8]uint8{0, 0, 1, 0, 0, 0, 0, 0}},
		{b: 64, expected: [8]uint8{0, 1, 0, 0, 0, 0, 0, 0}},
		{b: 255, expected: [8]uint8{1, 1, 1, 1, 1, 1, 1, 1}},
	}

	for _, i := range input {
		assert.EqualValues(t, i.expected, bits(i.b))
	}
}
