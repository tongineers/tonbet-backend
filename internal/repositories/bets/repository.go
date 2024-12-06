package bets

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/tongineers/dice-ton-api/internal/models"
)

type (
	Repository struct {
		db *gorm.DB
	}
)

const (
	DefaultReadLimit = 10
)

func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) MustMigrate() {
	if err := r.Migrate(); err != nil {
		panic(err)
	}
}

func (r *Repository) Migrate() error {
	return r.db.AutoMigrate(&models.Bet{})
}

func (r *Repository) Read() ([]*models.Bet, error) {
	bets := make([]*models.Bet, 0)
	res := r.db.Limit(DefaultReadLimit).
		Order("id desc").
		Find(&bets)
	return bets, res.Error
}

func (r *Repository) ReadByPlayerAddress(addr string) ([]*models.Bet, error) {
	bets := make([]*models.Bet, 0)
	res := r.db.Where("player_address = ?", addr).
		Limit(DefaultReadLimit).
		Order("id desc").
		Find(&bets)
	return bets, res.Error
}

func (r *Repository) ReadByStatus(status models.BetStatus) ([]*models.Bet, error) {
	bets := make([]*models.Bet, 0)
	res := r.db.Where("status = ?", status).Find(&bets)
	return bets, res.Error
}

func (r *Repository) ReadByIDs(ids ...int) ([]*models.Bet, error) {
	bets := make([]*models.Bet, 0)
	res := r.db.Where("id IN ?", ids).Find(&bets)
	return bets, res.Error
}

func (r *Repository) GetLastResolvedBetLT() (uint64, error) {
	var lastLT uint64
	res := r.db.Table("bets").
		Select("last_lt").
		Where("status = ?", models.BetStatusResolved).
		Order("last_lt desc").
		Limit(1).
		Find(&lastLT)

	return lastLT, res.Error
}

func (r *Repository) Update(bets ...*models.Bet) error {
	res := r.db.Clauses(
		clause.OnConflict{
			UpdateAll: true,
		},
	).Create(&bets)
	return res.Error
}
