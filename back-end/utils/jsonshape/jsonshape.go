// Package jsonshape checks decoded JSON against a contract written as code: required keys, types, enums
// and ranges. It reports every problem it finds instead of stopping at the first one, each as
// "path: what is wrong". Keys the contract does not mention are allowed, so a contract can grow without
// breaking older rows.
//
//	c, root := jsonshape.Parse(raw)
//	lap := root.Object("fastest_lap")
//	lap.Int("lap_number", 1, 500)
//	problems := c.Problems()
//
// Once a value is reported (missing, wrong type), checks on its children are skipped silently, so one
// bad object gives one problem and not one per field.
package jsonshape

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/google/uuid"
)

// Checker collects the problems found by the Object values it hands out.
type Checker struct {
	problems []string
}

// Problems returns what was found so far, in the order the checks ran; nil when everything matched.
func (c *Checker) Problems() []string {
	return c.problems
}

func (c *Checker) report(path, format string, v ...any) {
	if path == "" {
		path = "(root)"
	}
	c.problems = append(c.problems, path+": "+fmt.Sprintf(format, v...))
}

// Object is a JSON object under check. Its zero value (from a failed check) skips every check made on it.
type Object struct {
	c    *Checker
	path string
	m    map[string]any
}

// Parse decodes raw and returns its root, which must be an object. Numbers keep their literal form, so
// Int can tell 3 from 3.5.
func Parse(raw []byte) (*Checker, Object) {
	c := &Checker{}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		c.report("", "not valid JSON")
		return c, Object{}
	}
	m, ok := v.(map[string]any)
	if !ok {
		c.report("", "want an object, got %s", typeName(v))
		return c, Object{}
	}
	return c, Object{c: c, m: m}
}

func (o Object) child(key string) string {
	if o.path == "" {
		return key
	}
	return o.path + "." + key
}

// get returns the value under key, reporting it when absent or null.
func (o Object) get(key string) (any, bool) {
	if o.m == nil {
		return nil, false
	}
	v, ok := o.m[key]
	if !ok || v == nil {
		o.c.report(o.child(key), "missing")
		return nil, false
	}
	return v, true
}

// Object checks that key holds an object and returns it for further checks.
func (o Object) Object(key string) Object {
	v, ok := o.get(key)
	if !ok {
		return Object{}
	}
	m, ok := v.(map[string]any)
	if !ok {
		o.c.report(o.child(key), "want an object, got %s", typeName(v))
		return Object{}
	}
	return Object{c: o.c, path: o.child(key), m: m}
}

// Objects checks that key holds an array of objects and calls each for every element. An empty array is fine.
func (o Object) Objects(key string, each func(Object)) {
	v, ok := o.get(key)
	if !ok {
		return
	}
	items, ok := v.([]any)
	if !ok {
		o.c.report(o.child(key), "want an array, got %s", typeName(v))
		return
	}
	for i, item := range items {
		path := fmt.Sprintf("%s[%d]", o.child(key), i)
		m, ok := item.(map[string]any)
		if !ok {
			o.c.report(path, "want an object, got %s", typeName(item))
			continue
		}
		each(Object{c: o.c, path: path, m: m})
	}
}

// String checks that key holds a non-empty string and returns it ("" after a problem).
func (o Object) String(key string) string {
	v, ok := o.get(key)
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		o.c.report(o.child(key), "want a string, got %s", typeName(v))
		return ""
	}
	if strings.TrimSpace(s) == "" {
		o.c.report(o.child(key), "empty")
		return ""
	}
	return s
}

// StringMap checks that key holds an object whose values are all strings (it may be empty).
func (o Object) StringMap(key string) {
	obj := o.Object(key)
	for _, k := range slices.Sorted(maps.Keys(obj.m)) {
		if _, ok := obj.m[k].(string); !ok {
			o.c.report(obj.child(k), "want a string, got %s", typeName(obj.m[k]))
		}
	}
}

// Enum checks that key holds one of the given strings.
func (o Object) Enum(key string, values ...string) {
	s := o.String(key)
	if s != "" && !slices.Contains(values, s) {
		o.c.report(o.child(key), "%q is not one of %s", s, strings.Join(values, ", "))
	}
}

// UUID checks that key holds a UUID in its usual 36-character form.
func (o Object) UUID(key string) {
	s := o.String(key)
	if s == "" {
		return
	}
	if len(s) != 36 || uuid.Validate(s) != nil {
		o.c.report(o.child(key), "%q is not a UUID", s)
	}
}

// Match checks that key holds a string matching re; what describes the expected form in the problem.
func (o Object) Match(key string, re *regexp.Regexp, what string) {
	s := o.String(key)
	if s != "" && !re.MatchString(s) {
		o.c.report(o.child(key), "%q is not %s", s, what)
	}
}

// Int checks that key holds an integer in [lo, hi].
func (o Object) Int(key string, lo, hi int64) {
	n, ok := o.number(key)
	if !ok {
		return
	}
	i, err := n.Int64()
	if err != nil {
		o.c.report(o.child(key), "want an integer, got %s", n)
		return
	}
	if i < lo || i > hi {
		o.c.report(o.child(key), "%d is outside %d–%d", i, lo, hi)
	}
}

// Number checks that key holds a number in [lo, hi].
func (o Object) Number(key string, lo, hi float64) {
	n, ok := o.number(key)
	if !ok {
		return
	}
	f, err := n.Float64()
	if err != nil {
		o.c.report(o.child(key), "%s is not a finite number", n)
		return
	}
	if f < lo || f > hi {
		o.c.report(o.child(key), "%s is outside %g–%g", n, lo, hi)
	}
}

func (o Object) number(key string) (json.Number, bool) {
	v, ok := o.get(key)
	if !ok {
		return "", false
	}
	n, ok := v.(json.Number)
	if !ok {
		o.c.report(o.child(key), "want a number, got %s", typeName(v))
		return "", false
	}
	return n, true
}

func typeName(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case map[string]any:
		return "an object"
	case []any:
		return "an array"
	case string:
		return "a string"
	case json.Number:
		return "a number"
	case bool:
		return "a boolean"
	default:
		return fmt.Sprintf("%T", v)
	}
}
