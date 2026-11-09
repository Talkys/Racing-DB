package car

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// validParts follows spec 3.4.
const validParts = `{
	"gearbox": {"type": "SEQUENTIAL", "speeds": 8},
	"brakes": {"front": "Carbon-carbon discs", "rear": "Carbon-carbon discs"},
	"suspension": {"front": "Pushrod", "rear": "Pullrod"},
	"tyre_supplier": "Pirelli",
	"other": {"fuel_supplier": "Petronas", "hybrid_system": "MGU-K"}
}`

func TestPartsProblems(t *testing.T) {
	for _, tc := range []struct {
		name, old, new string // old → new edits validParts; old "" keeps it
		want           []string
	}{
		{name: "valid"},
		{name: "empty other", old: `{"fuel_supplier": "Petronas", "hybrid_system": "MGU-K"}`, new: `{}`},
		{name: "empty", old: validParts, new: `{}`, want: []string{
			"gearbox: missing", "brakes: missing", "suspension: missing", "tyre_supplier: missing", "other: missing",
		}},
		{name: "missing nested key", old: `, "speeds": 8`, want: []string{"gearbox.speeds: missing"}},
		{name: "wrong type", old: `"tyre_supplier": "Pirelli"`, new: `"tyre_supplier": ["Pirelli"]`,
			want: []string{"tyre_supplier: want a string, got an array"}},
		{name: "out of range", old: `"speeds": 8`, new: `"speeds": 0`, want: []string{"gearbox.speeds: 0 is outside 1–10"}},
		{name: "other must hold strings", old: `"hybrid_system": "MGU-K"`, new: `"hybrid_system": {"kw": 120}`,
			want: []string{"other.hybrid_system: want a string, got an object"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw := validParts
			if tc.old != "" {
				raw = strings.Replace(raw, tc.old, tc.new, 1)
				if raw == validParts {
					t.Fatalf("fragment %q not in validParts", tc.old)
				}
			}
			if got := partsProblems([]byte(raw)); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("problems =\n  %q\nwant\n  %q", got, tc.want)
			}
		})
	}
}

func TestCheckPartsFlagsAndKeepsTheParts(t *testing.T) {
	x := Car{ID: uuid.New(), Parts: json.RawMessage(`{"gearbox": {}}`)}
	if err := checkParts(t.Context(), &x); err != nil {
		t.Fatal(err)
	}
	if x.PartsValid || len(x.PartsErrors) == 0 || string(x.Parts) != `{"gearbox": {}}` {
		t.Errorf("invalid parts: valid=%v errors=%q parts=%s", x.PartsValid, x.PartsErrors, x.Parts)
	}

	x = Car{ID: uuid.New(), Parts: json.RawMessage(validParts)}
	if err := checkParts(t.Context(), &x); err != nil {
		t.Fatal(err)
	}
	if !x.PartsValid || x.PartsErrors != nil {
		t.Errorf("valid parts: valid=%v errors=%q", x.PartsValid, x.PartsErrors)
	}
}
