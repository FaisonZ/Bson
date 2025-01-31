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

### Notes

For a thoughtful write-up, [see "Post Thoughts"](notes/post-thoughts.md)

For other unexplained madness that I wrote and modified while developing this,
see the other files in the notes directory.

## Commands

### bson encode

```
cat file.json | bson encode > file.bson
```

encodes **file.json** into Bson and saves it into **file.bson**

### bson decode

```
cat file.bson | bson decode > file.json
```

decodes **file.bson** into JSON and saves it into **file.json**

### bson check

```
cat file.json | bson check
```

displays the size (in bytes) of file.json, the bson encoding and displays the
difference in size

## Compromises

I originally wanted to make sure all values that were encoded into Bson, would
be decoded and wrote back into JSON exactly how it was before encoding. Due to
how floats work, I decided to allow floats to come out however they come out.

If you have the value `1.23`, encoding then decoding and writing out JSON will
still get you `1.23`

However, if you have `1.0`, you'll end up with `1`

And because a floating point number is stored in binary as parts of an equation
instead of a flat number, the value you get after encoding and decoding *could*
end up with a rounding error.

Also, if it's good enough for Golang's json.Unmarshal(), then I'm not going to
worry too much on it.
