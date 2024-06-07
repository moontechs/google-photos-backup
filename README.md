* `docker build -t installer -f ./database/Dockerfile ./database`
* `docker run -it --rm --network=host -v $(pwd)/database:/data --env-file=.env.local installer task --dir /data install`
* `curlie -X POST -u "backend:change123" -H "NS: backup" -H "DB: backup" -H "Accept: application/json" http://localhost:8000/sql -d "select * from users;"`
* `docker run -it --rm --network=host -v $(pwd)/database:/data --env-file=.env.local installer task --dir /data create-user email=123@123.com password=change123`