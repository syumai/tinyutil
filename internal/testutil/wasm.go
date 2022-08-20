package testutil

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func RunWasm(t *testing.T, src string) (out io.Reader) {
	tmpdir, err := os.MkdirTemp("", "tinyutil")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	srcFile, err := os.CreateTemp(tmpdir, "*.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())

	_, err = io.Copy(srcFile, strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	err = srcFile.Close()
	if err != nil {
		t.Fatal(err)
	}

	wasmPath := filepath.Join(tmpdir, "test.wasm")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	buildCmd := exec.CommandContext(ctx, "tinygo", "build", "-o", wasmPath, "-target", "wasm", srcFile.Name())
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("err: %v\noutput: %s", err, string(output))
	}

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "deno", "run", "-A", "../testutil/lib/run_test.js", wasmPath)

	var outBuf, errBuf bytes.Buffer
	runCmd.Stdout = &outBuf
	runCmd.Stderr = &errBuf
	err = runCmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	if errBuf.Len() != 0 {
		t.Fatal(errBuf.String())
	}
	return &outBuf
}
