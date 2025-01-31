# My time making Bson (Post Thoughts)

Outline

* Wanted to print books from books in the creative commons
* Printing two to a page was fine, and came up with a way to sort the papges
* But margins were an issue
* Decided I needed to make an app that takes care of the page sorting and allows
for margin manipulation

* PDF has its own spec (major versions 1 and 2)
* Decided that I need some prep before even looking into that
* Needed a plan on how to think around binary specs and how to encode/decode a
binary file to a spec

* I know JSON and know it's one of the least optimal ways you can transfer data
from a bandwidth point of view
* Decided to make a binary encoding of JSON, and call it Bson (pronounced
Bison, sometimes Bee-sahn (like Nisan), dealer's choice)

* Started by counting out the types available in JSON
* The number can fit in 3 bits
* Started playing around with how to encode `{"foo":"bar"}`
* Decided null termination would be great
* Learned that null termination is not great
* Decided I need to use length-prefixing
* But how long? How about 5 bits, max 31
* After some though, decided I'd make a rope-like encoding:
  * `[length][bytes up to the prev length][next length][bytes up to the next length]...`
* Turns out I also should use length-prefixing for objects and arrays
* Got started

* BitBuilder first
* Then encoder

* Then decoder
* But I need a BitReader first
* Then decoder again

* When decoding an object, I ran into issues.
* Realized I skipped the value token for object keys
* So put a line to quick consume the bits, but ignore them
* Future improvement: All object keys are strings in JSON, so don't encode the
string token for object keys.

* Then I slapped together some JSON files, ran them forwards and backwards and
got the correct results both ways!
* Made a `check` command to show the amount of bytes the JSON and corresponding
Bson were, and the difference so I can see how much space I'm saving
* Something around 20%, but this is an excercise in working with binary and
specs

* Then I used some website to generate a really big JSON file with an array of
large objects
* This JSON file would encode, but decoding failed
* Then I remembered strings longer than 31 bytes
* Decided to do the following:
  * Get length of string in bytes
  * If > 31, then encode the length 31 and the first 31 bytes of the string
  * Then repeat the process with the rest of the string following those 31 bytes
  * Repeat until we have between 0 and 30 bytes, which is the last part encoded

* It worked!
* Did the same for Objects
* And the same for Arrays
* Now the array of large objects encodes and decodes successfully
* File reduction of 26%! Larger than I was expecting, since I'm leaving string
bytes exactly as they are (no compression)

* Had a thought, what if I collected strings and shoved them at the back of the
bson bytes?
* Flag the "id" of the string in the JSON encoded section
* And had an array of strings at the bottom
* Possibly help when object keys are repeated, or other things
* So just a thought
```
[bson version bits]
[bson encoded JSON bytes]...
... [String member of array, but instead of storing "applesauce", just 5]
...
[post bson encoded JSON bytes]
[encoded array of strings]
...[5th element: "applesauce"]
```
* Probably only worth doing if it doesn't increase bytes by a lot in JSON
without repeated strings or object key names
* And I need to make sure I can encode numbers well

* By the way. You might have noticed that I ignored numbers until now
* Started with ints
* Wanted to store a number in the smallest signed integer that it fits in
* as small as 8 bit, as large as 64 bit
* So had to store some bits to flag the size
* `0b00` for 8, `0b01` for 16, `0b10` for 32, and `0b11` for 64
* This gives me the fun of making a breaking future change if we get 128 bit
as a common integer size in the future.
* But Golang defaults to 64 bit when decoding JSON, so that works for me too
* First time had to look into Endianness.
* Decided to go with Big Endian, because network order
* It basically worked without issue
* One note, to follow in Golang's footsteps, I decided to return all decoded
ints as int64.
* The size int was still needed for decoding the bytes though, so no effort
wasted there

* Floats... There's a problem here
* If a JSON file has the value 3.0, my current code will encode that as an int
* If you then decode, the value will be 3
* On top of that, the max value of a int64 is larger than a float64, because of
fancy math reasons that I'm not going to read enough to understand
* So while looking into this issue, I learned about json.RawMessage
* If you Unmarshal json into a `map[string]json.RawMessage`, or into a
`[]json.RawMessage`, then you get the string as it was found in the JSON.
* Strings will be strings, numbers will be strings, arrays will be strings
including everything contained (including whitespace/new lines)

```go
jsonBlob := []byte{`{"a":1, "b": "bb", "c": {"foo":   "bar"}, "d": [  1],"e":1.23}`}
var d map[string]json.RawMessage
_ = json.Unmarshal(jsonBlob, &d)
for key, v := range d {
    fmt.Printf("[%s]: %s\n", key, v)
}
```

Outputs:
```
[a]: 1
[b]: "bb"
[c]: {"foo":   "bar"}
[d]: [  1]
[e]: 1.23
```

* Keeping the exact values is important to me
* If I encode a json file into bson, then decode that bson, I should end up with
the exact same JSON
  * Note: The order of object keys will not be guaranteed, because maps
* So I'm thinking now that I need to handle json Unmarshalling more closely
* Time to refactor before moving on to floats
