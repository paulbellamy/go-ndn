package encoding

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadNumber_EmptyReader(t *testing.T) {
	num, n, err := ReadNumber(bytes.NewReader([]byte{}))
	assert.EqualError(t, err, io.EOF.Error())
	assert.Equal(t, n, 0)
	assert.Equal(t, num, uint64(0))
}

func Test_ReadNumber_OneOctetValue(t *testing.T) {
	num, n, err := ReadNumber(bytes.NewReader([]byte{252}))
	assert.NoError(t, err)
	assert.Equal(t, n, 1)
	assert.Equal(t, num, uint64(252))
}

func Test_ReadNumber_TwoOctetValue(t *testing.T) {
	num, n, err := ReadNumber(bytes.NewReader([]byte{253, 0, 123}))
	assert.NoError(t, err)
	assert.Equal(t, n, 3)
	assert.Equal(t, num, uint64(123))
}

func Test_ReadNumber_TwoOctetValueUnderflow(t *testing.T) {
	num, n, err := ReadNumber(bytes.NewReader([]byte{253}))
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Equal(t, n, 1)
	assert.Equal(t, num, uint64(0))
}

func Test_ReadNumber_FourOctetValue(t *testing.T) {
	num, n, err := ReadNumber(bytes.NewReader([]byte{254, 0, 0, 0, 123}))
	assert.NoError(t, err)
	assert.Equal(t, n, 5)
	assert.Equal(t, num, uint64(123))
}

func Test_ReadNumber_FourOctetValueUnderflow(t *testing.T) {
	num, n, err := ReadNumber(bytes.NewReader([]byte{254}))
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Equal(t, n, 1)
	assert.Equal(t, num, uint64(0))
}

func Test_ReadNumber_EightOctetValue(t *testing.T) {
	num, n, err := ReadNumber(bytes.NewReader([]byte{255, 0, 0, 0, 0, 0, 0, 0, 123}))
	assert.NoError(t, err)
	assert.Equal(t, n, 9)
	assert.Equal(t, num, uint64(123))
}

func Test_ReadNumber_EightOctetValueUnderflow(t *testing.T) {
	num, n, err := ReadNumber(bytes.NewReader([]byte{255}))
	assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
	assert.Equal(t, n, 1)
	assert.Equal(t, num, uint64(0))
}
