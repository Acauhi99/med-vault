-- Remove admin user from default tenant
DELETE FROM user_tenants
WHERE user_id = (SELECT id FROM users WHERE email = 'acauhi.mateus@gmail.com')
  AND tenant_id = (SELECT id FROM tenants WHERE name = 'MedVault Demo');

-- Remove default tenant
DELETE FROM tenants WHERE name = 'MedVault Demo';
