### Overview

YTPL is a simple template render of YAML configuration files.  It uses [Go template langues](https://pkg.go.dev/text/template).

### Additional template functions

* [Sprig Function Documentation](http://masterminds.github.io/sprig/)

### Usage

```shell
ytpl -input ./config -output ./output
```

### Command line options

```
-input string
      a directory with yaml files which should be rendered (default ".")
-output string
      a folder with the result yaml files (default "./output")
```

### Example

[A sample yaml config](./test/dev/na/a.yaml) 
```yaml
a:
  b:
    c: {{ .abc }}
```

Yaml files which start with "_" prefix are used as a source of values for variables.

[A sample values holder](./test/_env.yaml)
```yaml
abc: Hello world
```

The output result will be

```yaml
a:
  b:
    c: Hello world
```