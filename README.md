# fuzzy-user-api

A simple registration/login service example using go gorm for storage and jwt for authentication.

## How to run?

The project is equipped with `docker-compose.yaml` file together with all the needed configuration in `dev.env` in order to be started:

The following command will start the fuzzy serivice as well as the DB locally in docker detached mode:
```
    make compose
```
The service will be listening on address `localhost:9205`.

When done one can run the below command in order to stop the running containers: 
```
    make decompose
```

## How to use? 

The project includes a `json` file (`Fuzzy-User-API.postman_collection.json`) that can be imported in Postman in order to load all the requests needed to make calls to the service.

## What about tests?

Tests can be run with the following command:

```
    make tests
```

