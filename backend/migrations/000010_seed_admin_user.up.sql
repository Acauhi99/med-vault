-- Create default tenant and add first admin user
-- This migration provisions the initial admin access for the demo

INSERT INTO tenants (id, name, status, created_at, updated_at)
SELECT gen_random_uuid(), 'MedVault Demo', 'active', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM tenants WHERE name = 'MedVault Demo');

-- Add the admin user to the default tenant
INSERT INTO user_tenants (user_id, tenant_id, role)
SELECT u.id, t.id, 'administrator'
FROM users u, tenants t
WHERE u.email = 'acauhi.mateus@gmail.com'
  AND t.name = 'MedVault Demo'
  AND NOT EXISTS (
    SELECT 1 FROM user_tenants ut
    WHERE ut.user_id = u.id AND ut.tenant_id = t.id
  );
