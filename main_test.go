package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parseWhitelistAddrs", func() {

	It("splits the comma separated string", func() {
		actual := parseWhitelistAddrs("10.0.0.1,192.168.42.0/24")
		Expect(actual).To(Equal([]string{"10.0.0.1", "192.168.42.0/24"}))
	})

	It("strips leading and trailing whitespace in each entry", func() {
		actual := parseWhitelistAddrs("10.0.0.1, 192.168.42.0/24")
		Expect(actual).To(Equal([]string{"10.0.0.1", "192.168.42.0/24"}))
	})

	It("returns an empty slice when given blank input", func() {
		actual := parseWhitelistAddrs("")
		Expect(actual).To(Equal([]string{}))
	})
})
