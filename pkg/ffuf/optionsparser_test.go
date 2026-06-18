package ffuf

import (
	"strings"
	"testing"
)

func TestTemplatePresent(t *testing.T) {
	template := "§"

	headers := make(map[string]string)
	headers["foo"] = "§bar§"
	headers["omg"] = "bbq"
	headers["§world§"] = "Ooo"

	goodConf := Config{
		Url:     "https://example.com/fooo/bar?test=§value§&order[§0§]=§foo§",
		Method:  "PO§ST§",
		Headers: headers,
		Data:    "line=Can we pull back the §veil§ of §static§ and reach in to the source of §all§ being?&commit=true",
	}

	if !templatePresent(template, &goodConf) {
		t.Errorf("Expected-good config failed validation")
	}

	badConfMethod := Config{
		Url:     "https://example.com/fooo/bar?test=§value§&order[§0§]=§foo§",
		Method:  "POST§",
		Headers: headers,
		Data:    "line=Can we pull back the §veil§ of §static§ and reach in to the source of §all§ being?&commit=§true§",
	}

	if templatePresent(template, &badConfMethod) {
		t.Errorf("Expected-bad config (Method) failed validation")
	}

	badConfURL := Config{
		Url:     "https://example.com/fooo/bar?test=§value§&order[0§]=§foo§",
		Method:  "§POST§",
		Headers: headers,
		Data:    "line=Can we pull back the §veil§ of §static§ and reach in to the source of §all§ being?&commit=§true§",
	}

	if templatePresent(template, &badConfURL) {
		t.Errorf("Expected-bad config (URL) failed validation")
	}

	badConfData := Config{
		Url:     "https://example.com/fooo/bar?test=§value§&order[§0§]=§foo§",
		Method:  "§POST§",
		Headers: headers,
		Data:    "line=Can we pull back the §veil of §static§ and reach in to the source of §all§ being?&commit=§true§",
	}

	if templatePresent(template, &badConfData) {
		t.Errorf("Expected-bad config (Data) failed validation")
	}

	headers["kingdom"] = "§candy"

	badConfHeaderValue := Config{
		Url:     "https://example.com/fooo/bar?test=§value§&order[§0§]=§foo§",
		Method:  "PO§ST§",
		Headers: headers,
		Data:    "line=Can we pull back the §veil§ of §static§ and reach in to the source of §all§ being?&commit=true",
	}

	if templatePresent(template, &badConfHeaderValue) {
		t.Errorf("Expected-bad config (Header value) failed validation")
	}

	headers["kingdom"] = "candy"
	headers["§kingdom"] = "candy"

	badConfHeaderKey := Config{
		Url:     "https://example.com/fooo/bar?test=§value§&order[§0§]=§foo§",
		Method:  "PO§ST§",
		Headers: headers,
		Data:    "line=Can we pull back the §veil§ of §static§ and reach in to the source of §all§ being?&commit=true",
	}

	if templatePresent(template, &badConfHeaderKey) {
		t.Errorf("Expected-bad config (Header key) failed validation")
	}
}

func TestProxyParsing(t *testing.T) {
	configOptions := NewConfigOptions()
	errorString := "Bad proxy url (-x) format. Expected http, https or socks5 url"

	// http should work
	configOptions.HTTP.ProxyURL = "http://127.0.0.1:8080"
	_, err := ConfigFromOptions(configOptions, nil, nil)
	if strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected http proxy string to work")
	}

	// https should work
	configOptions.HTTP.ProxyURL = "https://127.0.0.1"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected https proxy string to work")
	}

	// socks5 should work
	configOptions.HTTP.ProxyURL = "socks5://127.0.0.1"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected socks5 proxy string to work")
	}

	// garbage data should FAIL
	configOptions.HTTP.ProxyURL = "Y0 y0 it's GREASE"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if !strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected garbage proxy string to fail")
	}

	// Opaque URLs with the right scheme should FAIL
	configOptions.HTTP.ProxyURL = "http:sixhours@dungeon"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if !strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected opaque proxy string to fail")
	}

	// Unsupported protocols should FAIL
	configOptions.HTTP.ProxyURL = "imap://127.0.0.1"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if !strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected proxy string with unsupported protocol to fail")
	}
}

func TestReplayProxyParsing(t *testing.T) {
	configOptions := NewConfigOptions()
	errorString := "Bad replay-proxy url (-replay-proxy) format. Expected http, https or socks5 url"

	// http should work
	configOptions.HTTP.ReplayProxyURL = "http://127.0.0.1:8080"
	_, err := ConfigFromOptions(configOptions, nil, nil)
	if strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected http replay proxy string to work")
	}

	// https should work
	configOptions.HTTP.ReplayProxyURL = "https://127.0.0.1"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected https proxy string to work")
	}

	// socks5 should work
	configOptions.HTTP.ReplayProxyURL = "socks5://127.0.0.1"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected socks5 proxy string to work")
	}

	// garbage data should FAIL
	configOptions.HTTP.ReplayProxyURL = "Y0 y0 it's GREASE"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if !strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected garbage proxy string to fail")
	}

	// Opaque URLs with the right scheme should FAIL
	configOptions.HTTP.ReplayProxyURL = "http:sixhours@dungeon"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if !strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected opaque proxy string to fail")
	}

	// Unsupported protocols should FAIL
	configOptions.HTTP.ReplayProxyURL = "imap://127.0.0.1"
	_, err = ConfigFromOptions(configOptions, nil, nil)
	if !strings.Contains(err.Error(), errorString) {
		t.Errorf("Expected proxy string with unsupported protocol to fail")
	}
}

func TestDelayParsing(t *testing.T) {
	tests := []struct {
		name       string
		delayInput string
		wantErr    bool
		wantMin    float64
		wantMax    float64
		wantRange  bool
		wantHas    bool
	}{
		{
			name:       "single float",
			delayInput: "0.5",
			wantErr:    false,
			wantMin:    0.5,
			wantMax:    0,
			wantRange:  false,
			wantHas:    true,
		},
		{
			name:       "range float",
			delayInput: "0.5-2.0",
			wantErr:    false,
			wantMin:    0.5,
			wantMax:    2.0,
			wantRange:  true,
			wantHas:    true,
		},
		{
			name:       "range with integers",
			delayInput: "1-3",
			wantErr:    false,
			wantMin:    1.0,
			wantMax:    3.0,
			wantRange:  true,
			wantHas:    true,
		},
		{
			name:       "invalid single",
			delayInput: "abc",
			wantErr:    true,
		},
		{
			name:       "invalid range first",
			delayInput: "abc-2.0",
			wantErr:    true,
		},
		{
			name:       "invalid range second",
			delayInput: "0.5-xyz",
			wantErr:    true,
		},
		{
			name:       "too many parts",
			delayInput: "0.5-1.0-2.0",
			wantErr:    true,
		},
		{
			name:       "empty string",
			delayInput: "",
			wantErr:    false,
			wantMin:    0,
			wantMax:    0,
			wantRange:  false,
			wantHas:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configOptions := NewConfigOptions()
			configOptions.General.Delay = tt.delayInput
			configOptions.Input.Wordlists = []string{"test.txt"}
			configOptions.Input.InputMode = "clusterbomb"
			configOptions.HTTP.URL = "https://example.com/FUZZ"
			conf, err := ConfigFromOptions(configOptions, nil, nil)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for delay input %q, but got none", tt.delayInput)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for delay input %q: %v", tt.delayInput, err)
				return
			}

			if conf.Delay.Min != tt.wantMin {
				t.Errorf("Delay.Min = %v, want %v", conf.Delay.Min, tt.wantMin)
			}
			if conf.Delay.Max != tt.wantMax {
				t.Errorf("Delay.Max = %v, want %v", conf.Delay.Max, tt.wantMax)
			}
			if conf.Delay.IsRange != tt.wantRange {
				t.Errorf("Delay.IsRange = %v, want %v", conf.Delay.IsRange, tt.wantRange)
			}
			if conf.Delay.HasDelay != tt.wantHas {
				t.Errorf("Delay.HasDelay = %v, want %v", conf.Delay.HasDelay, tt.wantHas)
			}
		})
	}
}

func TestRecursionDepthParsing(t *testing.T) {
	tests := []struct {
		name               string
		recursion          bool
		recursionDepth     int
		wantDepth          int
	}{
		{
			name:           "recursion enabled with depth 3",
			recursion:      true,
			recursionDepth: 3,
			wantDepth:      3,
		},
		{
			name:           "recursion disabled with depth 0 (default)",
			recursion:      false,
			recursionDepth: 0,
			wantDepth:      0,
		},
		{
			name:           "recursion enabled with depth 0 (unlimited)",
			recursion:      true,
			recursionDepth: 0,
			wantDepth:      0,
		},
		{
			name:           "recursion enabled with depth 5",
			recursion:      true,
			recursionDepth: 5,
			wantDepth:      5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configOptions := NewConfigOptions()
			configOptions.HTTP.Recursion = tt.recursion
			configOptions.HTTP.RecursionDepth = tt.recursionDepth
			configOptions.Input.Wordlists = []string{"test.txt"}
			configOptions.Input.InputMode = "clusterbomb"
			configOptions.HTTP.URL = "https://example.com/FUZZ"

			conf, err := ConfigFromOptions(configOptions, nil, nil)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if conf.RecursionDepth != tt.wantDepth {
				t.Errorf("RecursionDepth = %v, want %v", conf.RecursionDepth, tt.wantDepth)
			}
			if conf.Recursion != tt.recursion {
				t.Errorf("Recursion = %v, want %v", conf.Recursion, tt.recursion)
			}
		})
	}
}
