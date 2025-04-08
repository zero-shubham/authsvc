
-- Migration for creating app_groups table
CREATE TABLE app_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    scopes TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for app_groups
CREATE INDEX idx_app_groups_org_id ON app_groups(org_id);

-- Migration for creating apps table
CREATE TABLE apps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL,
    app_grp_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_apps_app_groups FOREIGN KEY (app_grp_id) REFERENCES app_groups(id) ON DELETE SET NULL
);

-- Indexes for apps
CREATE INDEX idx_apps_org_id ON apps(org_id);
CREATE INDEX idx_apps_app_grp_id ON apps(app_grp_id);

-- Migration for creating users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    app_group_id UUID NOT NULL,
    org_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_users_app_groups FOREIGN KEY (app_group_id) REFERENCES app_groups(id) ON DELETE RESTRICT
);

-- Indexes for users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_app_group_id ON users(app_group_id);
CREATE INDEX idx_users_org_id ON users(org_id) WHERE org_id IS NOT NULL;