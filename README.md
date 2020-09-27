# MockRestServer

This server can be used to store and receive any (text) data via REST.
It's main purpose is to serve dummy-data for development or unit-tests where you need
a determined output.

You can run this service everywhere and use any data you want. The only thing you have to
do is to adjust your HTTP-Calls if you like to use this mocking-service in your tests.


## Download and start application

Download the current version from [Releases](https://github.com/mbedded/MockRestServer/releases). You can start
the executables via command-line.

```shell
# Allow execution of application
chmod +x MockRestServer_linux_amd64

# Start application with default parameters
./MockRestServer_linux_amd64

# Start application with custom parameters
./MockRestServer_linux_amd64 -port=1234 -database=myCustomDatabaseFile.db
```

The command line should show something like:

```
MockRestServer - Version: 1.0
2020/09/27 14:12:26 Webserver will be startet at http://localhost:8080
```

Now open http://localhost:8080 to access the application. If there is an error
it will be displayed in the console.

The application support different parameters so:

| Parameter | Description                                              | Default value     |
| --------- | -------------------------------------------------------- | ----------------- |
| -port     | Defines the local port of the http-webserver.            | 8080              |
| -database | Defines the path to the database file (sqlite).          | mockRestServer.db |
| -version  | If used, displays the version and exits the application. | -                 |


## Usage via Websites

If you open the homepage of this appliation you'll see links to create and display
your Mocks. On [Create Mock (localhost:8080)](http://localhost:8080/create) you can add new Mocks by
giving them a name and a content.

On the link [Show all (localhost:8080)](http://localhost:8080/showall) you can see all your created
Mocks and edit or delete them.


## Usage via API

The application is using REST-Endpoints to create, update, delete and show Mocks.
The model which has to be send and which is received by the api has the following layout:
```json
{
  "Key": "The unique key of the mock",
  "Content" : "Content of this mock"
}
```

This table is an overview of all endpoints. The actions are described
in detail below.

| Description     | HTTP Verb | Url                                      |
| --------------- | --------- | ---------------------------------------- |
| Create Mock     | POST      | http://localhost:8080/api/mock           |
| Update Mock     | PUT       | http://localhost:8080/api/mock           |
| Get Mock        | GET       | http://localhost:8080/api/mock/key/{key} |
| Delete Mock     | DELETE    | http://localhost:8080/api/mock/key/{key} |
| Get raw content | GET       | http://localhost:8080/raw/{key}          |
| Get all Mocks   | GET       | http://localhost:8080/api/mock/all       |


### Create new Mock

Creates a new Mock with the given key. If they key is empty a random key will be generated.

```bash
# Call via cURL
curl --request POST \
  --url http://localhost:8080/api/mock \
  --header 'content-type: application/json' \
  --data '{
	"Key": "",
	"Content": "my content"
}'

# Response
{
  "Key": "d784f08a-b732-40f9-aa11-6288ee733dea",
  "Content": "my content"
}

```

### Update existing Mock

Updates an existing Mock with the given key. If the key is wrong or missing you'll
receive an error.

```bash
# Call via cURL
curl --request PUT \
  --url http://localhost:5050/api/mock \
  --header 'content-type: application/json' \
  --data '{
  "Key": "d784f08a-b732-40f9-aa11-6288ee733dea",
  "Content": "new content"
}'

# Response
{
  "Key": "d784f08a-b732-40f9-aa11-6288ee733dea",
  "Content": "new content"
}
```

### Get existing Mock

Returns one Mock by its key. If the key is invalid you'll receive an error.

```bash
# Call via cURL
curl --request GET \
  --url http://localhost:5050/api/mock/key/d784f08a-b732-40f9-aa11-6288ee733dea

# Response
{
  "Id": 1,
  "Key": "d784f08a-b732-40f9-aa11-6288ee733dea",
  "Content": "new content"
}
```

**Hint:** You receive the internal database ID here but it's not used. If you need this
value somehow, feel fre to use it.


### Delete existing Mock

This call deletes a Mock from the database. If the key is invalid you'll receive an error.

```bash
# Call via cURL
curl --request DELETE \
  --url http://localhost:5050/api/mock/key/d784f08a-b732-40f9-aa11-6288ee733dea

# Response
[No content]
```


### Get all existing Mocks

This call returns a list of all existing Mocks in the database. If there are no Mocks
you'll receive an empty list.

```bash
# Call via cURL
curl --request GET \
  --url http://localhost:5050/api/mock/all

# Response
[
  {
    "Id": 2,
    "Key": "740bdb5b-c7bc-4e2f-8dc0-3378941f2be7",
    "Content": "Mock 1"
  },
  {
    "Id": 3,
    "Key": "b977d87e-fe71-4f19-a6f9-ccb9d5da198b",
    "Content": "Mock 2"
  }
]
```

### Get raw content of Mock

Just returns the `Content` Property of a Mock. It's identical to the
value you have set as you created the Mock.

```bash
# Call via cURL
curl --request GET \
  --url http://localhost:5050/raw/740bdb5b-c7bc-4e2f-8dc0-3378941f2be7

# Response
Mock 1
```


## How to compile

The GO-Compiler supports multiple operating systems out of the box.
With `env` you can set custom environment variables to your build process
which will be used and resetted after ([See here for details.](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04)).

Here you can see some values as example to compile a platform specific
executable:

```
# Compile for Linux
env GOOS=linux GOARCH=amd64 go build -o test_linux

# Compile for Windows
env GOOS=windows GOARCH=amd64 go build -o test_windows
```

To add a version number you have to call `go build -ldflags "-X main.version=1.0"`.
A complete call would look like this:

```
# Compile for Linux
env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=1.0" -o test_linux

# Compile for Windows
env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=1.0" -o test_windows
```

**Hint:** The name `main.version` is based on the variable in the *main.go*-file.
With this flags other values could be set, too.


## License

This project is licensed under MIT. Please see the License-file in this project
for detailed information.

