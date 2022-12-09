FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY parametros.txt .
RUN go mod download
COPY . .
RUN go BUILD -o ./out/resultado .
CMD ./out/resultado
