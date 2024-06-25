package cfglib

import (
	"log"
	"strconv"
	"strings"
)

const Separator = "__"

func parseEnvVars(envProvider func() []string) (envMap map[string]interface{}) {
	envMap = make(map[string]interface{})
	for _, line := range envProvider() {
		keyValue := strings.SplitN(line, "=", 2)
		if len(keyValue) != 2 {
			continue
		}
		keys := strings.Split(keyValue[0], "__")
		lastKeyIndex := len(keys) - 1
		currentMap := envMap
		var currentSlice []interface{}
		for i, key := range keys {
			if strings.HasPrefix(key, "_") || strings.HasSuffix(key, "_") {
				continue
			}
			if i == lastKeyIndex {
				if index, err := strconv.Atoi(key); err == nil {
					currentSlice[index] = keyValue[1]
					currentMap[keys[i-1]] = currentSlice
				} else {
					currentMap[key] = keyValue[1]
				}
			} else {

				if currentSlice != nil {
					index, err := strconv.Atoi(keys[i])
					if err != nil {
						log.Default().Println("Error: ", err)
						continue
					}
					if len(currentSlice) <= index {
						newSlice := make([]interface{}, index+1)
						copy(newSlice, currentSlice)
						currentSlice = newSlice
					}
					var nextMap map[string]interface{}
					var nextSlice []interface{}
					if _, err := strconv.Atoi(keys[i+1]); err == nil {
						nextSlice = make([]interface{}, 0)
						currentSlice[index] = nextSlice
					} else {
						nextMap = make(map[string]interface{})
						currentSlice[index] = nextMap
					}
					currentMap[keys[i-1]] = currentSlice
					currentSlice = nextSlice
					currentMap = nextMap
					continue
				}

				if _, ok := currentMap[key]; !ok {
					if _, err := strconv.Atoi(keys[i+1]); err == nil {
						currentMap[key] = make([]interface{}, 0)
					} else {
						currentMap[key] = make(map[string]interface{})
					}
				}
				switch current := currentMap[key].(type) {
				case map[string]interface{}:
					currentMap = current
				case []interface{}:
					index, _ := strconv.Atoi(keys[i+1]) // next key must be an index
					if len(current) <= index {
						newSlice := make([]interface{}, index+1)
						copy(newSlice, current)
						current = newSlice
					}
					// is the next index a last key?
					if i+1 < lastKeyIndex {
						// check the key over next key
						if _, err := strconv.Atoi(keys[i+2]); err == nil {
							current[index] = make([]interface{}, 0)
						} else {
							current[index] = make(map[string]interface{})
						}
					}
					currentSlice = current
					currentMap[key] = current
				}
			}
		}
	}
	return
}

func flattenMap(envMap map[string]interface{}) map[string]interface{} {
	flatMap := make(map[string]interface{})
	flatten("", envMap, flatMap)
	return flatMap
}

func flatten(prefix string, m map[string]interface{}, flatMap map[string]interface{}) {
	for k, v := range m {
		key := prefix + strings.ToLower(k)
		flatMap[key] = v
		switch val := v.(type) {
		case map[string]interface{}:
			flatten(key+".", val, flatMap)
		case []interface{}:
			for i, item := range val {
				switch current := item.(type) {
				case map[string]interface{}:
					flatten(key+"."+strconv.Itoa(i)+".", current, flatMap)
				default:
					flatMap[key+"."+strconv.Itoa(i)] = item
				}
			}
		default:
			flatMap[key] = val
		}
	}
}
