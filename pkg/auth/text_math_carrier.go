package auth

import (
	"fmt"
)

type textMapCarrier struct {
	headers map[string][]string
}

func (carrier *textMapCarrier) Get(key string) string {
	values := carrier.headers[key]
	if len(values) > 0 {
		fmt.Printf("Get key: %s\t%s\n", key, values[0])
		return values[0]
	} else {
		return ""
	}
}

func (carrier *textMapCarrier) Set(key string, value string) {
}

func (carrier *textMapCarrier) Keys() []string {
  var result []string
  for k, _ := range carrier.headers {
    result = append(result, k)
  }
  return result
}
