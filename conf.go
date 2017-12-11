package main

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"fmt"
)

type Conf struct {
	Command string `yaml:"command"`
}

func GetConf() *Conf {

	yamlFile, err := ioutil.ReadFile("conf/conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err  #%v ", err)
	}

	var c Conf
	yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	return &c
}
