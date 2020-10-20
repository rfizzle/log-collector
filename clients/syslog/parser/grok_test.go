package parser

import (
	"testing"
)

func TestParseGrok(t *testing.T) {
	grokLog1 := `2020-01-29T21:26:10Z 10.10.10.10 PulseSecure: - - - 2020-01-29 21:26:10 - ibos1 - [4.4.4.4] user1(Com1-Reliable)[Com1-Reliable-Grp-TST] - Login succeeded for user1/Com1-Reliable (session:00000000) from 5.5.5.5 with Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.30729; .NET CLR 3.5.30729; wbx 1.0.0; Zoom 3.6.0).`
	grokPattern1 := `%{TIMESTAMP_ISO8601:timestamp} %{PROG:pulse_appliance} %{PROG:software}: - - - %{TIMESTAMP_ISO8601:timestamp2} - %{WORD:w1} - \[%{NOTSPACE:ip}\] %{GREEDYDATA:user}\[%{GREEDYDATA:group}\] - %{GREEDYDATA:logmsg}`

	if _, err := ParseGrok(grokLog1, []string{grokPattern1}); err != nil {
		t.Fatalf("failed to parse Grok message: %v", err)
	}
}

func TestParseEventWithGrokPatterns(t *testing.T) {
	grokLog1 := `2020-01-29T21:26:10Z 10.10.10.10 PulseSecure: - - - 2020-01-29 21:26:10 - ibos1 - [4.4.4.4] user1(Com1-Reliable)[Com1-Reliable-Grp-TST] - Login succeeded for user1/Com1-Reliable (session:00000000) from 5.5.5.5 with Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.30729; .NET CLR 3.5.30729; wbx 1.0.0; Zoom 3.6.0).`
	grokPattern1 := `%{TIMESTAMP_ISO8601:timestamp} %{PROG:pulse_appliance} %{PROG:software}: - - - %{TIMESTAMP_ISO8601:timestamp2} - %{WORD:w1} - \[%{NOTSPACE:ip}\] %{GREEDYDATA:user}\[%{GREEDYDATA:group}\] - %{GREEDYDATA:logmsg}`
	grokExpectedKeyValuePair1 := []string{"software", "PulseSecure"}
	grokExpectedKeyValuePair2 := []string{"user", "user1(Com1-Reliable)"}
	grokExpectedKeyValuePair3 := []string{"group", "Com1-Reliable-Grp-TST"}
	grokExpectedKeyValues := [][]string{grokExpectedKeyValuePair1, grokExpectedKeyValuePair2, grokExpectedKeyValuePair3}

	if _, err := parseEventWithGrokPatterns(grokLog1, []string{grokPattern1}); err != nil {
		t.Fatalf("failed to parse Grok message: %v", err)
	}

	resultMap, err := parseEventWithGrokPatterns(grokLog1, []string{grokPattern1})

	if err != nil {
		t.Fatalf("failed to parse Grok message")
	}

	for _, v := range grokExpectedKeyValues {
		if resultMap[v[0]] != v[1] {
			t.Errorf(`resultMap["%s"] got %s; expected %s`, v[0], resultMap[v[0]], v[1])
		}
	}

}
