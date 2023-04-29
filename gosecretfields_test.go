package gosecretfields

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestAsSecret(t *testing.T) {
	type args struct {
		value         string
		redactedValue []string
	}
	tests := []struct {
		name string
		args args
		want Secret[string]
	}{
		{
			name: "With default redacted value",
			args: args{
				value:         "Hiro Protagonist",
				redactedValue: []string{},
			},
			want: Secret[string]{SecretValue: "Hiro Protagonist", redactedValue: ""},
		},
		{
			name: "With explicit redacted value",
			args: args{
				value:         "Hiro Protagonist",
				redactedValue: []string{"REDACTED"},
			},
			want: Secret[string]{SecretValue: "Hiro Protagonist", redactedValue: "REDACTED"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsSecret(tt.args.value, tt.args.redactedValue...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

type character struct {
	Name   Secret[string]
	Age    Secret[int]
	Friend *character
}

type testCase struct {
	name       string
	value      character
	redactJson bool
}

var sampleCharacterHiro = character{
	Name: AsSecret("Hiro Protagonist"),
	Age:  AsSecret(30),
}

var sampleCharacterYT = character{
	Name: AsSecret("YT"),
	Age:  AsSecret(15),
}

func containsSecrets(text string) bool {
	return strings.Contains(text, sampleCharacterHiro.Name.SecretValue) ||
		strings.Contains(text, fmt.Sprint(sampleCharacterHiro.Age.SecretValue)) ||
		strings.Contains(text, sampleCharacterYT.Name.SecretValue) ||
		strings.Contains(text, fmt.Sprint(sampleCharacterYT.Age.SecretValue))
}

func TestMarshallingAndStringer(t *testing.T) {

	hiroWithFriend := sampleCharacterHiro
	hiroWithFriend.Friend = &sampleCharacterYT

	tests := []testCase{
		{
			name:       "Redacting JSON",
			value:      sampleCharacterHiro,
			redactJson: true,
		},
		{
			name:       "NOT Redacting JSON",
			value:      sampleCharacterHiro,
			redactJson: false,
		},
		{
			name:       "Redacting JSON with references",
			value:      hiroWithFriend,
			redactJson: true,
		},
		{
			name:       "NOT Redacting JSON with references",
			value:      hiroWithFriend,
			redactJson: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if asString := fmt.Sprint(testCase.value); containsSecrets(asString) {
				t.Errorf("Stringer interface implementation must prevent secrets from leaking:\n%s", asString)
			}

			testCase.value.Name.CleartextJSON = !testCase.redactJson
			testCase.value.Age.CleartextJSON = !testCase.redactJson

			if maybeFriend := testCase.value.Friend; maybeFriend != nil {
				maybeFriend.Name.CleartextJSON = !testCase.redactJson
				maybeFriend.Age.CleartextJSON = !testCase.redactJson
			}

			var unmarshalled character

			marshalled, _ := json.Marshal(testCase.value)
			marshalledString := string(marshalled)

			if testCase.redactJson && containsSecrets(marshalledString) {
				t.Error("Marshalled JSON can't contain secret values when JSON redaction is enabled")
			}

			unmarshallingError := json.Unmarshal(marshalled, &unmarshalled)

			if unmarshallingError != nil {
				t.Errorf("Unmarshalling should not fail. Error: %s", unmarshallingError)
			}

			unmarshalled.Name.CleartextJSON = !testCase.redactJson
			unmarshalled.Age.CleartextJSON = !testCase.redactJson

			if maybeUnmarshalledFriend := unmarshalled.Friend; maybeUnmarshalledFriend != nil {
				maybeUnmarshalledFriend.Name.CleartextJSON = !testCase.redactJson
				maybeUnmarshalledFriend.Age.CleartextJSON = !testCase.redactJson
			}

			if !testCase.redactJson && !reflect.DeepEqual(unmarshalled, testCase.value) {
				t.Errorf("Unmarshalled value doesn't match the value that generated the JSON\nUnmarshalled value:\n%v\nValue:\n%v", unmarshalled, testCase.value)
			}

			if testCase.redactJson && reflect.DeepEqual(unmarshalled, testCase.value) {
				t.Errorf("It should not be possible to unmarshall the same value from secrets-hidden JSONs\n%s", marshalledString)
			}

		})
	}
}
