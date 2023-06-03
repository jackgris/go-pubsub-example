# Run the example of Pub/Sub with Nats

## Run our nats server

We will user Docker, this is the official documentation: [Docker Tutorial](https://docs.nats.io/running-a-nats-service/nats_docker/nats-docker-tutorial) and [Nats Docker](https://docs.nats.io/running-a-nats-service/nats_docker)


Basically we need to pull the docker image with this command:

```bash
docker pull nats
```

I choose use JetStream, for that reason we need to add the flag -js at the end.

And after that can start our server:

```bash
docker run --rm -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats-server -ti nats:latest -js
```

If you want to see the interactions in the terminal where you run Docker, you should use this command inside the nats folder.

```bash
 docker run --rm -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats-server -v "$(pwd)"/nats-config:/nats-config -ti nats:latest -c /nats-config/nats.conf -js
```

You can test the server is running right with this:

```bash
telnet localhost 4222
```

The output should be like this:

```
Trying ::1...
Connected to localhost.
Escape character is '^]'.
INFO {"server_id":"NDP7NP2P2KADDDUUBUDG6VSSWKCW4IC5BQHAYVMLVAJEGZITE5XP7O5J","version":"2.0.0","proto":1,"go":"go1.11.10","host":"0.0.0.0","port":4222,"max_payload":1048576,"client_id":13249}
```

## Run our server (subcriber)

First you have to run the subscriber, this will create a topic in the Nats Server running in the Docker container, this are inside the server folder.
For run this example, in this folder you only have to run this command:

```bash
go run ./main.go
```

You don't have to see anything until the client publish the messages. After that you can see the messages in this server and you can see the logs in the Docker container.

## Run our clients (publisher)

This application is in the client folder. For run this example, in this folder you only have to run this command:

```bash
go run ./main.go
```

This will publish 500 messages into the Nats Server running in the Docker container, and our server application will receive and show all the messages.
