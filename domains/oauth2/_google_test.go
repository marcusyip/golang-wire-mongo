package oauth2

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	people "google.golang.org/api/people/v1"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"project.scmp.tech/technology/newsroom-system/accounts/core/config"
	"project.scmp.tech/technology/newsroom-system/accounts/core/logger"
	acctmgr "project.scmp.tech/technology/newsroom-system/accounts/managers/account"
	"project.scmp.tech/technology/newsroom-system/accounts/mocks/mock_managers/mock_account"
	mock_repos "project.scmp.tech/technology/newsroom-system/accounts/mocks/mock_repositories"
	"project.scmp.tech/technology/newsroom-system/accounts/models"
	muser "project.scmp.tech/technology/newsroom-system/accounts/models/user"
)

var _ = Describe("GoogleOAuth2", func() {
	var (
		mockCtrl     *gomock.Controller
		mockUserRepo *mock_repos.MockUserRepository
		mockAcctMgr  *mock_account.MockManager
		googleOAuth2 *GoogleOAuth2
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserRepo = mock_repos.NewMockUserRepository(mockCtrl)
		mockAcctMgr = mock_account.NewMockManager(mockCtrl)
		googleOAuth2 = &GoogleOAuth2{
			logger:     logger.Get(),
			googleConf: config.GoogleConfig{},
			userRepo:   mockUserRepo,
			acctMgr:    mockAcctMgr,
		}
	})

	Context("when user with google ID exists", func() {
		It("should update and return the user", func() {
			mockUser := &models.User{
				ModelImpl: models.ModelImpl{
					ID: bson.ObjectIdHex("507f1f77bcf86cd799439011"),
				},
				OAuth2Token: &oauth2.Token{},
			}
			mockUserRepo.EXPECT().FindByGoogleID("some_google_id").
				Times(1).
				Return(mockUser, nil)
			mockUserRepo.EXPECT().UpdateByID(bson.ObjectIdHex("507f1f77bcf86cd799439011"), gomock.Any()).
				Times(1).
				Return(nil).
				Do(func(id bson.ObjectId, fields bson.M) {
					Expect(fields).Should(Equal(bson.M{
						"$set": bson.M{
							"oauth2_token": &oauth2.Token{},
							"display_name": "Tom Chan",
							"avatar_url":   "some_url",
							"department":   "Technology",
							"title":        "Web Developer",
						},
					}))
				})
			params := &AuthorizeParams{
				AuthCode: "some_code",
			}
			exchangeToken := func(authCode string) (*oauth2.Token, error) {
				return &oauth2.Token{}, nil
			}
			getPersonCall := func(token *oauth2.Token) (*people.Person, error) {
				Expect(token).NotTo(BeNil())

				person := &people.Person{
					Metadata: &people.PersonMetadata{
						Sources: []*people.Source{
							{
								Type: "PROFILE",
								Id:   "some_google_id",
							},
						},
					},
					Names: []*people.Name{
						{
							DisplayName: "Tom Chan",
						},
					},
					Organizations: []*people.Organization{
						{
							Metadata:   &people.FieldMetadata{Primary: false, Verified: true},
							Title:      "Internship",
							Department: "Product",
						},
						{
							Metadata:   &people.FieldMetadata{Primary: true, Verified: true},
							Title:      "Web Developer",
							Department: "Technology",
						},
						{
							Metadata:   &people.FieldMetadata{Primary: false, Verified: true},
							Title:      "Random Developer",
							Department: "Product",
						},
					},
					Photos: []*people.Photo{
						{
							Url: "some_url",
						},
					},
					EmailAddresses: []*people.EmailAddress{
						{
							Value: "tom.chan@scmp.com",
						},
					},
				}
				return person, nil
			}
			// mockUserRepo.EXPECT().FindByGoogleID(
			mockAcctMgr.EXPECT().ValidateEmail("tom.chan@scmp.com").
				Return(nil).
				Times(1)

			user, err := googleOAuth2.doAuthorize(params, exchangeToken, getPersonCall)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.ID).To(Equal(bson.ObjectIdHex("507f1f77bcf86cd799439011")))
			Expect(user.DisplayName).To(Equal("Tom Chan"))
			Expect(user.Department).To(Equal("Technology"))
			Expect(user.Title).To(Equal("Web Developer"))
			Expect(user.AvatarURL).To(Equal("some_url"))
		})
	})

	Context("when user with google ID not exists", func() {
		It("should create and return the user", func() {
			mockUser := &models.User{
				ModelImpl: models.ModelImpl{
					ID: bson.ObjectIdHex("507f1f77bcf86cd799439011"),
				},
				OAuth2Token: &oauth2.Token{},
				GoogleID:    "some_google_id",
				DisplayName: "Tom Chan",
				Department:  "Technology",
				Title:       "Web Developer",
				AvatarURL:   "some_url",
				Type:        muser.TypeUser,
				Email:       "tom.chan@scmp.com",
			}
			mockUserRepo.EXPECT().FindByGoogleID("some_google_id").
				Times(1).
				Return(nil, mgo.ErrNotFound)
			mockUserRepo.EXPECT().Create(gomock.Any()).
				Times(1).
				Return(mockUser, nil).
				Do(func(m models.Model) (models.Model, error) {
					user := m.(*models.User)
					Expect(user.GoogleID).Should(Equal("some_google_id"))
					return m, nil
				})
			params := &AuthorizeParams{
				AuthCode: "some_code",
			}
			exchangeToken := func(authCode string) (*oauth2.Token, error) {
				return &oauth2.Token{}, nil
			}
			getPersonCall := func(token *oauth2.Token) (*people.Person, error) {
				Expect(token).NotTo(BeNil())

				person := &people.Person{
					Metadata: &people.PersonMetadata{
						Sources: []*people.Source{
							{
								Type: "PROFILE",
								Id:   "some_google_id",
							},
						},
					},
					Names: []*people.Name{
						{
							DisplayName: "Tom Chan",
						},
					},
					Organizations: []*people.Organization{
						{
							Metadata:   &people.FieldMetadata{Primary: false, Verified: true},
							Title:      "Internship",
							Department: "Product",
						},
						{
							Metadata:   &people.FieldMetadata{Primary: true, Verified: true},
							Title:      "Web Developer",
							Department: "Technology",
						},
						{
							Metadata:   &people.FieldMetadata{Primary: false, Verified: true},
							Title:      "Random Developer",
							Department: "Product",
						},
					},
					Photos: []*people.Photo{
						{
							Url: "some_url",
						},
					},
					EmailAddresses: []*people.EmailAddress{
						{
							Value: "tom.chan@scmp.com",
						},
					},
				}
				return person, nil
			}

			// mockUserRepo.EXPECT().FindByGoogleID(
			mockAcctMgr.EXPECT().ValidateEmail("tom.chan@scmp.com").
				Return(nil).
				Times(1)

			user, err := googleOAuth2.doAuthorize(params, exchangeToken, getPersonCall)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.ID).To(Equal(bson.ObjectIdHex("507f1f77bcf86cd799439011")))
			Expect(user.DisplayName).To(Equal("Tom Chan"))
			Expect(user.Department).To(Equal("Technology"))
			Expect(user.Title).To(Equal("Web Developer"))
			Expect(user.AvatarURL).To(Equal("some_url"))
			Expect(user.Type).To(Equal(muser.TypeUser))
		})
	})

	Context("when user login with non-whitelisted email domain", func() {
		var (
			params        *AuthorizeParams
			exchangeToken func(authCode string) (*oauth2.Token, error)
			getPersonCall func(token *oauth2.Token) (*people.Person, error)
		)

		BeforeEach(func() {
			params = &AuthorizeParams{
				AuthCode: "some_code",
			}
			exchangeToken = func(authCode string) (*oauth2.Token, error) {
				return &oauth2.Token{}, nil
			}
			getPersonCall = func(token *oauth2.Token) (*people.Person, error) {
				Expect(token).NotTo(BeNil())

				person := &people.Person{
					Metadata: &people.PersonMetadata{
						Sources: []*people.Source{
							{
								Type: "PROFILE",
								Id:   "some_google_id",
							},
						},
					},
					Names: []*people.Name{
						{
							DisplayName: "Mike Chan",
						},
					},
					Organizations: []*people.Organization{
						{
							Metadata:   &people.FieldMetadata{Primary: true, Verified: true},
							Title:      "Web Developer",
							Department: "Technology",
						},
					},
					Photos: []*people.Photo{
						{
							Url: "some_url",
						},
					},
					EmailAddresses: []*people.EmailAddress{
						{
							Value: "mike.chan@somedomain.com",
						},
					},
				}
				return person, nil
			}

			mockAcctMgr.EXPECT().ValidateEmail("mike.chan@somedomain.com").Return(acctmgr.ErrInvalidEmailDomain)
		})

		It("should return error", func() {
			_, err := googleOAuth2.doAuthorize(params, exchangeToken, getPersonCall)
			Expect(err).To(Equal(acctmgr.ErrInvalidEmailDomain))
		})
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})
})
