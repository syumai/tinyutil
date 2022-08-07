package test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	const want = "want body"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("want: %s, got: %s", http.MethodGet, req.Method)
		}
		w.Write([]byte(want))
	}))
	defer srv.Close()

	src := fmt.Sprintf(`
package main

import (
	"io"
	"os"

	"github.com/syumai/tinyutil/httputil"
)

const url = "%s"

func main() {
	resp, err := httputil.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
`, srv.URL)

	tmpdir, err := os.MkdirTemp("", "tinyutil")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	srcFile, err := os.CreateTemp(tmpdir, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(srcFile.Name())
	defer srcFile.Close()

	_, err = io.Copy(srcFile, strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	wasmPath := filepath.Join(tmpdir, "test.wasm")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	buildCmd := exec.CommandContext(ctx, "tinygo", "build", "-o", wasmPath, "-target", "wasm", srcFile.Name())
	err = buildCmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "deno", "run", "-A")
}
