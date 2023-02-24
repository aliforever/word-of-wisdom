package transport

type writeItem struct {
	data  []byte
	err   chan error
	count chan int
}

func newWriteItem(data []byte) writeItem {
	return writeItem{
		data:  data,
		err:   make(chan error),
		count: make(chan int),
	}
}

func (w writeItem) writeResult(count int, err error) {
	go w.writeErr(err)
	go w.writeCount(count)
}

func (w writeItem) writeErr(err error) {
	w.err <- err
}

func (w writeItem) writeCount(count int) {
	w.count <- count
}
