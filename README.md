gobinary
===
package for binary header in go.

## install
`go get github.com/asxalex/gobinary`

## usage
```go
package main

import (
	"fmt"
	"github.com/asxalex/gobinary"
)

func main() {
	header := gobinary.NewBinaryHeader()
	header.AddField("version", 2)  // 2 bit for version
	header.AddField("MF", 1)       // 1 bit for M(ore)F(ragment)
	header.AddField("reserve1", 5) // 5 bit for reservation

	header.SetBitValue("version", 1) // set the version field to 1
	header.SetBitValue("MF", 1)      // set the MF to 1
	bin := header.ToBinary()         // get the binary according to the header
	fmt.Println(bin)

	// set the fields "versiong", "MF" and "reserve1"
	// accroding to the binary 0x80, that's to say, the
	// version field is set to 0b10, MF to 0b0, and reserve1 to 0b00000
	header.FromBinary([]byte{0x80})
}
```

## test and benchmark

```shell
$ go test -v -bench=.
=== RUN   TestBinary
--- PASS: TestBinary (0.00s)
goos: darwin
goarch: amd64
pkg: gobinary
BenchmarkBinaryConversion-4   	 5000000	       362 ns/op
PASS
ok  	gobinary	2.189s
```

## License
BSD
