# build stage
FROM golang:alpine AS build-env
WORKDIR go/src/github.com/vbansal/travel_mate_service/
ADD . .
RUN cd helpers && go install && cd ../services && go install && cd .. && go build -o travel_mate_app
#RUN cd /src && go build -o travel_mate_app

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/travel_mate_app /app/
ENTRYPOINT [ "travel_mate_app" ]


