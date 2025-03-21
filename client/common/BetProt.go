package common

import (
	"os"
	"fmt"
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"io"
	"bytes"
)

type BetProt struct {
	agency string
	name string
	surname string
	document string
	birthdate string
	number string
}

func NewBetProt(agency string, params []string) *BetProt {
	return &BetProt{
		agency: agency,
		name: params[0],
		surname: params[1],
		document: params[2],
		birthdate: params[3],
		number: params[4],
	}
}

func (b *BetProt) serialize() []byte {
	return []byte(fmt.Sprintf("%s;%s;%s;%s;%s;%s\n", b.agency, b.name, b.surname, b.document, b.birthdate, b.number))
}

func SendBatches(maxBatch int, agency string, writer *bufio.Writer) error {
	file, err := os.Open(fmt.Sprintf("agency-%s.csv", agency))
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	lineCount := 0
	var buffer bytes.Buffer 
	for {
		if lineCount == maxBatch {
			lineCount = 0
			if err := SendBatch(buffer, writer, false); err != nil {
				return err
			}
			buffer.Reset()
		}
		line, err := reader.Read()
		if err == io.EOF {
			if err := SendBatch(buffer, writer, true); err != nil {
				return err
			}
			break
		}
		if err != nil {
			return err
		}
		bet := NewBetProt(agency, line)
		buffer.Write(bet.serialize())
		lineCount++
	}
	return nil
}

func SendBatch(buffer bytes.Buffer, writer *bufio.Writer, lastBatch bool) error {
	msg := buffer.Bytes()
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(msg)))
	if lastBatch == true {
        writer.WriteByte(1)
    } else {
        writer.WriteByte(0)
    }
	writer.Write(lengthBytes)
	writer.Write(msg)
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}