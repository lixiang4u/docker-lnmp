FROM golang:latest AS stage1
WORKDIR /apps/api/
COPY ./api/go.mod ./
COPY ./api/go.sum ./
RUN go mod download
COPY ./api/ ./
RUN go build -o lnmp-cli .


FROM node:latest
WORKDIR /apps/web/
COPY ./web/package.json ./
COPY ./web/package-lock.json ./
RUN npm install
COPY ./web/ ./
RUN npm run build

COPY --from=stage1 apps/api/lnmp-cli ./dist

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["/apps/web/dist/lnmp-cli","s"]



