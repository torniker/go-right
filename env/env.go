package env

import (
	"fmt"
	"os"
	"strings"
)

// Key names for environment variable
const (
	EnvKey string = "ENV"
)

// Define environment constants
const (
	Local      string = "local"
	Testing    string = "testing"
	Staging    string = "staging"
	Production string = "production"
)

// server params
const (
	PortKey       = "PORT"
	BasePathKey   = "BASE_PATH"
	MainGoPathKey = "MAIN_GO_PATH"
)

var env *Environment

// Environment object
type Environment struct {
	service string
	vars    map[string]string
}

// New creates new environment for given service
func New(service string) {
	env = &Environment{
		service: service,
	}
	if strings.HasSuffix(os.Args[0], ".test") {
		os.Setenv(EnvKey, fmt.Sprint(Testing))
	}
	// set env vars into the map
	env.vars = make(map[string]string)
	for _, vars := range os.Environ() {
		pair := strings.SplitN(vars, "=", 2)
		env.Set(pair[0], pair[1])
	}
}

func Service() string {
	return env.Service()
}

func (e *Environment) Service() string {
	return e.service
}

func Get(key string) string {
	return env.Get(key)
}

func (e *Environment) Get(key string) string {
	if val, ok := e.vars[key]; ok {
		return val
	}
	return os.Getenv(key)
}

func SetForLocal(key, val string) {
	if IsLocal() {
		env.Set(key, val)
	}
}

func IsLocal() bool {
	return env.String() == Local
}

func Set(key, val string) {
	env.Set(key, val)
}

func (e *Environment) Set(key, val string) {
	e.vars[key] = val
}

func (e *Environment) String() string {
	if e == nil {
		panic("environment is nil. please use env.New() func before logging anything")
	}
	switch e.Get(EnvKey) {
	case Staging:
		return Staging
	case Production:
		return Production
	}
	return Local
}

func String() string {
	return env.String()
}
