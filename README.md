# Binary JSON (bson)

A learning experience with converting things into binary and then decoding back
from binary

## The idea

If you start with something like this (which shows as 21 bytes large):

```json
{
    "foo": "bar"
}
```

Could you encode that into a binary format that is smaller than the source
JSON?

If you hexdump the file:

```
$ xxd -b file.json
00000000: 01111011 00001010 00100000 00100000 00100000 00100000  {.    
00000006: 00100010 01100110 01101111 01101111 00100010 00111010  "foo":
0000000c: 00100000 00100010 01100010 01100001 01110010 00100010   "bar"
00000012: 00001010 01111101 00001010                             .}.
```

Here's how it lines up line by line:
```
1,1  `{\n`     01111011 00001010
2,1  `    `    00100000 00100000 00100000 00100000
2,5  `"foo":`  00100010 01100110 01101111 01101111 00100010 00111010
2,11 ` "bar"\n` 00100000 00100010 01100010 01100001 01110010 00100010 00001010
3,1  `}\n`     01111101 00001010
```

What if we used some specific binary "keys" to flag what is about to start?

| JSON Thing | Description | aaa |
| -- | -- | -- |
| Root Object | The root level object | tk |
| Root Array | The root level array | tk |
| Object Member | A member of an object ("<string>": <value>) | tk |
| Value | A value, either in a Root Array, Object Member, or Array Value | tk |
| Array | an ordered list of Values ([<value>, <value>, <value>]) | tk |
| String | A string, was surrounded by double quotes ("foobar")| tk |

Thinking with pseudo packing
```
[bson-version(4b)]
[root-object(1B)]
[object-member(1B)[string"foo"(1+3B)][value-string"bar"(1+3B)]
```
1. root
2. object
3. array
4. string
* number
    5. int
    6. float
7. true
8. false
9. null
10. object member
11. NB

1. 001 - object
    1. length (key count)
    2. object members
2. 010 - array
    1. length prefix
    2. values
3. 011 - string
    1. length prefix
    2. value
4. 100 - number
    1. num type
        1. 0 - int
        2. 1 - float
    2. value
5. 101 - boolean
    1. 0 - false
    2. 1 - true
6. 110 - null
7. 111 - unused

```
[bson-version]
[root]
[object token]
    [object-member token]
        [string"foo"][NB (closes string)]
        [string"bar"][NB (closes string and object-member)]
[NB (closes object token)]


0000        [bson-version]
001         [object token]
011 . 000       [string"foo"][NB (closes string)]
011 . 000           [string"bar"][NB (closes string and object-member)]
000         [NB (closes object token)]
```
01100110 01101111 01101111 - foo
01100010 01100001 01110010 - bar

```
0000                                    [bson-version]
001                                     [object]
0001                                    [object length (1)]
011 00000011 01100110 01101111 01101111 [string(3)"foo"]
011 00000011 01100010 01100001 01110010 [string(3)"bar"]


0000 001 0001 011 00000011 01100110 01101111 01101111 011 00000011 01100010 01100001 01110010
00000010 00101100 00001101 10011001 10111101 10111101 10000001 10110001 00110000 10111001 0
vs.
01111011 00100010 01100110 01101111 01101111 00100010 00111010 00100010 01100010 01100001 01110010 00100010 01111101



0000                                    [bson-version]
001 00001                               [object(1)]
011 00011 01100110 01101111 01101111    [string(3)"foo"]
011 00011 01100010 01100001 01110010    [string(3)"bar"]

0000 001 00001 011 00011 01100110 01101111 01101111 011 00011 01100010 01100001 01110010
\x02     \x16     \x36     \x66     \xF6     \xF6     \x36     \x26     \x17     /x20
\x02\x16\x36\x66\xF6\xF6\x36\x26\x17/x20
00000010 00010110 00110110 01100110 11110110 11110110 00110110 00100110 00010111 00101111 01111000 00110010
00000010 00010110 00110110 01100110 11110110 11110110 00110110 00100110 00010111 00101111 01111000 00110010 00110000
00000010 00010110 00110110 01100110 11110110 11110110 00110110 00100110 00010111 00100000
vs.
00000010 00101100 00001101 10011001 10111101 10111101 10000001 10110001 00110000 10111001 0
vs.
01111011 00100010 01100110 01101111 01101111 00100010 00111010 00100010 01100010 01100001 01110010 00100010 01111101

```

```json
{
    "version": 20,
    "names": [
        "bob",
        "bill",
        "frank"
    ],
    "office": {
        "street1": "123 test dr",
        "street2": null,
        "city": "Testington",
        "state": "WI",
        "zip": "53120"
    },
    "rate": 2.14,
    "active": true,
    "newsletterConsent": false
}
```

```
[bson-version]
[root]
[object token]
    [object-member token]
        [string"version"][NB (closes string)]
        [int"20"][NB (closes int and object-member)]
    [object-member token]
        [string"names"][NB (closes string)]
        [array]
            [string"bob"][NB (closes string)]
            [string"bill"][NB (closes string)]
            [string"frank"][NB (closes string)]
            [NB (closes array and object-member)]
    [object-member token]
        [string"office"][NB (closes string)]
        [object token]
            [object-member token]
                [string"street1"][NB (closes string)]
                [string"123 test dr"][NB (closes string and object-member)]
            [object-member token]
                [string"street2"][NB (closes string)]
                [null][NB (closes object-member)]
            [object-member token]
                [string"city"][NB (closes string)]
                [string"Testington"][NB (closes string and object-member)]
            [object-member token]
                [string"state"][NB (closes string)]
                [string"WI"][NB (closes string and object-member)]
            [object-member token]
                [string"zip"][NB (closes string)]
                [string"53120"][NB (closes string and object-member)]
            [NB (closes object token and object-member)]
    [object-member token]
        [string"rate"][NB (closes string)]
        [float"2.14"][NB (closes int and object-member)]
    [object-member token]
        [string"active"][NB (closes string)]
        [true][NB (closes object-member)]
    [object-member token]
        [string"newsletterConsent"][NB (closes string)]
        [false][NB (closes object-member)]
[NB (closes object token)]




0001                                                            [bson-version]
001 00110                                                       [object(6)]
011 00111 01110110 01100101 01110010 01110011 01101001          [string(7)"version"]
011 00010 00110010 00110000                                     [string(2)"20"]
011 00101 01101110 01100001 01101101 01100101 01110011          [string(5)"names"]
010 00011                                                       [array(3)]
011 00011 01100010 01101111 01100010                            [string(3)"bob"]
011 00100 01100010 01101001 01101100 01101100                   [string(4)"bill"]
011 00101 01100110 01110010 01100001 01101110 01101011          [string(5)"frank"]
011 00110 01101111 01100110 01100110 01101001 01100011 01100101 [string(6)"office"]
001 00101                                                       [object(5)]
011 00111 01110011 01110100 01110010 01100101 01100101 01110100 00110001    [string(7)"street1"]
011 01011 00110001 00110010 00110011 00100000 01110100 01100101 01110011 01110100 00100000 01100100 01110010   [string(11)"123 test dr"]
011 00111 01110011 01110100 01110010 01100101 01100101 01110100 00110010    [string(7)"street2"]
110                                                             [null]
011 00100 01100011 01101001 01110100 01111001                   [string(4)"city"]
011 01010 01010100 01100101 01110011 01110100 01101001 01101110 01100111 01110100 01101111 01101110 [string(10)"Testington"]
011 00101 01110011 01110100 01100001 01110100 01100101          [string(5)"state"]
011 00010 01010111 01001001                                     [string(2)"WI"]
011 00011 01111010 01101001 01110000                            [string(3)"zip"]
011 00101 00110101 00110011 00110010 00110001 00110000          [string(5)"53120"]
011 00100 01110010 01100001 01110100 01100101                   [string(4)"rate"]
011 00100 00110010 00101110 00110001 00110100                   [string(4)"2.14"]
011 00110 01100001 01100011 01110100 01101001 01110110 01100101 [string(6)"active"]
101 1                                                           [bool][true]
011 10001 01101110 01100101 01110111 01110011 01101100 01100101 01110100 01110100 01100101 01110010 01000011 01101111 01101110 01110100 01100101 01101110 01110100 [string(17)"newsletterConsent"]
101 0                                                           [bool][false]

0001 001 00110 011 00111 01110110 01100101 01110010 01110011 01101001 011 00010 00110010 00110000 011 00101 01101110 01100001 01101101 01100101 01110011 010 00011 011 00011 01100010 01101111 01100010 011 00100 01100010 01101001 01101100 01101100 011 00101 01100110 01110010 01100001 01101110 01101011 011 00110 01101111 01100110 01100110 01101001 01100011 01100101 001 00101 011 00111 01110011 01110100 01110010 01100101 01100101 01110100 00110001 011 01011 00110001 00110010 00110011 00100000 01110100 01100101 01110011 01110100 00100000 01100100 01110010 011 00111 01110011 01110100 01110010 01100101 01100101 01110100 00110010 110 011 00100 01100011 01101001 01110100 01111001 011 01010 01010100 01100101 01110011 01110100 01101001 01101110 01100111 01110100 01101111 01101110 011 00101 01110011 01110100 01100001 01110100 01100101 011 00010 01010111 01001001 011 00011 01111010 01101001 01110000 011 00101 00110101 00110011 00110010 00110001 00110000 011 00100 01110010 01100001 01110100 01100101 011 00100 00110010 00101110 00110001 00110100 011 00110 01100001 01100011 01110100 01101001 01110110 01100101 101 1 011 10001 01101110 01100101 01110111 01110011 01101100 01100101 01110100 01110100 01100101 01110010 01000011 01101111 01101110 01110100 01100101 01101110 01110100 101 0

\x12\x66\x77\x66\x57\x27\x36\x96
0001001001100110011101110110011001010111001001110011011010010110
\x23\x23\x06\x56\xE6\x16\xD6\x57
0010001100100011000001100101011011100110000101101101011001010111
\x34\x36\x36\x26\xF6\x26\x46\x26
0011010000110110001101100010011011110110001001100100011000100110
\x96\xC6\xC6\x56\x67\x26\x16\xE6
1001011011000110110001100101011001100111001001100001011011100110

1011011001100110111101100110011001100110100101100011011001010010

0101011001110111001101110100011100100110010101100101011101000011

0001011010110011000100110010001100110010000001110100011001010111

0011011101000010000001100100011100100110011101110011011101000111

0010011001010110010101110100001100101100110010001100011011010010

1110100011110010110101001010100011001010111001101110100011010010

1101110011001110111010001101111011011100110010101110011011101000

1100001011101000110010101100010010101110100100101100011011110100

1101001011100000110010100110101001100110011001000110001001100000

1100100011100100110000101110100011001010110010000110010001011100

0110001001101000110011001100001011000110111010001101001011101100

1100101101101110001011011100110010101110111011100110110110001100

1010111010001110100011001010111001001000011011011110110111001110

1000110010101101110011101001010
```

So values would not be packable. And if we ignore whitespace in the JSON, then
we can save a lot of bytes around things like spaces, tabs and newlines. Let
people use their own formatters!

## Commands

### bson encode

```
bson encode file.json file.bson
```

encodes **file.json** into **file.bson**

### bson decode

```
bson decode file.bson file.json
```

decodes **file.bson** into **file.json**

