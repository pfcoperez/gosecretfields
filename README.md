# gosecretfields
Simple library to represent secret or sensitive values in Go variables or fields.
It aims to provide default behaviours that prevent the contents of these fields information leaking, for example through logs.

## What this library does

It provides a wrapper which is not disimilar to a collection of just one element. Any string serialization of this collection results in a redacted value and JSON serializations may:

- Generate a JSON representation of the unique value contained in _the collection_, if its flag `ClearText` is set to `true` 

OR

- Generate a JSON representation of the redacted counterpart of the value it wraps.
_You can open the secret box, you need to open the secret box_.

This provides a safe framework to work with secrets: Printing, logging and serializing secret tagged values default to not showing them while they can **explicitly** accessed.
This shift from default disclosure to default hiding makes it easy for visual or automatic inspections to detect when a log or a JSON HTTP response might output information that needs to be handle with care while relieving the developer from the cognitive burden of having to worry about the code maybe leaking secrets around.

## What this library doesn't do

This small library:

- Doesn't cipher memory contents.
- Doesn't provide reverse encryption.

But the redacted value can be anything the developer wants: From default values for the type to encrypted versions of the secret value.

## Usage example

https://github.com/pfcoperez/gosecretfields/blob/f7bc50b57c4eb1d367946f2947e6eeb6436b7349/demo/app.go#L1-L56
