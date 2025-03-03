package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sarvochcha01/enlace-backend/internal/models"
)

type CommentRepository interface {
	CreateComment(*models.CreateCommentDTO) error
	GetComment(uuid.UUID) (*models.CommentDTO, error)
	UpdateComment(uuid.UUID, string) error
	DeleteComment(uuid.UUID) error

	GetAllComments(uuid.UUID) ([]models.CommentDTO, error)

	GetCommentCreator(uuid.UUID) (uuid.UUID, error)
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) CreateComment(commentDTO *models.CreateCommentDTO) error {

	queryString := `
		INSERT INTO comments (project_id, task_id, created_by, comment)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(queryString, commentDTO.ProjectID, commentDTO.TaskID, commentDTO.CreatedBy, commentDTO.Comment)

	if err != nil {
		return err
	}

	return nil
}

func (r *commentRepository) GetComment(commentID uuid.UUID) (*models.CommentDTO, error) {
	var commentDTO models.CommentDTO

	queryString := `
		SELECT id, project_id, task_id, created_by, comment, created_at, updated_at
		FROM comments
		WHERE id = $1	
	`

	err := r.db.QueryRow(queryString, commentID).Scan(&commentDTO.ID, &commentDTO.ProjectID, &commentDTO.TaskID, &commentDTO.CreatedBy, &commentDTO.Comment, &commentDTO.CreatedAt, &commentDTO.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &commentDTO, nil
}

func (r *commentRepository) UpdateComment(commentID uuid.UUID, newComment string) error {
	query := `
		UPDATE comments
		SET comment = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(query, newComment, commentID)
	return err
}

func (r *commentRepository) DeleteComment(commentID uuid.UUID) error {
	queryString := `
		DELETE FROM comments
		WHERE id = $1
	`

	result, err := r.db.Exec(queryString, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no comment found with the given ID")
	}

	return nil
}

func (r *commentRepository) GetCommentCreator(commentID uuid.UUID) (uuid.UUID, error) {
	var creatorID uuid.UUID
	err := r.db.QueryRow("SELECT created_by FROM comments WHERE id = $1", commentID).Scan(&creatorID)
	if err != nil {
		return uuid.Nil, err
	}
	return creatorID, nil
}

func (r *commentRepository) GetAllComments(taskID uuid.UUID) ([]models.CommentDTO, error) {
	var comments []models.CommentDTO

	queryString := `
		SELECT id, project_id, task_id, created_by, comment, created_at, updated_at
		FROM comments
		WHERE task_id = $1
	`

	rows, err := r.db.Query(queryString)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var comment models.CommentDTO
		if err := rows.Scan(&comment.ID, &comment.ProjectID, &comment.TaskID, &comment.CreatedBy, &comment.Comment, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil

}
