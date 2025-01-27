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
2,11 `"bar"\n` 00100000 00100010 01100010 01100001 01110010 00100010 00001010
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
[bson-version]
[root-object]
[object-member[string"foo"][value-string"bar"]
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

