package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type ValueMode string

const (
	ValueModeAuto = "auto"
	ValueModeHttp = "http"
)

func ParseConfig(file string) (*GeneratorConfig, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	config := &GeneratorConfig{}
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	for _, c := range config.Counters {
		if c.Number == 0 {
			c.Number = 1
		}
		if c.ValueMode == "" {
			c.ValueMode = ValueModeAuto
		}
	}
	for _, g := range config.Gauges {
		if g.Number == 0 {
			g.Number = 1
		}
		if g.Range.Upper == 0 {
			g.Range.Upper = 100
		}
		if g.ValueMode == "" {
			g.ValueMode = ValueModeAuto
		}

	}
	return config, nil
}

type GeneratorConfig struct {
	Counters []*CounterConfig `yaml: counters,omitempty`
	Gauges   []*GaugeConfig   `yaml: gauges,omitempty`
}

type CounterConfig struct {
	Prefix    string         `yaml: prefix`
	Number    int            `yaml: number,omitempty`
	Labels    []LabelSetting `yaml: labels,omitempty`
	ValueMode ValueMode      `yaml: valuemode,omitempty`
}

type GaugeConfig struct {
	Prefix    string         `yaml: prefix`
	Number    int            `yaml: number,omitempty`
	Labels    []LabelSetting `yaml: labels,omitempty`
	Range     GaugeRange     `yaml: range,omitempty`
	ValueMode ValueMode      `yaml: valuemode,omitempty`
}
type LabelSetting struct {
	Name     string   `yaml: name`
	ValueSet []string `yaml: valueset`
}
type GaugeRange struct {
	Upper int64 `yaml: upper`
	Lower int64 `yaml: lower`
}
