package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
	"github.com/sarvochcha01/enlace-backend/internal/repositories"
)

type CommentService interface {
	CreateComment(*models.CreateCommentDTO, string) error
	GetComment(uuid.UUID) (*models.CommentResponseDTO, error)
	UpdateComment(*models.UpdateCommentDTO, string) error
	DeleteComment(*models.DeleteCommentDTO, string) error

	GetAllCommentsForTask(taskID uuid.UUID, projectID uuid.UUID, firebaseUID string) ([]models.CommentResponseDTO, error)
}

type commentService struct {
	commentRepository    repositories.CommentRepository
	userService          UserService
	projectMemberService ProjectMemberService
}

func NewCommentService(cr repositories.CommentRepository, us UserService, pms ProjectMemberService) CommentService {
	return &commentService{commentRepository: cr, userService: us, projectMemberService: pms}
}

func (s *commentService) CreateComment(commentDTO *models.CreateCommentDTO, firebaseUID string) error {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)

	if err != nil {
		log.Println("UserID not found: ", err)
		return errors.New("UserID not found: " + err.Error())
	}

	var projectMemberID uuid.UUID
	projectMemberID, err = s.projectMemberService.GetProjectMemberID(userID, commentDTO.ProjectID)

	if err != nil {
		log.Println("Project Member not found: ", err)
		return errors.New("Project Member not found: " + err.Error())
	}

	commentDTO.CreatedBy = projectMemberID

	return s.commentRepository.CreateComment(commentDTO)
}

func (r *commentService) GetComment(commentID uuid.UUID) (*models.CommentResponseDTO, error) {
	return r.commentRepository.GetComment(commentID)
}

func (s *commentService) UpdateComment(updateCommentDTO *models.UpdateCommentDTO, firebaseUID string) error {

	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return errors.New("user not found")
	}

	projectMemberID, err := s.projectMemberService.GetProjectMemberID(userID, updateCommentDTO.ProjectID)
	if err != nil {
		return errors.New("Project Member not found: " + err.Error())
	}

	creatorID, err := s.commentRepository.GetCommentCreator(updateCommentDTO.CommentID)
	if err != nil {
		return errors.New("comment creator not found")
	}

	if creatorID != projectMemberID {
		return errors.New("unauthorized: you can only edit your own comments")
	}

	return s.commentRepository.UpdateComment(updateCommentDTO.CommentID, updateCommentDTO.Comment)
}

func (s *commentService) DeleteComment(deleteCommentDTO *models.DeleteCommentDTO, firebaseUID string) error {
	userID, err := s.userService.FindUserIDByFirebaseUID(firebaseUID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err.Error())
	}

	projectMemberID, err := s.projectMemberService.GetProjectMemberID(userID, deleteCommentDTO.ProjectID)
	if err != nil {
		return fmt.Errorf("project Member not found: %v", err.Error())
	}

	creatorID, err := s.commentRepository.GetCommentCreator(deleteCommentDTO.CommentID)
	if err != nil {
		return fmt.Errorf("comment creator not found: %v", err.Error())
	}

	if creatorID != projectMemberID {
		return fmt.Errorf("unauthorized: you can only edit your own comments")
	}

	return s.commentRepository.DeleteComment(deleteCommentDTO.CommentID)
}

func (s *commentService) GetAllCommentsForTask(taskID uuid.UUID, projectID uuid.UUID, firebaseUID string) ([]models.CommentResponseDTO, error) {

	_, err := s.projectMemberService.GetProjectMemberIDByFirebaseUID(firebaseUID, projectID)

	if err != nil {
		log.Println("Failed to get Comments. Only projects members can access comments", err)
		return nil, errors.New("failed to get Comments. Only projects members can access comments")
	}

	return s.commentRepository.GetAllCommentsForTask(taskID)
}
