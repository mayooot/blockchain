package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Int64ToBytes(in int64) []byte {
	var buff bytes.Buffer
	if err := binary.Write(&buff, binary.BigEndian, in); err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
