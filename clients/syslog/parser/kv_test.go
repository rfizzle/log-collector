package parser

import (
	"github.com/jjeffery/kv"
	"testing"
)

var kvMessage1 = []byte(`dvc=10.118.182.162 rt=1600239263565 cat=illusive{{COLON}}SYS`)
var kvMessage2 = []byte(`message this stuff dvc=10.118.182.162 rt=1600239263565 cat=illusive{{COLON}}SYS`)
var kvMessage3 = []byte(`dvc==10.118.182.162 rt==1600239263565 cat==illusive{{COLON}}SYS`)

func TestParseKV(t *testing.T) {
	kvExpectedLength := 6
	kvExpectedKeyValuePair1 := []string{"dvc", "10.118.182.162"}
	kvExpectedKeyValuePair2 := []string{"rt", "1600239263565"}
	kvExpectedKeyValuePair3 := []string{"cat", "illusive{{COLON}}SYS"}
	kvExpectedKeyValues := [][]string{kvExpectedKeyValuePair1, kvExpectedKeyValuePair2, kvExpectedKeyValuePair3}

	if text, list := kv.Parse(kvMessage1); text != nil {
		t.Fatalf("failed to parse KV message")
	} else if len(list) != kvExpectedLength {
		t.Fatalf("len(list) got %v; expected %v", len(list), kvExpectedLength)
	}

	resultMap, err := parseKeyValue(string(kvMessage1), false)

	if err != nil {
		t.Fatalf("failed to parse KV message")
	}

	for _, v := range kvExpectedKeyValues {
		if resultMap[v[0]] != v[1] {
			t.Errorf(`resultMap["%s"] got %s; expected %s`, v[0], resultMap[v[0]], v[1])
		}
	}
}

func TestParseKV2(t *testing.T) {
	if text, _ := kv.Parse(kvMessage2); text == nil {
		t.Errorf("failed to error on invalid KV message")
	}

	if text, _ := kv.Parse(kvMessage3); text == nil {
		t.Errorf("failed to error on invalid KV message")
	}
}
