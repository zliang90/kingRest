package conf

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"text/template"

	"github.com/go-playground/validator/v10"
	json "github.com/json-iterator/go"
	"github.com/zliang90/kingRest/internal/restful/errors"
	"github.com/zliang90/kingRest/pkg/log"
	"gopkg.in/yaml.v3"
)

var (
	config   *Config
	validate *validator.Validate

	lock = new(sync.RWMutex)

	funcs template.FuncMap

	leftDelim  = "{{"
	rightDelim = "}}"
)

type Config struct {
	// App operation mode: dev/test/prod
	Env string `validate:"oneof=dev test prod" yaml:"Env"`

	// restful api listen address, :8086
	WebServer `validate:"required" yaml:"WebServer"`

	// Log level: fatal, error, warning, info, debug
	LogLevel string `validate:"oneof=debug info warning error fatal" yaml:"LogLevel"`

	// Api error file
	ApiErrorFile string `validate:"required" yaml:"ApiErrorFile"`

	// database source
	DataSources map[string]DataSource `validate:"required" yaml:"DataSources"`
}

type WebServer struct {
	Addr              string `yaml:"Addr" validate:"required"`
	MaxHeaderBytes    int    `yaml:"MaxHeaderBytes" validate:"required"`
	ReadTimeout       int    `yaml:"ReadTimeout" validate:"required"`
	ReadHeaderTimeout int    `yaml:"ReadHeaderTimeout" validate:"required"`
	WriteTimeout      int    `yaml:"WriteTimeout" validate:"required"`
	IdleTimeout       int    `yaml:"IdleTimeout" validate:"required"`
}

type DataSource struct {
	Addr     string `yaml:"Addr"`
	IdleConn int    `yaml:"Idle"`
	MaxConn  int    `yaml:"Max"`
	Debug    bool   `yaml:"Debug"`
}

func init() {
	flag.StringVar(&leftDelim, "leftDelim", "[[", "config template left delim")
	flag.StringVar(&rightDelim, "rightDelim", "]]", "config template right delim")

	if validate == nil {
		validate = validator.New()
		// validate.SetTagName("validate")
	}

	funcs = template.FuncMap{
		"slice": func(args ...interface{}) []interface{} {
			return args
		},
	}
}

// Validate validate config file
func (c Config) Validate() (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			switch e := _err.(type) {
			case error:
				err = e
			case string:
				err = fmt.Errorf("config validate error, %s", e)
			default:
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	if err = validate.Struct(&c); err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	for _, e := range err.(validator.ValidationErrors) {
		if e != nil {
			if _, ok := e.(error); ok {
				err = e
				break
			}
		}
	}
	return err
}

func (c Config) String() string {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonBytes)
}

// LoadConfig loading config file to struct object
func LoadConfig(cfgPath string) error {
	text, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	var c Config

	// render config template
	buf := new(bytes.Buffer)
	t := template.Must(
		template.New("config").
			Delims(leftDelim, rightDelim).
			Funcs(funcs).
			Parse(string(text)))
	if err = t.Execute(buf, nil); err != nil {
		return err
	}
	if err = yaml.Unmarshal(buf.Bytes(), &c); err != nil {
		return err
	}
	if err = c.Validate(); err != nil {
		return err
	}

	// log level to uppercase
	c.LogLevel = strings.ToUpper(c.LogLevel)

	// error message
	if err = errors.LoadMessages(c.ApiErrorFile); err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()

	errors.SetEnv(c.Env)
	config = &c
	return nil
}

func GetConfig() *Config {
	lock.RLock()
	defer lock.RUnlock()

	return config
}
