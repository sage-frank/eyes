### go version required

> go>=1.20

### description

>
>  gin
> <br> session
> <br> jwt
>

### example

> * add-watermark-pdf
> * excel-export
> * asm

### depends

- elasticsearch
- mysql
- redis

### run && build

```
- Makefile 编译

launch:

- windows:

    docker compose -f docker-compose.yaml up -d
- linux:

    docker-compose -f docker-compose.yaml up -d
 - direct run:

    go run main.go

- build: 

    go build -ldflag='-w -s' -o eye .

```

