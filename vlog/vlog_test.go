package vlog_test

import (
	. "github.com/baotingfang/go-pivnet-client/vlog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"log"
)

var _ = Describe("Vlog", func() {
	Context("Test log instance methods", func() {
		var (
			logger    *Logger
			outBuffer *gbytes.Buffer
			errBuffer *gbytes.Buffer
		)

		BeforeEach(func() {
			outBuffer = gbytes.NewBuffer()
			errBuffer = gbytes.NewBuffer()
			logger = &Logger{
				OutLogger: log.New(outBuffer, "", log.LstdFlags),
				ErrLogger: log.New(errBuffer, "", log.LstdFlags),
				Prefix:    "T1",
			}
		})
		It("Test Error in vlog", func() {
			logger.LogLevel = ErrorLevel

			logger.Error("test error: %s", "error1")
			logger.Warn("test warn: %s", "warn1")
			logger.Info("test info: %s", "info1")
			logger.Debug("test debug: %s", "debug1")

			Expect(errBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[ERROR\] test error: error1`))
			Expect(outBuffer).ShouldNot(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[WARNING\] test warn: warn1`))
			Expect(outBuffer).ShouldNot(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[INFO\] test info: info1`))
			Expect(outBuffer).ShouldNot(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[DEBUG\] test debug: debug1`))
		})
		It("Test Warn in vlog", func() {
			logger.LogLevel = WarnLevel

			logger.Error("test error: %s", "error1")
			logger.Warn("test warn: %s", "warn1")
			logger.Info("test info: %s", "info1")
			logger.Debug("test debug: %s", "debug1")

			Expect(errBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[ERROR\] test error: error1`))
			Expect(outBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[WARNING\] test warn: warn1`))
			Expect(outBuffer).ShouldNot(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[INFO\] test info: info1`))
			Expect(outBuffer).ShouldNot(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[DEBUG\] test debug: debug1`))
		})
		It("Test Info in vlog", func() {
			logger.LogLevel = InfoLevel

			logger.Error("test error: %s", "error1")
			logger.Warn("test warn: %s", "warn1")
			logger.Info("test info: %s", "info1")
			logger.Debug("test debug: %s", "debug1")

			Expect(errBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[ERROR\] test error: error1`))
			Expect(outBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[WARNING\] test warn: warn1`))
			Expect(outBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[INFO\] test info: info1`))
			Expect(outBuffer).ShouldNot(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[DEBUG\] test debug: debug1`))
		})
		It("Test Debug in vlog", func() {
			logger.LogLevel = DebugLevel
			logger.Error("test error: %s", "error1")
			logger.Warn("test warn: %s", "warn1")
			logger.Info("test info: %s", "info1")
			logger.Debug("test debug: %s", "debug1")

			Expect(errBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[ERROR\] test error: error1`))
			Expect(outBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[WARNING\] test warn: warn1`))
			Expect(outBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[INFO\] test info: info1`))
			Expect(outBuffer).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T1\]\[DEBUG\] test debug: debug1`))
		})
	})

	Context("Test init log", func() {
		It("init log", func() {
			InitLog("T2", DebugLevel)
			Expect(Log.LogLevel).To(Equal(DebugLevel))
			Expect(Log.Prefix).To(Equal("T2"))
		})
	})

	Context("Test vlog global functions", func() {
		It("Test global Log instance", func() {
			output := gbytes.NewBuffer()
			errOutput := gbytes.NewBuffer()

			InitLog("T3", DebugLevel)
			Log.ErrLogger = log.New(errOutput, "", log.LstdFlags)
			Log.OutLogger = log.New(output, "", log.LstdFlags)

			Log.Error("test error: %s", "error1")
			Log.Warn("test warn: %s", "warn1")
			Log.Info("test info: %s", "info1")
			Log.Debug("test debug: %s", "debug1")

			Expect(errOutput).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T3\]\[ERROR\] test error: error1`))
			Expect(output).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T3\]\[WARNING\] test warn: warn1`))
			Expect(output).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T3\]\[INFO\] test info: info1`))
			Expect(output).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T3\]\[DEBUG\] test debug: debug1`))
		})

		It("Test vlog package functions", func() {
			output := gbytes.NewBuffer()
			errOutput := gbytes.NewBuffer()

			InitLog("T4", DebugLevel)
			Log.ErrLogger = log.New(errOutput, "", log.LstdFlags)
			Log.OutLogger = log.New(output, "", log.LstdFlags)

			Error("test error: %s", "error1")
			Warn("test warn: %s", "warn1")
			Info("test info: %s", "info1")
			Debug("test debug: %s", "debug1")

			Expect(errOutput).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T4\]\[ERROR\] test error: error1`))
			Expect(output).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T4\]\[WARNING\] test warn: warn1`))
			Expect(output).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T4\]\[INFO\] test info: info1`))
			Expect(output).Should(gbytes.Say(`\d+/\d+/\d+ \d+:\d+:\d+ \[T4\]\[DEBUG\] test debug: debug1`))
		})
	})
})
