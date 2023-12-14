package post_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPostService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Post Service Test Suite")
}
