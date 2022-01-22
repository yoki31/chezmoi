package progress

import (
	"io"
	"sync"
	"time"
)

type ProxyReader struct {
	sync.Mutex
	reader        io.Reader
	creationTime  time.Time
	firstByteTime time.Time
	bytesRead     int
	readCalls     int
}

type ReaderStats struct {
	BytesRead      int
	BytesPerSecond float64
	ReadCalls      int
}

type ReaderStatsFunc func() ReaderStats

func NewProxyReader(r io.Reader) *ProxyReader {
	return &ProxyReader{
		reader:       r,
		creationTime: time.Now(),
	}
}

func (r *ProxyReader) Read(p []byte) (int, error) {
	r.Lock()
	defer r.Unlock()
	n, err := r.reader.Read(p)
	if r.firstByteTime.IsZero() {
		r.firstByteTime = time.Now()
	}
	r.bytesRead += n
	r.readCalls++
	return n, err
}

func (r *ProxyReader) Stats() ReaderStats {
	r.Lock()
	defer r.Unlock()

	now := time.Now()

	var bytesPerSecond float64
	if !now.Equal(r.firstByteTime) {
		bytesPerSecond = float64(r.bytesRead) / now.Sub(r.firstByteTime).Seconds()
	}

	return ReaderStats{
		BytesRead:      r.bytesRead,
		BytesPerSecond: bytesPerSecond,
		ReadCalls:      r.readCalls,
	}
}
