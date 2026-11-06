package queries

import (
	"backend/utils/apperror"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var testCfg = EntitySearchConfig{
	BaseCondition: "visible",
	PageSize:      25,
	DefaultOrder:  "race_date DESC, id ASC",
	AllowedFilters: map[string]AllowedFilter{
		"id":     UUIDFilter("id"),
		"name":   TextFilter("name"),
		"season": IntFilter("season"),
		"active": BoolFilter("is_active"),
	},
	AllowedSorts: map[string]string{"name": "name", "id": "id", "season": "season"},
}

var testCols = []string{"id", "name"}

func build(t *testing.T, rawQuery string) (string, []any, error) {
	t.Helper()
	q, err := url.ParseQuery(rawQuery)
	if err != nil {
		t.Fatal(err)
	}
	sc, err := BuildCriteriaFromURL(q)
	if err != nil {
		return "", nil, err
	}
	return BuildSearchQuery(testCols, "races", testCfg, *sc)
}

func assertBadRequest(t *testing.T, err error) {
	t.Helper()
	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Status != http.StatusBadRequest {
		t.Fatalf("expected a 400 apperror, got %v", err)
	}
}

func TestBuildSearchQuery(t *testing.T) {
	a := uuid.MustParse("b11a43a8-4e36-47a8-8e68-301fc1cae6ef")
	b := uuid.MustParse("a22b54b9-5f47-48b9-9f79-412fd2dbf7f0")

	tests := []struct {
		name      string
		rawQuery  string
		wantQuery string
		wantArgs  []any
	}{
		{
			name:      "no params uses base condition, default order and first page",
			rawQuery:  "",
			wantQuery: "select id, name from races where (visible) order by race_date DESC, id ASC limit 26 offset 0",
			wantArgs:  []any{},
		},
		{
			name:      "typed int filter and page",
			rawQuery:  "filter=season,gte,2026&page=3",
			wantQuery: "select id, name from races where (visible) and season >= $1 order by race_date DESC, id ASC limit 26 offset 50",
			wantArgs:  []any{int64(2026)},
		},
		{
			name:      "ilike wraps and escapes the value",
			rawQuery:  "filter=name,ilike,50%25_off",
			wantQuery: "select id, name from races where (visible) and name ILIKE $1 order by race_date DESC, id ASC limit 26 offset 0",
			wantArgs:  []any{`%50\%\_off%`},
		},
		{
			name:      "filter value may contain commas",
			rawQuery:  "filter=name,eq,a,b",
			wantQuery: "select id, name from races where (visible) and name = $1 order by race_date DESC, id ASC limit 26 offset 0",
			wantArgs:  []any{"a,b"},
		},
		{
			name:      "in parses every id",
			rawQuery:  "filter=id,in," + a.String() + "|" + b.String(),
			wantQuery: "select id, name from races where (visible) and id = ANY($1) order by race_date DESC, id ASC limit 26 offset 0",
			wantArgs:  []any{[]uuid.UUID{a, b}},
		},
		{
			name:      "bool filter",
			rawQuery:  "filter=active,eq,true",
			wantQuery: "select id, name from races where (visible) and is_active = $1 order by race_date DESC, id ASC limit 26 offset 0",
			wantArgs:  []any{true},
		},
		{
			name:      "custom sort gets the id tiebreaker",
			rawQuery:  "sort=name,desc",
			wantQuery: "select id, name from races where (visible) order by name DESC, id ASC limit 26 offset 0",
			wantArgs:  []any{},
		},
		{
			name:      "repeated filters are AND-ed in order",
			rawQuery:  "filter=name,eq,Monaco&filter=season,gte,2020&filter=season,lte,2026",
			wantQuery: "select id, name from races where (visible) and name = $1 and season >= $2 and season <= $3 order by race_date DESC, id ASC limit 26 offset 0",
			wantArgs:  []any{"Monaco", int64(2020), int64(2026)},
		},
		{
			name:      "comma-separated sorts with optional directions",
			rawQuery:  "sort=season,DESC,name",
			wantQuery: "select id, name from races where (visible) order by season DESC, name ASC, id ASC limit 26 offset 0",
			wantArgs:  []any{},
		},
		{
			name:      "repeated sort params",
			rawQuery:  "sort=season&sort=name,desc",
			wantQuery: "select id, name from races where (visible) order by season ASC, name DESC, id ASC limit 26 offset 0",
			wantArgs:  []any{},
		},
		{
			name:      "an explicit id sort replaces the tiebreaker",
			rawQuery:  "sort=season,id,desc",
			wantQuery: "select id, name from races where (visible) order by season ASC, id DESC limit 26 offset 0",
			wantArgs:  []any{},
		},
		{
			name:      "trailing comma is ignored",
			rawQuery:  "sort=name,",
			wantQuery: "select id, name from races where (visible) order by name ASC, id ASC limit 26 offset 0",
			wantArgs:  []any{},
		},
		{
			name:      "sorting by id needs no tiebreaker",
			rawQuery:  "sort=id",
			wantQuery: "select id, name from races where (visible) order by id ASC limit 26 offset 0",
			wantArgs:  []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args, err := build(t, tt.rawQuery)
			if err != nil {
				t.Fatal(err)
			}
			if query != tt.wantQuery {
				t.Errorf("query\n got: %s\nwant: %s", query, tt.wantQuery)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("args\n got: %#v\nwant: %#v", args, tt.wantArgs)
			}
		})
	}
}

func TestBuildSearchQueryRejectsBadInput(t *testing.T) {
	for _, rawQuery := range []string{
		"page=0",
		"page=-1",
		"page=abc",
		"filter=name",
		"filter=nope,eq,1",
		"filter=season,ilike,20",
		"filter=season,eq,abc",
		"filter=active,eq,maybe",
		"filter=id,eq,not-a-uuid",
		"filter=id,in,b11a43a8-4e36-47a8-8e68-301fc1cae6ef|nope",
		"filter=name,eq,a&filter=season,eq,abc", // every filter is checked, not only the first
		"filter=name,eq,a&filter=nope,eq,1",
		"filter=name&filter=season,eq,1",
		strings.TrimPrefix(strings.Repeat("&filter=name,eq,a", maxFilters+1), "&"),
		"sort=round",
		"sort=name,sideways",
		"sort=desc",
		"sort=name,asc,desc",
		"sort=name,season,name",
		"sort=name&sort=name,desc",
	} {
		t.Run(rawQuery, func(t *testing.T) {
			_, _, err := build(t, rawQuery)
			assertBadRequest(t, err)
		})
	}
}

func TestBuildGetQuery(t *testing.T) {
	got := BuildGetQuery(testCols, "races", testCfg)
	want := "select id, name from races where id = $1 and (visible)"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}

	noBase := testCfg
	noBase.BaseCondition = ""
	if got := BuildGetQuery(testCols, "races", noBase); got != "select id, name from races where id = $1" {
		t.Errorf("got %s", got)
	}
}
