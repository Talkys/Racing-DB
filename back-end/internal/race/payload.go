package race

import (
	"context"
	"log/slog"
	"regexp"

	"backend/utils/jsonshape"
	"backend/utils/routes"
)

// The payload contracts of spec 3.3, one per category that has one. Other categories accept any object.
var payloadContracts = map[string]func(jsonshape.Object){
	"f1":     checkF1,
	"motogp": checkMotoGP,
	"wec":    checkWEC,
}

// lapTime is m:ss.sss, e.g. 1:27.511.
var lapTime = regexp.MustCompile(`^[0-9]{1,2}:[0-5][0-9]\.[0-9]{3}$`)

func checkF1(p jsonshape.Object) {
	lap := p.Object("fastest_lap")
	lap.UUID("driver_id")
	lap.Int("lap_number", 1, 500)
	lap.Match("lap_time", lapTime, "a lap time like 1:27.511")

	p.Objects("tyre_strategies", func(s jsonshape.Object) {
		s.UUID("driver_id")
		s.Objects("stops", func(stop jsonshape.Object) {
			stop.Int("lap", 1, 500)
			stop.Enum("compound", "SOFT", "MEDIUM", "HARD", "INTERMEDIATE", "WET")
		})
	})
	p.Int("safety_car_deployments", 0, 50)
	p.Number("track_temperature_celsius", -30, 80)
}

func checkMotoGP(p jsonshape.Object) {
	p.UUID("rider_of_the_day_id")
	tyres := p.Object("tyre_spec")
	tyres.Enum("front", "SOFT", "MEDIUM", "HARD")
	tyres.Enum("rear", "SOFT", "MEDIUM", "HARD")
	p.UUID("holeshot_winner_id")
	p.Enum("weather_conditions", "DRY", "WET")
}

func checkWEC(p jsonshape.Object) {
	p.Int("total_laps_completed", 0, 1000)
	winners := p.Object("class_winners")
	winners.UUID("hypercar")
	winners.UUID("lmgt3")
	p.Int("full_course_yellow_count", 0, 100)
	p.Int("red_flag_duration_minutes", 0, 1440)
}

// payloadProblems checks a payload against its category's contract; nil means it matches.
func payloadProblems(categoryID string, payload []byte) []string {
	c, root := jsonshape.Parse(payload)
	if contract, ok := payloadContracts[categoryID]; ok {
		contract(root)
	}
	return c.Problems()
}

// checkPayload is the Resource's Check: it flags a payload that breaks its contract and keeps it as stored,
// so one bad ETL row neither fails the request nor reaches clients as if it were valid.
func checkPayload(ctx context.Context, x *Race) error {
	x.PayloadErrors = payloadProblems(x.CategoryID, x.Payload)
	x.PayloadValid = len(x.PayloadErrors) == 0
	if !x.PayloadValid {
		slog.WarnContext(ctx, "race payload breaks its contract",
			"race_id", x.ID, "category_id", x.CategoryID, "problems", x.PayloadErrors, "request_id", routes.RequestID(ctx))
	}
	return nil
}
