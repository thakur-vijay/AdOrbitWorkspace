package workspacemember

import "context"

type Service interface {
	SendInvite(ctx context.Context, dto SendInviteDTO, userID string) error
}
