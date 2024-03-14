package gogesrtp

import (
	"errors"
	"strconv"
	"strings"
)

type DataItem struct {
	Address      string
	DataType     byte
	StartAddress uint16
	Length       uint16
	IsBit        bool
	//ByteAddress  uint16
	BitAddress uint16 //0-7
	Value      interface{}
}

func ParseDataItemForm(address string, valLength uint16, isBit bool) (DataItem, error) {
	item := DataItem{}
	item.Address = strings.ToUpper(address)
	item.Length = valLength
	item.IsBit = isBit
	var err interface{}
	switch item.Address[0:2] {
	case "AI":
		if isBit {
			err = errors.New("AI area notSupported Bit")
			break
		}
		item.DataType = AI
		item.StartAddress = uint16(address[2])
	case "AQ":
		if isBit {
			err = errors.New("AQ area notSupported Bit")
			break
		}
		item.DataType = AQ
		item.StartAddress = uint16(address[2])
	case "SA":
		if isBit {
			item.DataType = SA_BIT
		} else {
			item.DataType = SA_BYTE
		}
		item.StartAddress = uint16(address[2])
	case "SB":
		if isBit {
			item.DataType = SB_BIT
		} else {
			item.DataType = SB_BYTE
		}
		item.StartAddress = uint16(address[2])
	case "SC":
		if isBit {
			item.DataType = SC_BIT
		} else {
			item.DataType = SC_BYTE
		}
		item.StartAddress = uint16(address[2])
	default:
		num, err := strconv.ParseInt(address[1:], 10, 16)
		if err != nil {
			break
		}
		switch address[0:1] {
		case "R":
			if isBit {
				err = errors.New("R area notSupported Bit")
				break
			}
			item.DataType = R
			item.StartAddress = uint16(num)
		case "I":
			if isBit {
				item.DataType = I_BIT
			} else {
				item.DataType = I_BYTE
			}
			item.StartAddress = uint16(num)
		case "M":
			if isBit {
				item.DataType = M_BIT
			} else {
				item.DataType = M_BYTE
			}
			item.StartAddress = uint16(num)
		case "Q":
			if isBit {
				item.DataType = Q_BIT
			} else {
				item.DataType = Q_BYTE
			}
			item.StartAddress = uint16(num)
		case "T":
			if isBit {
				item.DataType = T_BIT
			} else {
				item.DataType = T_BYTE
			}
			item.StartAddress = uint16(num)
		case "G":
			if isBit {
				item.DataType = G_BIT
			} else {
				item.DataType = G_BYTE
			}
			item.StartAddress = uint16(num)
		default:
			err = errors.New("address error")
		}
	}
	if err != nil {
		return DataItem{}, err.(error)
	}
	if item.StartAddress > 0 {
		//item.ByteAddress = uint16(math.Ceil(dataitem.StartAddress / 8))
		item.BitAddress = uint16((item.StartAddress - 1) % 8)
	}
	return item, nil
}
