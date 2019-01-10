package json

import (
	"encoding/json"
	"fmt"
)

func Filter(obj map[string]interface{}, filter map[string]interface{}) (map[string]interface{}, error) {
	return applyFilter(obj, filter, "")
}

func applyFilter(obj map[string]interface{}, filter map[string]interface{}, path string) (map[string]interface{}, error) {
	if obj == nil {
		return nil, nil
	}

	out := make(map[string]interface{}, len(obj))
	for k, subFilter := range filter {
		v, ok := obj[k]
		if !ok {
			continue // ignore missing key or throw error
		}
		sf, ok := subFilter.(map[string]interface{})
		if !ok {
			out[k] = v // just keep it as is
		} else {
			// apply sub filter
			// if v is an map, apply sub filter directly
			// else if v is an array of objects, apply to sub filter to individual elements
			// else, throw an error (filter is trying to apply to non objects)

			switch u := v.(type) {
			case map[string]interface{}:
				subOut, err := applyFilter(u, sf, path+k+".")
				if err != nil {
					return nil, err
				}
				out[k] = subOut
			case []interface{}:
				for i := range u {
					entry, ok := u[i].(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("can't apply filter %s on %s%s[%d]: %v", toJson(sf), path, k, i, u[i]) // report the path to v
					}
					subOut, err := applyFilter(entry, sf, fmt.Sprintf("%s%s[%d].", path, k, i))
					if err != nil {
						return nil, err
					}
					u[i] = subOut
				}
				out[k] = u
			default:
				return nil, fmt.Errorf("can't apply filter %s on %s%s: %v", toJson(sf), path, k, v)
			}
		}
	}
	return out, nil
}

func toJson(v interface{}) string {
	str, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%q", v)
	}
	return string(str)
}
