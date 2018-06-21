package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Checker struct {
	Interval  int64           `json:"interval" yaml:"interval"`
	Debug     bool            `json:"debug" yaml:"debug"`
	Sites     map[string]Site `json:"sites" yaml:"sites"`
	Notifiers []Notifier      `json:"-" yaml:"-"`
}

func NewChecker(config string, notifiers []Notifier) (Checker, error) {
	c := Checker{}

	data, err := ioutil.ReadFile(config)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal([]byte(data), &c)
	c.Notifiers = notifiers
	return c, err
}

func (c Checker) Run() {
	log.Println("Started...")
	for {
		for i, s := range c.Sites {
			if c.Debug {
				log.Printf("Scanning '%s'\n", i)
			}
			err := s.Check(c.Notifiers)
			if err != nil {
				log.Println(err)
			}
			if c.Debug {
				log.Printf("%s => %s\n", c.Sites[i].Hash, s.Hash)
			}
			c.Sites[i] = s
		}

		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

type Site struct {
	URL       string `json:"url" yaml:"url"`
	Template  string `json:"template" yaml:"template"`
	Recipient string `json:"recipient" yaml:"recipient"`
	Hash      string `json:"-" yaml:"-"`
}

func (s *Site) Check(n []Notifier) error {
	oldHash := s.Hash
	response, err := http.Get(s.URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	s.Hash = GetMD5Hash(string(contents))
	if oldHash != "" && oldHash != s.Hash {
		for _, notifier := range n {
			err = notifier.Notify(s.Recipient, s.Template)
			if err != nil {
				log.Printf("Error while notifing: %s\n", err.Error())
			}
		}
	}
	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
