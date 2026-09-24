package tidyunits

import "testing"

func TestParseBytes(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"0", 0},
		{"512", 512},
		{"1kb", 1000},
		{"1 KB", 1000},
		{"1kib", 1024},
		{"1.5gb", 1500000000},
		{"2 MiB", 2 * 1024 * 1024},
		{"3bytes", 3},
	}

	for _, c := range cases {
		got, err := ParseBytes(c.in)
		if err != nil {
			t.Fatalf("ParseBytes(%q) returned error: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("ParseBytes(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestParseBytesErrors(t *testing.T) {
	for _, in := range []string{"", "gb", "-5gb", "5 elephants"} {
		if _, err := ParseBytes(in); err == nil {
			t.Errorf("ParseBytes(%q) expected error, got nil", in)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0 bytes"},
		{1, "1 byte"},
		{512, "512 bytes"},
		{1500, "1.5 kB"},
		{1500000, "1.5 MB"},
		{1000000000, "1 GB"},
	}

	for _, c := range cases {
		got := FormatBytes(c.in)
		if got != c.want {
			t.Errorf("FormatBytes(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}
