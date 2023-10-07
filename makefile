build-docker-image:
	docker build -t jungle:latest . \
	&& docker tag jungle benpeng/jungle \
	&& docker push benpeng/jungle

mysql:
	docker run --name mysql8 --network jungle-network -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root456789 -d mysql

migrateup:
	migrate -path db/migrations -database "mysql://root:root456789@tcp(localhost:3306)/jungle?multiStatements=true" -verbose up 1

migratedown:
	migrate -path db/migrations -database "mysql://root:root456789@tcp(localhost:3306)/jungle?multiStatements=true" -verbose down 1

migratenew:
	migrate create -ext sql -dir db/migrations -seq init_schema

createdb:
	docker exec -it mysql8 mysql -u root -proot456789 -e "CREATE DATABASE jungle;"

sqlc:
	sqlc generate

clear_db_user:
	docker exec -i mysql8 mysql -u root -proot456789 < ./db/query/clearUser.sql

.PHONY: build-docker-image mysql migrateup migratedown sqlc clear_db_user createdb