package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	connAttempts int
	connTimeout  time.Duration
	Conn         *gorm.DB
}

func New(cfg config.Config, opts ...Option) (*Postgres, error) {

	pg := &Postgres{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.PG.Username, cfg.PG.Password, cfg.PG.Host, cfg.PG.Port, cfg.PG.Dbname, cfg.PG.Sslmode)

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	var err error

	for pg.connAttempts > 0 {
		pg.Conn, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	// migration
	pg.Conn.AutoMigrate(&entity.Role{})
	pg.Conn.AutoMigrate(&entity.User{})

	pg.seedData(cfg)

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgres) seedData(cfg config.Config) {
	var roles = []entity.Role{
		{Name: "admin", Description: "Administrator role"},
		{Name: "customer", Description: "Authenticated customer role"},
		{Name: "anonymous", Description: "Unauthenticated customer role"}}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(cfg.Authen.AdminPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Println("Error hashing password:", err)
		return
	}

	var user = []entity.User{{
		Username: cfg.Authen.AdminUsername,
		Email:    cfg.Authen.AdminEmail,
		Password: string(encryptedPassword),
		RoleID:   1,
	}}

	// Upsert roles
	result := p.Conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},                   // Assuming uniqueness on Name
		DoUpdates: clause.AssignmentColumns([]string{"description"}), // Update only Description
	}).Create(&roles)
	if err := result.Error; err != nil {
		log.Println("Error upserting roles:", err)
		return
	}

	// Fetch role ID for "admin" role
	var adminRole entity.Role
	result = p.Conn.Where("name = ?", "admin").First(&adminRole)
	if err := result.Error; err != nil {
		log.Println("Error fetching admin role:", err)
		return
	}

	// Update user with the fetched role ID
	user[0].RoleID = adminRole.ID

	// Upsert user
	result = p.Conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}}, // Assuming uniqueness on Username
		DoUpdates: clause.AssignmentColumns([]string{"email", "password", "role_id"}),
	}).Create(&user)
	if err := result.Error; err != nil {
		log.Println("Error upserting user:", err)
		return
	}
}
