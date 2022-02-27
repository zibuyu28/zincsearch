package analysis

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/char"

	"github.com/prabhatsharma/zinc/pkg/errors"
	pluginanalysis "github.com/prabhatsharma/zinc/pkg/plugin/analysis"
	zincchar "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis/char"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func RequestCharFilter(data map[string]interface{}) (map[string]analysis.CharFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make(map[string]analysis.CharFilter)
	for name, options := range data {
		typ, err := zutils.GetStringFromMap(options, "type")
		if err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[char_filter] %s option [%s] should be exists", name, "type"))
		}
		filter, err := RequestCharFilterSingle(typ, options)
		if err != nil {
			return nil, err
		}
		filters[name] = filter
	}

	return filters, nil
}

func RequestCharFilterSlice(data []interface{}) ([]analysis.CharFilter, error) {
	if data == nil {
		return nil, nil
	}

	filters := make([]analysis.CharFilter, 0, len(data))
	for _, options := range data {
		var err error
		var filter analysis.CharFilter
		switch v := options.(type) {
		case string:
			filter, err = RequestCharFilterSingle(v, nil)
		case map[string]interface{}:
			var typ string
			typ, err = zutils.GetStringFromMap(options, "type")
			if err != nil {
				return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] option [type] should be exists")
			}
			filter, err = RequestCharFilterSingle(typ, options)
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, "[char_filter] option should be string or object")
		}
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}

	return filters, nil
}

func RequestCharFilterSingle(name string, options interface{}) (analysis.CharFilter, error) {
	name = strings.ToLower(name)
	switch name {
	case "ascii_folding", "asciifolding":
		return char.NewASCIIFoldingFilter(), nil
	case "html", "html_strip":
		return char.NewHTMLCharFilter(), nil
	case "zero_width_non_joiner":
		return char.NewZeroWidthNonJoinerCharFilter(), nil
	case "regexp", "pattern", "pattern_replace":
		return zincchar.NewRegexpCharFilter(options)
	case "mapping":
		return zincchar.NewMappingCharFilter(options)
	default:
		if v, ok := pluginanalysis.GetCharFilter(name); ok {
			return v, nil
		}
		return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[char_filter] unkown character filter [%s]", name))
	}
}
