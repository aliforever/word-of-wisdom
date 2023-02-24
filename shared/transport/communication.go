package transport

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

/*type headerLength int

const (
	HeaderLength16 headerLength = iota
	HeaderLength32
	HeaderLength64
)*/

type Transport struct {
	conn *net.TCPConn

	writeQueue chan writeItem
}

func NewTransport(conn *net.TCPConn) *Transport {
	t := &Transport{conn: conn, writeQueue: make(chan writeItem)}

	go t.processWriteQueue()

	return t
}

func (t *Transport) ReadContentLength() (int64, error) {
	contentLength := new(int64)

	err := binary.Read(t.conn, binary.LittleEndian, contentLength)
	if err != nil {
		return 0, err
	}

	return *contentLength, nil
}

func (t *Transport) ReadBytes(length int64) ([]byte, error) {
	data := make([]byte, length)

	read, err := io.ReadFull(t.conn, data)
	if err != nil {
		return nil, err
	}

	if int64(read) != length {
		return data, fmt.Errorf("invalid_content_length_given")
	}

	return data, nil
}

func (t *Transport) writeContentLength(length int64) error {
	return binary.Write(t.conn, binary.LittleEndian, length)
}

func (t *Transport) WriteNextBytes(data []byte) (int, error) {
	return t.writeBytes(data)
}

func (t *Transport) writeBytes(data []byte) (int, error) {
	nw := newWriteItem(data)

	go t.pushToWriteQueue(nw)

	return <-nw.count, <-nw.err
}

func (t *Transport) ReadNextBytes() ([]byte, error) {
	contentLength, err := t.ReadContentLength()
	if err != nil {
		return nil, err
	}

	return t.ReadBytes(contentLength)
}

func (t *Transport) ReceiveBytes(bytes chan<- []byte) error {
	var (
		data []byte
		err  error
	)

	for {
		data, err = t.ReadNextBytes()
		if err != nil {
			return err
		}

		go t.writeReceivedBytes(bytes, data)
	}
}

func (t *Transport) writeReceivedBytes(ch chan<- []byte, data []byte) {
	ch <- data
}

func (t *Transport) processWriteQueue() {
	for data := range t.writeQueue {
		t.write(data)
	}
}

func (t *Transport) pushToWriteQueue(data writeItem) {
	t.writeQueue <- data
}

func (t *Transport) write(item writeItem) {
	err := binary.Write(t.conn, binary.LittleEndian, int64(len(item.data)))
	if err != nil {
		item.writeResult(0, err)
		return
	}

	result, err := t.conn.Write(item.data)

	item.writeResult(result, err)
}
