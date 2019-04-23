## api-check

`api-check` is a cli tool that allows you to simply and quickly test restful APIs. Written in Go, `api-check` can be used to test APIs written in any language.

`api-check` works by making HTTP requests to your server and then asserting that the correct response is received as defined in a test definition files within your project.


## prerequisites

To install `api-check` from source do the following: 
```
go get github.com/JonathonGore/api-check
cd ${GOPATH}/src/github.com/JonathonGore/api-check && go install
```

`api-check` has no external dependencies!

**Note:** You need to ensure `${GOPATH}/bin` is in your `PATH` 

*Coming soon `api-check` available in Homebrew.*

## usage

`api-check` can be used as a standalone binary, or can be integrated directly with `go test`.

### Running Standalone

`api-check` looks for test definitions stored in `json` files with the `.ac.json` extension stored in any subdirectory of your project. 

You can run all test definitions in your project by running `$ api-check run` in the root of your project directory.

For more info on available commands you can run:

`$ api-check help`

### Integrating with go test

In the root of your project create a file `main_test.go` with the following contents.

```
package main

import (
	"testing"

	"github.com/JonathonGore/api-check/suite"
)

func TestMain(t *testing.T) {
	suite.Run(t)
}

```

The above will invoke `api-check run` and run all test definitions at or below the current directory.

## examples

Using api-check you can assert that your server produces exactly the correct JSON, by using the `json` key in the response body.

**users.ac.json:**

```
[{
  "hostname": "http://localhost:3000",
  "endpoint": "/users/Jack",
  "method": "get",
  "response": {
    "code": 200,
    "headers": {
        "Content-Type": "application/json"
    },
    "json": {
       "username": "Jack",
       "email": "jack@gmail.com",
       "first_name": "Jack",
       "last_name": "Gore"
    }
  }
}]
```

Additionally, instead of asserting an exact JSON match, `api-check` also allows you to assert the structure of a JSON response using the `ofType` key:


```
[{
  "hostname": "http://localhost:3000",
  "endpoint": "/users/Jack",
  "method": "get",
  "response": {
    "code": 200,
    "headers": {
        "Content-Type": "application/json"
    },
    "ofType": {
       "username": "string",
       "email": "string",
       "user_id": "number"
       "first_name": "string",
       "last_name": "string",
       "aliases": [
           "string"
       ]
    }
  }
}]
```

The above is a test files each contain a single test definition.

These test definitions will make a `GET` request to `http://localhost:3000/users/Jack`. It will assert that it receives the response are specified in the `response` key.

### Configuring api-check

`api-check` can be configured by placing a file named `.ac.json` in the directory where you will run your `api-check` commands.

The `.ac.json` file is a plain JSON (see `examples/` for an exmaple of this file) file that supports the following keys:

* `setup-script`
    * The name of a bash script to be ran before executing any of the test suites.
* `cleanup-script`
    * The name of a bash script to be ran after the execution of all tests.
* `hostname`
    * The default hostname to be used in your test definitions, allows you to not have to specify hostname in each test definition.


