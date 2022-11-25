Lambda Sim
==========

Lambda Sim, as the name implies, is a simulator for AWS Lambda. Its puspose is to allow a user to configure the lambda function code to be
initialised and run as it would in a Lambda runtime. The runtimes themselves must be configured with the local runtime API URLs. E.g., to
run a custom, compiled runtime to interact with Lambda Sim:

```shell
AWS_LAMBDA_RUNTIME_API=localhost:8999/z_test ./build/z_test
```

Events can then be passed using the CLI or a web server listening on a configurable port.
