
all: export CREATE_DB=$(shell cat database/create_db)
all:
		@go build -o pass_manager -ldflags '-X "main.db_CREATE_COMMAND=$(CREATE_DB)"'

dependencies:
		go get -u github.com/lib/pq
