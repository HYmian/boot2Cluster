package conf

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Register    []Register `yaml:"register"`
	BootCommand string     `yaml:"bootCommand"`
	LiveCommand string     `yaml:"liveCommand"`
}

type Register struct {
	Regular string `yaml:"regular"`
	File    string `yaml:"file"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg *Config
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}

type Boot struct {
	Templates   []*template.Template
	Targets     []string
	BootCommand string
	LiveCommand string
	Nodes       []Node

	mutex *sync.Mutex
}

func NewBoot(cfg *Config, clusterNum uint) *Boot {
	var tmpls []*template.Template
	var files []string

	for _, v := range cfg.Register {
		data, err := ioutil.ReadFile(v.File)
		if err != nil {
			log.Println(err.Error())
			return nil
		}

		t := template.Must(
			template.New(v.File).Parse(
				fmt.Sprintf(
					"%s{{ range .}}%s\n{{ end }}",
					data,
					v.Regular,
				),
			),
		)

		tmpls = append(tmpls, t)
		files = append(files, v.File)
	}

	return &Boot{
		Templates:   tmpls,
		Targets:     files,
		BootCommand: cfg.BootCommand,
		LiveCommand: cfg.LiveCommand,
		Nodes:       make([]Node, 0, clusterNum),
		mutex:       &sync.Mutex{},
	}
}

func (b *Boot) Entry() error {
	for i := 0; i < len(b.Targets); i++ {
		file, err := os.OpenFile(b.Targets[i], os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("open target file error: %s", err.Error())
		}
		defer file.Close()

		if err = b.Templates[i].Execute(file, b.Nodes); err != nil {
			return err
		}
	}
	return nil
}

func (b *Boot) ExecBootCommand() error {
	return Exec(b.BootCommand)
}

func (b *Boot) ExecLiveCommand() error {
	return Exec(b.LiveCommand)
}

func (b *Boot) GetNodes() Nodes {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.Nodes
}
