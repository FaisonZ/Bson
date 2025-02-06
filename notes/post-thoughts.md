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
way better than me. So now I'll just pay him to print and bind books for me
instead.

But he has a problem: The software he uses doesn't let him adjust margins for
the book. So in some books, he'll have a lot of wasted space around the text,
and sometimes the text gets too close to the middle fold and becomes unreadable.

"Do you think you could figure this out?" he asks.

"Sure," I says. And so I start looking for the PDF specification.

## One small favour

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
head ([view my brainstorming notes at your own risk](/notes/brainstorming.md)).

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

Finally, to be safe, I decided to start the Bson encoding with 4 bits for the
Bson version used for encoding.

This felt like a good start to this learning project.

## Don't let UTF-8 multi-byte you

Something in the back of my mind started to worry me. And mind you, I have yet
to code at this point, so it's the best time to worry!

UTF-8 characters can span multiple bytes. What stops a byte inside of a UTF-8
character from being a null byte?

As best as I could find, nothing. Which means I can't depend on null bytes to
terminate a string. So back to my plan!

After that null byte research, I determined that length prefixing is the way to
go for strings, objects, and arrays.

And now things look and feel more realistic:
```
001 [length 1] [encoding of key "foo"] [encoding of value "bar"]
```

I'll know exactly how many array items or object fields follow in the encoding,
so decoding Bson should be a breeze!

But how many bits should I use for my length? Well, I chose 5. Not because 31
seems like a reasonable max length for most strings (think of object keys). I
purely chose it because when added to the 3 bits for the Value type, it's a
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

And so I started coding and I began with the encoder. But before starting the
encoder, I needed to be able to write binary with bits.

I wrote out how I wanted to accomplish this in my [notes on Bit Builder](/notes/bit-builder.md)
(Note, I have no idea why I wrote a section on long strings in that file, but if
you skip the Strings section, you get a full picture of my plan for Bit Builder).

Most of the bits I plan to write are less than a byte, so I had to learn proper
bit shifting. For example, to write my version bits `0001`, I need to shift 4
to the left:

```
// Starts with a full byte:
0b0000_0001

// Then shift to the left:
0b0000_0001 << 4
0b0001_0000

// Finally, get the bites into a slice of bytes
bytes[0] |= 0b0001_0000

// Keep track of the next bit available for writing, bit 5
0b0001_0000
//     ^ That one
```

And the next bits I add will always be either the Object or Array value type bits

```
// The object type bits (001)
0b0000_0001

// We need the three bits to be written starting at the 5th bit
0b0000_0001 << 1
0b0000_0010

// Write those bits into the same byte the version bits are in:
bytes[0] |= 0b0000_0010

// Track the next available bit, bit 8
0b0001_0010
//        ^ That one
```

This was a fun task to accomplish. I had to keep track of my position in a Golang
slice of bytes (`[]byte`), and my bit position in that byte. And I had to figure
out how to properly write a collection of bits between two bytes when it happens
to span the existing byte and a new one.

Take a look at the code if you're interested. I was pretty happy with the results!

With Bit Builder in hand, I started writing the encoder.

In this, I learned a lot about how Go "Unmarshals" JSON into maps, array, strings,
etc.. And there's something interesting about writing code to turn unknown types
into a proper Go type before deciding what to do with it. Maybe not the best use
of a strongly typed language, but it wasn't too painful.

Either way, I took poor notes on what happened in this phase, but I just remember
avoiding numbers like the plague. When I got to the point that I could take
`{"foo": "bar"}` and spit out binary, I had a line of code that would print `"I
haven't implemented numbers yet!"` to the console if you tried to encode a number.

## Uno reverso

Now that I was able to encode some JSON into Bson, I wanted to turn the tables
and decode Bson back into JSON.

And since I'm working with bits and not full bytes, I had to develop a way to
read bit by bit. That's how I ended up with Bit Reader.

Bit Reader was a lot like Bit Builder, in that I had to keep track of what byte
I'm looking at and what bit in that byte. The main difference is that I already
have a full list of bytes.

So honestly, it went really smoothly. The one thing I had to work out was how to
avoid getting extra bits from a byte when I only want a few bits. To solve this,
I had to create a mask that I could `&` with the byte I'm reading from:

```
// The current byte, but I want 3 bits starting at 5
0b0001_0010

// Shift my bits so the 3 bits are all the way to the right
0b0001_0010 >> 1
0b0000_1001

// Now if I just returned this, I would have `1001 instead of `001`
// So create a mask of `111` to get rid of the extra 1
(2^3 - 1) = 0b0000_0111
0b0000_1001 & 0b0000_0111
0b0000_0001

// And now the bits are properly extracted
```

Another fun taks!

So on to the decoder.

It was kind of easy.

I used the Bit Reader to get bits, interpret them, and then read values into
maps, array, strings, etc.. Doing the encoding earlier helped me to have a good
idea on how to tackle it.

I did run into a small bug that highlighted a design issue in my spec: In
reading bits for an Object field, I went straight to reading the string and
skipped the 3 bits that declared it a string.

It was obvious at that point: JSON object fields are always strings, which means
I don't have to encode the Value type bits for a string when encoding object
field names. So I put it on my mental list for later.

## Test data

It was time to stop hardcoding test json in my code and make some commands to
encode and decode json from a file.

To keep things simple for me, meaning not have to deal with reading from and
writing to files, I decided to take JSON from Standard In and write Bson to
Standard Out.

```
cat stuff.json | bson encode > stuff.bson

cat stuff.bson | bson decode > stuff.json
```

This was fairly simple and you can look at `main.go` and probably understand it
well enough if you don't know Golang.

Then I put together a bunch of JSON files into a `jsons` directory. Some I
crafted by hand, others with [JSON Generator](https://json-generator.com/). But
when ready, I started running them forwards and backwards through Bson:

```
$ cat jsons/object-simple.json 
{"active":false, "foo":"bar", "num":10, "thing":null}
$ cat jsons/object-simple.json | bson encode > a.bson
$ cat a.bson | bson encode > a.json
$ cat a.json
{"active":false,"foo":"bar","num":10,"thing":null}
```

Success! And I could see that the Bson file was smaller than the JSON file. So
Double Success!!

But I got tired of manually determining how many bytes were saved when encoding
some JSON, so I made another command:

```
$ cat jsons/object-simple.json | bson check
Json size: 54
Bson size: 29
diff: 25
```

Life was good. And so I decided to try out some JSON files with large data:

```
$ cat jsons/array-large.json | bson encode | bson decode
Invalid value token: 111
```

What? Why does that file fail when all the others worked?

Well remember that thought I had one night while walking my baby at 10pm? I
forgot to write any code to handle strings longer than 31 bytes when encoding or
decoding.

So it did come back later to bite me.

I went over the general plan already, so I won't discuss it again here. If you
want to, you can see my initial notes that I put in my
[notes on Bit Builder](/notes/bit-builder.md), for whatever reason.

Either way, I tackled it bit by bit:

* Encode long string
* Decode long string
* Run JSON with long strings through encode and decode
* Encode long array
* Decode long array
* Run JSON with a long array through encode and decode
* Encode long object
* Decode long object
* Run JSON with a long object through encode and decode

And with that, I could encode large JSON values and properly decode the Bson
back into JSON.

For fun, here's the Bson check for the file that previously failed:

```
$ cat jsons/array-large.json | bson check
Json size: 7751
Bson size: 5706
diff: 2045
```

That's a 26% reduction! Which probably suprised me more than it did you.

## Cow thoughts

While doing cow chores at night, I sometimes have thoughts that I need to
remember for later:

"What if I collected strings, shoved them at the back of the Bson file, and then
referenced those strings by ID instead of duplicating string bytes?"

Could give a lot of savings for arrays of homogeneous objects.

```
[bson version bits]
[bson encoded JSON bytes]
...
... [String member of array, but instead of storing "applesauce", just 5]
...
[post bson encoded JSON bytes]
[encoded array of strings]
...[5th element: "applesauce"]...
```

Either way, write it down, consider it later.

## Is this your number?

If you open the file [/jsons/array-large.json](/jsons/array-large.json), you'll
notice something missing: numbers.

I can't consider Bson complete until I have a solution for numbers, so it's time
to bite the bullet and learn more about things like signed/unsigned integers of
varying bit sizes, and why 0.1 + 0.2 = 0.30000000000000004

But before learning what a mantissa is, I looked at the easier Integers first.

No decimals, and for the most part (ignoring negatives), it's just plain binary.
And after learning about 2's compliment, even the negatives made sense.

Let's get Integers in the Spec.

I want Integers to be encoded in their smallest amount of bytes possible. If a
number fits in a signed 8 bit integer, shove it in a byte. If it fits into a
32 bit integer, 3 bytes.

Since there are only 4 integer sizes that I can tell (8, 16, 32, 64), I ditched
the 5 bit length encoding and used a 2 bit flag instead:

* `00` - 8 bit
* `01` - 16 bit
* `10` - 32 bit
* `11` - 64 bit

If we end up in a world with 128 bit integers, I'll have to make Bson V2 and
bump up the bits to 3.

```
001 [size flag] [Integer bytes]
```

Before encoding though, I needed to convert a signed integer into bytes. And
while looking through the Go standard library, I realized this will require some
decisions around endianness.

So I decided Big and moved on.

Wrote the encoding with some effort. The decoding went smoother, and I followed
in Go's footsteps by returning a 64 integer when decoding. Added new JSON files
and revelled in the success!

```
$ cat jsons/array-objects-with-numbers.json | bson check
Json size: 9339
Bson size: 6566
diff: 2773
```

## We'll all float on, alright!

Floats are interesting. Pretty incredible really. Best I can understand from
what I read, a float basically stores two parts of an equation, and a signed bit.

Way beyond me, so I'm not going to try and do anything more than store the bytes
that Go unmarshals out of the JSON string.

But how do I make sure that I'm dealing with a float and not an integer when all
that Go gives me is a 64 bit float?

And on top of that, apparently a 64 bit integer can store a much larger integer
than a 64 bit float can. So how do I get Go to return me integers instead of
floats when appropriate?

That's what lead me to finding `json.RawMessage`.

When you blindly unmarshal JSON in Go, you get the following types:

| JSON Type | Go Type |
| -- | -- |
| boolean | bool |
| numbers | float64 |
| string | string |
| array | []any |
| object | map[string]any |
| null | nil |

But if you tell Go to unmarshal a JSON object into `map[string]json.RawMessage`,
you get something different:

```go
jsonBlob := []byte{`{"a":1, "b": "bb", "c": {"foo":   "bar"}, "d": [  1],"e":1.23}`}
var d map[string]json.RawMessage
json.Unmarshal(jsonBlob, &d)
for key, v := range d {
    fmt.Printf("[%s]: %s\n", key, v)
}
```

The output of that code is this:
```
[a]: 1
[b]: "bb"
[c]: {"foo":   "bar"}
[d]: [  1]
[e]: 1.23
```

Which means I can look at the start of each value to determine what type it is.
As I found on a thread somewhere online:

* `{` means an object
* `[` means array
* `"` means a string
* `t` means true
* `f` means false
* `n` means null
* Everything else is a number

Numbers in JSON are interesting:

* Integers: 10, -10, etc.
* Floats: 0.1, 1.1, -1.1
* Floats?: 1e5, -1E-3

And with that in mind, if I don't find a period `.` or an `e` or `E`, then it's
an integer! So I put together some functions to take the JSON unmarshalling one
value at a time and was able to put Integer values into Integers and every other
number into 64 bit floats.

Finally, I generated a JSON file with floats in it, and the encoding and
decoding worked!

```
$ cat jsons/array-objects-with-floats.json | bson check
Json size: 6653
Bson size: 4564
diff: 2089
```

## Everything led up to this moment

I had a goal: come up with a binary encoding of JSON and see if I can make the
results fewer bytes than the source JSON.

Not only did I succeed, but the savings were between 25% and 55% on the test
files I used!

And even better, I learned how to deal with binary, came up with a specification
and wrote code that follows the specs.

Now it was time to reward myself and move on to the next leg of my journey. I
found that the PDF standard can be downloaded freely from the PDF Association!
I was now ready for whatever binary encoding Adobe chose for PDF files.

Let's take a quick look at some examples from the PDF 1.7 standard:

```
Integer objects
123 43445 +17 -98 0

Real objects
34.5 -3.62 +123.6 4. -.002 0.0

Strings
(This is a string)

Arrays
[549 3.14 false (Ralph) /SomeName]

Dictionaries
<< /Type /Example
    /Subtype /DictionaryExample
    /Version 0.01
    /IntegerItem 12
    /StringItem (a string)
    /Subdictionary <<   /Item1 0.4
                        /Item2 true
                        /LastItem (not!)
                        /VeryLastItem (OK)
                   >>
>>
```

... Well, I guess PDF is a plain text standard...

## When the cows come home

I finished writing this whole thing and decided to remove the 3 bit value type
in front of object keys and see what happens.

Here's the check on the largest JSON test file before the change:

```
$ cat jsons/array-objects-with-floats.json | bson check
Json size: 6653
Bson size: 4564
diff: 2089
```

And here's the check after the update:
```
$ cat jsons/array-objects-with-floats.json | bson check
Json size: 6653
Bson size: 4512
diff: 2141
```

So the change saved 52 bytes for that file, an extra 0.8% gained.

Small as it is, it's still an improvement! And either way, it feels good to not
include bits that aren't needed.

The spec is updated now and so are the tests.
