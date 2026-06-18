package filter

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ffuf/ffuf/v2/pkg/ffuf"
)

const AllStatuses = 0

type StatusFilter struct {
	Value      []ffuf.ValueRange
	Regexps    []*regexp.Regexp
	regexpRaw  []string
}

func NewStatusFilter(value string) (ffuf.FilterProvider, error) {
	var intranges []ffuf.ValueRange
	var regexps []*regexp.Regexp
	var regexpRaw []string

	for _, sv := range strings.Split(value, ",") {
		sv = strings.TrimSpace(sv)
		if sv == "" {
			continue
		}
		// Check if this is a regex pattern (enclosed in /, e.g., /pattern/ or /pattern/i)
		if strings.HasPrefix(sv, "/") && strings.Count(sv, "/") >= 2 {
			// Extract pattern and flags
			lastSlash := strings.LastIndex(sv, "/")
			if lastSlash > 0 {
				pattern := sv[1:lastSlash]
				flags := sv[lastSlash+1:]
				
				// Handle flags
				if strings.Contains(flags, "i") {
					pattern = "(?i)" + pattern
				}
				
				re, err := regexp.Compile(pattern)
				if err != nil {
					return &StatusFilter{}, fmt.Errorf("Status filter or matcher (-fc / -mc): invalid regex pattern %s: %s", sv, err)
				}
				regexps = append(regexps, re)
				regexpRaw = append(regexpRaw, sv)
				continue
			}
		}
		
		if sv == "all" {
			intranges = append(intranges, ffuf.ValueRange{Min: AllStatuses, Max: AllStatuses})
		} else {
			vr, err := ffuf.ValueRangeFromString(sv)
			if err != nil {
				return &StatusFilter{}, fmt.Errorf("Status filter or matcher (-fc / -mc): invalid value %s", sv)
			}
			intranges = append(intranges, vr)
		}
	}
	return &StatusFilter{Value: intranges, Regexps: regexps, regexpRaw: regexpRaw}, nil
}

func (f *StatusFilter) MarshalJSON() ([]byte, error) {
	value := make([]string, 0)
	for _, v := range f.Value {
		if v.Min == 0 && v.Max == 0 {
			value = append(value, "all")
		} else {
			if v.Min == v.Max {
				value = append(value, strconv.FormatInt(v.Min, 10))
			} else {
				value = append(value, fmt.Sprintf("%d-%d", v.Min, v.Max))
			}
		}
	}
	value = append(value, f.regexpRaw...)
	return json.Marshal(&struct {
		Value string `json:"value"`
	}{
		Value: strings.Join(value, ","),
	})
}

func (f *StatusFilter) Filter(response *ffuf.Response) (bool, error) {
	statusMatched := false
	
	// Check status codes if any are defined
	if len(f.Value) > 0 {
		for _, iv := range f.Value {
			if iv.Min == AllStatuses && iv.Max == AllStatuses {
				// Handle the "all" case
				statusMatched = true
				break
			}
			if iv.Min <= response.StatusCode && response.StatusCode <= iv.Max {
				statusMatched = true
				break
			}
		}
	} else {
		// No status codes defined, only regex matching
		statusMatched = true
	}
	
	// If status didn't match, return early
	if !statusMatched {
		return false, nil
	}
	
	// Check regex patterns if any are defined
	if len(f.Regexps) > 0 {
		// Combine headers and body for matching
		matchheaders := ""
		for k, v := range response.Headers {
			for _, iv := range v {
				matchheaders += k + ": " + iv + "\r\n"
			}
		}
		matchdata := []byte(matchheaders)
		matchdata = append(matchdata, response.Data...)
		
		for _, re := range f.Regexps {
			if re.Match(matchdata) {
				return true, nil
			}
		}
		// None of the regex patterns matched
		return false, nil
	}
	
	// No regex patterns, return status match result
	return statusMatched, nil
}

func (f *StatusFilter) Repr() string {
	var strval []string
	for _, iv := range f.Value {
		if iv.Min == AllStatuses && iv.Max == AllStatuses {
			strval = append(strval, "all")
		} else if iv.Min == iv.Max {
			strval = append(strval, strconv.Itoa(int(iv.Min)))
		} else {
			strval = append(strval, strconv.Itoa(int(iv.Min))+"-"+strconv.Itoa(int(iv.Max)))
		}
	}
	strval = append(strval, f.regexpRaw...)
	return strings.Join(strval, ",")
}

func (f *StatusFilter) ReprVerbose() string {
	if len(f.Regexps) > 0 {
		return fmt.Sprintf("Response status: %s, Regexp: %s", f.Repr(), strings.Join(f.regexpRaw, ","))
	}
	return fmt.Sprintf("Response status: %s", f.Repr())
}
