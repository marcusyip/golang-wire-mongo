package oauth2

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOauth2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oauth2 Suite")
}
