CREATE TABLE projects (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    name NVARCHAR(255) NOT NULL UNIQUE,
    description NVARCHAR(MAX) NOT NULL CONSTRAINT DF_projects_description DEFAULT '',
    created_at DATETIMEOFFSET NOT NULL CONSTRAINT DF_projects_created_at DEFAULT SYSUTCDATETIME(),
    updated_at DATETIMEOFFSET NOT NULL CONSTRAINT DF_projects_updated_at DEFAULT SYSUTCDATETIME()
);

CREATE TABLE services (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    project_id UNIQUEIDENTIFIER NOT NULL,
    name NVARCHAR(255) NOT NULL,
    image NVARCHAR(2048) NOT NULL,
    replicas INT NOT NULL CONSTRAINT DF_services_replicas DEFAULT 1,
    created_at DATETIMEOFFSET NOT NULL CONSTRAINT DF_services_created_at DEFAULT SYSUTCDATETIME(),
    updated_at DATETIMEOFFSET NOT NULL CONSTRAINT DF_services_updated_at DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_services_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT UQ_services_project_name UNIQUE (project_id, name)
);

CREATE TABLE deployments (
    id UNIQUEIDENTIFIER PRIMARY KEY,
    service_id UNIQUEIDENTIFIER NOT NULL,
    version NVARCHAR(255) NOT NULL,
    status NVARCHAR(50) NOT NULL CONSTRAINT DF_deployments_status DEFAULT 'pending',
    created_at DATETIMEOFFSET NOT NULL CONSTRAINT DF_deployments_created_at DEFAULT SYSUTCDATETIME(),
    updated_at DATETIMEOFFSET NOT NULL CONSTRAINT DF_deployments_updated_at DEFAULT SYSUTCDATETIME(),
    CONSTRAINT FK_deployments_service FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE INDEX idx_services_project_id
    ON services(project_id);

CREATE INDEX idx_deployments_service_id
    ON deployments(service_id);
