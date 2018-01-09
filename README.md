![https://github.com/monstar-lab/fr-circle-api/](https://cdn-images-1.medium.com/max/600/1*bnyJ9a-DxAKV-qJXKDAZkQ.png)

# Golang docker Env (Minio, Postgres, Redis) 
In this setup we will have five docker containers as follows.

* 1 Go container : This is the main container which will host our Go app.
* 1 Minio container : Object storage server with Amazon S3 compatible API.
* 1 Postgres container : Postgres DB server.
* 1 Redis container : Key-value store that functions as a data structure server.
* 1 Nginx container : This will be used as a reverse proxy server.

## Prerequisites

You have to have docker, docker composer and Go installed on your local machine.

[Docker-composer installation guide](https://docs.docker.com/compose/install/)

[Go installation guide](https://golang.org/dl/)

And also you have to add GOPATH to your local system path. 

* You can confirm the GOPATH using `go env GOPATH` command.
* [Check this to learn how to add new path to your system path](https://github.com/golang/go/wiki/SettingGOPATH)


## Download, Build, Run, Code and Test 

### Download and build

On your local machine, clone this repository and build docker images issuing following commands.

```
go get github.com/monstar-lab/fr-circle-api 
cd $(go env GOPATH)/src/github.com/monstar-lab/fr-circle-api
docker-compose build
```

It will take some time to complete these three commands as it has to pull all the images needed at first run. 

### Run

#### Foreground mode

```
docker-compose up
```

By Issuing above command, we can have our containers up and running in the foreground. Hence we can see all the logs in our terminal window. 

*(You will have to kill any applications which is using port 80.)*

We can use `CTRL + C` to stop containers and get the access back to our terminal. 

#### Background mode

```
docker-compose up -d
```

Using the flag -d, we can have our containers running in the background. However we cannot see logs in our terminal window.


In either mode we can now access our go application on [http://localhost/](http://localhost/)

#### Stop and remove

Issue the following command to stop and remove all running containers.

```
docker-compose down
``` 

### Code

When we start up our containers, our go application will start to serve and be watching for any file change.  

Each time a go source file is changed, Our app will incrementally be rebuilt. So we don't need to worry about building our app again and again. 

We can just refresh our browser to see the changes we make.

#### Editor

We can use any editor or IDE as our wish. However we must set up [gometalinter](https://github.com/alecthomas/gometalinter) and any editor related plugins that makes go coding fast and easy.  

Some of the common editors are ....

* VScode 
* GoLand
* Vim
* Emacs
* Atom

#### Connecting to other containers from go container

Since go container is linked to DB, Redis and Minio containers, We can connect to those service from our go source as follows.

##### DB

	- DBMS : postgres
	- Host : db
	- Port : 5432 
	- User : postgres
	- Password : mypass
	- DBname : sample

##### Redis

	- Host : redis
	- Port : 6379 

##### Minio

	- Endpoint : s3:9000
	- Access key : AKIAIOSFODNN7EXAMPLE 	
	- Secret Key : wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY



#### Add new library
If we want to import a new library to our code, 
we can issue a command as follows while containers are running.
It will download the library to the project. 

```
docker exec -it go-cont dep ensure -add github.com/foo/bar
```
*(If you have go dep installed locally, you can do the same by `dep ensure -add github.com/foo/bar`)* 

then we can import that library in our source code like `import "github.com/foo/bar"`

*** ***Do not use `go get` in this environment to import libraries.***

### Test

While the docker containers are up and running, use the following command to run all the go test files. 

```
docker exec -it go-cont go test ./...
```

If docker containers are not running use the following command. (This will start up our containers by itself to run the test.)

```
docker-compose run go go test ./...
```

## Container Names

* Go container : go-cont
* Minio container : storage-cont
* Postgres container : db-cont
* Redis container : kvs-cont
* Nginx container : proxy-cont

## Useful commands

### Access shell inside a container

Issue the Following command while the container is running.

Replace go-cont with any of the container name above.

```
docker exec -it go-cont sh
```
### Ensure dependencies
To analyze and download missing dependencies issue the following command while containers are running.

```
docker exec -it go-cont dep ensure
```

*(If you have go dep installed locally, you can do the same by `dep ensure`)*  
 
## About Files

### `Dockerfile`

This file is responsible to create a new docker image by copying the Go source in `go-app` directory, build and run it in the
container. 

### `docker-compose.yml`

This file is useful for defining and running multi-container Docker applications like ours. We use this file to configure our application’s services. Then, by `docker-compose build` command, we can create and start all the services. 

### `Gopkg.toml`, `Gopkg.lock`
These files are used by [dep](https://github.com/golang/dep/) to track and manage dependencies for our project.

### `.nginx.conf`
This is a basic configuration file grabbed from the stock nginx docker image. 
Server block is added to this file to get web requests (on ports 80,443) forwarded to the Go container. 

### Simple go code samples to undestand how to use this setup

### `main.go`

This runs a http server that listens on port `:8080` inside go container, so that we can access it via web browser on `http://localhost:8080`. Since this is the starting point of our application, this file is responsible to route a url to the correct httphandler function.

*we can access it via port `http://localhost:80` as well since the proxy server is configured to forward all the requests to port 80 and 443 to the go container.* 

### `redis.go`

This file includes a sample code to understand how to deal with redis storage in this environment. 
[Redigo](https://github.com/garyburd/redigo/) is used in this sample code, 
  which is a Go client for the Redis database.


### `db.go`

This file includes a sample code to understand how to connect to postgres database in this environment. 

This file is responsible for the followings. 

* `http://localhost/db/`: lists all the records in a table called dummytable.
* `http://localhost/db/add/{some text}`: inserts a new records to a table called dummytable with the value of {some text}.
 
### `s3.go`

This file includes a sample code to understand how to deal with minio object storage in this environment. 

This file is responsible for the followings. 

* `http://localhost/s3/`: lists all the objects in a bucket called testbucket.
* `http://localhost/s3/triggeraput/`: uploads the `s3_upload_test_file.txt` to the testbucket.

We can access the minio service from the minio server directly via `http://localhost:9000` and check our objects.

Read [official minio guide](https://docs.minio.io/docs/golang-client-quickstart-guide
) to learn more. 

### `main_test.go`

This file includes a simple test. 

This test will run when we issue the `docker-compose run go go test ./...` command which is mentioned in the test section above.


### Happy coding...!
