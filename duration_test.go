package tidyunits

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	cases := []struct {
		in   string
		want time.Duration
	}{
		{"90s", 90 * time.Second},
		{"1h30m", time.Hour + 30*time.Minute},
		{"2 hours 15 mins", 2*time.Hour + 15*time.Minute},
		{"1.5 days", 36 * time.Hour},
		{"1.5h", 90 * time.Minute},
		{"500ms", 500 * time.Millisecond},
		{"1us", time.Microsecond},
		{"1µs", time.Microsecond},
		{"1ns", time.Nanosecond},
		{"1d, 2h, 30m", 24*time.Hour + 2*time.Hour + 30*time.Minute},
		{"  45s  ", 45 * time.Second},
	}

	for _, c := range cases {
		got, err := ParseDuration(c.in)
		if err != nil {
			t.Fatalf("ParseDuration(%q) returned error: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("ParseDuration(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseDurationErrors(t *testing.T) {
	for _, in := range []string{"", "  ", "90", "90 fortnights", "h", "1.5"} {
		if _, err := ParseDuration(in); err == nil {
			t.Errorf("ParseDuration(%q) expected error, got nil", in)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		in   time.Duration
		want string
	}{
		{0, "0ns"},
		{500 * time.Nanosecond, "500ns"},
		{1500 * time.Nanosecond, "1.5µs"},
		{1500 * time.Microsecond, "1.5ms"},
		{1500 * time.Millisecond, "1.5s"},
		{90 * time.Second, "1.5m"},
		{90 * time.Minute, "1.5h"},
		{36 * time.Hour, "1.5d"},
	}

	for _, c := range cases {
		got := FormatDuration(c.in)
		if got != c.want {
			t.Errorf("FormatDuration(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParseFormatDurationRoundTrip(t *testing.T) {
	for _, in := range []string{"45s", "1.5m", "500ms", "1.5h", "1.5d"} {
		d, err := ParseDuration(in)
		if err != nil {
			t.Fatalf("ParseDuration(%q) returned error: %v", in, err)
		}
		if got := FormatDuration(d); got != in {
			t.Errorf("FormatDuration(ParseDuration(%q)) = %q, want %q", in, got, in)
		}
	}
}
