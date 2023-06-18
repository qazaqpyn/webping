# webping
A program that checks the list of sites for availability.

Once a minute, it checks whether the sites from the list are available with worker pool in the background and locally notes the time of access to them.
Users have 3 query options (endpoints):
1. Get access time to a specific site.
2. Get the name of the site with the minimum access time.
3. Get the name of the site with the maximum access time.

Administrators get statistics on user requests for the three above endpoints.

## Technologies 
- JWT Authentication 
- Swagger
- Worker Pool
- Graceful Shutdown

## Before run 
1. start mongodb in docker 

```
Docker run -it --rm --name mongo_db -e MONGODB_DATABASE=main -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=qwerty -p 27019:27017 -d mongo
```

2. rename .env-example to .env

## To run 
go build -o app cmd/main.go && ./app


