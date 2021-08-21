APP_NAME := libipfs-example

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

default: run

clean:
	rm -rf .libipfs* libipfs* test-download.jpg

build:
	cp -rf ../bin/libipfs-linux* .
	g++ -pthread -o libipfs-example -ldl example.cpp libipfs-linux.a

test: 
	./libipfs-example

run: clean \
	build \
	test 