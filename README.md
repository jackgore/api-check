## api-check

`api-check` is a tool that allows you to simply and quickly test restful APIs. Written in Go, `api-check` can be used to test APIs written in any language.

`api-check` works by creating test definitions which are plain JSON files. Test definitions describe the expected ins and outs of API endpoints. The test definitions are then ran against your running API to verify the results.

## prerequisites

To install `api-check` from source do the following: 
```
go get github.com/JonathonGore/api-check
cd ${GOPATH}/src/github.com/JonathonGore/api-check && go install
```

`api-check` has no external dependencies!

**Note:** You need to ensure `${GOPATH}/bin` is in your `PATH` 

** Coming soon `api-check` available in docker.

## usage

`api-check` can be used as a standalone binary, or can be integrated directly with `go test`.

### Running Standalone

`api-check` looks for test definitions stored in `json` files with the `.ac.json` extension stored in any subdirectory in your project. 

You can run all test definitions in your project by running `api-check` in the root of your project directory.

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

The above will invoke `api-check` and run all test definitions at or below the current directory.

## examples

An example test definition for testing a users endpoint of a webserver would look like this:

**users.ac.json:**

```
[{
  "hostname": "http://localhost:3000",
  "method": "post",
  "endpoint": "/users",
  "request": {
    "json": {
       "username": "Jack",
       "password": "password",
       "email": "jack@gmail.com",
       "first_name": "Jack",
       "last_name": "Gore"
    }
  },
  "response": {
    "code": 200
  }
},
{
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

The above is a test definition file containing two test definitions.

The first test definition will make a `POST` request to the url specified in the `hostname` and`endpoint` paramaters (`http://localhost:3000/users`). It will attach the provided `json` struct in the request body. Then it will `assert` that it receives a `200` statuscode and a `Content-Type` header from the webserver.

The second test definition will make a `GET` request to `http://localhost:3000/users/Jack`. It will assert that it receives the attached `json` struct in the response body, as well as ensuring it receives a `200` status code.

Running `api-check` in the same directory as `users.ac.json` will cause the tests to be ran and if successful result in the following output:

```
Running go api-check

API Check Test for: http://localhost:3001/users succeeded
API Check Test for: http://localhost:3001/users/Jack succeeded

2 tests ran. 2 successful. 0 failures.
```

