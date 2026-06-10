package main

import (
	"bytes"
	"testing"
)

func env(kv map[string]string) func(string) string {
	return func(key string) string { return kv[key] }
}

func TestRun_Success(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := run(env(map[string]string{
		"GITHUB_ACTIONS": "true",
		"GITHUB_TOKEN":   "token",
	}), &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if stderr.String() != "plugin_schema_version=1\n" {
		t.Fatalf("expected no stderr output, got %q", stderr.String())
	}
}

func TestRun_Failure(t *testing.T) {
	t.Parallel()

	var stderr bytes.Buffer
	code := run(env(map[string]string{}), &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if stderr.Len() == 0 {
		t.Fatal("expected stderr output")
	}
}
