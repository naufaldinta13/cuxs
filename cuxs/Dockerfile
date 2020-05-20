FROM qasico/go:latest

# Apps Dir
RUN mkdir -p "$GOPATH/src/github.com/naufaldinta13/cuxs/cuxs"
ADD . "$GOPATH/src/github.com/naufaldinta13/cuxs/cuxs"
WORKDIR "$GOPATH/src/github.com/naufaldinta13/cuxs/cuxs"

# Installing Dependency
RUN glide update

# Entry Point
ENTRYPOINT "$GOPATH/src/github.com/naufaldinta13/cuxs/cuxs"

CMD ["sleep","3600"]
