FROM golang:1.6

EXPOSE 9110

RUN go get github.com/mesosphere/mesos-exporter

ENTRYPOINT [ "/go/bin/mesos-exporter" ]
