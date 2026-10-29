package articlesearch

import (
	"reflect"
	"strings"
	"testing"
)

func ptr(s string) *string { return &s }

func TestBuildPreview(t *testing.T) {
	tests := []struct {
		name        string
		raw         *string
		title       *string
		highlighted bool
		want        []PreviewSegment
	}{
		{
			name: "no text",
			want: []PreviewSegment{},
		},
		{
			name:        "highlights become segments",
			raw:         ptr("a drive from pole, with a single"),
			highlighted: true,
			want:        []PreviewSegment{{Text: "a drive from "}, {Text: "pole", Match: true}, {Text: ", with a "}, {Text: "single", Match: true}},
		},
		{
			name:        "leading title is dropped, even when highlighted",
			raw:         ptr("Verstappen masters Monaco: A controlled drive"),
			title:       ptr("verstappen masters monaco"),
			highlighted: true,
			want:        []PreviewSegment{{Text: "A controlled drive"}},
		},
		{
			name:  "a preview that is only the title is kept",
			raw:   ptr("Max Verstappen"),
			title: ptr("Max Verstappen"),
			want:  []PreviewSegment{{Text: "Max Verstappen"}},
		},
		{
			name:  "title elsewhere is kept",
			raw:   ptr("Intro. Max Verstappen"),
			title: ptr("Max Verstappen"),
			want:  []PreviewSegment{{Text: "Intro. Max Verstappen"}},
		},
		{
			name:        "unpaired markers never leak",
			raw:         ptr("broken marker"),
			highlighted: true,
			want:        []PreviewSegment{{Text: "broken marker"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPreview(tt.raw, tt.title, tt.highlighted)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v\nwant %#v", got, tt.want)
			}
		})
	}
}

func TestBrowsePreviewIsTruncatedAtAWord(t *testing.T) {
	long := strings.Repeat("Vitória sem sustos em Monte Carlo. ", 20)
	got := buildPreview(&long, nil, false)
	if len(got) != 1 || !strings.HasSuffix(got[0].Text, " …") {
		t.Fatalf("got %#v", got)
	}
	text := strings.TrimSuffix(got[0].Text, " …")
	if n := len([]rune(text)); n > previewRunes || n < previewRunes/2 {
		t.Errorf("length %d out of range", n)
	}
	if !strings.HasPrefix(long, text) {
		t.Errorf("not a prefix: %q", text)
	}
}
