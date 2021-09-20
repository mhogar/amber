FROM golang:latest

# set working directory 
ENV APP_PATH /go/src/app
WORKDIR ${APP_PATH}

# copy source files
COPY . ${APP_PATH}

# build the application
RUN go build -v

# run the server
CMD ${APP_PATH}/amber
