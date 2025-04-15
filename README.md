# gollet

Hello there... If you're reading this, then you will see how silly my attempt on tryna write a Golang REST Api. T-T

In case you're interested to run this project in your local machine. Here's a Step by step guide on how to do so.

## Step 1

Firstly, we need to make sure that we have already installed the correct version of the golang compiler in our local machine.
Assuming that you've already `git clone` this repo and you're already in the root path of the project.

Now, let's start by making sure we have the correct golang compiler. If you're using ASDF then, you can easily install them by doing so.

```
# This will install the golang plugin
asdf plugin-add golang

# install the golang version based on the .tool-versions
asdf install
```

## Step 2

Since, we are using docker to help make things easier for anyone to setup and run the server. Therefore, we can easily use the `docker-compose.yaml` script to run the server in docker instances.

We can simply setup and destroy the docker instance by running these command

```
# To create the docker instance
$ docker-compose up -d --build

# To destroy the docker instance
$ docker-compose down
```

After running the `docker-compose up -d --build` command. You will notice there will be a new instances created in your docker. Since I'm using OrbStack to manage my docker instances. Therefore, I added a screenshot on what you will see after we run the docker compose command.


[<img src="public/orbstack.png">](#)

The REST api are now expose on our localhost port 8080.

But, if you're more interested to run the server locally. You can do so easily by running this command in your CLI.

```
go ./cmd/gollet/main.go
```

This command will setup the database for you (Assuming you already have postgres installed) if it doesn't exist and run the migration. Then it will start the server at port 8080.

## Step 3

Now, after we have successfully run the server on our docker. We can now test the backend service by sending a request to these path

```
# For user signup and login
POST http://localhost:8080/api/users/signup

{
    "email": "email@example.com",
    "password": "password123"
}

POST http://localhost:8080/api/users/login

{
    "email": "email@example.com",
    "password": "password123"
}

# For wallet interactions
POST http://localhost:8080/api/users/:user_id/wallets/deposit

{
    "amount": 50.00,
    "currency": "USD"
}

POST http://localhost:8080/api/users/:user_id/wallets/withdraw

{
    "amount": 15.00,
    "currency": "USD"
}

POST http://localhost:8080/api/users/:user_id/wallets/transfer

{
    "amount": 5.00,
    "recipient_user_id": 3
}

GET http://localhost:8080/api/users/:user_id/wallets/balance

GET http://localhost:8080/api/users/:user_id/wallets/transactions/history
```

## Step 4

Now, to run the test for the services. We can simply do it my running

```
go test ./internal/services
```