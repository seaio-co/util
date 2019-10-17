package bench

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *A) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "BirthDay":
			z.BirthDay, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "Phone":
			z.Phone, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Siblings":
			z.Siblings, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "Spouse":
			z.Spouse, err = dc.ReadUint8()
			if err != nil {
				return
			}
		case "Money":
			z.Money, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}