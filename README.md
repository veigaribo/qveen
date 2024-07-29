# Qveen

Generate files from templates.

If you are looking for an usage example, look no further &mdash;
Qveen uses itself: <https://github.com/veigaribo/qveen/tree/main/qveen>.

## Parameters

Generation is done based on parameter files. A parameter file should
contain TOML and optionally a table with metadata, named `meta` by
default.

Files will be generated from templates with access to the values
provided in the parameter file.

The `meta` table, if present, may contain the following entries:

- `template`: A path to the Go template file to be expanded;
- `output`: The file path in which to store the resulting file;
- `pairs`: Pairs of templates and outputs;
- `prompts`: A list of values to be provided interactively.

You can either provide `template` and `output` to process a single
template or `pairs` to process multiple templates with the same data.
Or both, in which case the root `template` and `output` constitute what
is considered the first pair for reporting purposes.

`template` may be a path to a local file or a URL. `output` may be a
path to a local file or `-`, which will make the file be output to
stdout.

Additionally, `meta.output` may be a directory ending with `/`. In that
case, it will function as a prefix, and the remainder of the path must
be provided as a flag when invoking the tool.

The `prompts` key allows for values to be interactively provided and is
expected to contain an array of tables with the following keys:

- `kind`: Determines the type of prompt to present. Allowed values are
  `input`, for single line texts; `text`, for potentially multiline
  texts; `confirm`, for a boolean true or false; and `select`, for
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

Values for prompts may also be provided as flags. This will be required
if not running in an interactive terminal.

Almost\* every value in the parameter file may reference others using
template syntax, however expansion is done only once, i.e. not
recursively, and in three steps: first for the `meta.prompts` values,
before actually performing the prompts, then for regular values outside
of the `meta` table, and the for the remaining values in the `meta`
table.

\* `meta.prompts[].kind` is currently an exception and does not expand.

Example:

``` toml
# Available as {{.language}}
language = "en_US"

# Using `meta.{template,output}`.
[meta]
template = "templates/route.go.tmpl"
output = "routes/{{snakecase .name}}_route.go"

# Using `meta.pairs[].{template,output}` at the same time.
# Both `output` files will be generated.
[[meta.pairs]]
template = "templates/route_test.go.tmpl"
output = "routes/{{snakecase .name}}_route_test.go"

# Available as {{.name}}
[[meta.prompts]]
name = "name"
kind = "text"
title = "Name:"
```

## Arguments and flags

Parameter files shall be provided as positional arguments for the
`qveen` executable. They may be a path to a local file, an URL, or `-`,
which means the contents should come from stdin.

`qveen` also accepts the following flags:

- `--template` / `-t`: Defines the path to the Go template file to be
  expanded. Will override `meta.template` if and only if processing a
  single template \* output pair;
- `--output` / `-o`: Defines the file path in which to store the
  resulting file. May be a prefix, like `meta.output`. If it isn't, it
  will override `meta.output` if and only if processing a single
  template \* output pair;
- `--prompt-value` / `-o`: Provides a value for a prompt. Example:
  `-p name="value"` will set the value `value` to the prompt named
  `name`. The value should be valid for the respective kind of prompt.
  If providing a value for a `select` prompt, use the option's `title`
  on the right-hand side: `-p selection="Option's title"`;
- `--meta-key` / `-m`: Changes the key in which to look for metadata
  from the default of `meta`. Must be a top-level field;
- `--left-delim` / `-l`: Changes the string to use as the left
  delimiter of actions in the templates from the default `{{`;
- `--right-delim` / `-r`: Changes the string to use as the right
  delimiter of actions in the templates from the default `}}`.
- `--overwrite` / `-y`: Skips the confirmation that Qveen would
  normally require before writing over existing files.
- `--help` / `-h`: Displays information and immediately exits.

`--template` and `--output` will be expanded as templates in the same
way as `meta.template` and `meta.output`.

Example:

``` shell
qveen -l '<%=' -r '%>' -t templates/controller.ts.tmpl -o 'src/controllers/<%= kebabcase (lowercase .name) %>.ts' -p name=Brenda qveen/auth.toml
```

## Templates

Templates are regular Go template files. The `.` object will be a map
containing the values in the parameter file.

Various functions and templates are provided for convenience. See
<https://github.com/veigaribo/qveen/blob/main/docs/template_lib.md>
for more.

> You can use `jq`!

Example:

``` c
#include <{{dotcase "stdio h"}}>

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
