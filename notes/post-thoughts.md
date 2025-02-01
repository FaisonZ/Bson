# My time making Bson (Post Thoughts)

This is the story of how I ended up designing a specification for a binary
encoding of JSON.

## Old books and artisanal book binding

For some years now, I've spent time browsing the Internet Archive for old and
interesting books.

I make cheddar cheese based on a book from 1900, learned some woodcraft skills
from old proto Boy Scout organization handbooks, and most recently I started
learning Old English from a textbook in the Creative Commons.

Study for me has always been easier with a physical book, especially when you
need to flip between the front and back of a book frequently. Try doing that
with a PDF on your phone like I did, and you'll likely find that you've stopped
trying.

And so began my simple task: Figure out how to print a PDF as a book.

So on a weekend, I considered how pages could be printed, 2 to a side of paper,
and folded to make a book that is 5.5"x8.5". First think of how many pieces of
paper can be reasonably folded in half. I settled on 8. Then think up the
arrangement of the pages on front and back. I figured this out and made a script
that spat out the pages for front printings and back printings. After that,
print the page sets to separate PDFs that could then be printed 2 to a page.
Finally fold them, sew them, and feel good about life.

That was a fun solid day of effort there.

Well it turns out a friend of mine already has a book binding system and does it
way better than me. So now I'll just pay him to print and bind books for me.

But he has a problem: The software he uses doesn't let him adjust margins for
the book. So in some books, he'll have a lot of wasted space around the text,
and sometimes the text gets too close to the middle fold and becomes unreadable.

"Do you think you could figure this out?" he asks.

"Sure," I says. And so I start looking for the PDF specification.

## Sometimes a simple task is blocked by other simple tasks

I never really watched the show Malcolm in the Middle, but long ago I caught one
scene that frequently comes to mind when I'm doing things:

Hal, the father of the family, needs to replace a lightbulb. When he goes to get
a light bulb from a closet shelf, the shelf wobbles due to a broken support. So
he goes to get a screw driver from a drawer, which squeeks. On and on it goes
until he's working under his car and his wife asks him if he's replaced the
lightbulb yet. So Hal comes out from under the car and yells "What does it look
like I'm doing?!"

[Hal fixing a light bulb](https://www.youtube.com/watch?v=AbSehcT19u0)

Well I looked for the PDF specification, which Adobe requires you to buy for $0,
then I realized that I have never worked with file specifications outside of
HTML and JSON. And those are plain text, which I highly doubt PDFs are.

And that's the moment I decided to learn how to develop a binary encoding, and
write a program that can encode and decode it.

You know, instead of "buying" the PDF spec and reading it.

## A simple binary encoding project

I mentioned JSON before. It has a really straight-forward standard and I work
with it frequently. It is a plain text representation of data, and so it's a
little bulky. Maybe I could put together rules to encode JSON into binary and
save some bytes.

And so I started brainstorming over a Binary encoding of JSON, known as Bson.

## Artisanal bit packing

I hit the ground running and wrote notes that half covered the thoughts in my
head ([view my brainstorming notes at your own risk](/notes/brainstorming.md)].

There were a lot of thoughts and I started thinking with bytes to encode types
and lengths, but I realized quickly that whole bytes would probably make the
binary encoding **larger** than the original JSON.

So I looked at how many types of values there were, 7, and how many bits could
cover each type, 3. If you look at the V1 Spec, you'll see all the Value types,
but to give an example here, an object is `001` and a string is `011`.

And then I started considering a common form of JSON: `{"foo": "bar"}`.

I already have types defined for objects and strings, so the next step was to
define how the values were encoded. Well I learned some C a long time ago, so
obviously I can just use null to terminate a string, right?

For those who know before reading on, you know we'll come back to this later.

So for a string, the binary encoding will look something like this:
```
011 [bytes for the string] [Null Byte]
```

Cool. On to objects. Objects have an unknown number of fields in them. I suppose
I can also null terminate that. Then we can just encode a string for the object
key followed by the encoding of that field's value:

```
001 [encoding of "foo"] [encoding of "bar"] [Null Byte]
```

And arrays can work the same way. `["foo", "bar"]`:
```
010 [encoding of "foo"] [encoding of "bar"] [Null Byte]
```

This felt like a good start to this learning project.

## Don't let UTF-8 multi-byte you on the butt

Something in the back of my mind started to worry me. And mind you, I have yet
to code at this point, so it's the best time to worry!

UTF-8 characters can span multiple bytes. What stops a byte inside of a UTF-8
character from being a null byte?

As best as I could find, nothing. Which means I can't depend on null bytes to
terminate a string. So back to my plan!

After that null byte research, I decided to go with length prefixing for strings,
objects, and arrays.

And now things look and feel more realistic:
```
001 [length 1] [encoding of key "foo"] [encoding of value "bar"]
```

I'll know exactly how many array items or object fields follow in the encoding,
so decoding Bson should be a breeze!

But how many bits should I use for my length? Well, I decided on 5. Not because
31 seems like a reasonable max length for most strings (think of object keys).
I purely chose it because when added to the 3 bits for the Value type, it's a
whole byte.

Purely asthetics, and it definitley won't come back to bite me later.

Except later was later that night when I was walking my baby around the house.
What happens if you have a string that's 32 bytes long?

Luckily you can get a lot of thinking in when a baby is too restless to sleep
and demands a 20 minute walk with papa at 10 pm.

The solution was simple: Encode the string in chunks of 31 bytes max. Just look
at the length of the string, if it's 31, then expect another length encoded
after the 31 bytes of string data, followed by more bytes of the string. Keep
doing that until the length is less than 31 and then you have the whole string.

So a small string would look like this:
```
011 [length < 31] [String bytes matching length]
```

And a big string would look like this:
```
011 [length = 31] [31 bytes of string] [length < 31] [The final bytes of the string]
```

While figuring this out, it became obvious that I have the exact same problem
with Object fields and Array elements. So we just do the exact same thing.

Object with a lot of fields:
```
001 [length = 31] [key-value pair of 31 fields] [length < 31] [key-value pair of the rest]
```

Array with a lot of elements:
```
010 [length = 31] [31 items] [length < 31] [the rest of the items]
```

## Laying the foundation

Outline


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

* Decided to keep floats as 64 bit

* Wrote the spec

* The end
