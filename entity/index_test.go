package entity_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/brynbellomy/ginkgo-reporter"
)

func TestTexture(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithCustomReporters(t, "Entity Suite", []Reporter{
		&reporter.TerseReporter{Logger: &reporter.DefaultLogger{}},
	})
}
