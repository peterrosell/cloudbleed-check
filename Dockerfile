FROM scratch
MAINTAINER Peter Rosell<peter.rosell@gmail.com>
ADD cloudbleed-check cloudbleed-check
ENTRYPOINT ["/cloudbleed-check", "c"]