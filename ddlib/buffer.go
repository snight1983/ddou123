package ddlib

var bufCache *Queue = NewQueue()

// CBuffer .
type CBuffer struct {
	Buf  []byte
	Size uint32
}

// NBuffer .
func NBuffer() *CBuffer {
	buf := bufCache.Pop()
	if nil == buf {
		buf = &CBuffer{}
	}
	return buf.(*CBuffer)
}

// FBuffer .
func FBuffer(obj *CBuffer) {
	obj.Buf = obj.Buf[:0]
	obj.Size = 4
	bufCache.Push(obj)
}

// AppendInt32 添加 AppendInt32
func (buf *CBuffer) AppendInt32(value int32) {
	//Buf.AppendInt32
	//binary.Write(buf.Buf[buf.Size:], binary.BigEndian, &value)
}
