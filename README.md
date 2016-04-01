[![GoDoc](https://godoc.org/github.com/cybozu-go/transocks?status.png)][godoc]
[![Build Status](https://travis-ci.org/cybozu-go/transocks.png)](https://travis-ci.org/cybozu-go/transocks)

transocks - a transparent SOCKS5/HTTP proxy
===========================================

**transocks** is a background service to redirect TCP connections
transparently to a SOCKS5 server or a HTTP proxy server like [Squid][].

Currently, transocks supports only Linux iptables with DNAT/REDIRECT target.

Features
--------

* IPv4 and IPv6

    Both IPv4 and IPv6 are supported.
    Note that `nf_conntrack_ipv4` or `nf_conntrack_ipv6` kernel modules
    must be loaded beforehand.

* SOCKS5 and HTTP proxy (CONNECT)

    We recommend using SOCKS5 server if available.
    Looking for a good SOCKS5 server?  Take a look at our [usocksd][]!

    HTTP proxies often prohibits CONNECT method to make connections
    to ports other than 443.  Make sure your HTTP proxy allows CONNECT
    to the ports you want.

* Library and executable

    transocks comes with a handy executable.
    You may use the library to create your own.

Usage
-----

`transocks [-h] [-f CONFIG]`

The default configuration file path is `/usr/local/etc/transocks.toml`.

`transocks` does not have *daemon* mode.  Use systemd or upstart to
run it on your background.

Install
-------

Use Go 1.5 or better.

```
go get github.com/cybozu-go/transocks/cmd/transocks
```

Configuration file format
-------------------------

`transocks.toml` is a [TOML][] file.

`listen` and `proxy_url` are mandatory.
Other items are optional.

```
# listening address of transocks.
listen = "localhost:1081"

proxy_url = "socks5://10.20.30.40:1080"  # for SOCKS5 server
#proxy_url = "http://10.20.30.40:3128"   # for HTTP proxy server

log_level = "info"
log_file = "/var/log/transocks.log"
```

Redirecting connections by iptables
-----------------------------------

Use DNAT or REDIRECT target in OUTPUT chain of the `nat` table.

Save the following example to a file, then execute:
`sudo iptables-restore < FILE`

```
*nat
:PREROUTING ACCEPT [0:0]
:INPUT ACCEPT [0:0]
:OUTPUT ACCEPT [0:0]
:POSTROUTING ACCEPT [0:0]
:TRANSOCKS - [0:0]
-A OUTPUT -p tcp -j TRANSOCKS
-A TRANSOCKS -d 0.0.0.0/8 -j RETURN
-A TRANSOCKS -d 10.0.0.0/8 -j RETURN
-A TRANSOCKS -d 127.0.0.0/8 -j RETURN
-A TRANSOCKS -d 169.254.0.0/16 -j RETURN
-A TRANSOCKS -d 172.16.0.0/12 -j RETURN
-A TRANSOCKS -d 192.168.0.0/16 -j RETURN
-A TRANSOCKS -d 224.0.0.0/4 -j RETURN
-A TRANSOCKS -d 240.0.0.0/4 -j RETURN
-A TRANSOCKS -p tcp -j REDIRECT --to-ports 1081
COMMIT
```

Use *ip6tables* to redirect IPv6 connections.

Library usage
-------------

Read [the documentation][godoc].

License
-------

[MIT](https://opensource.org/licenses/MIT)

Author
------

[@ymmt2005][]

[godoc]: https://godoc.org/github.com/cybozu-go/transocks
[Squid]: http://www.squid-cache.org/
[usocksd]: https://github.com/cybozu-go/usocksd
[TOML]: https://github.com/toml-lang/toml
[@ymmt2005]: https://github.com/ymmt2005
