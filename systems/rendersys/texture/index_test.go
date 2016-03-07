package texture_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/brynbellomy/ginkgo-reporter"
)

func TestTexture(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithCustomReporters(t, "Texture Suite", []Reporter{
		&reporter.TerseReporter{Logger: &reporter.DefaultLogger{}},
	})
}
