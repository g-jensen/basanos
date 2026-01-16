package spec

import (
	"gopkg.in/yaml.v3"
)

type Hook struct {
	Run     string `yaml:"run"`
	Timeout string `yaml:"timeout"`
}

type RunBlock struct {
	Command string `yaml:"command"`
	Timeout string `yaml:"timeout"`
}

type Assertion struct {
	Command string `yaml:"command"`
	Timeout string `yaml:"timeout"`
}

type Scenario struct {
	ID         string            `yaml:"id"`
	Name       string            `yaml:"name"`
	Env        map[string]string `yaml:"env"`
	OnFailure  string            `yaml:"on_failure"`
	Before     *Hook             `yaml:"before"`
	After      *Hook             `yaml:"after"`
	BeforeEach *Hook             `yaml:"before_each"`
	AfterEach  *Hook             `yaml:"after_each"`
	Run        *RunBlock         `yaml:"run"`
	Assertions []Assertion       `yaml:"assertions"`
	Scenarios  []Scenario        `yaml:"scenarios"`
}

type Context struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Env         map[string]string `yaml:"env"`
	OnFailure   string            `yaml:"on_failure"`
	Before      *Hook             `yaml:"before"`
	After       *Hook             `yaml:"after"`
	BeforeEach  *Hook             `yaml:"before_each"`
	AfterEach   *Hook             `yaml:"after_each"`
	Scenarios   []Scenario        `yaml:"scenarios"`
}

func ParseContext(data []byte) (*Context, error) {
	var ctx Context
	err := yaml.Unmarshal(data, &ctx)
	if err != nil {
		return nil, err
	}
	return &ctx, nil
}
