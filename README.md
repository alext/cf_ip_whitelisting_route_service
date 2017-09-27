# CloudFoundry IP Restricting Route Service

This is a Proof-of-Concept CloudFoundry app that implements a
[route-service](https://docs.cloudfoundry.org/services/route-services.html) to
add IP whitelisting to an application.

If your CF deployment has a self-signed SSL certificate, set the
`SKIP_SSL_VALIDATION` environment variable to avoid SSL errors when proxying to
the backend.
