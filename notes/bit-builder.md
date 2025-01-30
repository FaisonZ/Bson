# Binary Builder Thinkings

## Initial Thoughts

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

## Strings

I've decided that strings will be length-prefixed, since UTF-8 makes it
impossible to depend on null terminating.

This leads to a possibility that strings will be encoded and end up larger than
if the string was just stored as plain text.

To get around this, I'm thinking of basically doing a sort of Rope:

```
[String Token][length bits][String bytes, up to max length][length bits][more string bytes]
```

At the moment, I'm using 5 bits to encode length for several things:
* Objects (how many keys)
* Object member keys (the strings)
* Arrays (how many elements)
* Strings (how many bytes)

A 5 bit integer has a max value of `2^5 - 1 = 31`. This works for the majority
of JSON object keys that I've seen, and will likely do well for the rest. But
the following string would break the current encoding:

```
abcdefghijklmnopqrstuvwxyz012345
```

The length of 32 will end up writing a length 0 (`0b0_0000`) and the string
bytes will then be interpretted as something other than strings.

### The Approach

For a string less than 31 bytes large:

* Get the length, which will be between 0 (`0b0_0000`) and (`0b1_1110`)
* Write the length bytes to the bson
* Write the string bytes to the bson
* Move along

For a string more than or equal to 31 bytes:

* Write a max string chunk value for the length to the bson: `0b1_1111`
* Write 31 bytes of the string to the bson
* Repeat, starting with the 32nd byte of the string
  * This could result in following the <31 bytes rules
  * It could result in the >=31 bytes rules

As an example, here's how the alphanumeric string pans out:
```
// abcdefghijklmnopqrstuvwxyz012345

011 11111   // String, 31 bytes in length
.........   // the first 31 bytes of the string (a...z0...4)
00001       // second length of the same string, 1 byte in length
.........   // The final byte of the string (5)
```

Then when decoding the string, it will need to follow these rules:
1. Find the string token
2. Get the length bits
3. Get the bytes equal to the length bits
4. If length bits was 31
  1. Go back to 2
  2. Append the next set of string bytes to the current set
5. Return the full string

And something will be needed for Objects and Arrays as well

## A bunch of binary thinking

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
0b0001_0010, 0b0001_0110, 0b0011_0110, 0b0110_0110, 0b1111_0110, 0b1111_0110, 0b0011_0110, 0b0010_0110, 0b0001_0111, 0b0010_0000
// b.bytes = [0b0001_0010, 0b0001_0110, 0b0011_0110, 0b0110_0110, 0b1111_0110, 0b1111_0110, 0b0011_0110, 0b0010_0110, 0b0001_0111, 0b0010_0000]; currByteLen = 4 ; currByte = 10

b.WriteTo(someFileWriter)
// Literally just passes b.bytes into someFileWriter.Write()
// Should return (10, nil)

```

01100110 01101111 01101111 - foo
01100010 01100001 01110010 - bar

