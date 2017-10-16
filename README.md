# CloudFoundry IP Restricting Route Service

This is a Proof-of-Concept CloudFoundry app that implements a
[route-service](https://docs.cloudfoundry.org/services/route-services.html) to
add IP whitelisting to an application. It does this by comparing the address in
the `X-Forwarded-For` header with the whitelist.

This accepts a set of IP addresses or CIDRs as a comma-separated list in the
`WHITELIST_ADDRS` environment variable.

Set `TRUSTED_ROUTERS` to a comma-separated list of IPs of the routers in your
stack. When processing the `X-Forwarded-For` header, the last entry that
doesn't match one of these will be used as the client IP.

If your CF deployment has a self-signed SSL certificate, set the
`SKIP_SSL_VALIDATION` environment variable to avoid SSL errors when proxying to
the backend.
