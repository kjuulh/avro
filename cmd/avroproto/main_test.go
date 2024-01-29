package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "Update golden files")

func TestAvroProto_RequiredFlags(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantExitCode int
	}{}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := realMain(test.args, io.Discard, io.Discard)

			assert.Equal(t, test.wantExitCode, got)
		})
	}
}

func TestAvroProto_GeneratesSchemaStdout(t *testing.T) {
	var buf bytes.Buffer

	args := []string{"avroproto", "-p", "testpkg", "testdata/schema.avsc"}
	gotCode := realMain(args, &buf, io.Discard)
	require.Equal(t, 0, gotCode)

	want, err := os.ReadFile("testdata/golden.proto")
	require.NoError(t, err)
	assert.Equal(t, want, buf.Bytes())
}

func TestAvroProto_GeneratesSchema(t *testing.T) {
	path, err := os.MkdirTemp("./", "avroproto")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(path) })

	file := filepath.Join(path, "test.go")
	args := []string{"avroproto", "-p", "testpkg", "-o", file, "testdata/schema.avsc"}
	gotCode := realMain(args, io.Discard, io.Discard)
	require.Equal(t, 0, gotCode)

	got, err := os.ReadFile(file)
	require.NoError(t, err)

	if *update {
		err = os.WriteFile("testdata/golden.proto", got, 0600)
		require.NoError(t, err)
	}

	want, err := os.ReadFile("testdata/golden.proto")
	require.NoError(t, err)
	assert.Equal(t, want, got)

	assert.Equal(t, string(want), string(got))
}
