CREATE TABLE
    workspaces (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        name VARCHAR(150) NOT NULL,
        slug VARCHAR(180) NOT NULL UNIQUE,
        type VARCHAR(20) NOT NULL CHECK (type IN ('advertiser', 'publisher')),
        country CHAR(2) NOT NULL,
        timezone VARCHAR(100) NOT NULL,
        currency CHAR(3) NOT NULL,
        industry VARCHAR(100),
        website TEXT,
        description TEXT,
        status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'deleted')),
        created_by UUID NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_workspaces_created_by ON workspaces (created_by);

CREATE INDEX idx_workspaces_type ON workspaces (type);

CREATE INDEX idx_workspaces_status ON workspaces (status);