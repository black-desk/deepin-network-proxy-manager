package config_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	. "github.com/black-desk/deepin-network-proxy-manager/internal/config"
	. "github.com/black-desk/deepin-network-proxy-manager/internal/test/ginkgo-helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {
	ContextTable("load from valid configuration (%s)", func(path string) {
		var (
			err     error
			file    *os.File
			content []byte
		)

		BeforeEach(func() {
			if file, err = os.Open(path); err != nil {
				Fail(fmt.Sprintf("Failed to open configuration %s: %s", path, err.Error()))
			}

			if content, err = io.ReadAll(file); err != nil {
				Fail(fmt.Sprintf("Failed to read configuration %s: %s", path, err.Error()))
			}

			_, err = Load(content)
		})
		AfterEach(func() {
			file.Close()
		})
		It("should success.", func() {
			Expect(err).To(BeNil())
		})
	},
		ContextTableEntry("../../test/data/example_config.yaml"),
	)
})

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configuration Suite")
}
