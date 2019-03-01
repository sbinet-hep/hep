package main

import (
	"flag"
	"log"

	"go-hep.org/x/hep/groot"
)

func main() {
	flag.Parse()
	f, err := groot.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("evt")
	if err != nil {
		log.Fatal(err)
	}

	evt := o.(*Event)
	log.Printf("evt: %#v", *evt)
}
