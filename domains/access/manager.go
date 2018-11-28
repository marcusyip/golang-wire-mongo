package access

import (
	"context"
	"time"

	"github.com/marcusyip/golang-wire-mongo/core/rand"
	"github.com/marcusyip/golang-wire-mongo/models"
	repos "github.com/marcusyip/golang-wire-mongo/repositories"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	logger     *logrus.Logger
	userRepo   *repos.UserRepository
	accessRepo *repos.AccessRepository
}

func (mgr *Manager) Authenticate(ctx context.Context, accessToken string) (*models.User, error) {
	// log := mgr.getLogger("Authenticate")
	filter := bson.D{
		{"access_token", accessToken},
		{"expire_at", bson.D{
			{"$gt", time.Now()},
		}},
	}
	access, err := mgr.accessRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	user, err := mgr.userRepo.FindByID(ctx, access.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (mgr *Manager) Resolve(ctx context.Context, accessToken string) (*models.Access, error) {
	filter := bson.D{
		{"access_token", accessToken},
		{"expire_at", bson.D{
			{"$gt", time.Now()},
		}},
	}
	access, err := mgr.accessRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return access, nil
}

func (mgr *Manager) CreateForUser(ctx context.Context, user *models.User, provider models.AccessProvider) (*models.Access, error) {
	// log := mgr.getLogger("CreateForUser")
	access := &models.Access{
		Model:        *models.NewModel(),
		UserID:       user.ID,
		AccessToken:  rand.RandString(48),
		RefreshToken: rand.RandString(48),
		Provider:     provider,
		ExpireAt:     time.Now().AddDate(1, 0, 0),
	}
	err := mgr.accessRepo.Create(ctx, access)
	if err != nil {
		return nil, err
	}
	return access, nil
}

func (mgr *Manager) getLogger(method string) *logrus.Entry {
	return mgr.logger.WithFields(logrus.Fields{"domain": "access", "method": method})
}

func NewManager(
	logger *logrus.Logger,
	userRepo *repos.UserRepository,
	accessRepo *repos.AccessRepository,
) *Manager {
	return &Manager{
		logger:     logger,
		userRepo:   userRepo,
		accessRepo: accessRepo,
	}
}
