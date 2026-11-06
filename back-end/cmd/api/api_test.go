package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"testing"
)

const pageSize = 25

// Well-known seed rows (same ids as scripts/smoke.sh).
const (
	raceMonaco     = "c39a7593-36c5-4ceb-8d07-6bc82f808381"
	raceMugello    = "d40b8604-47d6-5dfc-9e18-7cd93f919492"
	driverVer      = "b11a43a8-4e36-47a8-8e68-301fc1cae6ef"
	driverBag      = "a22b54b9-5f47-48b9-9f79-412fd2dbf7f0"
	articleMonaco  = "75000000-0000-4000-8000-000000000001"
	articleF1      = "75000000-0000-4000-8000-000000000002"
	articleDraft   = "75000000-0000-4000-8000-000000000003"
	articleFuture  = "75000000-0000-4000-8000-000000000006" // dated 2099: must stay in the future
	unknownUUID    = "00000000-0000-4000-8000-000000000099"
	visibleArticle = "published_at IS NOT NULL AND published_at <= now()"
)

// resource describes one list/get endpoint. The expectations are written as SQL against the seeded
// database, independently of the repositories, so a broken repository query shows up as a mismatch.
type resource struct {
	path         string
	table        string
	where        string // rows the API serves; empty = every row
	defaultOrder string // ORDER BY the list uses without sort=
	ids          [2]string
	missing      string    // well-formed id that does not exist (404)
	badID        string    // malformed id (400)
	filter       [3]string // column, eq, value: must match ids[0]
	sorts        []string  // every sortable column
	paged        bool      // the seed holds more than one page
}

var resources = []resource{
	{
		path: "categories", table: "categories", defaultOrder: "name ASC, id ASC",
		ids: [2]string{"f1", "motogp"}, missing: "xyz", badID: "no!",
		filter: [3]string{"is_active", "eq", "true"}, sorts: []string{"id", "name", "is_active"},
	},
	{
		path: "circuits", table: "circuits", defaultOrder: "name ASC, id ASC",
		ids:     [2]string{"23fa3b45-df13-4f9e-a864-44aa4cf81190", "6fa85f64-5717-4562-b3fc-2c963f66afa6"},
		missing: unknownUUID, badID: "not-a-uuid",
		filter: [3]string{"country", "eq", "Monaco"}, sorts: []string{"name", "location", "country", "length_km"},
	},
	{
		path: "drivers", table: "drivers", defaultOrder: "last_name ASC, first_name ASC, id ASC",
		ids: [2]string{driverVer, driverBag}, missing: unknownUUID, badID: "not-a-uuid",
		filter: [3]string{"nationality", "eq", "Dutch"}, sorts: []string{"last_name", "first_name", "code", "nationality"}, paged: true,
	},
	{
		path: "races", table: "races", defaultOrder: "race_date DESC, id ASC",
		ids: [2]string{raceMonaco, raceMugello}, missing: unknownUUID, badID: "not-a-uuid",
		filter: [3]string{"category_id", "eq", "f1"}, sorts: []string{"race_date", "name", "category_id", "season", "round"}, paged: true,
	},
	{
		path: "teams", table: "teams", defaultOrder: "name ASC, id ASC",
		ids:     [2]string{"71000000-0000-4000-8000-000000000001", "71000000-0000-4000-8000-000000000002"},
		missing: unknownUUID, badID: "not-a-uuid",
		filter: [3]string{"country", "eq", "United Kingdom"}, sorts: []string{"name", "short_name", "country", "founded_year"},
	},
	{
		path: "engines", table: "engines", defaultOrder: "manufacturer ASC, name ASC, id ASC",
		ids:     [2]string{"72000000-0000-4000-8000-000000000001", "72000000-0000-4000-8000-000000000002"},
		missing: unknownUUID, badID: "not-a-uuid",
		filter: [3]string{"aspiration", "eq", "HYBRID_TURBO"}, sorts: []string{"manufacturer", "name", "configuration", "displacement_cc", "cylinders", "aspiration", "power_kw"},
	},
	{
		path: "chassis", table: "chassis", defaultOrder: "manufacturer ASC, name ASC, id ASC",
		ids:     [2]string{"73000000-0000-4000-8000-000000000001", "73000000-0000-4000-8000-000000000002"},
		missing: unknownUUID, badID: "not-a-uuid",
		filter: [3]string{"material", "eq", "Carbon fibre monocoque"}, sorts: []string{"manufacturer", "name", "material", "wheelbase_mm", "weight_kg"},
	},
	{
		path: "cars", table: "cars", defaultOrder: "season DESC, name ASC, id ASC",
		ids:     [2]string{"74000000-0000-4000-8000-000000000001", "74000000-0000-4000-8000-000000000002"},
		missing: unknownUUID, badID: "not-a-uuid",
		filter: [3]string{"season", "eq", "2026"}, sorts: []string{"name", "season", "category_id"}, paged: true,
	},
	{
		path: "articles", table: "articles", where: visibleArticle, defaultOrder: "published_at DESC, id ASC",
		ids: [2]string{articleMonaco, articleF1}, missing: articleDraft, badID: "not-a-uuid",
		filter: [3]string{"entity_type", "eq", "race"}, paged: true, sorts: []string{"published_at", "created_at", "entity_type", "kind"},
	},
	{
		path: "article-texts", table: "article_texts", where: "article_id IN (SELECT id FROM articles WHERE " + visibleArticle + ")",
		defaultOrder: "language ASC, id ASC",
		ids:          [2]string{"76000000-0000-4000-8000-000000000001", "76000000-0000-4000-8000-000000000002"},
		missing:      "76000000-0000-4000-8000-000000000004", badID: "not-a-uuid", // ...04 is the draft's text
		filter: [3]string{"language", "eq", "en"}, sorts: []string{"title", "language"},
	},
}

func TestResources(t *testing.T) {
	requireDB(t)
	for _, res := range resources {
		t.Run(res.path, func(t *testing.T) {
			t.Run("list pages through every row in default order", func(t *testing.T) {
				want := res.oracle(t, "", res.defaultOrder)
				if res.paged && len(want) <= pageSize {
					t.Fatalf("seed has %d rows, want more than one page (%d)", len(want), pageSize)
				}

				var got []string
				for page := 1; ; page++ {
					rows := getList(t, res.path, url.Values{"page": {fmt.Sprint(page)}})
					if len(rows) > pageSize+1 {
						t.Fatalf("page %d: %d rows, want at most %d", page, len(rows), pageSize+1)
					}
					if len(rows) <= pageSize {
						got = append(got, ids(rows)...)
						break
					}
					got = append(got, ids(rows[:pageSize])...)
					// The extra row only signals a next page, which must start with it.
					next := getList(t, res.path, url.Values{"page": {fmt.Sprint(page + 1)}})
					if len(next) == 0 || next[0]["id"] != rows[pageSize]["id"] {
						t.Fatalf("page %d look-ahead row %v does not open page %d", page, rows[pageSize]["id"], page+1)
					}
				}
				assertIDs(t, got, want)

				pastEnd := len(want)/pageSize + 2
				if rows := getList(t, res.path, url.Values{"page": {fmt.Sprint(pastEnd)}}); len(rows) != 0 {
					t.Errorf("page %d past the end: %d rows, want none", pastEnd, len(rows))
				}
			})

			t.Run("filter", func(t *testing.T) {
				col, op, val := res.filter[0], res.filter[1], res.filter[2]
				rows := getList(t, res.path, url.Values{"filter": {col + "," + op + "," + val}})
				assertIDs(t, ids(rows), res.oracle(t, fmt.Sprintf("%s::text = '%s'", col, val), res.defaultOrder, pageSize+1))
				if !slices.Contains(ids(rows), res.ids[0]) {
					t.Errorf("filter result lacks %s", res.ids[0])
				}
			})

			t.Run("sort", func(t *testing.T) {
				for _, col := range res.sorts {
					for _, dir := range []string{"asc", "desc"} {
						rows := getList(t, res.path, url.Values{"sort": {col + "," + dir}})
						assertIDs(t, ids(rows), res.oracle(t, "", col+" "+dir+", id ASC", pageSize+1))
					}
				}
			})

			t.Run("get matches the list row", func(t *testing.T) {
				var got map[string]any
				getJSON(t, "/"+res.path+"/"+res.ids[0], nil, http.StatusOK, &got)
				listed := getList(t, res.path, url.Values{"filter": {"id,eq," + res.ids[0]}})
				if len(listed) != 1 || !reflect.DeepEqual(got, listed[0]) {
					t.Errorf("get = %v\nlist = %v", got, listed)
				}
			})

			t.Run("getMany with id,in", func(t *testing.T) {
				rows := getList(t, res.path, url.Values{"filter": {"id,in," + res.ids[0] + "|" + res.ids[1] + "|" + res.missing}})
				got := ids(rows)
				slices.Sort(got)
				want := []string{res.ids[0], res.ids[1]}
				slices.Sort(want)
				assertIDs(t, got, want)
			})

			t.Run("errors", func(t *testing.T) {
				assertError(t, "/"+res.path+"/"+res.missing, nil, http.StatusNotFound)
				assertError(t, "/"+res.path+"/"+url.PathEscape(res.badID), nil, http.StatusBadRequest)
				for _, q := range []url.Values{
					{"filter": {"nope,eq,1"}},
					{"filter": {res.filter[0] + ",bogus,1"}},
					{"filter": {res.filter[0]}},
					{"filter": {"id,in,not-an-id|" + res.ids[0]}},
					{"filter": {"id,eq," + res.ids[0], "nope,eq,1"}},
					{"sort": {"nope,asc"}},
					{"sort": {res.sorts[0] + ",sideways"}},
					{"sort": {res.sorts[0], res.sorts[0] + ",desc"}},
					{"page": {"0"}},
					{"page": {"x"}},
				} {
					if res.path == "categories" && strings.HasPrefix(q.Get("filter"), "id,in,") {
						continue // category ids are free-form slugs
					}
					assertError(t, "/"+res.path, q, http.StatusBadRequest)
				}
			})
		})
	}
}

// One case per filter kind and operator, beyond the eq filters in TestResources.
func TestFilterOperators(t *testing.T) {
	requireDB(t)
	cases := []struct {
		res    string
		filter string
		where  string // same condition in SQL
	}{
		{"drivers", "last_name,ilike,VER", "last_name ILIKE '%ver%'"},
		{"drivers", "last_name,ilike,%", "last_name LIKE '%\\%%'"}, // % is matched literally, not as a wildcard
		{"drivers", "code,neq,VER", "code <> 'VER'"},
		{"races", "season,gte,2026", "season >= 2026"},
		{"races", "round,lt,5", "round < 5"},
		{"circuits", "length_km,gt,10", "length_km > 10"},
		{"chassis", "weight_kg,lte,768", "weight_kg <= 768"},
		{"categories", "is_active,eq,false", "NOT is_active"},
		{"categories", "id,in,f1|wec|nope", "id IN ('f1', 'wec')"},
		{"races", "winner_id,in," + driverVer + "|" + driverBag, fmt.Sprintf("winner_id IN ('%s', '%s')", driverVer, driverBag)},
		{"article-texts", "article_id,eq," + articleMonaco, fmt.Sprintf("article_id = '%s'", articleMonaco)},
		{"article-texts", "article_id,eq," + articleDraft, "false"},
		{"article-texts", "title,ilike,draft", "false"}, // the draft's text is hidden even when it matches
	}
	byPath := map[string]resource{}
	for _, r := range resources {
		byPath[r.path] = r
	}
	for _, c := range cases {
		t.Run(c.res+"?"+c.filter, func(t *testing.T) {
			res := byPath[c.res]
			rows := getList(t, res.path, url.Values{"filter": {c.filter}})
			assertIDs(t, ids(rows), res.oracle(t, c.where, res.defaultOrder, pageSize+1))
		})
	}
}

// Several filters and sorts in one request.
func TestCombinedQueries(t *testing.T) {
	requireDB(t)
	cases := []struct {
		res     string
		query   url.Values
		where   string
		orderBy string
	}{
		{"races", url.Values{"filter": {"category_id,eq,f1", "season,eq,2026"}}, "category_id = 'f1' AND season = 2026", "race_date DESC, id ASC"},
		{"races", url.Values{"filter": {"season,eq,2025", "round,gte,5", "round,lte,12"}}, "season = 2025 AND round BETWEEN 5 AND 12", "race_date DESC, id ASC"},
		{"races", url.Values{"filter": {"category_id,eq,f1", "category_id,eq,wec"}}, "false", "race_date DESC, id ASC"},
		{"races", url.Values{"sort": {"category_id,season,desc,round"}}, "", "category_id ASC, season DESC, round ASC, id ASC"},
		{"races", url.Values{"sort": {"season,desc", "name"}, "filter": {"category_id,eq,motogp"}}, "category_id = 'motogp'", "season DESC, name ASC, id ASC"},
		{"articles", url.Values{"filter": {"entity_type,eq,race", "entity_id,eq," + raceMonaco}}, fmt.Sprintf("entity_type = 'race' AND entity_id = '%s'", raceMonaco), "published_at DESC, id ASC"},
		{"engines", url.Values{"sort": {"aspiration,desc,power_kw"}}, "", "aspiration DESC, power_kw ASC, id ASC"},
		{"categories", url.Values{"sort": {"is_active,id,desc"}}, "", "is_active ASC, id DESC"},
	}
	byPath := map[string]resource{}
	for _, r := range resources {
		byPath[r.path] = r
	}
	for _, c := range cases {
		t.Run(c.res+"?"+c.query.Encode(), func(t *testing.T) {
			res := byPath[c.res]
			rows := getList(t, res.path, c.query)
			assertIDs(t, ids(rows), res.oracle(t, c.where, c.orderBy, pageSize+1))
		})
	}

	t.Run("too many filters", func(t *testing.T) {
		filters := make([]string, 11)
		for i := range filters {
			filters[i] = "season,gte,2000"
		}
		assertError(t, "/races", url.Values{"filter": filters}, http.StatusBadRequest)
	})
}

func TestHiddenArticles(t *testing.T) {
	requireDB(t)
	for _, id := range []string{articleDraft, articleFuture} {
		assertError(t, "/articles/"+id, nil, http.StatusNotFound)
		if rows := getList(t, "article-texts", url.Values{"filter": {"article_id,eq," + id}}); len(rows) != 0 {
			t.Errorf("texts of hidden article %s are served: %v", id, rows)
		}
	}
}

func TestNotFoundRoute(t *testing.T) {
	requireDB(t)
	assertError(t, "/nope", nil, http.StatusNotFound)
}

type searchResult struct {
	ID       string  `json:"id"`
	Language *string `json:"language"`
	Title    *string `json:"title"`
	Preview  []struct {
		Text  string `json:"text"`
		Match bool   `json:"match"`
	} `json:"preview"`
}

func (r searchResult) previewText() string {
	var b strings.Builder
	for _, s := range r.Preview {
		b.WriteString(s.Text)
	}
	return b.String()
}

func TestArticleSearch(t *testing.T) {
	requireDB(t)

	search := func(t *testing.T, q url.Values) []searchResult {
		t.Helper()
		var results []searchResult
		getJSON(t, "/articles/search", q, http.StatusOK, &results)
		return results
	}
	only := func(t *testing.T, results []searchResult, id string) searchResult {
		t.Helper()
		if len(results) != 1 || results[0].ID != id {
			t.Fatalf("results = %v, want only %s", results, id)
		}
		return results[0]
	}
	find := func(t *testing.T, results []searchResult, id string) searchResult {
		t.Helper()
		i := slices.IndexFunc(results, func(r searchResult) bool { return r.ID == id })
		if i < 0 {
			t.Fatalf("results = %v, want %s among them", results, id)
		}
		return results[i]
	}

	t.Run("marks matched words", func(t *testing.T) {
		r := only(t, search(t, url.Values{"q": {"pole"}, "lang": {"en"}}), articleMonaco)
		if !slices.ContainsFunc(r.Preview, func(s struct {
			Text  string `json:"text"`
			Match bool   `json:"match"`
		}) bool {
			return s.Match && strings.EqualFold(s.Text, "pole")
		}) {
			t.Errorf("preview %v does not mark 'pole'", r.Preview)
		}
	})

	t.Run("prefers the requested language", func(t *testing.T) {
		r := find(t, search(t, url.Values{"q": {"verstappen"}, "lang": {"pt-BR"}}), articleMonaco)
		if r.Language == nil || *r.Language != "pt-br" {
			t.Errorf("language = %v, want pt-br", r.Language)
		}
		r = find(t, search(t, url.Values{"q": {"verstappen"}, "lang": {"en"}}), articleMonaco)
		if r.Title == nil || *r.Title != "Verstappen masters Monaco" {
			t.Errorf("title = %v, want the English one", r.Title)
		}
	})

	t.Run("phrases", func(t *testing.T) {
		only(t, search(t, url.Values{"q": {`"single pit stop"`}}), articleMonaco)
		if r := search(t, url.Values{"q": {`"stop single pit"`}}); len(r) != 0 {
			t.Errorf("out-of-order phrase matched: %v", r)
		}
	})

	t.Run("hidden articles never match", func(t *testing.T) {
		if r := search(t, url.Values{"q": {"draft or scheduled"}}); len(r) != 0 {
			t.Errorf("results = %v, want none", r)
		}
	})

	t.Run("previews are Markdown-free", func(t *testing.T) {
		r := only(t, search(t, url.Values{"q": {"controlled"}}), articleMonaco)
		if strings.Contains(r.previewText(), "**") {
			t.Errorf("preview %q still has Markdown", r.previewText())
		}
	})

	t.Run("browse lists every visible article newest first", func(t *testing.T) {
		results := search(t, url.Values{"lang": {"en"}})
		got := make([]string, len(results))
		for i, r := range results {
			got[i] = r.ID
		}
		articles := resources[slices.IndexFunc(resources, func(r resource) bool { return r.path == "articles" })]
		assertIDs(t, got, articles.oracle(t, "", "published_at DESC, id ASC", pageSize+1))
	})

	t.Run("bad input", func(t *testing.T) {
		for _, q := range []url.Values{
			{"lang": {"not_a_lang"}},
			{"page": {"0"}},
			{"q": {strings.Repeat("a", 201)}},
		} {
			assertError(t, "/articles/search", q, http.StatusBadRequest)
		}
	})
}

// oracle returns the ids the endpoint should serve, straight from SQL: rows matching res.where and the
// optional extra condition, in the given order, at most limit[0] of them.
func (res resource) oracle(t *testing.T, extra, orderBy string, limit ...int) []string {
	t.Helper()
	var conds []string
	for _, c := range []string{res.where, extra} {
		if c != "" {
			conds = append(conds, "("+c+")")
		}
	}
	query := "SELECT id::text FROM " + res.table
	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}
	query += " ORDER BY " + orderBy
	if len(limit) > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit[0])
	}

	rows, err := testDB.QueryContext(t.Context(), query)
	if err != nil {
		t.Fatalf("oracle %q: %v", query, err)
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return ids
}

func request(t *testing.T, path string, query url.Values) *httptest.ResponseRecorder {
	t.Helper()
	target := "/api/1" + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}
	rec := httptest.NewRecorder()
	testHandler.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, target, nil))
	return rec
}

func getJSON(t *testing.T, path string, query url.Values, wantStatus int, dst any) {
	t.Helper()
	rec := request(t, path, query)
	if rec.Code != wantStatus {
		t.Fatalf("GET %s?%s: status %d, want %d; body %s", path, query.Encode(), rec.Code, wantStatus, rec.Body)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), dst); err != nil {
		t.Fatalf("GET %s?%s: %v; body %s", path, query.Encode(), err, rec.Body)
	}
}

func getList(t *testing.T, path string, query url.Values) []map[string]any {
	t.Helper()
	var rows []map[string]any
	getJSON(t, "/"+path, query, http.StatusOK, &rows)
	if rows == nil {
		t.Fatalf("GET /%s?%s: got null, want a JSON array", path, query.Encode())
	}
	return rows
}

// assertError checks the status and the JSON error body.
func assertError(t *testing.T, path string, query url.Values, wantStatus int) {
	t.Helper()
	var body struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	}
	getJSON(t, path, query, wantStatus, &body)
	if body.Status != wantStatus || body.Error == "" {
		t.Errorf("GET %s?%s: error body %+v", path, query.Encode(), body)
	}
}

func ids(rows []map[string]any) []string {
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i], _ = row["id"].(string)
	}
	return out
}

func assertIDs(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("ids differ\n got: %v\nwant: %v", got, want)
	}
}
