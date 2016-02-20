//go:generate hel

package testhelpers_test

import (
	"github.com/apoydence/eachers"
	"github.com/apoydence/eachers/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AlwaysReturn", func() {
	var (
		args []interface{}
	)

	Context("given a struct with channels", func() {
		type sampleStruct struct {
			A chan int
			B chan string
		}

		var (
			receiver sampleStruct
		)

		BeforeEach(func() {
			receiver = sampleStruct{
				A: make(chan int, 100),
				B: make(chan string, 100),
			}
		})

		Context("proper args", func() {
			BeforeEach(func() {
				args = eachers.With(99, "some-arg")
			})

			It("keeps all the channels populated with the given arguments", func() {
				testhelpers.AlwaysReturn(receiver, args...)

				Eventually(receiver.A).Should(HaveLen(cap(receiver.A)))
				Eventually(receiver.B).Should(HaveLen(cap(receiver.B)))
			})

			It("sends the expected argument", func() {
				testhelpers.AlwaysReturn(receiver, args...)

				Eventually(receiver.A).Should(Receive(Equal(args[0])))
				Eventually(receiver.B).Should(Receive(Equal(args[1])))
			})
		})

		Context("invalid args", func() {
			It("not enough arguments", func() {
				f := func() {
					testhelpers.AlwaysReturn(receiver, args[:1])
				}
				Expect(f).To(Panic())
			})

			It("wrong argument type", func() {
				f := func() {
					testhelpers.AlwaysReturn(receiver, 99, 100)
				}
				Expect(f).To(Panic())
			})
		})
	})

	Context("given a channel", func() {
		var (
			channel chan int
		)

		BeforeEach(func() {
			channel = make(chan int, 100)
		})

		Context("with arguments", func() {

			BeforeEach(func() {
				args = eachers.With(99)
			})

			It("keeps the channel populated with the given argument", func() {
				testhelpers.AlwaysReturn(channel, args...)

				Eventually(channel).Should(HaveLen(cap(channel)))
			})

			It("sends the expected argument", func() {
				testhelpers.AlwaysReturn(channel, args...)

				Eventually(channel).Should(Receive(Equal(args[0])))
			})
		})

		Context("invalid arguments", func() {
			It("panics for no arguments", func() {
				f := func() {
					testhelpers.AlwaysReturn(channel)
				}
				Expect(f).To(Panic())
			})

			It("panics for more than one argument", func() {
				f := func() {
					testhelpers.AlwaysReturn(channel, 1, 2)
				}
				Expect(f).To(Panic())
			})

			It("panics if the argument type and channel type don't match", func() {
				f := func() {
					testhelpers.AlwaysReturn(channel, "invalid")
				}
				Expect(f).To(Panic())
			})
		})
	})

	Context("given something other than a channel or struct full of channels", func() {
		BeforeEach(func() {
			args = eachers.With(1, 2, 3)
		})

		It("panics for non channel or struct", func() {
			f := func() {
				testhelpers.AlwaysReturn(8, args)
			}

			Expect(f).To(Panic())
		})

		It("panics for read only channel", func() {
			var invalid <-chan int
			f := func() {
				testhelpers.AlwaysReturn(invalid, args)
			}

			Expect(f).To(Panic())
		})
	})
})
