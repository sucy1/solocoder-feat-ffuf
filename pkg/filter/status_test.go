package filter

import (
	"strings"
	"testing"

	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

func TestNewStatusFilter(t *testing.T) {
	f, _ := NewStatusFilter("200,301,400-410,500")
	statusRepr := f.Repr()
	if !strings.Contains(statusRepr, "200,301,400-410,500") {
		t.Errorf("Status filter was expected to have 4 values")
	}
}

func TestNewStatusFilterError(t *testing.T) {
	_, err := NewStatusFilter("invalid")
	if err == nil {
		t.Errorf("Was expecting an error from errenous input data")
	}
}

func TestStatusFiltering(t *testing.T) {
	f, _ := NewStatusFilter("200,301,400-498,500")
	for i, test := range []struct {
		input  int64
		output bool
	}{
		{200, true},
		{301, true},
		{500, true},
		{4, false},
		{399, false},
		{400, true},
		{444, true},
		{498, true},
		{499, false},
		{302, false},
	} {
		resp := ffuf.Response{StatusCode: test.input}
		filterReturn, _ := f.Filter(&resp)
		if filterReturn != test.output {
			t.Errorf("Filter test %d: Was expecing filter return value of %t but got %t", i, test.output, filterReturn)
		}
	}
}

func TestStatusFilterRegex(t *testing.T) {
	f, err := NewStatusFilter("/admin/")
	if err != nil {
		t.Fatalf("Failed to create status filter with regex: %s", err)
	}
	sf, ok := f.(*StatusFilter)
	if !ok {
		t.Fatalf("Expected *StatusFilter, got %T", f)
	}
	if len(sf.Regexps) != 1 {
		t.Errorf("Expected 1 regexp, got %d", len(sf.Regexps))
	}
}

func TestStatusFilterRegexCaseInsensitive(t *testing.T) {
	f, err := NewStatusFilter("/Admin/i")
	if err != nil {
		t.Fatalf("Failed to create status filter with case-insensitive regex: %s", err)
	}
	sf, ok := f.(*StatusFilter)
	if !ok {
		t.Fatalf("Expected *StatusFilter, got %T", f)
	}
	if len(sf.Regexps) != 1 {
		t.Errorf("Expected 1 regexp, got %d", len(sf.Regexps))
	}
}

func TestStatusFilterStatusAndRegex(t *testing.T) {
	f, err := NewStatusFilter("200,/(admin|login)/i")
	if err != nil {
		t.Fatalf("Failed to create status filter with status and regex: %s", err)
	}
	sf, ok := f.(*StatusFilter)
	if !ok {
		t.Fatalf("Expected *StatusFilter, got %T", f)
	}
	if len(sf.Value) != 1 {
		t.Errorf("Expected 1 status range, got %d", len(sf.Value))
	}
	if len(sf.Regexps) != 1 {
		t.Errorf("Expected 1 regexp, got %d", len(sf.Regexps))
	}
}

func TestStatusFilterFilteringWithRegex(t *testing.T) {
	f, _ := NewStatusFilter("/test/")
	for i, test := range []struct {
		body     string
		status   int64
		expected bool
	}{
		{"this is a test body", 200, true},
		{"no match here", 200, false},
		{"TEST uppercase", 200, false},
		{"", 200, false},
	} {
		resp := ffuf.Response{StatusCode: test.status, Data: []byte(test.body)}
		result, _ := f.Filter(&resp)
		if result != test.expected {
			t.Errorf("Filter test %d: Expected %t but got %t for body: %q", i, test.expected, result, test.body)
		}
	}
}

func TestStatusFilterFilteringWithStatusAndRegex(t *testing.T) {
	f, _ := NewStatusFilter("200,/(admin|login)/i")
	for i, test := range []struct {
		body     string
		status   int64
		expected bool
	}{
		{"admin page", 200, true},
		{"Admin Page", 200, true},
		{"LOGIN form", 200, true},
		{"other page", 200, false},
		{"admin page", 404, false},
		{"other page", 200, false},
	} {
		resp := ffuf.Response{StatusCode: test.status, Data: []byte(test.body)}
		result, _ := f.Filter(&resp)
		if result != test.expected {
			t.Errorf("Filter test %d: Expected %t but got %t for body: %q, status: %d", i, test.expected, result, test.body, test.status)
		}
	}
}
