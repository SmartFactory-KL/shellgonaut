package id_short_path

import (
	"strings"
	"testing"
)

func TestGatherStepsFromIdShortPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		idShortPath string
		want        []*IDShortPathStep
		wantErr     string
	}{
		{
			name:        "rejects empty path",
			idShortPath: "   ",
			wantErr:     "idShortPath cannot be empty",
		},
		{
			name:        "parses plain dotted path",
			idShortPath: "submodel.collection.property",
			want: []*IDShortPathStep{
				{StepType: StepTypeIDShort, NextIdShort: "submodel"},
				{StepType: StepTypeIDShort, NextIdShort: "collection"},
				{StepType: StepTypeIDShort, NextIdShort: "property"},
			},
		},
		{
			name:        "trims surrounding whitespace",
			idShortPath: "  submodel.collection  ",
			want: []*IDShortPathStep{
				{StepType: StepTypeIDShort, NextIdShort: "submodel"},
				{StepType: StepTypeIDShort, NextIdShort: "collection"},
			},
		},
		{
			name:        "rejects empty path parts",
			idShortPath: "submodel..property",
			wantErr:     "idShortPath cannot have empty parts",
		},
		{
			name:        "rejects missing closing bracket",
			idShortPath: "list[2",
			wantErr:     "missing ]",
		},
		{
			name:        "rejects empty index",
			idShortPath: "list[].property",
			wantErr:     "cannot be empty",
		},
		{
			name:        "parses index steps",
			idShortPath: "list[2].property",
			want: []*IDShortPathStep{
				{StepType: StepTypeIDShort, NextIdShort: "list"},
				{StepType: StepTypeIndex, NextIndex: 2},
				{StepType: StepTypeIDShort, NextIdShort: "property"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GatherStepsFromIdShortPath(tt.idShortPath)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got %d steps, want %d", len(got), len(tt.want))
			}

			for i := range tt.want {
				if got[i] == nil {
					t.Fatalf("step %d is nil", i)
				}
				if got[i].StepType != tt.want[i].StepType {
					t.Fatalf("step %d type mismatch: got %v want %v", i, got[i].StepType, tt.want[i].StepType)
				}
				if got[i].NextIdShort != tt.want[i].NextIdShort {
					t.Fatalf("step %d idShort mismatch: got %q want %q", i, got[i].NextIdShort, tt.want[i].NextIdShort)
				}
				if got[i].NextIndex != tt.want[i].NextIndex {
					t.Fatalf("step %d index mismatch: got %d want %d", i, got[i].NextIndex, tt.want[i].NextIndex)
				}
			}
		})
	}
}
