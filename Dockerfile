FROM public.ecr.aws/bitnami/golang:1.20 as build-image

WORKDIR /go/src
COPY go.mod go.sum cmd/aws-watch/*.go ./

RUN go build -o /aws-watch

FROM public.ecr.aws/lambda/go:1

COPY --from=build-image /aws-watch /var/task/.

# Command can be overwritten by providing a different command in the template directly.
CMD ["aws-watch"]
