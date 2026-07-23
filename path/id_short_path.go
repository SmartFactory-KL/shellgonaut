package id_short_path

import (
	"fmt"
	"strconv"
	"strings"
)

type IDShortPathStepType int

const (
	StepTypeIDShort = iota
	StepTypeIndex
)

type IDShortPathStep struct {
	StepType IDShortPathStepType

	NextIdShort string
	NextIndex   int
}

// GatherStepsFromIdShortPath creates a list of path stepts to go based on idShortPath
func GatherStepsFromIdShortPath(idShortPath string) ([]*IDShortPathStep, error) {
	idShortPath = strings.TrimSpace(idShortPath)

	if len(idShortPath) == 0 {
		return nil, fmt.Errorf("idShortPath cannot be empty")
	}

	// first split by . to get basic elements
	// then check for [] pairs to get index steps
	stepList := make([]*IDShortPathStep, 0, 20)
	pathByPoints := strings.SplitSeq(idShortPath, ".")

	for pathItem := range pathByPoints {
		if len(strings.TrimSpace(pathItem)) == 0 {
			return nil, fmt.Errorf("idShortPath cannot have empty parts: %s", idShortPath)
		}

		pathByIdx := strings.Split(pathItem, "[")
		for idx2, rawIdxItem := range pathByIdx {
			trimmedIdxItem := strings.TrimSpace(rawIdxItem)

			if idx2 == 0 {
				// the first item is always the "idShortPath"-Part of it.
				// for [123] that would simply be empty.
				if len(trimmedIdxItem) > 0 {
					stepList = append(stepList, &IDShortPathStep{
						StepType:    StepTypeIDShort,
						NextIdShort: trimmedIdxItem,
					})
				}
			} else {
				// nextIndex must be a valid int within []
				if !strings.HasSuffix(trimmedIdxItem, "]") {
					return nil, fmt.Errorf("invalid index on idShortPath element %d: %s - missing ]", idx2, trimmedIdxItem)
				}

				idxItem := strings.TrimSuffix(trimmedIdxItem, "]")

				if len(idxItem) == 0 {
					return nil, fmt.Errorf("invalid index on idShortPath element %d: %s - cannot be empty", idx2, idxItem)
				}

				nextIndex, err := strconv.ParseInt(idxItem, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid index on idShortPath element %d: %s - not a valid number", idx2, idxItem)
				}
				if nextIndex < 0 {
					return nil, fmt.Errorf("invalid index on idShortPath element %d: %s cannot be negative", idx2, idxItem)
				}

				stepList = append(stepList, &IDShortPathStep{
					StepType:  StepTypeIndex,
					NextIndex: int(nextIndex),
				})
			}

		}
	}

	return stepList, nil
}
