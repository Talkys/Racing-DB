package queries

import (
	"backend/utils/constants"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type AllowedFilter struct {
	DBColumn   string
	AllowedOps map[string]string // "eq" -> "=", "ilike" -> "ILIKE"
}

type EntitySearchConfig struct {
	AllowedFilters map[string]AllowedFilter
	AllowedSorts   map[string]string // api_name -> db_column
	DefaultSortCol string
	DefaultSortDir string
	PageSize       int
}

type SearchCriteria struct {
	Page      int
	FilterCol string
	FilterOp  string
	FilterVal string
	SortCol   string
	SortDir   string
}

func BuildCriteriaFromURL(query url.Values) (*SearchCriteria, error) {
	// Handle page with a default value (e.g., 1) if omitted,
	// or keep strict parsing if 'page' is strictly mandatory.
	page := 1
	if pStr := query.Get("page"); pStr != "" {
		var err error
		page, err = strconv.Atoi(pStr)
		if err != nil {
			return nil, err
		}
	}

	var filterCol, filterOp, filterVal string
	if filterBlock := query.Get("filter"); filterBlock != "" {
		parts := strings.SplitN(filterBlock, ",", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid filter parameters: %s", filterBlock)
		}
		filterCol, filterOp, filterVal = parts[0], parts[1], parts[2]
	}

	// strings.Cut handles missing commas safely out-of-the-box
	sortCol, sortDir, _ := strings.Cut(query.Get("sort"), ",")

	return &SearchCriteria{
		Page:      page,
		FilterCol: filterCol,
		FilterOp:  filterOp,
		FilterVal: filterVal,
		SortCol:   sortCol,
		SortDir:   sortDir,
	}, nil
}

func BuildSearchQuery(cols []string, table string, cfg EntitySearchConfig, sc SearchCriteria) (string, []any, error) {
	var queryBuilder strings.Builder
	args := make([]any, 0, 3)
	argIdx := 1

	columns := strings.Join(cols, ", ")

	queryBuilder.WriteString("select ")
	queryBuilder.WriteString(columns)
	queryBuilder.WriteString(" from ")
	queryBuilder.WriteString(table)

	if sc.FilterCol != "" {
		filterCfg, exists := cfg.AllowedFilters[sc.FilterCol]
		if !exists {
			return "", nil, fmt.Errorf("invalid filter column: %s", sc.FilterCol)
		}

		sqlOp, opExists := filterCfg.AllowedOps[sc.FilterOp]
		if !opExists {
			return "", nil, fmt.Errorf("invalid operator '%s' for column '%s'", sc.FilterOp, sc.FilterCol)
		}

		val := sc.FilterVal
		if sc.FilterOp == constants.Contains {
			val = "%" + val + "%"
		}

		fmt.Fprintf(&queryBuilder, " WHERE %s %s $%d", filterCfg.DBColumn, sqlOp, argIdx)
		args = append(args, val)
		argIdx++
	}

	sortCol := cfg.DefaultSortCol
	hasCustomSort := false
	if sc.SortCol != "" {
		if mappedCol, exists := cfg.AllowedSorts[sc.SortCol]; exists {
			sortCol = mappedCol
			hasCustomSort = true
		} else {
			return "", nil, fmt.Errorf("invalid sort column: %s", sc.SortCol)
		}
	}

	sortDir := cfg.DefaultSortDir
	if strings.ToLower(sc.SortDir) == "desc" {
		sortDir = "DESC"
	}

	if hasCustomSort && sortCol != cfg.DefaultSortCol {
		fmt.Fprintf(&queryBuilder, " ORDER BY %s %s, %s %s", sortCol, sortDir, cfg.DefaultSortCol, sortDir)
	} else {
		fmt.Fprintf(&queryBuilder, " ORDER BY %s %s", sortCol, sortDir)
	}

	limit := cfg.PageSize + 1
	offset := (sc.Page - 1) * cfg.PageSize
	fmt.Fprintf(&queryBuilder, " LIMIT %d OFFSET %d", limit, offset)

	return queryBuilder.String(), args, nil
}
