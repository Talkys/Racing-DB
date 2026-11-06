package queries

import (
	"backend/utils/apperror"
	"backend/utils/constants"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Kind is the type of a filterable column. Filter values are parsed into it before reaching the
// database, so a bad value is a 400 instead of a database error.
type Kind int

const (
	KindText Kind = iota
	KindInt
	KindDecimal
	KindBool
	KindUUID
)

// Operator sets per column type: api operator -> SQL operator.
var (
	TextOps   = map[string]string{constants.Equal: "=", constants.NotEqual: "<>", constants.Contains: "ILIKE"}
	NumberOps = map[string]string{constants.Equal: "=", constants.Greater: ">", constants.GreaterEqual: ">=", constants.Less: "<", constants.LessEqual: "<="}
	BoolOps   = map[string]string{constants.Equal: "="}
	IDOps     = map[string]string{constants.Equal: "=", constants.In: "= ANY"}
)

// tiebreaker keeps pagination stable when the sort column has duplicates. Every table has an id.
const tiebreaker = "id"

// maxInValues bounds `in` lists; a page never needs more related ids than this.
const maxInValues = 100

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

type AllowedFilter struct {
	DBColumn   string
	Kind       Kind
	AllowedOps map[string]string // "eq" -> "=", "ilike" -> "ILIKE"
}

func TextFilter(col string) AllowedFilter {
	return AllowedFilter{DBColumn: col, Kind: KindText, AllowedOps: TextOps}
}
func IntFilter(col string) AllowedFilter {
	return AllowedFilter{DBColumn: col, Kind: KindInt, AllowedOps: NumberOps}
}
func DecimalFilter(col string) AllowedFilter {
	return AllowedFilter{DBColumn: col, Kind: KindDecimal, AllowedOps: NumberOps}
}
func BoolFilter(col string) AllowedFilter {
	return AllowedFilter{DBColumn: col, Kind: KindBool, AllowedOps: BoolOps}
}
func UUIDFilter(col string) AllowedFilter {
	return AllowedFilter{DBColumn: col, Kind: KindUUID, AllowedOps: IDOps}
}

type EntitySearchConfig struct {
	BaseCondition  string // SQL condition always applied (search and get by id), e.g. hiding draft articles
	AllowedFilters map[string]AllowedFilter
	AllowedSorts   map[string]string // api_name -> db_column
	DefaultOrder   string            // ORDER BY used when no sort is requested; should end with the id tiebreaker
	PageSize       int
}

// Filter is one `column,operator,value` expression.
type Filter struct {
	Col, Op, Val string
}

// Sort is one ordering column; Dir is "", "asc" or "desc" as the client wrote it.
type Sort struct {
	Col, Dir string
}

type SearchCriteria struct {
	Page    int
	Filters []Filter // AND-ed together
	Sorts   []Sort   // in priority order
}

// maxFilters bounds the filters of one request; more than a handful is never a real search.
const maxFilters = 10

// ParsePage reads the 1-based `page` param, defaulting to 1.
func ParsePage(query url.Values) (int, error) {
	pStr := query.Get("page")
	if pStr == "" {
		return 1, nil
	}
	page, err := strconv.Atoi(pStr)
	if err != nil || page < 1 {
		return 0, apperror.BadRequest("invalid page: '%s' (must be an integer >= 1)", pStr)
	}
	return page, nil
}

// BuildCriteriaFromURL reads `page`, any number of `filter=column,operator,value` params, and `sort` as a
// comma-separated list of columns, each optionally followed by its direction (`season,desc,name`), given
// once or repeated (`sort=season,desc&sort=name`).
func BuildCriteriaFromURL(query url.Values) (*SearchCriteria, error) {
	page, err := ParsePage(query)
	if err != nil {
		return nil, err
	}

	filterBlocks := query["filter"]
	if len(filterBlocks) > maxFilters {
		return nil, apperror.BadRequest("too many filters: %d (at most %d)", len(filterBlocks), maxFilters)
	}
	filters := make([]Filter, 0, len(filterBlocks))
	for _, block := range filterBlocks {
		parts := strings.SplitN(block, ",", 3)
		if len(parts) != 3 {
			return nil, apperror.BadRequest("invalid filter: '%s' (expected column,operator,value)", block)
		}
		filters = append(filters, Filter{Col: parts[0], Op: parts[1], Val: parts[2]})
	}

	var sorts []Sort
	for _, block := range query["sort"] {
		for _, token := range strings.Split(block, ",") {
			switch lower := strings.ToLower(token); {
			case token == "":
				// tolerate "name," and empty params
			case (lower == "asc" || lower == "desc") && len(sorts) > 0 && sorts[len(sorts)-1].Dir == "":
				sorts[len(sorts)-1].Dir = lower
			default:
				sorts = append(sorts, Sort{Col: token})
			}
		}
	}

	return &SearchCriteria{Page: page, Filters: filters, Sorts: sorts}, nil
}

func BuildSearchQuery(cols []string, table string, cfg EntitySearchConfig, sc SearchCriteria) (string, []any, error) {
	var queryBuilder strings.Builder
	args := make([]any, 0, len(sc.Filters))
	conditions := make([]string, 0, len(sc.Filters)+1)

	if cfg.BaseCondition != "" {
		conditions = append(conditions, "("+cfg.BaseCondition+")")
	}

	for _, f := range sc.Filters {
		filterCfg, exists := cfg.AllowedFilters[f.Col]
		if !exists {
			return "", nil, apperror.BadRequest("invalid filter column: %s", f.Col)
		}

		sqlOp, opExists := filterCfg.AllowedOps[f.Op]
		if !opExists {
			return "", nil, apperror.BadRequest("invalid operator '%s' for column '%s'", f.Op, f.Col)
		}

		val, err := parseFilterValue(filterCfg, f.Op, f.Val)
		if err != nil {
			return "", nil, apperror.BadRequest("invalid value '%s' for column '%s'", f.Val, f.Col)
		}

		args = append(args, val)
		if f.Op == constants.In {
			conditions = append(conditions, fmt.Sprintf("%s %s($%d)", filterCfg.DBColumn, sqlOp, len(args)))
		} else {
			conditions = append(conditions, fmt.Sprintf("%s %s $%d", filterCfg.DBColumn, sqlOp, len(args)))
		}
	}

	orderBy, err := buildOrderBy(cfg, sc)
	if err != nil {
		return "", nil, err
	}

	queryBuilder.WriteString("select ")
	queryBuilder.WriteString(strings.Join(cols, ", "))
	queryBuilder.WriteString(" from ")
	queryBuilder.WriteString(table)
	if len(conditions) > 0 {
		queryBuilder.WriteString(" where ")
		queryBuilder.WriteString(strings.Join(conditions, " and "))
	}
	queryBuilder.WriteString(" order by ")
	queryBuilder.WriteString(orderBy)

	// One extra row tells the client whether another page exists.
	limit := cfg.PageSize + 1
	offset := (sc.Page - 1) * cfg.PageSize
	fmt.Fprintf(&queryBuilder, " limit %d offset %d", limit, offset)

	return queryBuilder.String(), args, nil
}

// BuildGetQuery selects one row by id, still honouring the config's BaseCondition.
func BuildGetQuery(cols []string, table string, cfg EntitySearchConfig) string {
	query := fmt.Sprintf("select %s from %s where %s = $1", strings.Join(cols, ", "), table, tiebreaker)
	if cfg.BaseCondition != "" {
		query += " and (" + cfg.BaseCondition + ")"
	}
	return query
}

func buildOrderBy(cfg EntitySearchConfig, sc SearchCriteria) (string, error) {
	if len(sc.Sorts) == 0 {
		return cfg.DefaultOrder, nil
	}

	terms := make([]string, 0, len(sc.Sorts)+1)
	seen := make(map[string]bool, len(sc.Sorts))
	for _, s := range sc.Sorts {
		sortCol, exists := cfg.AllowedSorts[s.Col]
		if !exists {
			return "", apperror.BadRequest("invalid sort column or direction: %s", s.Col)
		}
		if seen[sortCol] {
			return "", apperror.BadRequest("duplicate sort column: %s", s.Col)
		}
		seen[sortCol] = true

		sortDir := constants.Ascendent
		if s.Dir == "desc" {
			sortDir = constants.Descendent
		}
		terms = append(terms, sortCol+" "+sortDir)
	}

	if !seen[tiebreaker] {
		terms = append(terms, tiebreaker+" "+constants.Ascendent)
	}
	return strings.Join(terms, ", "), nil
}

func parseFilterValue(f AllowedFilter, op string, raw string) (any, error) {
	switch op {
	case constants.Contains:
		return "%" + likeEscaper.Replace(raw) + "%", nil
	case constants.In:
		parts := strings.Split(raw, "|")
		if len(parts) > maxInValues {
			return nil, fmt.Errorf("too many values: %d", len(parts))
		}
		switch f.Kind {
		case KindInt:
			return parseList(parts, parseInt)
		case KindDecimal:
			return parseList(parts, parseDecimal)
		case KindBool:
			return parseList(parts, strconv.ParseBool)
		case KindUUID:
			return parseList(parts, uuid.Parse)
		default:
			return parts, nil
		}
	}

	switch f.Kind {
	case KindInt:
		return parseInt(raw)
	case KindDecimal:
		return parseDecimal(raw)
	case KindBool:
		return strconv.ParseBool(raw)
	case KindUUID:
		return uuid.Parse(raw)
	default:
		return raw, nil
	}
}

func parseInt(s string) (int64, error)       { return strconv.ParseInt(s, 10, 64) }
func parseDecimal(s string) (float64, error) { return strconv.ParseFloat(s, 64) }

func parseList[T any](parts []string, parse func(string) (T, error)) ([]T, error) {
	values := make([]T, 0, len(parts))
	for _, p := range parts {
		v, err := parse(p)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}
