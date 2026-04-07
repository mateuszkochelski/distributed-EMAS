BINARY := bin/tsp-emas

NUM_AGENTS ?= 10000
NUM_ISLANDS ?= 100
INITIAL_ENERGY ?= 10
ENERGY_TRANSFER ?= 1
REPRODUCTION_THRESHOLD ?= 15
CHILD_ENERGY ?= 5
DEATH_THRESHOLD ?= 0
SAME_ISLAND_PROBABILITY ?= 0.999

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

SIM_FLAGS = \
	-num-agents=$(NUM_AGENTS) \
	-num-islands=$(NUM_ISLANDS) \
	-initial-energy=$(INITIAL_ENERGY) \
	-energy-transfer=$(ENERGY_TRANSFER) \
	-reproduction-threshold=$(REPRODUCTION_THRESHOLD) \
	-child-energy=$(CHILD_ENERGY) \
	-death-threshold=$(DEATH_THRESHOLD) \
	-same-island-probability=$(SAME_ISLAND_PROBABILITY)

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

.PHONY: build tsp-file tsp-random continuous clean

build:
	mkdir -p bin
	go build -o $(BINARY) .

tsp-file: build
	$(BINARY) -mode=tsp $(SIM_FLAGS) $(TSP_FLAGS) -tsp-source=file

tsp-random: build
	$(BINARY) -mode=tsp $(SIM_FLAGS) $(TSP_FLAGS) -tsp-source=random

continuous: build
	$(BINARY) -mode=continuous $(SIM_FLAGS) $(CONT_FLAGS)

clean:
	rm -rf bin
