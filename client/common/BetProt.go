package common

import (
	"os"
	"fmt"
	"bufio"
	"encoding/binary"
)

type BetProt struct {
	agency string
	name string
	surname string
	document string
	birthdate string
	number string
}

func NewBetProt(agency string) *BetProt {
	return &BetProt{
		agency: agency,
		name: os.Getenv("NOMBRE"),
		surname: os.Getenv("APELLIDO"),
		document: os.Getenv("DOCUMENTO"),
		birthdate: os.Getenv("NACIMIENTO"),
		number: os.Getenv("NUMERO"),
	}
}

func (b *BetProt) SendBet(writer *bufio.Writer) error {
	msg := []byte(fmt.Sprintf("%s;%s;%s;%s;%s;%s", b.agency, b.name, b.surname, b.document, b.birthdate, b.number))
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(msg)))
	writer.Write(lengthBytes)
	writer.Write(msg)
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}


