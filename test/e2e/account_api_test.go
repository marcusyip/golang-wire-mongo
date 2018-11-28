package e2e

import (
	// . "github.com/marcusyip/golang-wire-mongo/web"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/marcusyip/golang-wire-mongo/models"
	"github.com/marcusyip/golang-wire-mongo/test/factories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account API", func() {
	var (
		ctxUser *models.User
	)

	BeforeEach(func() {
		ctxUser = ctn.UserFactory.Create(
			factories.UserWithAccess,
		)
	})

	var _ = Describe("Show account", func() {
		Context("Valid access token", func() {
			It("response my own account", func() {
				req, _ := http.NewRequest("GET", "/api/v1/account", nil)
				AddAccessToken(req, ctxUser)
				w := httptest.NewRecorder()

				ctn.Engine.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var b map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &b)
				Expect(b["object"]).To(Equal("account"))
				Expect(b["id"]).To(Equal(ctxUser.ID.Hex()))
				Expect(b["email"]).To(Equal(ctxUser.Email))
				Expect(b["username"]).To(Equal(ctxUser.Username))
				Expect(b["display_name"]).To(Equal(ctxUser.DisplayName))
			})
		})

		Context("Missing access token", func() {
			It("response unauthorize", func() {
				req, _ := http.NewRequest("GET", "/api/v1/account", nil)
				w := httptest.NewRecorder()

				ctn.Engine.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))

				var b map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &b)
				Expect(b["object"]).To(Equal("error"))
				Expect(b["code"]).To(Equal(float64(10101)))
				Expect(b["message"]).To(Equal("Unauthorized"))
				// Expect(b["message"]).To(Equal("Missing access token"))
			})
		})
	})
})
