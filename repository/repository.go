package repository

import (
	"fmt"
	"gass/domain"
	"html"
	"net/mail"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GassRepository struct {
	dbClient *gorm.DB
}

func Initialize() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Europe/Istanbul", host, username, password, databaseName, port)
	dbClient, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return dbClient, err
}

func (r *GassRepository) CreateUser(user *domain.User) error {
	email, err := mail.ParseAddress(user.Email)
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Email = html.EscapeString(email.Address)

	err = r.dbClient.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GassRepository) FindUserById(id uint) (*domain.User, error) {
	var user *domain.User
	err := r.dbClient.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GassRepository) FindUploadById(id uint) (*domain.Upload, error) {
	var upload *domain.Upload
	err := r.dbClient.First(&upload, id).Error
	if err != nil {
		return nil, err
	}
	return upload, nil
}

func (r *GassRepository) FindUserByEmail(email string) (*domain.User, error) {
	var user *domain.User
	err := r.dbClient.Where("email=?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GassRepository) GetUserUploads(user *domain.User) ([]*domain.Upload, error) {
	var userUploads []*domain.Upload
	err := r.dbClient.Model(&user).Association("Uploads").Find(&userUploads)
	if err != nil {
		return nil, err
	}
	return userUploads, nil
}

func (r *GassRepository) GetUserAnalyses(user *domain.User) ([]*domain.Analysis, error) {
	var userAnalyses []*domain.Analysis
	err := r.dbClient.Model(&user).Association("Analyses").Find(&userAnalyses)
	if err != nil {
		return nil, err
	}
	return userAnalyses, nil
}

func (r *GassRepository) CreateAnalysis(analysis *domain.Analysis) error {
	err := r.dbClient.Create(&analysis).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GassRepository) GetAnalysisUpload(analysis *domain.Analysis) (*domain.Upload, error) {
	var upload *domain.Upload
	err := r.dbClient.First(&upload, analysis.UploadId).Error
	if err != nil {
		return nil, err
	}
	return upload, nil
}

func (r *GassRepository) CreateUpload(upload *domain.Upload) error {
	err := r.dbClient.Create(&upload).Error
	if err != nil {
		return err
	}
	return nil
}

func NewGassRepository() (*GassRepository, error) {
	dbClient, err := Initialize()
	if err != nil {
		return nil, err
	}

	dbClient.AutoMigrate(&domain.User{})
	dbClient.AutoMigrate(&domain.Upload{})
	dbClient.AutoMigrate(&domain.Analysis{})

	return &GassRepository{dbClient: dbClient}, nil
}
