package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Basic Auth proxy", func() {

	var (
		proxy    *AuthProxy
		backend  *ghttp.Server
		req      *http.Request
		response *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		ipset, err := NewIPSet([]string{"10.0.0.0/24"})
		Expect(err).NotTo(HaveOccurred())
		proxy = NewAuthProxy(ipset, 0)
		backend = ghttp.NewServer()
		backend.AllowUnhandledRequests = true
		backend.UnhandledRequestStatusCode = http.StatusOK
	})

	AfterEach(func() {
		backend.Close()
	})

	Context("with a request from route-services", func() {
		BeforeEach(func() {
			req = httptest.NewRequest("GET", "http://proxy.example.com/", nil)
			req.Header.Set("X-CF-Forwarded-Url", backend.URL())
			req.Header.Set("X-CF-Proxy-Signature", "Stub signature")
			req.Header.Set("X-CF-Proxy-Metadata", "Stub metadata")
		})

		JustBeforeEach(func() {
			response = httptest.NewRecorder()
			proxy.ServeHTTP(response, req)
		})

		Context("with an IP in the whitelist", func() {
			BeforeEach(func() {
				req.Header.Set("X-Forwarded-For", "10.0.0.3")
			})

			It("should proxy the request to the backend", func() {
				Expect(response.Code).To(Equal(http.StatusOK))

				Expect(backend.ReceivedRequests()).To(HaveLen(1))

				headers := backend.ReceivedRequests()[0].Header
				Expect(headers.Get("X-CF-Proxy-Signature")).To(Equal("Stub signature"))
				Expect(headers.Get("X-CF-Proxy-Metadata")).To(Equal("Stub metadata"))
			})

			It("preserves the Host header from the forwarded URL", func() {
				url, err := url.Parse(backend.URL())
				Expect(err).NotTo(HaveOccurred())

				beReq := backend.ReceivedRequests()[0]
				Expect(beReq.Host).To(Equal(url.Host))
			})

			Context("with a path and query in the forwarded URL", func() {
				BeforeEach(func() {
					req.Header.Set("X-CF-Forwarded-Url", backend.URL()+"/foo/bar?a=b")
				})
				It("preserves the path and query from the forwarded URL", func() {
					beReq := backend.ReceivedRequests()[0]

					Expect(beReq.URL.Path).To(Equal("/foo/bar"))
					Expect(beReq.URL.RawQuery).To(Equal("a=b"))
				})
			})
		})

		Context("with a non whitelisted IP", func() {
			BeforeEach(func() {
				req.Header.Set("X-Forwarded-For", "192.168.0.1")
			})

			It("returns a 403", func() {
				Expect(response.Code).To(Equal(403))
			})

			It("does not make a request to the backend", func() {
				Expect(backend.ReceivedRequests()).To(HaveLen(0))
			})
		})

		DescribeTable("authorization rules",
			func(xffHeader string, offset int, expectAuthorized bool) {
				req.Header.Set("X-Forwarded-For", xffHeader)
				proxy.xffOffset = offset

				response = httptest.NewRecorder()
				proxy.ServeHTTP(response, req)

				if expectAuthorized {
					Expect(response.Code).To(Equal(200))
				} else {
					Expect(response.Code).To(Equal(403))
				}
			},
			Entry("an IP in the whitelist", "10.0.0.1", 0, true),
			Entry("an IP in the whitelist with additional entries", "192.0.2.5, 10.0.0.1", 0, true),
			Entry("an IP not in the whitelist", "192.168.0.1", 0, false),
			Entry("an IP not in the whitelist with additional entries", "192.0.2.5, 192.168.0.1", 0, false),
			Entry("whitelisted IP in the wrong place in XFF", "10.0.0.1, 192.168.0.1", 0, false),
			Entry("whitelisted IP with an offset", "10.0.0.1, 192.168.0.1", 1, true),
			Entry("whitelisted IP with an offset and additional entries", "192.0.2.5, 10.0.0.1, 192.168.0.1", 1, true),
			Entry("offset beyond start of XFF header", "10.0.0.1", 1, false),
			Entry("offset well beyond start of XFF header", "10.0.0.1", 4, false),
			Entry("empty XFF header", "", 0, false),
			Entry("extra spaces in XFF header", "10.0.0.1,  192.168.0.1", 1, true),
		)
	})
})
