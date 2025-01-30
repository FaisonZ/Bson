# Numbers

My main professional languages have been PHP and Javascript, so I've never
really had to worry about integers and floats, let alone 8 bit, 16 bit, 32 bit,
or even 64 bit.

Which means I need to learn about these things deeply and come up with a way to
encode numbers in a JSON file in a way that isn't wasteful in bytes, but doesn't
lose things around float precision (if that even is a worry).

## Starting thoughts

Ideally I can encode something as an int or a float, and somehow call out the
byte size (8/16/32/64).

Golang unmarshalls all JSON numbers as float64, which means I need to come up
with a way to determine if it's an int or not.

I think I could confidently encode an int with a specific byte size. I'm not as
confident with floats though, so I might just leave them at 64 bits unless I
learn something better.

### Psuedo code

* Have float64, and tested to only be an int
* Check what size int it fits into (only signed ints)
* Convert the float64 into the proper sized int
* Add "Number" token bits
* Add "Int<bit size>" token bits
* Add the int's bytes

* Have float64, and tested to not be an int
* Add "Number" token bits
* Add "Float64" token bits
* Add the float's bytes

Tokens:
* Int (`0b100`)
  * Int8  (`0b00`)
  * Int16 (`0b01`)
  * Int32 (`0b10`)
  * Int64 (`0b11`)
* Float (`0b111`)
  * Float32 (`0b00`) (not used yet)
  * Float64 (`0b01`)

Examples?
| Number | Type | Type Token | Size Token | Number bytes | Combined | JSON bytes | Bson Bytes |
| - | - | - | - | - | - | - | - |
| 0 | Int8 | `0b100` | `0b00` |  `0b0000_0000` | `0b1000_0000, 0b0000_0000` | 1 | 2 |
| 10 | Int8 | `0b100` | `0b00` |  `0b0000_1010` | `0b1000_0000, 0b0101_0000` | 2 | 2 |
| -10 | Int8 | `0b100` | `0b00` | `0b1111_0110` | `0b1000_0111, 0b1011_0000` | 3 | 2 |
| 32021 | Int16 | `0b100` | `0b01` | `0b0111_1101, 0b0001_0101` | `0b1000_1011, 0b1110_1000, 0b1010_1000` | 5 | 3 |

Since I'm too lazy to figure out the bits myself and the website I used won't
work above 53 bits or with floats, I'm leaving the examples at that.

JSON is plain text, meaning every character used in a number is a byte. So as
the number gets larger, the more bytes bson saves.

* 0 is 1 byte in JSON, but 2 in Bson
* 10 is 2 bytes in JSON, but 2 in Bson
* -500 is 4 bytes in JSON, but 3 in Bson
* 65000 is 5 bytes in JSON, but 3 in Bson
* -65000 is 6 bytes in JSON, but 3 in Bson
* 1503238552 is 10 bytes in JSON, but 5 in Bson
* 5000000000000000000 is 19 bytes in JSON, but 9 bytes in Bson

