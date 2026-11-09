package jsonshape

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// contract exercises every check once.
func contract(root Object) {
	inner := root.Object("inner")
	inner.String("name")
	inner.UUID("id")
	root.Enum("kind", "A", "B")
	root.Int("count", 0, 10)
	root.Number("ratio", -1, 1)
	root.Match("code", regexp.MustCompile(`^[A-Z]{3}$`), "a three-letter code")
	root.StringMap("tags")
	root.Objects("items", func(item Object) {
		item.Int("n", 1, 5)
	})
}

func problems(raw string) []string {
	c, root := Parse([]byte(raw))
	contract(root)
	return c.Problems()
}

const valid = `{
  "inner": {"name": "x", "id": "b11a43a8-4e36-47a8-8e68-301fc1cae6ef"},
  "kind": "A", "count": 10, "ratio": -0.5, "code": "VER",
  "tags": {"a": "1"}, "items": [{"n": 1}, {"n": 5}],
  "extra": "unknown keys are allowed"
}`

func TestValid(t *testing.T) {
	if got := problems(valid); got != nil {
		t.Fatalf("problems = %q, want none", got)
	}
	if got := problems(`{"inner": {"name": "x", "id": "b11a43a8-4e36-47a8-8e68-301fc1cae6ef"}, "kind": "B", "count": 0,
		"ratio": 1, "code": "HAM", "tags": {}, "items": []}`); got != nil {
		t.Fatalf("empty map and array: problems = %q, want none", got)
	}
}

func TestProblems(t *testing.T) {
	for _, tc := range []struct {
		name, raw string
		want      []string
	}{
		{"not JSON", `{`, []string{"(root): not valid JSON"}},
		{"root not an object", `[]`, []string{"(root): want an object, got an array"}},
		{"empty object", `{}`, []string{
			"inner: missing", "kind: missing", "count: missing", "ratio: missing", "code: missing", "tags: missing", "items: missing",
		}},
		{"null counts as missing", replace(`"kind": "A"`, `"kind": null`), []string{"kind: missing"}},
		{"bad object skips its fields", replace(`"inner": {"name": "x", "id": "b11a43a8-4e36-47a8-8e68-301fc1cae6ef"}`, `"inner": 3`),
			[]string{"inner: want an object, got a number"}},
		{"nested missing key", replace(`"name": "x"`, `"nam": "x"`), []string{"inner.name: missing"}},
		{"empty string", replace(`"name": "x"`, `"name": " "`), []string{"inner.name: empty"}},
		{"wrong type", replace(`"name": "x"`, `"name": 1`), []string{"inner.name: want a string, got a number"}},
		{"bad uuid", replace(`"b11a43a8-4e36-47a8-8e68-301fc1cae6ef"`, `"b11a43a84e3647a88e68301fc1cae6ef"`),
			[]string{`inner.id: "b11a43a84e3647a88e68301fc1cae6ef" is not a UUID`}},
		{"enum", replace(`"kind": "A"`, `"kind": "a"`), []string{`kind: "a" is not one of A, B`}},
		{"int range", replace(`"count": 10`, `"count": 11`), []string{"count: 11 is outside 0–10"}},
		{"int fraction", replace(`"count": 10`, `"count": 2.5`), []string{"count: want an integer, got 2.5"}},
		{"int as string", replace(`"count": 10`, `"count": "10"`), []string{"count: want a number, got a string"}},
		{"number range", replace(`"ratio": -0.5`, `"ratio": -1.5`), []string{"ratio: -1.5 is outside -1–1"}},
		{"number overflow", replace(`"ratio": -0.5`, `"ratio": 1e400`), []string{"ratio: 1e400 is not a finite number"}},
		{"pattern", replace(`"code": "VER"`, `"code": "ver"`), []string{`code: "ver" is not a three-letter code`}},
		{"string map", replace(`"tags": {"a": "1"}`, `"tags": {"b": true, "a": 1}`),
			[]string{"tags.a: want a string, got a number", "tags.b: want a string, got a boolean"}},
		{"array type", replace(`"items": [{"n": 1}, {"n": 5}]`, `"items": {}`), []string{"items: want an array, got an object"}},
		{"array elements", replace(`"items": [{"n": 1}, {"n": 5}]`, `"items": [{"n": 0}, "x", {}]`),
			[]string{"items[0].n: 0 is outside 1–5", "items[1]: want an object, got a string", "items[2].n: missing"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := problems(tc.raw); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("problems =\n  %q\nwant\n  %q", got, tc.want)
			}
		})
	}
}

// replace swaps one fragment of the valid document, failing loudly if the fragment is not there.
func replace(old, new string) string {
	out := strings.Replace(valid, old, new, 1)
	if out == valid {
		panic("fragment not found: " + old)
	}
	return out
}
