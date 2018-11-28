package e2e

import (
	"net/http"
	"net/http/httptest"

	"github.com/marcusyip/golang-wire-mongo/models"
	"github.com/marcusyip/golang-wire-mongo/test/options"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Access API", func() {
	var (
		ctxUser *models.User
	)

	Context("Works", func() {
		BeforeEach(func() {
			ctxUser = ctn.UserFactory.Create(
				options.Field{"password", "thispassword123"},
			)
		})
		It("response the access entity with account", func() {
			body := jsonBody(map[string]interface{}{
				"username": ctxUser.Username,
				"password": "thispassword123",
			})
			req, _ := http.NewRequest("POST", "/api/v1/access", body)
			w := httptest.NewRecorder()
			ctn.Engine.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
			b := parseBody(w.Body)
			Expect(b["object"]).To(Equal("access"))
			Expect(len(b["access_token"].(string))).To(Equal(48))
		})
	})

	Context("Invalid password", func() {
		BeforeEach(func() {
			ctxUser = ctn.UserFactory.Create(
				options.Field{"password", "thispassword123"},
			)
		})
		It("response the access entity with account", func() {
			body := jsonBody(map[string]interface{}{
				"username": ctxUser.Username,
				"password": "wrongpassword",
			})
			req, _ := http.NewRequest("POST", "/api/v1/access", body)
			w := httptest.NewRecorder()
			ctn.Engine.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
			b := parseBody(w.Body)
			Expect(b["object"]).To(Equal("error"))
			Expect(b["message"]).To(Equal("Invalid credential"))
		})
	})

	Context("Not exists username", func() {
		It("response the access entity with account", func() {
			body := jsonBody(map[string]interface{}{
				"username": "notexists",
				"password": "wrongpassword",
			})
			req, _ := http.NewRequest("POST", "/api/v1/access", body)
			w := httptest.NewRecorder()
			ctn.Engine.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
			b := parseBody(w.Body)
			Expect(b["object"]).To(Equal("error"))
			Expect(b["message"]).To(Equal("Invalid credential"))
		})
	})
})
