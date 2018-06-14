## api-check

`api-check` is a tool that allows you simply and quickly test your restful api's.

`api-check` works by creating test definitions in plain JSON, then the test definitions are ran against your web-server to verify the results.

## prerequisites

To install `api-check` install it by running 
```
go get github.com/JonathonGore/api-check
cd ${GOPATH}/src/github.com/JonathonGore/api-check && go install
```

`api-check` has no external dependencies!

## usage

`api-check` looks for test definitions stored in `json` files with the `.ac.json` extension stored in any subdirectory in your project. 


You can run all test definitions by running `api-check` in the root of your project directory.

## examples

An example test definition for testing a users endpoint of a webserver would look like this.

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
    "json": {
       "username": "Jack",
       "email": "jack@gmail.com",
       "first_name": "Jack",
       "last_name": "Gore"
    }
  }
}]
```
The first will test making a `post` request to the url specified in the `hostname` and`endpoint` paramaters (`http://localhost:3000/users`). It will attach the provided `json` struct in the request body and it will `assert` that it receives a `200` statuscode from the webserver.

The second test definition will make a `get` request to `http://localhost:3000/users/Jack`. It will assert that it receives the attach `json` struct in the response body, as well as ensuring it receives a `200` status code.

Running `api-check` in the same directory as `users.ac.json` will cause the tests to be and if successful result in the following output:

```
Running go api-check

API Check Test for: http://localhost:3001/users succeeded
API Check Test for: http://localhost:3001/users/Jack succeeded

2 tests ran. 2 successful. 0 failures.
```
