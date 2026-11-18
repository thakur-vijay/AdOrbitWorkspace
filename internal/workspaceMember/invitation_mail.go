package workspacemember

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/thakur-vijay/adorbit-go/mailservice"
)

const invitationTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="color-scheme" content="light">
    <title>Workspace Invitation - AdOrbit</title>
</head>

<body style="margin:0;padding:0;background-color:#f3f5f9;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Arial,sans-serif;">

<table width="100%" cellpadding="0" cellspacing="0" border="0"
    style="background-color:#f3f5f9;padding:40px 12px;">
<tr>
<td align="center">

<table width="100%" cellpadding="0" cellspacing="0" border="0"
    style="max-width:520px;background-color:#ffffff;border-radius:16px;overflow:hidden;">

    <!-- Header -->
    <tr>
        <td align="center" style="background-color:#111827;padding:32px 24px;">
            <h1 style="margin:0;color:#ffffff;font-size:27px;font-weight:750;letter-spacing:-1px;">
                AdOrbit<span style="color:#6366f1;">.</span>
            </h1>
            <p style="margin:9px 0 0;color:#9ca3af;font-size:12px;letter-spacing:2px;text-transform:uppercase;">
                Advertising Intelligence Platform
            </p>
        </td>
    </tr>

    <!-- Main Content -->
    <tr>
        <td style="padding:38px 32px 20px;">
            <h2 style="margin:0 0 14px;color:#111827;font-size:24px;font-weight:700;text-align:center;">
                You're invited!
            </h2>

            <p style="margin:0;color:#6b7280;font-size:15px;line-height:1.7;text-align:center;">
                <strong style="color:#111827;">{{.InviterName}}</strong>
                has invited you to join
                <strong style="color:#111827;">{{.WorkspaceName}}</strong>
                on AdOrbit.
            </p>
        </td>
    </tr>

    <!-- Workspace Details -->
    <tr>
        <td style="padding:20px 32px 25px;">
            <table width="100%" cellpadding="0" cellspacing="0" border="0"
                style="background-color:#f9fafb;border:1px solid #e5e7eb;border-radius:12px;">
                <tr>
                    <td style="padding:18px;">
                        <p style="margin:0 0 5px;color:#9ca3af;font-size:11px;font-weight:600;letter-spacing:1px;text-transform:uppercase;">
                            Workspace
                        </p>
                        <p style="margin:0 0 16px;color:#111827;font-size:16px;font-weight:600;">
                            {{.WorkspaceName}}
                        </p>

                        <p style="margin:0 0 5px;color:#9ca3af;font-size:11px;font-weight:600;letter-spacing:1px;text-transform:uppercase;">
                            Invited by
                        </p>
                        <p style="margin:0 0 16px;color:#111827;font-size:14px;">
                            {{.InviterName}}
                        </p>

                        <p style="margin:0 0 5px;color:#9ca3af;font-size:11px;font-weight:600;letter-spacing:1px;text-transform:uppercase;">
                            Your role
                        </p>
                        <p style="margin:0;color:#4f46e5;font-size:14px;font-weight:600;">
                            {{.Role}}
                        </p>
                    </td>
                </tr>
            </table>
        </td>
    </tr>

    <!-- CTA -->
    <tr>
        <td align="center" style="padding:5px 32px 30px;">
            <a href="{{.InvitationLink}}"
                style="display:inline-block;background-color:#4f46e5;color:#ffffff;text-decoration:none;padding:15px 35px;border-radius:9px;font-size:14px;font-weight:600;">
                Accept Invitation
            </a>
        </td>
    </tr>

    <!-- Expiry Notice -->
    <tr>
        <td style="padding:0 32px 28px;">
            <table width="100%" cellpadding="0" cellspacing="0" border="0"
                style="background-color:#fffbeb;border:1px solid #fde68a;border-radius:10px;">
                <tr>
                    <td style="padding:15px 16px;">
                        <p style="margin:0;color:#92400e;font-size:13px;line-height:1.6;">
                            <strong>Invitation expiry</strong><br>
                            This invitation expires on
                            <strong>{{.ExpiresAt}}.</strong>
                            Please accept it before then.
                        </p>
                    </td>
                </tr>
            </table>
        </td>
    </tr>

    <!-- Fallback Link -->
    <tr>
        <td style="padding:0 32px 30px;">
            <p style="margin:0 0 10px;color:#6b7280;font-size:12px;line-height:1.6;">
                If the button doesn't work, copy and paste this link into your browser:
            </p>
            <p style="margin:0;word-break:break-all;font-size:12px;line-height:1.6;">
                <a href="{{.InvitationLink}}" style="color:#4f46e5;">
                    {{.InvitationLink}}
                </a>
            </p>
        </td>
    </tr>

    <!-- Security Notice -->
    <tr>
        <td style="padding:0 32px 35px;border-bottom:1px solid #f1f2f5;">
            <p style="margin:0;color:#6b7280;font-size:13px;line-height:1.7;text-align:center;">
                If you weren't expecting this invitation, you can safely ignore this email.
            </p>
        </td>
    </tr>

    <!-- Footer -->
    <tr>
        <td align="center" style="padding:25px 24px 30px;">
            <p style="margin:0 0 8px;color:#111827;font-size:13px;font-weight:600;">
                AdOrbit
            </p>
            <p style="margin:0;color:#9ca3af;font-size:11px;line-height:1.6;">
                © 2026 AdOrbit. All rights reserved.<br>
                This is an automated email. Please do not reply.
            </p>
        </td>
    </tr>

</table>

</td>
</tr>
</table>

</body>
</html>
`

type InvitationMailService struct {
	mailService mailservice.Service
}

type InvitationEmailDTO struct {
	Email          string
	WorkspaceName  string
	InviterName    string
	Role           string
	InvitationLink string
	ExpiresAt      time.Time
}

func NewInvitationMailService(
	mailService mailservice.Service,
) *InvitationMailService {
	return &InvitationMailService{
		mailService: mailService,
	}
}

func (s *InvitationMailService) Send(
	ctx context.Context,
	dto InvitationEmailDTO,
) error {

	tmpl, err := template.New("invitation").Parse(invitationTemplate)
	if err != nil {
		return fmt.Errorf("parse invitation template: %w", err)
	}

	var htmlBody bytes.Buffer

	data := struct {
		WorkspaceName  string
		InviterName    string
		Role           string
		InvitationLink string
		ExpiresAt      string
	}{
		WorkspaceName:  dto.WorkspaceName,
		InviterName:    dto.InviterName,
		Role:           dto.Role,
		InvitationLink: dto.InvitationLink,
		ExpiresAt:      dto.ExpiresAt.UTC().Format("02 Jan 2006, 03:04 PM MST"),
	}

	if err := tmpl.Execute(&htmlBody, data); err != nil {
		return fmt.Errorf("render invitation template: %w", err)
	}

	return s.mailService.Send(ctx, mailservice.Message{
		To:       []string{dto.Email},
		Subject:  "You're invited to join a workspace - AdOrbit",
		HTMLBody: htmlBody.String(),
		TextBody: fmt.Sprintf(
			"%s has invited you to join %s on AdOrbit.\n\n"+
				"Role: %s\n\n"+
				"Accept your invitation: %s\n\n"+
				"This invitation expires on %s.",
			dto.InviterName,
			dto.WorkspaceName,
			dto.Role,
			dto.InvitationLink,
			data.ExpiresAt,
		),
	})
}
