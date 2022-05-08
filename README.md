# Adobo

Adobo generates types from a jsonschema (draft 2019-09) in different languages (currently only
Go).

Adobo is written using [arr.ai](https://arr.ai)

**WARNING**: adobo is very much experimental and it also use a very experimental
language. Use at your own risk.

## Usage

Simplest usage is to just run it by pointing to your jsonschema file

```sh
adobo path/to/your/schema.json
```

The previous command will output Go code to your stdout.

For more information, you can run adobo with the `--help` flag.

## Future Updates :crossed_fingers:

### Other languages!

adobo wasn't meant to just work with Golang, it should also output to other
languages. To do this, it will require A LOT of work.

### Other jsonschema

There's a lot of jsonschema versions and it should support other versions. Not
to mention, there's more features to support in the 2019-09 draft like enums,
different patternProperties, and many more. The suffering never ends.

### More [arr.ai](https://arr.ai)

[arr.ai](https://arr.ai) is the best language ever and this is a totally
unbiased comment. Because it's the best language ever, the CLI should be written
in [arr.ai](https://arr.ai).
