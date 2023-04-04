FROM golang:latest AS stage1
WORKDIR /apps/
COPY ./api/go.mod ./
RUN go mod download
COPY ./api/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lnmp-cli-linux64 . && \
    CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o lnmp-cli-linux . && \
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o lnmp-cli-win64.exe . && \
    CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o lnmp-cli-win.exe . && \
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o lnmp-cli-darwin64


FROM node:latest
WORKDIR /apps/
COPY ./web/package.json ./
RUN npm install
COPY ./web/ ./
RUN npm run build

COPY --from=stage1 /apps/lnmp-cli* .
COPY ./common/ ./common
RUN chmod +x ./common/script/entrypoint.sh

ENTRYPOINT ["/bin/bash", "-c", "/apps/common/script/entrypoint.sh"]

# docker build --no-cache --tag docker-lnmp:latest .
# docker run --rm -v E:\\d:/home docker-lnmp:latest

