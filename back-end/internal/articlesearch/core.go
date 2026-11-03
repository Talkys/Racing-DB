package articlesearch

import (
	"context"
	"strings"
	"unicode"
	"unicode/utf8"
)

// previewRunes caps browse previews (search previews are sized by ts_headline).
const previewRunes = 240

type ArticleSearchCore struct {
	repo *ArticleSearchRepository
}

func NewArticleSearchCore() *ArticleSearchCore {
	return &ArticleSearchCore{
		repo: NewArticleSearchRepository(),
	}
}

func (c *ArticleSearchCore) SearchArticles(ctx context.Context, sc SearchCriteria) ([]ArticleSearchResult, error) {
	rows, err := c.repo.SearchArticles(ctx, sc)
	if err != nil {
		return nil, err
	}

	results := make([]ArticleSearchResult, 0, len(rows))
	for _, row := range rows {
		result := row.ArticleSearchResult
		result.Preview = buildPreview(row.Preview, row.Title, sc.Query != "")
		results = append(results, result)
	}
	return results, nil
}

func buildPreview(raw *string, title *string, highlighted bool) []PreviewSegment {
	if raw == nil {
		return []PreviewSegment{}
	}

	segments := splitHighlights(*raw)
	if title != nil {
		// Bodies usually open with their title as a heading; the client already shows the title.
		segments = trimPrefix(segments, *title)
	}
	if !highlighted {
		segments = truncate(segments, previewRunes)
	}
	return segments
}

// splitHighlights turns "a b c" into [{a } {b match} { c}].
func splitHighlights(s string) []PreviewSegment {
	segments := []PreviewSegment{}
	for {
		start := strings.Index(s, matchStart)
		if start < 0 {
			break
		}
		stop := strings.Index(s[start:], matchStop)
		if stop < 0 {
			break
		}
		stop += start
		if start > 0 {
			segments = append(segments, PreviewSegment{Text: s[:start]})
		}
		segments = append(segments, PreviewSegment{Text: s[start+len(matchStart) : stop], Match: true})
		s = s[stop+len(matchStop):]
	}
	// Unpaired markers cannot come from the source text, but never leak them.
	s = strings.NewReplacer(matchStart, "", matchStop, "").Replace(s)
	if s != "" {
		segments = append(segments, PreviewSegment{Text: s})
	}
	return segments
}

// trimPrefix drops the title (case-insensitive) from the start of the preview, plus the separators after it.
// The preview is left untouched when it is nothing but the title.
func trimPrefix(segments []PreviewSegment, title string) []PreviewSegment {
	var plain strings.Builder
	for _, seg := range segments {
		plain.WriteString(seg.Text)
	}
	text := plain.String()
	cut, ok := foldPrefixLen(text, title)
	if title == "" || !ok {
		return segments
	}

	for _, r := range text[cut:] {
		if !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			break
		}
		cut += len(string(r))
	}
	if cut == len(text) {
		return segments
	}

	trimmed := []PreviewSegment{}
	for _, seg := range segments {
		if cut >= len(seg.Text) {
			cut -= len(seg.Text)
			continue
		}
		seg.Text = seg.Text[cut:]
		cut = 0
		trimmed = append(trimmed, seg)
	}
	return trimmed
}

// foldPrefixLen reports whether text starts with prefix, case-insensitively, and the byte length of the match
// in text. It compares rune by rune, since folded runes can differ in UTF-8 length (ß is 2 bytes, ẞ is 3).
func foldPrefixLen(text, prefix string) (int, bool) {
	n := 0
	for prefix != "" {
		if n == len(text) {
			return 0, false
		}
		want, wantSize := utf8.DecodeRuneInString(prefix)
		got, gotSize := utf8.DecodeRuneInString(text[n:])
		if !strings.EqualFold(string(want), string(got)) {
			return 0, false
		}
		prefix = prefix[wantSize:]
		n += gotSize
	}
	return n, true
}

// truncate shortens an unhighlighted preview to about max runes, at a word boundary.
func truncate(segments []PreviewSegment, max int) []PreviewSegment {
	if len(segments) != 1 {
		return segments
	}
	runes := []rune(segments[0].Text)
	if len(runes) <= max {
		return segments
	}
	cut := string(runes[:max])
	if i := strings.LastIndexFunc(cut, unicode.IsSpace); i > max/2 {
		cut = cut[:i]
	}
	return []PreviewSegment{{Text: strings.TrimRightFunc(cut, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsPunct(r) }) + " …"}}
}
