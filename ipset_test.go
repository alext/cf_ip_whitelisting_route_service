package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("IPSet", func() {

	Describe("parsing list of networks", func() {

		It("handles a list of IPv4 CIDRs", func() {
			input := []string{
				"10.0.1.0/24",
				"10.4.5.3/18",
				"192.168.45.0/28",
			}
			ipset, err := NewIPSet(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(ipset.networks).To(HaveLen(3))

			Expect(ipset.networks[0].String()).To(Equal("10.0.1.0/24"))
			Expect(ipset.networks[1].String()).To(Equal("10.4.0.0/18"))
			Expect(ipset.networks[2].String()).To(Equal("192.168.45.0/28"))
		})

		It("converts IPv4 addresses into /32 CIDRs", func() {
			input := []string{
				"10.0.1.0",
				"10.4.5.3",
				"192.168.45.0",
			}
			ipset, err := NewIPSet(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(ipset.networks).To(HaveLen(3))

			Expect(ipset.networks[0].String()).To(Equal("10.0.1.0/32"))
			Expect(ipset.networks[1].String()).To(Equal("10.4.5.3/32"))
			Expect(ipset.networks[2].String()).To(Equal("192.168.45.0/32"))
		})

		It("handles a list of IPv6 CIDRs", func() {
			input := []string{
				"2001:db8:a0b:12f0::1/64",
				"2001:db8:e5b3::1/48",
			}
			ipset, err := NewIPSet(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(ipset.networks).To(HaveLen(2))

			Expect(ipset.networks[0].String()).To(Equal("2001:db8:a0b:12f0::/64"))
			Expect(ipset.networks[1].String()).To(Equal("2001:db8:e5b3::/48"))
		})

		It("converts IPv6 addresses to /128 CIDRs", func() {
			input := []string{
				"2001:db8:a0b:12f0::4",
				"2001:db8:e5b3::1",
			}
			ipset, err := NewIPSet(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(ipset.networks).To(HaveLen(2))

			Expect(ipset.networks[0].String()).To(Equal("2001:db8:a0b:12f0::4/128"))
			Expect(ipset.networks[1].String()).To(Equal("2001:db8:e5b3::1/128"))
		})

		It("handles a list with a mixture of formats", func() {
			input := []string{
				"10.0.1.0/24",
				"10.4.5.3",
				"2001:db8:e5b3::1/48",
				"2001:db8:a0b:12f0::4",
			}
			ipset, err := NewIPSet(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(ipset.networks).To(HaveLen(4))

			Expect(ipset.networks[0].String()).To(Equal("10.0.1.0/24"))
			Expect(ipset.networks[1].String()).To(Equal("10.4.5.3/32"))
			Expect(ipset.networks[2].String()).To(Equal("2001:db8:e5b3::/48"))
			Expect(ipset.networks[3].String()).To(Equal("2001:db8:a0b:12f0::4/128"))
		})

		It("errors with an invalid entry", func() {
			input := []string{
				"10.0.1.0/24",
				"not-an-ip",
				"10.4.5.3",
			}
			_, err := NewIPSet(input)
			Expect(err).To(HaveOccurred())
		})
	})

	DescribeTable("testing an IP",
		func(ip string, expected bool) {
			ipset, err := NewIPSet([]string{
				"10.0.1.0/24",
				"10.4.5.3",
				"2001:db8:e5b3::1/48",
				"2001:db8:a0b:12f0::4",
			})
			Expect(err).NotTo(HaveOccurred())

			actual := ipset.Contains(ip)
			Expect(actual).To(Equal(expected))
		},
		Entry("a whitelisted IPv4", "10.0.1.54", true),
		Entry("a non-whitelisted IPv4", "10.4.5.4", false),
		Entry("a whitelisted IPv6", "2001:db8:e5b3::123:1", true),
		Entry("a non-whitelisted IPv6", "2001:db8:a0b:12f0::3:4", false),
		Entry("an invalid IP", "not-an-ip", false),
	)
})
