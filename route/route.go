package route

import (
	"gopkg.in/yaml.v2"
)

type Route struct {
	Path string	`yaml:"path"`
	Method string	`yaml:"method"`
	Params []string	`yaml:"params"`
	Cache int	`yaml:"cache"`
}

func New(data []byte) (map[string]Route, error) {
	var routeMap map[string]Route
	err := yaml.Unmarshal(data, &routeMap)
	if err != nil {
		return nil, err
	}
	return routeMap, nil
}