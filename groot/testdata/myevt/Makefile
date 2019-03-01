.PHONY: all myevt

all: run myevt

myevt: myevt.h o.root
	/bin/rm -f evt.go
	root-gen-type -p main -o evt.go -t Event,P3 ./o.root
	go build -o ./myevt

o.root: run.C myevt.h
	root -b -q ./run.C
	root-ls -sinfos ./o.root

run: myevt
	./myevt ./o.root

