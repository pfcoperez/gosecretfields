package main

import (
	"encoding/json"
	"fmt"

	. "github.com/pfcoperez/gosecretfields"
)

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

func main() {

	// This prints the character with redacted fields
	fmt.Println(sampleCharacterHiro)

	// This prints Hiro's unredacted age, the developer needs to "open" the secret box
	fmt.Println(sampleCharacterHiro.Age.SecretValue)

	// This serializes and prints the book with all secrets redacted
	marshalledAllSecrets, _ := json.Marshal(snowCrash)
	fmt.Println(string(marshalledAllSecrets))

	// But we can set some fields open for JSON serialization
	changedSettings := NewImmutableSettings(true)

	snowCrash.Characters[1].Name.Settings = changedSettings
	snowCrash.Characters[1].Age.Settings = changedSettings

	marshalledSomeSecrets, _ := json.Marshal(snowCrash)
	fmt.Println(string(marshalledSomeSecrets))
}
