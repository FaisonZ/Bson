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

