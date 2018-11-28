package e2e

import (
	"bytes"
	"encoding/json" // . "github.com/marcusyip/golang-wire-mongo/web"
	"net/http"
	"net/http/httptest"

	"github.com/marcusyip/golang-wire-mongo/test/factories"
	"github.com/marcusyip/golang-wire-mongo/test/options"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func parseBody(body *bytes.Buffer) map[string]interface{} {
	var b map[string]interface{}
	err := json.Unmarshal(body.Bytes(), &b)
	if err != nil {
		panic(err)
	}
	return b
}

var _ = Describe("Register API", func() {
	Context("Valid input", func() {
		It("response the registered account", func() {
			body := jsonBody(map[string]interface{}{
				"username":     "gary123",
				"password":     "thisispassword123",
				"email":        "gary@sample.com",
				"display_name": "Gary",
			})
			req, _ := http.NewRequest("POST", "/api/v1/register", body)
			w := httptest.NewRecorder()
			ctn.Engine.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))

			b := parseBody(w.Body)
			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(b["object"]).To(Equal("account"))
			Expect(len(b["id"].(string))).To(Equal(24))
			Expect(b["email"]).To(Equal("gary@sample.com"))
			Expect(b["username"]).To(Equal("gary123"))
			Expect(b["display_name"]).To(Equal("Gary"))
		})
	})

	Context("Username already exists", func() {
		BeforeEach(func() {
			_ = ctn.UserFactory.Create(
				options.Field{"username", "alex123"},
				factories.UserWithAccess,
			)
		})
		It("response error object", func() {
			body := jsonBody(map[string]interface{}{
				"username":     "alex123",
				"password":     "thisispassword123",
				"email":        "alex123@sample.com",
				"display_name": "Gary",
			})
			req, _ := http.NewRequest("POST", "/api/v1/register", body)
			w := httptest.NewRecorder()

			ctn.Engine.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			b := parseBody(w.Body)
			Expect(b["object"]).To(Equal("error"))
			Expect(b["code"]).To(Equal(float64(10102)))
			Expect(b["message"]).To(Equal("Username already exists"))
		})
	})

	Context("Email already exists", func() {
		BeforeEach(func() {
			_ = ctn.UserFactory.Create(
				options.Field{"email", "lolo@sample.com"},
				factories.UserWithAccess,
			)
		})
		It("response error object", func() {
			body := jsonBody(map[string]interface{}{
				"username":     "lolo123",
				"password":     "thisispassword123",
				"email":        "lolo@sample.com",
				"display_name": "Gary",
			})
			req, _ := http.NewRequest("POST", "/api/v1/register", body)
			w := httptest.NewRecorder()

			ctn.Engine.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))

			b := parseBody(w.Body)
			Expect(b["object"]).To(Equal("error"))
			Expect(b["code"]).To(Equal(float64(10103)))
			Expect(b["message"]).To(Equal("Email already exists"))
		})
	})
})
