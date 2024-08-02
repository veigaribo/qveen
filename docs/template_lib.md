This file documents the template utilities made available for users of
Qveen.

In this document, a "character" mean a grapheme representable as a
single Unicode code point.

A "word" is defined as a sequence of characters that does not contain a
space.

If you don't know what a character being "title case" means, read it as
"upper case".

`=>` indicates the string output of the template. `$>` indicates
output to the terminal.

# Functions

## Arithmetic

### add :: ...int -> int
### mul :: ...int -> int

Adds or multiplies, respectively, the given integers.

```
{{mul 5 4 3 2 1}}

=> 120
```

### sub :: int -> ...int -> int
### div :: int -> ...int -> int

Successively subtracts or divides, respectively, the first argument
by the other ones.

```
{{sub 21 13 8}}

=> 0
```

### rem :: int -> int -> int

Returns the remainder of the division of the first argument by the
second.

```
{{rem 25 7}}

=> 4
```

## Text

### join :: []string -> string -> string

Returns a new string where each of the strings in the first argument
have been interspersed with the string in the second argument.

```
{{join (list "a" "b" "c") ", "}}

=> a, b, c
```

### uppercase :: string -> string
### lowercase :: string -> string

Transforms the string to upper and lower case respectively.

```
{{uppercase "in the early days computers were much simpler"}}

=> IN THE EARLY DAYS COMPUTERS WERE MUCH SIMPLER
```

### titlecase :: string -> string

Makes the first character of each word title case. Other characters are
not affected.

Note that this is not smart enough to ignore words such as "the" and
"is".

```
{{titlecase "in the early days computers were much simpler"}}

=> In The Early Days Computers Were Much Simpler
```

### pascalcase :: string -> string

Makes the first character of each word title case and removes every
space. Other characters are not affected.

```
{{pascalcase "in the early days computers were much simpler"}}

=> InTheEarlyDaysComputersWereMuchSimpler
```

### camelcase :: string -> string

Makes the first character of each word, except the first, title case,
and removes every space. Other characters are not affected.

```
{{camelcase "in the early days computers were much simpler"}}

=> inTheEarlyDaysComputersWereMuchSimpler
```

### snakecase :: string -> string

Substitutes every space character for an `_`. Other characters are not
affected.

```
{{snakecase "in the early days computers were much simpler"}}

=> in_the_early_days_computers_were_much_simpler
```

### kebabcase :: string -> string

Substitutes every space character for a `-`. Other characters are not
affected.

```
{{kebabcase "in the early days computers were much simpler"}}

=> in-the-early-days-computers-were-much-simpler
```

### constcase :: string -> string

Transforms every character to upper case and substitutes every space
character for an `_`. Same as `uppercase` then `snakecase`.

```
{{constcase "in the early days computers were much simpler"}}

=> IN_THE_EARLY_DAYS_COMPUTERS_WERE_MUCH_SIMPLER
```

### dotcase :: string -> string

Substitutes every space character for a `.`. Other characters are not
affected.

```
{{dotcase "in the early days computers were much simpler"}}

=> in.the.early.days.computers.were.much.simpler
```

### sentencecase :: string -> string

Makes the first character title case. Other characters are not affected.

```
{{sentencecase "in the early days computers were much simpler"}}

=> In the early days computers were much simpler
```

### escapebackslash :: string -> string -> string

Will return a new string similar to the one in the second argument but
with a backslash preceding every character contained in the first
string.

```
{{escapebackslash "\"\\" "The title of this paper is an homage to David Goldberg’s classic paper \"What Every Computer Scientist Should
Know About Floating-Point Arithmetic\""}}

=> The title of this paper is an homage to David Goldberg’s classic paper \"What Every Computer Scientist Should
Know About Floating-Point Arithmetic\"
```

Note that the backslashes in the input strings are there just as per
Go template syntax and do not actually exist.

### escapedouble :: string -> string -> string

Will return a new string similar to the one in the second argument but
with very character contained in the first string being preceded by
itself.

```
{{escapebackslash "\"\\" "The title of this paper is an homage to David Goldberg’s classic paper \"What Every Computer Scientist Should Know About Floating-Point Arithmetic\""}}

=> The title of this paper is an homage to David Goldberg’s classic paper ""What Every Computer Scientist Should Know About Floating-Point Arithmetic""
```

### escapehtml :: string -> string

HTML escapes the input string. Should do the same as the builtin `html`
function.

```
{{escapehtml "#include <emmintrin.h>"}}

=> #include &lt;emmintrin.h&gt;
```

### repl :: string -> string -> string -> string

Replaces all non-overlapping occurrences of the first string in the
third string with the second string.

```
{{repl "\n" "\\n" "- RAM hardware design (speed and parallelism).\n- Memory controller designs.\n- CPU caches.\n- Direct memory access (DMA) for devices"}}

=> - RAM hardware design (speed and parallelism).\n- Memory controller designs.\n- CPU caches.\n- Direct memory access (DMA) for devices
```

## Containers

### map :: ...any -> map[string]any

Shall be called with alternating keys and values. Returns a map that
associates every key with its accompanying value.

```
{{range $k, $v := map "key1" "value1" "key2" "value2"}}{{$k}}: {{$v}}{{end}}

=> key1: value1key2: value2
```

### list :: ...any -> []any

Creates a slice with the argument values.

```
{{range list "first" "last"}}{{.}}~{{end}}

=> first~last~
```

### set :: ([]any | map[string]any) -> (int | string) -> any -> ()

The first argument should be either a slice or a map. The second
argument should be something capable of indexing the container in the
first argument (an int for a slice or a string for a map). Then, the
value at the given index of the container will be set to the value in
the third argument. The slice may grow to fit such index.

```
{{- $xs := map -}}
{{- set $xs "one" 1 -}}
{{ dump $xs }}

$> {"one": 1}
```

### append :: []any -> any -> ()

Appends the second argument to the slice in the first.

```
{{- $xs := list -}}
{{- append $xs (map "one" 1) -}}
{{ dump $xs }}

$> [{"one": 1}]
```

### slice :: []any -> int -> ()
### slice :: []any -> int -> int -> ()

Behaves similarly to the builtin `slice` function.

- If given one integer `n`, will return a slice that skips the first
  `n` elements of slice in the first argument.
- If given two integers `n` and `m`, will return a slice that skips the
  first `n` elements and ends at the `m`th element of the slice in the
	first argument. The `m`th element itself is not included.

```
{{- $slice := (list 10 20 30 40 50) 1 3 -}}
{{ dump $slice }}

=> [20, 30]
```

> The type notation broke a little here.

## Assertions

### ismap :: any -> bool
### isstr :: any -> bool
### isint :: any -> bool
### isarr :: any -> bool

Returns if the argument is, respectively, a map, a string, an integer,
or a slice.

## jq

Qveen allows the usage of `jq` queries in templates. Under the hood,
<https://github.com/itchyny/gojq> is being used.

### jq1 :: string -> any -> any

Applies the query in the first argument to the object in the second.
Returns the first result found, or `nil` if there was none.

```
{{- $data := (list (map "n" 1) (map "n" 2) (map "n" 3)) -}}
{{jq1 ".[].n | select(. > 1)" $data}}

=> 2
```

### jqn :: string -> any -> []any

Applies the query in the first argument to the object in the second.
Returns a slice containing every result found.

```
{{- $data := (list (map "n" 1) (map "n" 2) (map "n" 3)) -}}
{{jqn ".[].n" $data}}

=> [1 2 3]
```

## Miscellaneous

### err :: string -> ⊥

Fails with the given reason.

```
{{err "This is wrogn."}}

$> Failed to execute template: template: qveen:1:2: executing "qveen" at <err "This is wrogn.">: error calling err: This is wrogn.
```

### dump :: ...any -> ()

Prints the arguments to stderr for inspection. The format strives to be
similar to JSON.

### probe :: any -> any

Prints the argument to stderr in the same manner as `dump`, but also
returns the argument untouched. Useful for inspecting values in
expressions.

```
{{add (mul 2 2 | probe) (mul 5 5)}}

=> 29
$> 4
```

## TOML, YAML and JSON

### toml :: any -> string
### yaml :: any -> string
### json :: any -> string

Returns the argument encoded in the given format.

> `json` is very useful for creating quoted strings.

```
{{toml (map "sometable" (map "value1" 1 "value2" 2))}}

=> [sometable]
value1 = 1
value2 = 2
```

# Templates

### join

Similar to the `join` function, but applies a template to each item
of the array. The argument should be a map containing the following
keys and values:

- `tmpl`: A string containing the name of the template that will be
  renderer on each item. That template will receive its corresponding
  item as its argument;
- `els`: The list of values from which to render;
- `sep`: The value used to separate each item;
- `pre`: If and only if there are any items in `els`, this will be
  output before anything else. Surprisingly useful.

```
{{- $items := (list "orange" "pear" "apple") -}}
{{- define "item" -}}{{.}}{{- end -}}
{{template "join" (map "tmpl" "item" "els" $items "sep" " - " "pre" "FRUITS\n")}}

=> FRUITS
orange - pear - apple
```

# Trimmers

Qveen utilises a modified template parser that recognizes the following
extra trimmers:

### Horizontal trimmer `~`

This trimmer will only trim horizontal space characters, so not
carriage return or line feed.

```
s := `
	{{- "In the early days computers were much simpler."}}
	{{~ "The various components of a system were developed together "}}
	{{- "and, as a result, were quite balanced in their performance."}}
`
```

### Vertical trimmer `#`

- If this trimmer is on the left, it will seek whitespace until it
finds something that is not whitespace. Then, if it was on a different
line from where it started, it will leave that line's line break
behind and consume the rest. If there was trailing whitespace in that
line, it will be kept too.

- If this trimmer is on the right, it will seek whitespace until it
finds something that is not whitespace. Then, if it was on a different
line from where it started, it will leave that line's whitespace behind
and consume the rest.

Made mainly to be used on control actions, where no output is produced.
Allows for cleanly joining lines without much worry for whitespaces.

```
import (
	{{# if $usesfmt #}}
	"fmt"
	{{# end #}}
	"github.com/veigaribo/qveen/utils"
	{{# if $usesfmt #}}
	"github.com/veigaribo/qveen/prompts"
	{{# end #}}
	"strings"
)
```

# Other features

You are allowed to invoke templates with variable names.

```
{{- define "n" -}}{{.}}{{- end -}}
{{- $n := "n" -}}

{{template $n 3 }}

=> 3
```

The following shorthands are accepted:

- `{{def ...}}` is the same as `{{define ...}}`;
- `{{t ...}}` is the same as `{{template ...}}`.

```
{{- def "one" -}} I {{.}} {{- end -}}
{{- $tmpl := "one" -}}

{{t $tmpl "u"}}

=> I u
```
