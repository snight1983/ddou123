package vpr

//https://github.com/liuxp0827/govpr.git
import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"os"
)

type verFile struct {
	file       *os.File
	readwriter *bufio.ReadWriter
}

func newVPerile(filename string) (*verFile, error) {

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	vprFile := &verFile{
		file:       file,
		readwriter: bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(file)),
	}

	return vprFile, nil
}

func (f *verFile) putInt(v int) (int, error) {
	var intBuf [4]byte

	data := intBuf[:4]
	data[0] = byte(v & 0xff)
	data[1] = byte((v >> 8) & 0xff)
	data[2] = byte((v >> 16) & 0xff)
	data[3] = byte((v >> 24) & 0xff)
	return f.readwriter.Write(data)
}

func (f *verFile) putByte(v byte) error {
	return f.readwriter.WriteByte(v)
}

func (f *verFile) putFloat64(v float64) (int, error) {
	var float64Buf [8]byte
	data := float64Buf[:8]
	putFloat64LE(data, v)
	return f.readwriter.Write(data)
}

func (f *verFile) getInt() (int, error) {

	var v uint32
	binary.Read(f.readwriter, binary.LittleEndian, &v)
	return int(v), nil
}

func (f *verFile) getByte() (byte, error) {
	return f.readwriter.ReadByte()
}

func (f *verFile) getFloat64() (float64, error) {
	var float64Buf [8]byte

	data := float64Buf[:8]
	_, err := io.ReadFull(f.readwriter, data)
	if err != nil {
		return .0, err
	}

	return getFloat64LE(data), nil
}

func (f *verFile) getFloat32() (float32, error) {
	var floatBuf [4]byte

	data := floatBuf[:4]
	_, err := io.ReadFull(f.readwriter, data)
	if err != nil {
		//log.Error(err)
		return .0, err
	}

	return getFloat32LE(data), nil
}

func (f *verFile) close() error {
	err := f.readwriter.Flush()
	if err != nil {
		return err
	}
	return f.file.Close()
}

func getUint16LE(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func putUint16LE(b []byte, v uint16) {
	binary.LittleEndian.PutUint16(b, v)
}

func getUint16BE(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func putUint16BE(b []byte, v uint16) {
	binary.BigEndian.PutUint16(b, v)
}

func getUint32LE(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func putUint32LE(b []byte, v uint32) {
	binary.LittleEndian.PutUint32(b, v)
}

func getUint32BE(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func putUint32BE(b []byte, v uint32) {
	binary.BigEndian.PutUint32(b, v)
}

func getUint64LE(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func putUint64LE(b []byte, v uint64) {
	binary.LittleEndian.PutUint64(b, v)
}

func getUint64BE(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func putUint64BE(b []byte, v uint64) {
	binary.BigEndian.PutUint64(b, v)
}

func getFloat32BE(b []byte) float32 {
	return math.Float32frombits(getUint32BE(b))
}

func putFloat32BE(b []byte, v float32) {
	putUint32BE(b, math.Float32bits(v))
}

func getFloat32LE(b []byte) float32 {
	return math.Float32frombits(getUint32LE(b))
}

func putFloat32LE(b []byte, v float32) {
	putUint32LE(b, math.Float32bits(v))
}

func getFloat64BE(b []byte) float64 {
	return math.Float64frombits(getUint64BE(b))
}

func putFloat64BE(b []byte, v float64) {
	putUint64BE(b, math.Float64bits(v))
}

func getFloat64LE(b []byte) float64 {
	return math.Float64frombits(getUint64LE(b))
}

func putFloat64LE(b []byte, v float64) {
	putUint64LE(b, math.Float64bits(v))
}

func uvarintSize(x uint64) int {
	i := 0
	for x >= 0x80 {
		x >>= 7
		i++
	}
	return i + 1
}

func varintSize(x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return uvarintSize(ux)
}

func getUvarint(b []byte) (uint64, int) {
	return binary.Uvarint(b)
}

func putUvarint(b []byte, v uint64) int {
	return binary.PutUvarint(b, v)
}

func getVarint(b []byte) (int64, int) {
	return binary.Varint(b)
}

func putVarint(b []byte, v int64) int {
	return binary.PutVarint(b, v)
}

func readUvarint(r io.ByteReader) (uint64, error) {
	return binary.ReadUvarint(r)
}

func readVarint(r io.ByteReader) (int64, error) {
	return binary.ReadVarint(r)
}
