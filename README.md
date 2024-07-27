# Qveen

Generate files from templates.

## Parameters

Generation is done based on parameter files. A parameter file should
contain TOML and optionally include a table named `meta`.

Files will be generated from templates, and the values those templates
will have access to are going to be those defined in that file.

The `meta` table, if present, may contain the following entries:

- `template`: A path to the Go template file to be expanded;
- `output`: The file path in which to store the resulting file;
- `prompt`: A list of values to be provided interactively.

The prompt key is expected to contain an array of tables with the
following keys:

- `kind`: Determines the type of prompt to present. Allowed values are
  `input`, for single line texts; `text`, for potentially multiline
  texts; `confirm`, for a boolean true / false; and `select`, for
  selection amongst a set of options;
- `name`: Name of the variable in which to bind;
- `title`: Text to show when prompting.

If `kind` is set to `select`, an `options` field is also expected to
exist and contain the options from which the user will select. Each
option may be either a simple string or a table containing `title` and,
optionally, `value`. In the latter case, `title` is the text the user
will see, and `value` is the value that will become available in the
context of the template. If the option is a string, or if `value` is
omitted from the table, then the value will be a string with the same
contents as the title, which, in the former case, is just the value of
the string.

Additionaly, `meta.output` may be a directory ending with `/`. In that
case, it will function as a prefix, and the remainder of the path must
be provided as a flag when invoking the tool.

Almost\* every value in the parameter file may reference others using
template syntax, however expansion is done only once, i.e. not
recursively, and in three steps: first for the `meta.prompt` values,
before actually performing the prompts, then for regular values outside
of the `meta` table, and the for the remaining values in the `meta` table.

\* `meta.prompt[#].kind` is currently an exception and does not expand.

Example:

``` toml
# Available as {{.language}}
language = "en_US"

[meta]
template = "templates/route.go.tmpl"
output = "routes/{{snakecase .name}}_route.go"

# Available as {{.name}}
[[meta.prompt]]
name = "name"
kind = "text"
```

## Arguments and flags

Parameter files shall be provided as positional arguments for the
`qveen` executable.

`qveen` also accepts the following flags:

- `--template` / `-t`: A path to the Go template file to be
expanded;
- `--output` / `-o`: The file path in which to store the resulting
file. May be a prefix, like `meta.output`;
- `--meta-key` / `-m`: The key of the field to look for metadata in
instead of `meta`. Must be a top-level field;
- `--left-delim` / `-l`: String to use as the left delimiter of actions
in the templates, whose default is `{{`;
- `--right-delim` / `-r`: String to use as the right delimiter of
actions in the templates, whose default is `}}`.

`--template` and `--output` will be expanded as templates in the same
way as `meta.template` and `meta.output`.

Options provided both in the parameter file and as a flags will assume
the value assigned in the flag, unless it is the case of a pair of
`outputs` where one of them is a prefix, in which case both values 
will be combined.

`--template` will only be used if a single parameter file was provided.
If there were multiple, each one of them should have their template
file specified in `meta`. `--output` will only be considered for
multiple files if it defines a prefix.

Example:

``` shell
# Single file.
qveen -t templates/controller.ts.tmpl -o 'src/controllers/{{kebabcase .name}}.ts' qveen/auth.toml

# Multiple files (potentially).
qveen -o src/controllers/ qveen/controllers/*.toml
```

## Templates

Templates are regular Go template files. The `.` object will be a map
containing the values in the parameter file.

The utility functions:

- `uppercase`
- `lowercase`
- `titlecase`
- `pascalcase`
- `camelcase`
- `snakecase`
- `kebabcase`
- `constcase`
- `dotcase`
- `sentencecase`

are also provided. They work on the assumption that words are separated
by spaces, and work best when the string is all lower case.

They should all work fine for non-ASCII graphemes.

Example:

``` c
#include <stdio.h>

int main() {
    printf("{{uppercase .text}}!\n");
    return {{.status}};
}
```

## Disclaimer

This project is in early development and is thus likely to contain bugs.
Feel free to report them.

### TODO

- Allow usage of this project as a library;
- Add more functions to templates, in particular for escaping;
- Generate man page;
- Document installation;
- Include examples of usage patterns;
