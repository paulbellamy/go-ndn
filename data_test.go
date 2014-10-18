package ndn

import (
	"bytes"
	"testing"
	"time"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/stretchr/testify/assert"
)

func Test_Data_Name(t *testing.T) {
	subject := &Data{}

	// when unspecified
	assert.Equal(t, subject.GetName().Size(), 0)

	name := Name{Component{"a"}, Component{"b"}}
	subject.SetName(name)
	assert.Equal(t, subject.GetName(), name)

	// check it is copied
	name.Clear()
	assert.Equal(t, subject.GetName(), Name{Component{"a"}, Component{"b"}})
}

func Test_Data_FreshnessPeriod(t *testing.T) {
	subject := &Data{}

	// when unspecified
	assert.Equal(t, subject.GetFreshnessPeriod(), -1)
	subject.SetFreshnessPeriod(2 * time.Minute)
	assert.Equal(t, subject.GetFreshnessPeriod(), 2*time.Minute)
}

func Test_Data_Content(t *testing.T) {
	subject := &Data{}

	// when unspecified
	assert.Equal(t, subject.GetContent(), []byte{})
	content := []byte("hello world")
	subject.SetContent(content)
	assert.Equal(t, subject.GetContent(), content)
}

func Test_Data_Signature(t *testing.T) {
	subject := &Data{}

	assert.Nil(t, subject.GetSignature())
	signature := Sha256WithRSASignature{}
	subject.SetSignature(signature)
	assert.Equal(t, subject.GetSignature(), signature)
}

func Test_Data_WriteTo(t *testing.T) {
	expected, err := encoding.ParentTLV{
		T: encoding.DataType,
		V: []encoding.TLV{
			encoding.ParentTLV{
				T: encoding.NameType,
				V: []encoding.TLV{
					encoding.ByteTLV{
						T: encoding.NameComponentType,
						V: []byte{'a'},
					},
				},
			},
		},
	}.MarshalBinary()
	assert.NoError(t, err)

	buf := &bytes.Buffer{}
	subject := &Data{
		name: Name{Component{"a"}},
	}
	n, err := subject.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, n, int64(len(buf.Bytes())))
	assert.Equal(t, buf.Bytes(), expected)
}

func Test_Data_WriteTo_WithoutAName(t *testing.T) {
	buf := &bytes.Buffer{}

	subject := &Data{}
	n, err := subject.WriteTo(buf)
	assert.EqualError(t, err, ErrNameRequired.Error())
	assert.Equal(t, n, int64(0))
	assert.Nil(t, buf.Bytes(), "Expected nothing to be written, but data was found.")
}
