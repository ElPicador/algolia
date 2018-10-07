NAME=api

all: $(NAME)

$(NAME):
	@mkdir -p ./bin || echo ""
	go build -o ./bin/$(NAME) ./cmd/api

clean:
	rm -rf ./bin/$(NAME)

test:
	go test -race -count=1 ./...
