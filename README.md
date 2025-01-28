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

