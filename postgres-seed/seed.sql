INSERT INTO users (id, email, password_hash, role)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'admin@dummy.com', '', 'admin'),
    ('00000000-0000-0000-0000-000000000002', 'user@dummy.com', '', 'user')
ON CONFLICT (id) DO NOTHING;
