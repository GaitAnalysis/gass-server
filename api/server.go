package api

import (
	"gass/domain"
	v1 "gass/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	svc *v1.Service
}

func (s *Server) Register(c *gin.Context) {
	var input *domain.AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{Email: input.Email, Password: input.Password}
	if err := s.svc.RegisterUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (s *Server) Login(c *gin.Context) {
	var input *domain.AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.svc.LoginUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": token})
}

func (s *Server) AddUpload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	upload, err := s.svc.SaveUpload(c, fileHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"upload": upload})
}

func (s *Server) GetUploads(c *gin.Context) {
	uploads, err := s.svc.GetUploads(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uploads": uploads})
}

func (s *Server) AddAnalysis(c *gin.Context) {
	var analysisInput *domain.AnalysisInput
	if err := c.ShouldBindJSON(&analysisInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analysis, err := s.svc.AddAnalysis(analysisInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analysis": analysis})
}

func (s *Server) GetAnalyses(c *gin.Context) {
	analyses, err := s.svc.GetAnalyses(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analyses": analyses})
}

func NewServer(svc *v1.Service) *Server {
	return &Server{
		svc: svc,
	}
}
