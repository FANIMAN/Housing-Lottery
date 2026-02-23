package domain

import "time"

type AuditLogResponse struct {
	ID          string     `json:"id"`
	AdminID     *string    `json:"admin_id"`
	AdminEmail  *string    `json:"admin_email"`
	Action      string     `json:"action"`
	EntityType  *string    `json:"entity_type"`
	EntityID    *string    `json:"entity_id"`
	HTTPStatus  *int       `json:"http_status"`
	IPAddress   *string    `json:"ip_address"`
	UserAgent   *string    `json:"user_agent"`
	ErrorMessage *string   `json:"error_message"`
	CreatedAt   time.Time  `json:"created_at"`
}