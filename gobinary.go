package gobinary

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"
)

var (
	ErrHeaderExists    = errors.New("header name exists")
	ErrHeaderNotExists = errors.New("header name does not exist")
	ErrValueNotMatch   = errors.New("value not match")
	ErrLengthNotMatch  = errors.New("length not match")
)

type Bit struct {
	Value uint32
	Digit uint8
}

type Header struct {
	BitMap   map[string]*Bit
	BitSlice []*Bit
	Length   int16
	Mu       *sync.RWMutex
}

func (h *Header) FromBinary(b []byte) error {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	if int(h.Length) != len(b)*8 {
		return ErrLengthNotMatch
	}
	index := 0
	for _, v := range h.BitSlice {
		f := 32 - v.Digit
		tmp := make([]byte, 4)
		for i := 0; i < int(v.Digit); i++ {
			first := int(f) + int(i)
			firstByte := first / 8
			firstBit := first % 8
			indexByte := (index + i) / 8
			indexBit := (index + i) % 8
			f := tmp[firstByte]
			bit := (b[indexByte] >> (7 - uint(indexBit))) & 0x1
			if bit == 1 {
				f |= (0x1 << (7 - uint(firstBit)))
				tmp[firstByte] = f
			}
		}
		v.Value = ByteToUint32(tmp)
		index += int(v.Digit)
	}
	return nil
}

func (h *Header) GetValue(field string) (uint32, bool) {
	h.Mu.RLock()
	defer h.Mu.RUnlock()
	if v, ok := h.BitMap[field]; ok {
		return v.Value, true
	}
	return 0, false
}

func (h *Header) ToBinary() []byte {
	h.Mu.RLock()
	defer h.Mu.RUnlock()
	bitLength := 0
	byteLength := h.Length / 8
	if h.Length%8 != 0 {
		byteLength++
	}
	result := make([]byte, byteLength)
	for _, v := range h.BitSlice {
		first := bitLength
		b := v.GetBinary()
		for i := 32 - v.Digit; i < 32; i++ {
			firstByte := first / 8
			firstBit := first % 8
			byteIndex := i / 8
			bitIndex := i % 8
			f := result[firstByte]
			if (b[byteIndex]>>(7-bitIndex))&0x1 == 1 {
				f |= 0x1 << (7 - uint(firstBit))
			} else {
				f &= (0xff ^ (0x1 << (7 - uint(firstBit))))
			}
			result[firstByte] = f
			first++
		}
		bitLength += int(v.Digit)
	}
	return result
}

func DumpByte(b []byte) {
	fmt.Println("dump byte: ")
	for i := 0; i < len(b)*8; i++ {
		byteIdx := i / 8
		bitIdx := i % 8
		bb := b[byteIdx]
		value := (bb & (0x1 << (7 - uint(bitIdx)))) >> (7 - uint(bitIdx))
		fmt.Printf("%d ", value)
	}
	fmt.Println()
}

func (h *Header) DumpValue() {
	h.Mu.RLock()
	defer h.Mu.RUnlock()
	for idx, v := range h.BitSlice {
		fmt.Printf("%02dth(%d bit): ", idx, v.Digit)
		b := v.GetBinary()
		for i := 32 - v.Digit; i < 32; i++ {
			byteidx := i / 8
			bitidx := i % 8
			fmt.Printf("%d", (b[byteidx]&(0x1<<(7-bitidx)))>>(7-bitidx))
		}
		fmt.Println()
	}
}

func (h *Header) SetBitValue(name string, value uint64) error {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	if b, ok := h.BitMap[name]; ok {
		if float64(value) >= math.Pow(2, float64(b.Digit)) {
			return ErrValueNotMatch
		}
		b.Value = uint32(value)
		return nil
	}
	return ErrHeaderNotExists
}

func (h *Header) AddField(name string, bitlength uint8) error {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	if _, ok := h.BitMap[name]; ok {
		return ErrHeaderExists
	}
	b := &Bit{
		Value: 0,
		Digit: bitlength,
	}
	h.Length += int16(bitlength)
	h.BitMap[name] = b
	h.BitSlice = append(h.BitSlice, b)
	return nil
}

func NewBinaryHeader() *Header {
	header := &Header{
		BitMap:   make(map[string]*Bit),
		BitSlice: []*Bit{},
		Length:   0,
		Mu:       new(sync.RWMutex),
	}
	return header
}

func MakeBit(value uint32, digit uint8) Bit {
	return Bit{
		Value: value,
		Digit: digit,
	}
}

type MyBit struct {
	BB Bit
}

func ByteToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func (b *Bit) GetBinary() []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, b.Value)
	return result
}

func main() {
	header := NewBinaryHeader()
	header.AddField("version", 2)
	header.SetBitValue("version", 1)
	header.AddField("MF", 1)
	header.AddField("ack", 1)
	header.AddField("resend", 1)
	header.AddField("reserve1", 3)
	header.AddField("length", 5)
	header.AddField("reserve2", 3)
	header.AddField("id", 7)
	header.AddField("check", 1)
	header.FromBinary([]byte{0x80, 0, 0})
	header.DumpValue()
	b := header.ToBinary()
	fmt.Println(b)
	DumpByte(header.ToBinary())
}
