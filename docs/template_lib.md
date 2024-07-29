This file documents the functions and templates made available for users
of Qveen.

In this document, a "character" means graphemes representable as a
single Unicode code point.

A "word" is defined as something that follows a space (as defined by
Unicode) or the start of the string.

If you don't know what a character being "title case" means, read it as
"upper case".

# Functions

## Text

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

=> Failed to execute template: template: qveen:1:2: executing "qveen" at <err "This is wrogn.">: error calling err: This is wrogn.
```

### tmpl :: string -> any -> string

Invokes the template with the name given in the first argument and the
argument in the second. Is basically the same as `{{template}}`, but,
being a function, allows for arbitrary expressions / pipelines as
arguments.

```
{{- define "n" -}}{{.}}{{- end -}}
{{tmpl "n" 3 }}

=> 1
```

# Templates

### join

Join intersperses template segments, so that every element is separated
by something, but that something does not appear before or after the
ends. The parameter should be a map containing the following keys and
values:

- `tmpl`: A string containing the name of the template that will be
renderer on each item. That template will receive its corresponding item
as its argument;
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
