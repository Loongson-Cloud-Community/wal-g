package internal

import (
	"github.com/wal-g/wal-g/internal/walparser/parsingutil"
	"github.com/wal-g/wal-g/utility"
	"io"
)

type WalPartDataType uint8

const (
	PreviousWalHeadType WalPartDataType = 0
	WalTailType         WalPartDataType = 1
	WalHeadType         WalPartDataType = 2
)

type WalPart struct {
	dataType WalPartDataType
	id       uint8
	data     []byte
}

func NewWalPart(dataType WalPartDataType, id uint8, data []byte) *WalPart {
	return &WalPart{dataType, id, data}
}

func (part *WalPart) Save(writer io.Writer) error {
	_, err := writer.Write([]byte{byte(part.dataType), part.id})
	if err != nil {
		return err
	}
	dataLen := uint32(len(part.data))
	_, err = writer.Write(utility.ToBytes(&dataLen))
	if err != nil {
		return err
	}
	_, err = writer.Write(part.data)
	return err
}

func saveWalParts(parts []WalPart, writer io.Writer) error {
	for _, part := range parts {
		err := part.Save(writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadWalPart(reader io.Reader) (*WalPart, error) {
	var dataType WalPartDataType
	var partId uint8
	var dataLen uint32
	err := parsingutil.ParseMultipleFieldsFromReader([]parsingutil.FieldToParse{
		{Field: &dataType, Name: "part data type"},
		{Field: &partId, Name: "part number"},
		{Field: &dataLen, Name: "part data len"},
	}, reader)
	if err != nil {
		return nil, err
	}
	data := make([]byte, dataLen)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}
	return &WalPart{dataType, partId, data}, nil
}