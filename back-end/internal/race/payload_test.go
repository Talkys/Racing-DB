package race

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const driver = "b11a43a8-4e36-47a8-8e68-301fc1cae6ef"

// Valid payloads, one per contracted category (spec 3.3).
var validPayloads = map[string]string{
	"f1": `{
		"fastest_lap": {"driver_id": "` + driver + `", "lap_number": 44, "lap_time": "1:27.511"},
		"tyre_strategies": [{"driver_id": "` + driver + `", "stops": [{"lap": 18, "compound": "MEDIUM"}]},
		                    {"driver_id": "` + driver + `", "stops": []}],
		"safety_car_deployments": 2,
		"track_temperature_celsius": 38.5
	}`,
	"motogp": `{
		"rider_of_the_day_id": "` + driver + `", "tyre_spec": {"front": "HARD", "rear": "SOFT"},
		"holeshot_winner_id": "` + driver + `", "weather_conditions": "DRY"
	}`,
	"wec": `{
		"total_laps_completed": 368, "class_winners": {"hypercar": "` + driver + `", "lmgt3": "` + driver + `"},
		"full_course_yellow_count": 4, "red_flag_duration_minutes": 0
	}`,
}

func TestPayloadProblems(t *testing.T) {
	for _, tc := range []struct {
		name, category, old, new string // old → new edits the category's valid payload; old "" keeps it
		want                     []string
	}{
		{name: "f1 valid", category: "f1"},
		{name: "motogp valid", category: "motogp"},
		{name: "wec valid", category: "wec"},
		{name: "category without contract takes any object", category: "indycar"},

		{name: "f1 empty", category: "f1", old: validPayloads["f1"], new: `{}`, want: []string{
			"fastest_lap: missing", "tyre_strategies: missing", "safety_car_deployments: missing", "track_temperature_celsius: missing",
		}},
		{name: "f1 missing nested key", category: "f1", old: `"lap_number": 44, `, want: []string{"fastest_lap.lap_number: missing"}},
		{name: "f1 bad lap time", category: "f1", old: `"1:27.511"`, new: `"87.5"`,
			want: []string{`fastest_lap.lap_time: "87.5" is not a lap time like 1:27.511`}},
		{name: "f1 bad compound", category: "f1", old: `"MEDIUM"`, new: `"SUPERSOFT"`,
			want: []string{`tyre_strategies[0].stops[0].compound: "SUPERSOFT" is not one of SOFT, MEDIUM, HARD, INTERMEDIATE, WET`}},
		{name: "f1 wrong type", category: "f1", old: `"safety_car_deployments": 2`, new: `"safety_car_deployments": "2"`,
			want: []string{"safety_car_deployments: want a number, got a string"}},
		{name: "f1 out of range", category: "f1", old: `38.5`, new: `380`,
			want: []string{"track_temperature_celsius: 380 is outside -30–80"}},
		{name: "f1 driver id not a uuid", category: "f1", old: `"` + driver + `", "lap_number"`, new: `"VER", "lap_number"`,
			want: []string{`fastest_lap.driver_id: "VER" is not a UUID`}},

		{name: "motogp missing key", category: "motogp", old: `"weather_conditions": "DRY"`, new: `"weather": "DRY"`,
			want: []string{"weather_conditions: missing"}},
		{name: "motogp bad enum", category: "motogp", old: `"rear": "SOFT"`, new: `"rear": "WET"`,
			want: []string{`tyre_spec.rear: "WET" is not one of SOFT, MEDIUM, HARD`}},
		{name: "motogp wrong type", category: "motogp", old: `{"front": "HARD", "rear": "SOFT"}`, new: `["HARD", "SOFT"]`,
			want: []string{"tyre_spec: want an object, got an array"}},

		{name: "wec fraction", category: "wec", old: `368`, new: `368.5`, want: []string{"total_laps_completed: want an integer, got 368.5"}},
		{name: "wec out of range", category: "wec", old: `"full_course_yellow_count": 4`, new: `"full_course_yellow_count": -1`,
			want: []string{"full_course_yellow_count: -1 is outside 0–100"}},
		{name: "wec null winner", category: "wec", old: `"lmgt3": "` + driver + `"`, new: `"lmgt3": null`,
			want: []string{"class_winners.lmgt3: missing"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, ok := validPayloads[tc.category]
			if !ok {
				raw = `{"anything": [1, "two"]}`
			}
			if tc.old != "" {
				edited := strings.Replace(raw, tc.old, tc.new, 1)
				if edited == raw {
					t.Fatalf("fragment %q not in the %s payload", tc.old, tc.category)
				}
				raw = edited
			}
			if got := payloadProblems(tc.category, []byte(raw)); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("problems =\n  %q\nwant\n  %q", got, tc.want)
			}
		})
	}
}

func TestCheckPayloadFlagsAndKeepsThePayload(t *testing.T) {
	x := Race{ID: uuid.New(), CategoryID: "f1", Payload: json.RawMessage(`{}`)}
	if err := checkPayload(t.Context(), &x); err != nil {
		t.Fatal(err)
	}
	if x.PayloadValid || len(x.PayloadErrors) == 0 || string(x.Payload) != `{}` {
		t.Errorf("invalid payload: valid=%v errors=%q payload=%s", x.PayloadValid, x.PayloadErrors, x.Payload)
	}

	x = Race{ID: uuid.New(), CategoryID: "motogp", Payload: json.RawMessage(validPayloads["motogp"])}
	if err := checkPayload(t.Context(), &x); err != nil {
		t.Fatal(err)
	}
	if !x.PayloadValid || x.PayloadErrors != nil {
		t.Errorf("valid payload: valid=%v errors=%q", x.PayloadValid, x.PayloadErrors)
	}
	if body, _ := json.Marshal(x); strings.Contains(string(body), "payload_errors") {
		t.Errorf("payload_errors should be omitted when valid: %s", body)
	}
}
