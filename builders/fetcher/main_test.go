package main

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func errCmp(e1, e2 error) bool {
	return errors.Is(e1, e2) || errors.Is(e2, e1)
}

func Test_extractBuilderRef(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		ref  string
		err  error
		tag  string
	}{
		{
			name: "valid tag",
			ref:  "refs/tags/v1.2.3",
			tag:  "v1.2.3",
		},
		{
			name: "invalid ref",
			ref:  "heads/tags/v1.2.3",
			err:  errorInvalidRef,
		},
		{
			name: "invalid tag prerelease",
			ref:  "refs/tags/v1.2.3-alpha",
			err:  errorInvalidRef,
		},
		{
			name: "invalid tag build",
			ref:  "refs/tags/v1.2.3+123",
			err:  errorInvalidRef,
		},
		{
			name: "invalid tag metadata",
			ref:  "refs/tags/v1.2.3-aplha+123",
			err:  errorInvalidRef,
		},
	}

	for _, tt := range tests {
		tt := tt // Re-initializing variable so it is not changed while executing the closure below
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag, err := extractBuilderRef(tt.ref)
			if !errCmp(err, tt.err) {
				t.Errorf(cmp.Diff(err, tt.err, cmpopts.EquateErrors()))
			}

			if err != nil {
				return
			}

			if tag != tt.tag {
				t.Errorf(cmp.Diff(tag, tt.tag))
			}
		})
	}
}
