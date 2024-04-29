## OpenTelemetry SDK Configuration Validator

This application will validate a yaml or json file against the [OpenTelemetry
SDK Configuration schema]().

### Use

```
$ go get
$ go build

$ ./otel_config_validator -o out.json examples/kitchen-sink.yaml
Valid OpenTelemetry Configuration!
```

### Testing

Run the Go unit tests:

```
$ go test .
```

Running the tests of the compiled CLI requires
[shelltest](https://github.com/simonmichael/shelltestrunner),
[jq](https://github.com/jqlang/jq/) and [yq](https://github.com/mikefarah/yq):

```
$ shelltest -c --diff --all shelltests/*.test
```

