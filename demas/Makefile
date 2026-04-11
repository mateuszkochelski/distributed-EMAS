BINARY := tsp-emas

RUN_SECONDS ?= 30

NUM_AGENTS ?= 10000
NUM_ISLANDS ?= 100
INITIAL_ENERGY ?= 10
ENERGY_TRANSFER ?= 1
REPRODUCTION_THRESHOLD ?= 15
CHILD_ENERGY ?= 5
DEATH_THRESHOLD ?= 0
SAME_ISLAND_PROBABILITY ?= 0.999
LOG_MEETINGS ?= false

TSP_SOURCE ?= file
TSP_INPUT_FILE ?= ../testdata/tsplib/a280.tsp
TSP_NUM_CITIES ?= 100
TSP_MIN_COORD ?= -100
TSP_MAX_COORD ?= 100

CONT_DIM ?= 30
CONT_LOWER ?= -5.12
CONT_UPPER ?= 5.12
CONT_SIGMA ?= 0.2
CONT_OBJECTIVE ?= rastrigin

LOCAL_BIN := bin/$(BINARY)
X86_BIN := bin/x86/$(BINARY)
ARM_BIN := bin/arm/$(BINARY)

COMMON_FLAGS = \
	-run-seconds=$(RUN_SECONDS) \
	-num-agents=$(NUM_AGENTS) \
	-num-islands=$(NUM_ISLANDS) \
	-initial-energy=$(INITIAL_ENERGY) \
	-energy-transfer=$(ENERGY_TRANSFER) \
	-reproduction-threshold=$(REPRODUCTION_THRESHOLD) \
	-child-energy=$(CHILD_ENERGY) \
	-death-threshold=$(DEATH_THRESHOLD) \
	-same-island-probability=$(SAME_ISLAND_PROBABILITY) \
	-log-meetings=$(LOG_MEETINGS)

TSP_FLAGS = \
	-tsp-source=$(TSP_SOURCE) \
	-tsp-input-file=$(TSP_INPUT_FILE) \
	-tsp-num-cities=$(TSP_NUM_CITIES) \
	-tsp-min-coord=$(TSP_MIN_COORD) \
	-tsp-max-coord=$(TSP_MAX_COORD)

CONT_FLAGS = \
	-cont-dim=$(CONT_DIM) \
	-cont-lower=$(CONT_LOWER) \
	-cont-upper=$(CONT_UPPER) \
	-cont-sigma=$(CONT_SIGMA) \
	-cont-objective=$(CONT_OBJECTIVE)

.PHONY: build build-local build-x86 build-arm build-all \
	tsp-file tsp-random continuous \
	run-x86-tsp-file run-x86-tsp-random run-x86-continuous \
	run-arm-tsp-file run-arm-tsp-random run-arm-continuous \
	clean

build: build-local

build-local:
	mkdir -p bin
	go build -o $(LOCAL_BIN) .

build-x86:
	mkdir -p bin/x86
	GOOS=linux GOARCH=amd64 go build -o $(X86_BIN) .

build-arm:
	mkdir -p bin/arm
	GOOS=linux GOARCH=arm64 go build -o $(ARM_BIN) .

build-all: build-x86 build-arm

tsp-file: build-local
	$(LOCAL_BIN) -mode=tsp $(COMMON_FLAGS) $(TSP_FLAGS) -tsp-source=file

tsp-random: build-local
	$(LOCAL_BIN) -mode=tsp $(COMMON_FLAGS) $(TSP_FLAGS) -tsp-source=random

continuous: build-local
	$(LOCAL_BIN) -mode=continuous $(COMMON_FLAGS) $(CONT_FLAGS)

clean:
	rm -rf bin
