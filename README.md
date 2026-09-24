# tidy-units

Byte sizes and durations show up as strings from all kinds of places: config
files, API responses, command-line flags, spreadsheets someone hand-typed.
The spellings never agree - `1.5GB`, `1.5 gb`, `1500000`, `1.5GiB` might all
mean roughly the same thing to the human who wrote them, but they aren't the
same number, and code that only handles one spelling breaks the moment
someone pastes in another.

`tidy-units` is a small Go package for turning that mess into an exact value,
and turning an exact value back into something a person would write. Every
exported function is a pure function of its input: no globals, no I/O, no
hidden state. That makes them trivial to unit test and safe to call from
anywhere, including concurrently.

## Install

```
go get github.com/xedelgado/tidy-units
```

## Usage

```go
package main

import (
	"fmt"

	tidyunits "github.com/xedelgado/tidy-units"
)

func main() {
	bytes, err := tidyunits.ParseBytes("1.5 GiB")
	if err != nil {
		panic(err)
	}
	fmt.Println(bytes)                    // 1610612736
	fmt.Println(tidyunits.FormatBytes(bytes)) // "1.61 GB"

	d, err := tidyunits.ParseDuration("2 hours 15 mins")
	if err != nil {
		panic(err)
	}
	fmt.Println(d)                          // 2h15m0s
	fmt.Println(tidyunits.FormatDuration(d)) // "2.25h"
}
```

## What's supported

`ParseBytes` accepts a number followed by an optional unit: `b`/`byte`/
`bytes`, decimal units `kb`/`mb`/`gb`/`tb`/`pb`, and binary units `kib`/
`mib`/`gib`/`tib`/`pib`. Matching is case-insensitive and whitespace between
the number and unit is optional. `FormatBytes` always renders back using
decimal units, since that's what most readers expect from a plain byte
count.

`ParseDuration` accepts one or more number+unit pairs, optionally separated
by spaces or commas: `ns`, `us`/`µs`, `ms`, `s`/`sec`/`seconds`, `m`/`min`/
`minutes`, `h`/`hr`/`hours`, `d`/`day`/`days` (treated as a fixed 24h).
`FormatDuration` picks the single largest unit that keeps the number
readable, rather than mixing units the way `time.Duration.String` does.

## Status

Early. The parsing rules above cover common cases but not every locale or
abbreviation in the wild - see the issues for what's planned next.

## License

MIT, see [LICENSE](LICENSE).
