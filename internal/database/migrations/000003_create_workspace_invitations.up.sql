CREATE TABLE
    workspace_invitations (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        workspace_id UUID NOT NULL REFERENCES workspaces (id) ON DELETE CASCADE,
        email VARCHAR(255) NOT NULL,
        role VARCHAR(20) NOT NULL DEFAULT 'member',
        token_hash VARCHAR(255) NOT NULL UNIQUE,
        status VARCHAR(20) NOT NULL DEFAULT 'pending',
        invited_by UUID NOT NULL,
        expires_at TIMESTAMPTZ NOT NULL,
        accepted_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT chk_workspace_invitations_role CHECK (role IN ('admin', 'member')),
        CONSTRAINT chk_workspace_invitations_status CHECK (
            status IN ('pending', 'accepted', 'expired', 'revoked')
        )
    );

CREATE INDEX idx_workspace_invitations_workspace_id ON workspace_invitations (workspace_id);

CREATE INDEX idx_workspace_invitations_status ON workspace_invitations (status);

CREATE INDEX idx_workspace_invitations_expires_at ON workspace_invitations (expires_at);

CREATE UNIQUE INDEX uq_workspace_invitations_workspace_email ON workspace_invitations (workspace_id, LOWER(email));