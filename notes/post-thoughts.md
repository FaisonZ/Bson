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
* Decided to make a binary encoding of JSON, and call it Bson (pronounced Bison)

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

