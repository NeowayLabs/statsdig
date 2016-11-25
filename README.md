# StatsDig

Yeah, it is a pretty cheesy name for the idea of mixing up
the StatsD protocol and Sysdig Cloud services.

So why another statsd Golang client ? Well, actually there
is not a lot of reason, but there is some.

Depending on how you send UDP packets on Go it will work perfectly
fine when there is someone listening on the port where the packets
are being sent, but when they are not, your metrics won't be
collected correctly by the sysdig agent (because Go will not even send them).

Why would I send packets to ports where no one is listening ?
Well Sysdig is that magic, it will use kernel introspection to
collect the UDP packets independent from where they are being sent.
It is a very common use case to just sent to localhost on the
default port.

We optimize for this use case and make things as simple as possible.
Also we added support to Sysdig extension of tags as a first class
citizen on the API.

Also we add here some tools that helped us to debug problems like
metrics disappearing (usually our fault), etc.

More on Sysdig StatsD magic:

* [Metrics Integration: StatsD](https://support.sysdigcloud.com/hc/en-us/articles/204376099-Metrics-integrations-StatsD)
* [StatsD Teleportation](https://support.sysdigcloud.com/hc/en-us/articles/204470339)

Sysdig magic has a know limitation that required the udp to be under IPv4,
so the library is hardcoded to always use **udp4** as the protocol.


## API Reference

[API Reference docs](https://godoc.org/github.com/NeowayLabs/statsdig)


## Debugging Metric Collection

To make it easier to run on a cluster for tests we added a
docker image with all the commands embedded on them.

You can just run:

```
docker run neowaylabs/statsdig:latest ./sender
```

To send a lot of default test metrics to your sysdig agent
(of course you must run this on a host that has the agent running).

If you are feeling hardcore, just build the project, it is
plain Go code with no dependencies besides Go itself:

```
go run cmd/sender/main.go
```

If you are using Kubernetes you can just run:

```
kubectl apply -f tests/k8s/sender.yml
```

Using our test example that will run a job on your cluster that will
generate metrics for you.

In the past we had problems running directly on host X on docker X on
Kubernetes, that is why we have all these debugging options.
