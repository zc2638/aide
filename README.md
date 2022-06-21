# aide

Command line installer help framework.

The installation can be performed according to the installation configuration,
and you can also quickly generate custom installation tools based on `pipeline` and `stage`.

## How to use for golang library

> go get -u github.com/zc2638/aide

## Examples

[**参考**](https://github.com/zc2638/aide/blob/main/examples)

## Pipeline

Supported prompt types:

- Input
- Password
- Text
- Confirm
- Select
- MultiSelect

### 1. Installation tool

#### For normal

Download from [Releases](https://github.com/zc2638/aide/releases)

#### For gopher

```shell
go install github.com/zc2638/aide/cmd/aide@latest
```

### 2. Create the pipeline config file

You can define your own pipeline configuration.

pipeline.yaml

```yaml
apiVersion: v1
kind: Pipeline
metadata:
  name: test
  label:
    project: aide
    version: v1
spec:
  prompts:
    - name: custom_name
      type: Input
      message: What's your name?
    - name: gender
      type: Select
      message: What's your gender?
      enum: [ "male", "female", "unknown" ]
  steps:
    - name: "step1"
      render:
        src: examples/file/test.in
        dest: testdata/test.out
    - command: env
    - command: echo $custom_name\($gender\)
```

### 3. Execute the pipeline according to the config file

```shell
aide apply -f pipeline.yaml
```