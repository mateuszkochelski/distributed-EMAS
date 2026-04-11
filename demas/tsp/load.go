package tsp

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadFromFile(path string) (Problem, error) {
	file, err := os.Open(path)
	if err != nil {
		return Problem{}, fmt.Errorf("open %q: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var (
		inNodeSection  bool
		dimension      int
		edgeWeightType string
		problemType    string
		cities         []City
		lineNo         int
	)

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if line == "EOF" {
			break
		}

		if !inNodeSection {
			if line == "NODE_COORD_SECTION" {
				inNodeSection = true

				if problemType != "" && problemType != "TSP" {
					return Problem{}, fmt.Errorf("line %d: unsupported TYPE %q", lineNo, problemType)
				}

				if edgeWeightType != "EUC_2D" {
					return Problem{}, fmt.Errorf("line %d: unsupported EDGE_WEIGHT_TYPE %q", lineNo, edgeWeightType)
				}

				if dimension <= 0 {
					return Problem{}, fmt.Errorf("line %d: invalid DIMENSION %d", lineNo, dimension)
				}

				cities = make([]City, 0, dimension)
				continue
			}

			key, value, ok := parseHeaderLine(line)
			if !ok {
				continue
			}

			switch key {
			case "TYPE":
				problemType = value
			case "DIMENSION":
				parsed, err := strconv.Atoi(value)
				if err != nil {
					return Problem{}, fmt.Errorf("line %d: parse DIMENSION: %w", lineNo, err)
				}
				dimension = parsed
			case "EDGE_WEIGHT_TYPE":
				edgeWeightType = value
			}

			continue
		}

		city, err := parseNodeCoordLine(line, lineNo)
		if err != nil {
			return Problem{}, err
		}

		cities = append(cities, city)

		if len(cities) == dimension {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return Problem{}, fmt.Errorf("scan %q: %w", path, err)
	}

	if !inNodeSection {
		return Problem{}, fmt.Errorf("missing NODE_COORD_SECTION")
	}

	if len(cities) != dimension {
		return Problem{}, fmt.Errorf("expected %d cities, got %d", dimension, len(cities))
	}

	return NewProblem(cities), nil
}

func parseHeaderLine(line string) (key, value string, ok bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	key = strings.TrimSpace(parts[0])
	value = strings.TrimSpace(parts[1])

	if key == "" || value == "" {
		return "", "", false
	}

	return key, value, true
}

func parseNodeCoordLine(line string, lineNo int) (City, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return City{}, fmt.Errorf("line %d: expected 'id x y', got %q", lineNo, line)
	}

	_, err := strconv.Atoi(fields[0])
	if err != nil {
		return City{}, fmt.Errorf("line %d: parse node id: %w", lineNo, err)
	}

	x, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return City{}, fmt.Errorf("line %d: parse x: %w", lineNo, err)
	}

	y, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return City{}, fmt.Errorf("line %d: parse y: %w", lineNo, err)
	}

	return City{
		X: x,
		Y: y,
	}, nil
}
