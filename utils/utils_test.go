package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatToBigInt(t *testing.T) {
	f := FloatToBigInt(1.2, 12)
	assert.Equal(t, "1200000000000", f.String())
}
