package ndn

import (
	"testing"

	"github.com/paulbellamy/go-ndn/encoding"
	"github.com/stretchr/testify/assert"
)

func Test_NameFromTLV(t *testing.T) {
	tlv := encoding.ParentTLV{
		T: encoding.NameType,
		V: []encoding.TLV{
			encoding.ByteTLV{
				T: encoding.NameComponentType,
				V: []byte("foo"),
			},
		},
	}
	subject, err := NameFromTLV(tlv)
	assert.NoError(t, err)
	assert.Equal(t, subject, Name{Component{"foo"}})
}

func Test_NameFromTLV_WrongType(t *testing.T) {
	tlv := encoding.ParentTLV{T: encoding.InterestType}
	subject, err := NameFromTLV(tlv)
	assert.EqualError(t, err, "TLV is not a name")
	assert.Nil(t, subject)
}

func Test_Name_Copy(t *testing.T) {
	subject := Name{Component{"a"}, Component{"b"}}
	subject2 := subject.Copy()
	assert.Equal(t, subject, subject2)

	// check it was copied
	(&subject).Clear()
	assert.Equal(t, subject2, Name{Component{"a"}, Component{"b"}})
}

func Test_Name_Append(t *testing.T) {
	subject := Name{}
	subject = subject.AppendBytes([]byte("a"))
	subject = subject.AppendString("b☃")
	//AppendBlob(??).
	//AppendComponent(Component{}).
	//AppendName(Name{"e", "f"}).
	//AppendSegment(0x12).
	//AppendSegmentOffset(4).
	//AppendSequenceNumber(4).
	//AppendTimestamp(now).
	//AppendVersion(123).
}

func Test_Name_Clear(t *testing.T) {
	subject := Name{Component{"a"}, Component{"b"}}
	subject.Clear()
	assert.Equal(t, subject, Name{})
}

func Test_Name_Compare(t *testing.T) {
	subject := Name{Component{"foo"}, Component{"bar"}}
	assert.Equal(t, subject.Compare(Name{Component{"foo"}}), 1)
	assert.Equal(t, subject.Compare(Name{Component{"foo"}, Component{"ba"}}), 1)
	assert.Equal(t, subject.Compare(Name{Component{"foo"}, Component{"aar"}}), 1)
	assert.Equal(t, subject.Compare(Name{Component{"foo"}, Component{"bar"}}), 0)
	assert.Equal(t, subject.Compare(Name{Component{"foo"}, Component{"car"}}), -1)
	assert.Equal(t, subject.Compare(Name{Component{"foo"}, Component{"bara"}}), -1)
	assert.Equal(t, subject.Compare(Name{Component{"foo"}, Component{"bar"}, Component{"baz"}}), -1)
}

func Test_Name_Equals(t *testing.T) {
	subject := Name{Component{"foo"}, Component{"bar"}}
	assert.False(t, subject.Equals(Name{Component{"foo"}}))
	assert.False(t, subject.Equals(Name{Component{"foo"}, Component{"baz"}}))
	assert.True(t, subject.Equals(Name{Component{"foo"}, Component{"bar"}}))
}

func Test_Name_Get(t *testing.T) {
	component1 := Component{"a"}
	component2 := Component{"b"}
	subject := Name{component1, component2}
	assert.Equal(t, subject.Get(0), component1)
	assert.Equal(t, subject.Get(1), component2)
	assert.Equal(t, subject.Get(-1), component2)
	assert.Equal(t, subject.Get(-2), component1)
}

func Test_Name_GetPrefix(t *testing.T) {
	component1 := Component{"a"}
	component2 := Component{"b"}
	component3 := Component{"c"}
	subject := Name{component1, component2, component3}
	assert.Equal(t, subject.GetPrefix(0), Name{})
	assert.Equal(t, subject.GetPrefix(1), Name{component1})
	assert.Equal(t, subject.GetPrefix(2), Name{component1, component2})
	assert.Equal(t, subject.GetPrefix(3), Name{component1, component2, component3})
	assert.Equal(t, subject.GetPrefix(4), Name{component1, component2, component3})
	assert.Equal(t, subject.GetPrefix(-1), Name{component1, component2})
	assert.Equal(t, subject.GetPrefix(-2), Name{component1})
	assert.Equal(t, subject.GetPrefix(-3), Name{})
}

func Test_Name_GetSubName(t *testing.T) {
	component1 := Component{"a"}
	component2 := Component{"b"}
	component3 := Component{"c"}
	subject := Name{component1, component2, component3}
	assert.Equal(t, subject.GetSubName(0, 0), Name{})
	assert.Equal(t, subject.GetSubName(0, 2), Name{component1, component2})
	assert.Equal(t, subject.GetSubName(1, 1), Name{component2})
	assert.Equal(t, subject.GetSubName(1, 3), Name{component2, component3})
	assert.Equal(t, subject.GetSubName(2, 2), Name{component3})
	assert.Equal(t, subject.GetSubName(3, 3), Name{})
	assert.Equal(t, subject.GetSubName(-1, 1), Name{component3})
	assert.Equal(t, subject.GetSubName(-2, 2), Name{component2, component3})
	assert.Equal(t, subject.GetSubName(-3, 2), Name{component1, component2})
}

func Test_Name_Match(t *testing.T) {

	component1 := Component{"a"}
	component2 := Component{"b"}
	subject := Name{component1, component2}
	assert.True(t, subject.Match(Name{component1, component2}))
	assert.True(t, subject.Match(Name{component1, component2, Component{"foo"}}))
	assert.False(t, subject.Match(Name{component1}))
	assert.False(t, subject.Match(Name{component1, Component{"foo"}, Component{"bar"}}))

	// always true for empty names
	assert.True(t, Name{}.Match(Name{component1}))
}

func Test_Name_Size(t *testing.T) {
	subject := Name{Component{"a"}, Component{"b"}, Component{"c"}}
	assert.Equal(t, subject.Size(), 3)
}

func Test_Name_IsEmpty(t *testing.T) {
	assert.True(t, Name{}.IsEmpty())
	assert.False(t, Name{Component{"a"}, Component{"b"}, Component{"c"}}.IsEmpty())
}

func Test_Name_IsBlank(t *testing.T) {
	assert.True(t, Name{}.IsBlank())
	assert.False(t, Name{Component{"a"}}.IsBlank())
}
