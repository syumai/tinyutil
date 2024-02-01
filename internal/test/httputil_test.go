package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/syumai/tinyutil/internal/testutil"
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
	"fmt"
	"io"

	"github.com/syumai/tinyutil/httputil"
)

func main() {
	resp, err := httputil.Get(%q)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
`, srv.URL)

	out := testutil.RunWasm(t, src)
	b, err := io.ReadAll(out)
	if err != nil {
		t.Fatal(err)
	}

	gotBody := strings.TrimSpace(string(b))
	if want != gotBody {
		t.Fatalf("want: %s, got: %s", want, gotBody)
	}
}

func TestPost(t *testing.T) {
	const (
		wantResBody     = "want res body"
		wantReqBody     = "want req body"
		wantContentType = "text/plain"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			t.Fatalf("want: %s, got: %s", http.MethodPost, req.Method)
		}
		b, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		gotReqBody := string(b)
		if gotReqBody != wantReqBody {
			t.Fatalf("want: %s, got: %s", wantReqBody, gotReqBody)
		}
		gotContentType := req.Header.Get("Content-Type")
		if gotContentType != wantContentType {
			t.Fatalf("want: %s, got: %s", wantContentType, gotContentType)
		}
		w.Write([]byte(wantResBody))
	}))
	defer srv.Close()

	src := fmt.Sprintf(`
package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/syumai/tinyutil/httputil"
)

func main() {
	resp, err := httputil.Post(%q, %q, strings.NewReader(%q))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
`, srv.URL, wantContentType, wantReqBody)

	out := testutil.RunWasm(t, src)
	b, err := io.ReadAll(out)
	if err != nil {
		t.Fatal(err)
	}

	gotBody := strings.TrimSpace(string(b))
	if wantResBody != gotBody {
		t.Fatalf("want: %s, got: %s", wantResBody, gotBody)
	}
}

func TestPostForm(t *testing.T) {
	const (
		wantResBody     = "want res body"
		wantReqBody     = "bar=baz&foo=quux"
		wantContentType = "application/x-www-form-urlencoded"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			t.Fatalf("want: %s, got: %s", http.MethodPost, req.Method)
		}
		b, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		gotReqBody := string(b)
		if gotReqBody != wantReqBody {
			t.Fatalf("want: %s, got: %s", wantReqBody, gotReqBody)
		}
		gotContentType := req.Header.Get("Content-Type")
		if gotContentType != wantContentType {
			t.Fatalf("want: %s, got: %s", wantContentType, gotContentType)
		}
		w.Write([]byte(wantResBody))
	}))
	defer srv.Close()

	src := fmt.Sprintf(`
package main

import (
	"fmt"
	"io"
	"net/url"

	"github.com/syumai/tinyutil/httputil"
)

func main() {
	data := url.Values{
		"bar": []string{"baz"},
		"foo": []string{"quux"},
	}
	resp, err := httputil.PostForm(%q, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
`, srv.URL)

	out := testutil.RunWasm(t, src)
	b, err := io.ReadAll(out)
	if err != nil {
		t.Fatal(err)
	}

	gotBody := strings.TrimSpace(string(b))
	if wantResBody != gotBody {
		t.Fatalf("want: %s, got: %s", wantResBody, gotBody)
	}
}

func TestHead(t *testing.T) {
	const (
		wantHeaderKey   = "want header key"
		wantHeaderValue = "want header value"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodHead {
			t.Fatalf("want: %s, got: %s", http.MethodHead, req.Method)
		}
		w.Header().Set(wantHeaderKey, wantHeaderValue)
	}))
	defer srv.Close()

	src := fmt.Sprintf(`
package main

import (
	"fmt"

	"github.com/syumai/tinyutil/httputil"
)

func main() {
	resp, err := httputil.Head(%q)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Header.Get(%q))
}
`, srv.URL, wantHeaderKey)

	out := testutil.RunWasm(t, src)
	b, err := io.ReadAll(out)
	if err != nil {
		t.Fatal(err)
	}

	gotHeaderValue := strings.TrimSpace(string(b))
	if wantHeaderValue != gotHeaderValue {
		t.Fatalf("want: %s, got: %s", wantHeaderValue, gotHeaderValue)
	}
}
