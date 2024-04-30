package content

import (
	"net/http"
	"strconv"
	"strings"
)

type AcceptRange struct {
	Type       string
	Subtype    string
	Weight     float64
	Parameters map[string]string
	raw        string
}

func NegotiateContentType(r *http.Request, offers ...ContentType) string {
	if len(offers) == 0 {
		offers = append(offers, ContentTypeText)
	}

	accepts := AcceptMediaTypes(r)
	offerRanges := make([]AcceptRange, 0, len(offers))

	for _, offer := range offers {
		offerRanges = append(offerRanges, ParseAcceptRange(string(offer)))
	}

	return negotiateContentType(accepts, offerRanges)
}

func AcceptMediaTypes(r *http.Request) []AcceptRange {
	var result []AcceptRange

	for _, value := range r.Header["Accept"] {
		result = append(result, ParseAcceptRanges(value)...)
	}

	return result
}

func ParseAcceptRanges(accepts string) []AcceptRange {
	var result []AcceptRange
	remaining := accepts

	for {
		var accept string

		accept, remaining = extractFieldAndSkipToken(remaining, ',')
		result = append(result, ParseAcceptRange(accept))

		if len(remaining) == 0 {
			break
		}
	}

	return result
}

func ParseAcceptRange(accept string) AcceptRange {
	typeAndSubtype, rawParams := extractFieldAndSkipToken(accept, ';')
	t, st := extractFieldAndSkipToken(typeAndSubtype, '/')
	params := extractParams(rawParams)
	weight := extractWeight(params)

	return AcceptRange{
		Type:       t,
		Subtype:    st,
		Parameters: params,
		Weight:     weight,
		raw:        accept,
	}
}

func extractFieldAndSkipToken(accept string, sep rune) (string, string) {
	field, rest := extractField(accept, sep)
	if len(rest) > 0 {
		rest = rest[1:]
	}

	return field, rest
}

func extractField(entry string, sep rune) (string, string) {
	field := entry
	rest := ""

	for index, char := range entry {
		if char == sep {
			field = strings.TrimSpace(entry[:index])
			rest = strings.TrimSpace(entry[index:])
		}
	}

	return field, rest
}

func extractParams(rawParams string) map[string]string {
	params := make(map[string]string)
	rest := rawParams

	for {
		var param string

		param, rest = extractFieldAndSkipToken(rest, ';')
		if len(param) > 0 {
			key, value := extractFieldAndSkipToken(param, '=')
			params[key] = value
		}
		if len(rest) == 0 {
			break
		}
	}

	return params
}

func extractWeight(params map[string]string) float64 {
	if weight, ok := params["q"]; ok {
		result, err := strconv.ParseFloat(weight, 64)
		if err == nil {
			return result
		}
	}

	return 1
}

func compareParams(a, b map[string]string) int {
	var count int

	for key, value := range a {
		if v, ok := b[key]; ok && value == v {
			count++
		}
	}

	return count
}

func negotiateContentType(accepts []AcceptRange, offers []AcceptRange) string {
	best := offers[0].raw
	bestWeight := offers[0].Weight
	bestParams := 0

	for _, offer := range offers {
		for _, accept := range accepts {
			booster := float64(0)

			if accept.Type != "*" {
				booster++

				if accept.Subtype != "*" {
					booster++
				}
			}

			switch {
			case bestWeight > (accept.Weight + booster):
				continue

			case accept.Type == "*" && accept.Subtype == "*" && ((accept.Weight + booster) > bestWeight):
				best = offer.raw
				bestWeight = accept.Weight + booster

			case accept.Subtype == "*" && accept.Type == offer.Type && ((accept.Weight + booster) > bestWeight):
				best = offer.raw
				bestWeight = accept.Weight + booster

			case accept.Type == offer.Type && accept.Subtype == offer.Subtype:
				paramCount := compareParams(accept.Parameters, offer.Parameters)
				if paramCount >= bestParams {
					best = offer.raw
					bestWeight = accept.Weight + booster
					bestParams = paramCount
				}
			}
		}
	}

	return best
}
