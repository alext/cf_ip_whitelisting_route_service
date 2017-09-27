# CloudFoundry IP Restricting Route Service

This is a Proof-of-Concept CloudFoundry app that implements a
[route-service](https://docs.cloudfoundry.org/services/route-services.html) to
add IP whitelisting to an application. It does this by comparing the address in
the `X-Forwarded-For` header with the whitelist.

This accepts a set of IP addresses or CIDRs as a comma-separated list in the
`WHITELIST_ADDRS` environment variable.

To control which entry in the `X-Forwarded-For` header to look at you can set
the `XFF_OFFSET` environment variable to a number indicating how many entries
to strip off the end before taking the IP. If unset, this defaults to 0.

If your CF deployment has a self-signed SSL certificate, set the
`SKIP_SSL_VALIDATION` environment variable to avoid SSL errors when proxying to
the backend.
