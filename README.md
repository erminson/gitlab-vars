# glvars
## Description
`glvars` tool for import and export Gitlab CI/CD variables
## Setup
- Using environment variables: 
 ```bash
  export GLVARS_PRIVATE_TOKEN="your-gitlab-private-token"
  export GLVARS_PROJECT_ID=1 # Your gitlab project id
  ```
- Using config file `.glvars.yaml` (in `$HOME` directory): 
```bash
  private-token: your-gitlab-private-token
  project-id: 1
```
- Using flag `-p` or `--project` to set project id:
```bash
  glvars [command] -p 1
  glvars [command] --project 1
```
## Use cases and examples
- Export variables in json forma
```bash
  glvars export
```
or
```bash
  glvars export -p 1
```
- Import variables. To import variables into a project use command `import` with `-f` or `--filename` flag with value path to file with array of variables
```bash
  glvars import -f ./vars.json
```
  or
```bash
  glvars import --filename ./vars.json
```

## Install
You can build project from source
```bash
  go build -o glvars cmd/main.go
  project-id: 1
```

Or download precompiled binaries from [release](https://github.com/erminson/gitlab-vars/actions/workflows/release.yml)