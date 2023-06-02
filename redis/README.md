# Running the example (Pub/Sub with Redis)

## Run Redis

First of all, you need Redis running on your computer, you can install Redis following these instructions [Redis: Getting Started](https://redis.io/docs/getting-started/) or like I prefer using Docker.

If you want to run Redis from Docker (I believe this is the best option). If you have Docker installed and the daemon running.
Only need this command:

```bash
docker run -d --name redis-stack-server -p 6379:6379 redis/redis-stack-server:latest
```
Or you can follow the docs: [Redis Docker](https://redis.io/docs/stack/get-started/install/docker/)

## Run our server and clients

After we have Redis running and listening in the `PORT 6379`, only need to ran the server and the clients. (NOTE: if you never run `go mod tidy` in this project, do that first.)

Inside the server folder run this command:

```bash
go run ./main.go
```
That will start our server and will be listening in the `PORT 3000`

After that, you can run in multiple terminals a couple of clients running the same command:

```bash
go run ./main.go
```

## Test all, sending some messages

You can test the API using some program like Postman, here is an example using curl from the terminal:

```bash
curl -X POST http://localhost:3000 -H 'Content-Type: application/json' -d '{"Name":"john", "Email":"hola@gmail.com"}'
```

If all is right, you have to see in every terminal the same message. I hope this is useful.
