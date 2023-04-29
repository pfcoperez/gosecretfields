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

````go
type Character struct {
	Name   Secret[string]
	Age    Secret[int]
	Friend *Character
}

type Book struct {
	Tittle     string
	Author     string
	Characters []Character
}

var sampleCharacterHiro = Character{
	Name: AsSecret("Hiro Protagonist"),
	Age:  AsSecret(30),
}

var sampleCharacterYT = Character{
	Name: AsSecret("YT"),
	Age:  AsSecret(15),
}

var snowCrash = Book{
	Tittle:     "Snow Crash",
	Author:     "Neal Stephenson",
	Characters: []Character{sampleCharacterHiro, sampleCharacterYT},
}

// This prints the character with redacted fields
fmt.Println(sampleCharacterHiro)

// This prints Hiro's unredacted age, the developer needs to "open" the secret box
fmt.Println(sampleCharacterHiro.Age.SecretValue)

// This serializes and prints the book with all secrets redacted
marshalledAllSecrets, _ := json.Marshal(snowCrash)
fmt.Println(string(marshalledAllSecrets))

// But we can set some fields open for JSON serialization
snowCrash.Characters[1].Name.CleartextJSON = true
snowCrash.Characters[1].Age.CleartextJSON = true

marshalledSomeSecrets, _ := json.Marshal(snowCrash)
fmt.Println(string(marshalledSomeSecrets))
```
