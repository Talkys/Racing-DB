package car

import (
	"context"
	"log/slog"

	"backend/utils/jsonshape"
	"backend/utils/routes"
)

// partsProblems checks parts against the contract of spec 3.4, the same for every category; nil means it matches.
func partsProblems(parts []byte) []string {
	c, p := jsonshape.Parse(parts)

	gearbox := p.Object("gearbox")
	gearbox.String("type")
	gearbox.Int("speeds", 1, 10)
	for _, key := range []string{"brakes", "suspension"} {
		pair := p.Object(key)
		pair.String("front")
		pair.String("rear")
	}
	p.String("tyre_supplier")
	p.StringMap("other")

	return c.Problems()
}

// checkParts is the Resource's Check: it flags parts that break the contract and keeps them as stored,
// so one bad ETL row neither fails the request nor reaches clients as if it were valid.
func checkParts(ctx context.Context, x *Car) error {
	x.PartsErrors = partsProblems(x.Parts)
	x.PartsValid = len(x.PartsErrors) == 0
	if !x.PartsValid {
		slog.WarnContext(ctx, "car parts break their contract",
			"car_id", x.ID, "problems", x.PartsErrors, "request_id", routes.RequestID(ctx))
	}
	return nil
}
