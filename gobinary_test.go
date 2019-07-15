package gobinary

import (
	"testing"
)

var (
	binaryCode = []byte{0x40, 0x80, 0xff}
)

func BenchmarkBinaryConversion(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		header.FromBinary(binaryCode)
		header.ToBinary()
	}
}

func TestBinary(t *testing.T) {
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
	header.FromBinary(binaryCode)
	b := header.ToBinary()
	if len(b) != len(binaryCode) {
		t.Log("length error")
		t.FailNow()
	}
	for idx, _ := range b {
		if b[idx] != binaryCode[idx] {
			t.Logf("%d mismatch: %d != %d", idx, b[idx], binaryCode[idx])
			t.FailNow()
		}
	}
}
