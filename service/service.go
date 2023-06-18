package v1

import (
	"bytes"
	"fmt"
	"gass/domain"
	"gass/pkg/jwt"
	"gass/repository"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *repository.GassRepository
	httpClient *http.Client
}

func (s *Service) RegisterUser(user *domain.User) error {
	return s.repository.CreateUser(user)
}

func (s *Service) LoginUser(input *domain.AuthInput) (token string, err error) {
	user, err := s.repository.FindUserByEmail(input.Email)
	if err != nil {
		return "", err
	}

	err = s.validatePassword(user, input.Password)
	if err != nil {
		return "", err
	}

	token, err = jwt.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) SaveUpload(c *gin.Context, fileHeader *multipart.FileHeader) (*domain.Upload, error) {
	userID, err := jwt.CurrentUserID(c)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.FindUserById(userID)
	if err != nil {
		return nil, err
	}

	userUploads, err := s.repository.GetUserUploads(user)
	if err != nil {
		return nil, err
	}

	fileRef := "uploads/" + strconv.Itoa(int(user.ID)) + "_" + strconv.Itoa(len(userUploads))
	upload := &domain.Upload{
		FileRef: fileRef,
		Size:    fileHeader.Size,
		UserId:  userID,
		User:    *user,
	}

	err = s.repository.CreateUpload(upload)
	if err != nil {
		return nil, err
	}

	err = c.SaveUploadedFile(fileHeader, fileRef)
	if err != nil {
		return nil, err
	}

	err = s.requestAnalysis(userID, upload.ID, fileRef)
	return upload, err
}

func (s *Service) GetUploads(c *gin.Context) ([]*domain.Upload, error) {
	userID, err := jwt.CurrentUserID(c)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.FindUserById(userID)
	if err != nil {
		return nil, err
	}

	uploads, err := s.repository.GetUserUploads(user)
	if err != nil {
		return nil, err
	}
	return uploads, nil
}

func (s *Service) requestAnalysis(userID uint, uploadID uint, uploadPath string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	err := bodyWriter.WriteField("user_id", strconv.Itoa(int(userID)))
	if err != nil {
		return err
	}

	err = bodyWriter.WriteField("video_id", strconv.Itoa(int(uploadID)))
	if err != nil {
		return err
	}

	file, err := os.Open(uploadPath)
	if err != nil {
		return err
	}
	defer file.Close()

	videoWriter, err := bodyWriter.CreateFormFile("video", uploadPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(videoWriter, file)
	if err != nil {
		return err
	}
	bodyWriter.Close()

	server_addr := fmt.Sprintf("%s/classify", os.Getenv("INFERENCE_SERVER"))
	req, err := http.NewRequest("POST", server_addr+"/classify", bodyBuf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(os.Stdout, res.Body)
	return err
}

func (s *Service) AddAnalysis(input *domain.AnalysisInput) (*domain.Analysis, error) {
	user, err := s.repository.FindUserById(input.UserId)
	if err != nil {
		return nil, err
	}

	upload, err := s.repository.FindUploadById(input.UploadId)
	if err != nil {
		return nil, err
	}

	analysis := &domain.Analysis{
		UserId:   input.UserId,
		User:     *user,
		UploadId: input.UploadId,
		Upload:   *upload,
		Result:   input.Result,
	}

	err = s.repository.CreateAnalysis(analysis)
	if err != nil {
		return nil, err
	}

	return analysis, nil
}

func (s *Service) GetAnalyses(c *gin.Context) ([]*domain.Analysis, error) {
	userID, err := jwt.CurrentUserID(c)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.FindUserById(userID)
	if err != nil {
		fmt.Println("a")
		return nil, err
	}

	uploads, err := s.repository.GetUserUploads(user)
	if err != nil {
		fmt.Println("b")
		return nil, err
	}

	uploadsMap := make(map[uint]*domain.Upload, len(uploads))
	for _, upload := range uploads {
		uploadsMap[upload.ID] = upload
	}

	analyses, err := s.repository.GetUserAnalyses(user)
	if err != nil {
		fmt.Println("c")
		return nil, err
	}

	for _, analysis := range analyses {
		upload, ok := uploadsMap[analysis.UploadId]
		if ok {
			analysis.Upload = *upload
		}
	}
	return analyses, nil
}

func (s *Service) validatePassword(user *domain.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func NewService(repository *repository.GassRepository, httpClient *http.Client) *Service {
	return &Service{
		repository: repository,
		httpClient: httpClient,
	}
}
