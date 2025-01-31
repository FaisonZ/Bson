# Bson V1 Specification

**Bson** - Bson is a binary encoding of JSON, hence the name Bson.

The general form of Bson is shown below:
```
[Bson version bits][Root Value][Value][Value][Value]....
```

If all the bits of a Bson encoding does not complete the final byte, zeros are
added to complete the last byte.

## Example with Explanations inline

Bson encoding of `{"foo":"bar"}` without explanation:
```
00010010 00010110 00110110 01100110 11110110
11110110 00110110 00100110 00010111 00100000
```

Bson encoding for `{"foo":"bar"}` with explanation
```
0001
^^^^ Version
001     00001
^Object ^length
  011     00011   01100110 01101111 01101111
  ^String ^Length ^------- "foo" ----------^
  011     00011   01100010 01100001 01110010
  ^String ^Length ^------- "bar" ----------^
0000
^ Extra bits to fill out the last byte
```

## Full breakdown

### Version

Each Bson file starts with 4 bits for the version.

```
0001    // Version bits (V1)
```

Using 4 bits means that this spec can have at most 16 versions (15, since I'm
skipping `0000`), but this is a learning project, so I don't expect any version
aside from version 1

### Root value

Following the version bits is the Root value.

Bson must have a root value that is either an Object or an Array. So the next
bits will be the encoding for that type of value

### Values

There are 7 Value types:
* Boolean
* Null
* String
* Array
* Object
* Float
* Integer

Each value follows the same general form:

* 3 bits for the Value type (defined in each Value type subsection)
* 5 bits for length if the Value type encodes length
* Bytes for the Value
  * If the Value has length bits preceeding it, then the next bytes matching
    that length will be the value
  * If the length preceeding a value is the max value of a 5 bit unsigned
    integer, then an additional set of length bits and value bytes follows

Each Value type subsection contains example encodings

#### Boolean

* Value type: `101`
* Value: `0` for false, or `1` for true

Encoding for true:
```
101 1   // Boolean token and true value
```

Encoding for false:
```
101 0   // Boolean token and false value
```

#### Null

* Value type: `110`

Null is only ever null, so the 3 bit Value type token is all that is needed.

Encoding for null:
```
110     // Null token
```

#### String

* Value type: `011`
* Length: `00000` - `11111`
* Value: "Length" number of bytes of the String
* Repeated Length and Value until the Length is less than 31 (`11111`)

Encoding for `"Cheese"`:
```
011 00110                   // String token and length(6)
01000011 01101000 01100101  // The 6 bytes of "Cheese"
01100101 01110011 01100101
```

Encoding for a string larger than 31 bytes.

`"abcdefghijklmnopqrstuvwxyz012345"`:
```
011 11111                   // String token and length(31)
01100001 01100010 01100011  // The first 31 bytes of the string
01100100 01100101 01100110 01100111 01101000 01101001 01101010 01101011
01101100 01101101 01101110 01101111 01110000 01110001 01110010 01110011
01110100 01110101 01110110 01110111 01111000 01111001 01111010 00110000
00110001 00110010 00110011 00110100
00001                       // length(1) of the next part of the string
00110101                    // The final byte of the string "5"
```

#### Array

* Value type: `010`
* Length: `00000` - `11111`
* Value: "Length" number of Values belonging to the array
* Repeated Length and Value until the Length is less than 31 (`11111`)

Encoding for `["a", "b", "c"]`
```
010 00011           // Array token and length(3)
011 00001 01100001  // [0]: "a"
011 00001 01100010  // [1]: "b"
011 00001 01100011  // [2]: "c"
```

#### Object

* Value type: `001`
* Length: `00000` - `11111`
* Value: "Length" number of pairs of String Value for object key and Value

As a convention, object Key-Value pairs are encoded in ascending order based on
the key

Before Bson V1 is finalized, I will consider removing the String token from
Object keys. Object keys in JSON are always strings, so might not need the token

Encoding for `{"a": null, "b": "bar"}`
```
001 00010                       // Object token and length(2)
011 00001 01100001              // Object key "a"
110                             // Value of "a": null
011 00001 01100010              // Object key "b"
011 00011 01100010 01100010     // Value of "b": "bb"
```

#### Float

* Value type: `111`
* Value: 8 bytes for a 64 bit float

All Floats are encoded in Network Order (Big Endian)

Encoding for `123.456`
```
111                 // Float token
01000000 01011110   // The 8 bytes for a 64 bit Float
11011101 00101111
00011010 10011111
10111110 01110111
```

Encoding for `123.45e+2`
```
111                 // Float token
01000000 11001000   // The 8 bytes for a 64 bit Float
00011100 10000000
00000000 00000000
00000000 00000000
```

#### Integer

* Value type: `100`
* Int Size Flag: `00` - 8 bit, `01` - 16 bit, `10` - 32 bit, `11` - 64 bit
* Value: 1, 2, 4, or 8 bytes for a signed integer corresponding to the size flag

All integers larger than 1 byte are encoded in Network Order (Big Endian)

Encoding for `-10`
```
100 00      // Integer token and 8 bit size flag
11110110    // The 8 bits for the signed integer -10
```

Encoding for `32021`
```
100 01              // Integer token and 16 bit size flag
01111101 00010101   // The 16 bits for the signed integer 32021
```
