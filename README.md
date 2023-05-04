[![License: MIT](https://img.shields.io/github/license/erminson/gitlab-vars)](https://github.com/erminson/gitlab-vars/blob/master/LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/erminson/gitlab-vars)](https://goreportcard.com/report/github.com/erminson/gitlab-vars)
[![Tests](https://github.com/erminson/gitlab-vars/actions/workflows/test.yaml/badge.svg)](https://github.com/erminson/gitlab-vars/actions/workflows/test.yaml)
[![Codecov](https://codecov.io/github/erminson/gitlab-vars/branch/master/graph/badge.svg?token=MCOD2CJ0T1)](https://codecov.io/github/erminson/gitlab-vars)
[![Release](https://img.shields.io/github/v/release/erminson/gitlab-vars)](https://github.com/erminson/gitlab-vars/releases/latest)

# glvars
## Description
`glvars` CLI tool for import and export project-level Gitlab CI/CD Variables.
## Setup
- Using environment variables add your personal access token (`api` scope) and project id (optional): 
 ```bash
  export GLVARS_PRIVATE_TOKEN="your-personal-access-token"
  export GLVARS_PROJECT_ID=278964 # Your gitlab project id
  ```
- Or you can create config file `.glvars.yaml` (in `$HOME` directory): 
```bash
  host: gitlab.example.com # For gitlab server on own domain. By default: https://gitlab.com
  private-token: your-gitlab-private-token
  project-id: 278964
```
- You also can specify your project id using flag `-p` or `--project`:
```bash
  glvars [command] -p 278964
```
## Use cases and examples
- Export variables in JSON format from project specified in config file or env variable:
```bash
  glvars export
```
or you can explicitly specify project id using flag `-p`:   
```bash
  glvars export -p 278964
```
- To import variables into a project use command `import` with `-f` or `--filename` flag with value path to file with array of variables
```bash
  glvars import -f ./vars.json
```
JSON file example:
```json
[
  {
    "variable_type": "env_var",
    "key": "MYSQL_HOST",
    "value": "127.0.0.1",
    "protected": false,
    "masked": false,
    "raw": true,
    "environment_scope": "production"
  },
  {
    "variable_type": "env_var",
    "key": "MYSQL_USER",
    "value": "admin",
    "protected": false,
    "masked": false,
    "raw": true,
    "environment_scope": "production"
  }
]
```

## Installation

Download a binary suitable for your OS at the [release page](https://github.com/erminson/gitlab-vars/releases/latest).

You can build project from source:

```bash
  go build -o glvars cmd/main.go
```