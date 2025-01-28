# Binary Builder Thinkings

In Go, we have things like the StringBuilder

If you run the following code:

```go
func main() {
	var b strings.Builder
	for i := 3; i >= 1; i-- {
		fmt.Fprintf(&b, "%d...", i)
	}
	b.WriteString("ignition")
	fmt.Println(b.String())

}
```

You get the following output:

```
3...2...1...ignition
```

So perhaps this a good paradigm to emulate in gluing partial bytes together into
full bytes.

Psuedo code:
```
type BitBuilder struct {
    bytes []byte
    currByteLen uint8
    currByte int
}

var b bitBuilder
// b.bytes = [0b0000_0000]; currByteLen = 0 ; currByte = 0
b.addBits(0b0001, 4) // bson v 1
// b.bytes = [0b0001_0000]; currByteLen = 4 ; currByte = 0
b.addBits(0b001, 3) // Object
// b.bytes = [0b0001_0010]; currByteLen = 7 ; currByte = 0
b.addBits(0b00001, 5) // Object length
// b.bytes = [0b0001_0010, 0b0001_0000]; currByteLen = 4 ; currByte = 1
b.addBits(0b011, 3) // String
// b.bytes = [0b0001_0010, 0b0001_0110]; currByteLen = 7 ; currByte = 1
b.addBits(0b00011, 5) // String length 3
// b.bytes = [0b0001_0010, 0b0001_0110, 0b0011_0000]; currByteLen = 4 ; currByte = 2
b.addBytes([]byte("foo")) // "foo" = 0b01100110 0b01101111 0b01101111 
// b.bytes = [0b0001_0010, 0b0001_0110, 0b0011_0110, 0b0110_0110, 0b1111_0110, 0b1111_0000]; currByteLen = 4 ; currByte = 6
b.addBits(0b011, 3) // String
// b.bytes = [0b0001_0010, 0b0001_0110, 0b0011_0110, 0b0110_0110, 0b1111_0110, 0b1111_0110]; currByteLen = 7 ; currByte = 6
b.addBits(0b00011, 5) // String length 3
// b.bytes = [0b0001_0010, 0b0001_0110, 0b0011_0110, 0b0110_0110, 0b1111_0110, 0b1111_0110, 0b0011_0000]; currByteLen = 4 ; currByte = 7
b.addBytes([]byte("bar")) // "bar" = 0b01100010 0b01100001 0b01110010 
// b.bytes = [0b0001_0010, 0b0001_0110, 0b0011_0110, 0b0110_0110, 0b1111_0110, 0b1111_0110, 0b0011_0110, 0b0010_0110, 0b0001_0111, 0b0010_0000]; currByteLen = 4 ; currByte = 10

b.WriteTo(someFileWriter)
// Literally just passes b.bytes into someFileWriter.Write()
// Should return (10, nil)

```

01100110 01101111 01101111 - foo
01100010 01100001 01110010 - bar

