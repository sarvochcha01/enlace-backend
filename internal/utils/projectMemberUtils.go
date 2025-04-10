package utils

import "github.com/sarvochcha01/enlace-backend/internal/models"

func HasEditPrivileges(projectMember *models.ProjectMemberResponseDTO) bool {
	return projectMember.Role == models.RoleOwner || projectMember.Role == models.RoleEditor
}
