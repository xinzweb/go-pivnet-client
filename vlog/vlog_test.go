package vlog_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"log"

	. "github.com/baotingfang/go-pivnet-client/vlog"
)

var _ = Describe("Vlog", func() {
	var (
		logger    *Logger
		outBuffer *gbytes.Buffer
		errBuffer *gbytes.Buffer
	)

	BeforeEach(func() {
		outBuffer = gbytes.NewBuffer()
		errBuffer = gbytes.NewBuffer()
		logger = &Logger{
			OutLogger: log.New(outBuffer, "T1 ", log.LstdFlags),
			ErrLogger: log.New(errBuffer, "T1 ", log.LstdFlags),
		}
	})
	It("Test Error in vlog", func() {
		logger.LogLevel = ErrorLevel
		logger.Error("test error: %s", "error1")
		Expect(errBuffer).Should(gbytes.Say(`^T1 .+ \[ERROR\] test error: error1`))
	})
	It("Test Info in vlog", func() {
		logger.LogLevel = InfoLevel
		logger.Info("test info: %s", "info1")
		Expect(outBuffer).Should(gbytes.Say(`^T1 .+ \[INFO\] test info: info1`))
	})
	It("Test Debug in vlog", func() {
		logger.LogLevel = DebugLevel
		logger.Debug("test debug: %s", "debug1")
		Expect(outBuffer).Should(gbytes.Say(`^T1 .+ \[DEBUG\] test debug: debug1`))
	})
	It("Test Warn in vlog", func() {
		logger.LogLevel = WarnLevel
		logger.Warn("test warn: %s", "warn1")
		Expect(outBuffer).Should(gbytes.Say(`^T1 .+ \[WARNING\] test warn: warn1`))
	})
})
