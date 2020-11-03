# storemanager

## requirements 
- go1.15

### Steps

Start db and table creation
```
docker compose up -d db
```

Set env
```
export CONFIG_PATH=/home/arjun/ARJUN/projects/storemanager/config/local
```

Migrate Store data
```
go run dbmigration/store/main.go
```

Start App
```
go run cmd/app/main.go
```

[Swagger Doc](https://github.com/arjunksofficial/storemanager/blob/main/api/swagger.json)

[Postman Collection](https://github.com/arjunksofficial/storemanager/blob/main/api/SM.postman_collection.json)