CREATE TABLE IF NOT EXISTS users (
    acc_tag TEXT PRIMARY KEY,
    reg_country TEXT,
    reg_city TEXT,
    first_email TEXT,
    phone TEXT,
    first_device TEXT,
    is_donator BOOLEAN DEFAULT FALSE,
    reg_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT PRIMARY KEY,
    acc_tag TEXT REFERENCES users(acc_tag) ON DELETE CASCADE,
    session_ip TEXT,
    device_id TEXT,
    asn TEXT,
    country TEXT,
    city TEXT,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id TEXT PRIMARY KEY,
    acc_tag TEXT REFERENCES users(acc_tag) ON DELETE CASCADE,
    amount NUMERIC(15, 2),
    currency TEXT,
    payment_method TEXT,
    payment_provider TEXT,
    device_id TEXT,
    device_model TEXT,
    ip TEXT,
    country TEXT,
    city TEXT,
    asn TEXT,
    session_id TEXT REFERENCES sessions(session_id),
    timestamp TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS tickets (
    ticket_id BIGSERIAL PRIMARY KEY,
    acc_tag TEXT REFERENCES users(acc_tag) ON DELETE CASCADE,
    claimant_tag TEXT,
    device_id TEXT,
    final_percentage TEXT,
    knowledge_score NUMERIC(10, 4),
    penalty_score NUMERIC(10, 4),
    ip_info JSONB, 
    created_at TIMESTAMPTZ DEFAULT NOW()
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ticket_details (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT REFERENCES tickets(ticket_id) ON DELETE CASCADE,
    name TEXT,
    code TEXT,
    value NUMERIC(15, 4),
    weight NUMERIC(15, 4),
    result NUMERIC(15, 4),
    comment TEXT,
    status TEXT
);

CREATE TABLE IF NOT EXISTS system_info (
    ticket_id BIGINT PRIMARY KEY REFERENCES tickets(ticket_id) ON DELETE CASCADE,
    os TEXT,
    platform TEXT,
    arch TEXT,
    kernel TEXT,
    cpu_model TEXT,
    cpu_cores INTEGER,
    total_ram BIGINT,
    hostname TEXT,
    machine_id TEXT NOT NULL,
    username TEXT
);

CREATE TABLE IF NOT EXISTS dbrecord (
    acc_tag TEXT PRIMARY KEY,
    reg_country TEXT,
    reg_city TEXT,
    first_email TEXT,
    phone TEXT,
    first_device TEXT,
    is_donator BOOLEAN DEFAULT FALSE,
    reg_date TIMESTAMPTZ,
    
    devices JSONB,     
    first_transaction JSONB,   
    user_history JSONB,    
    
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dbrecord_email ON dbrecord(first_email);
CREATE INDEX IF NOT EXISTS idx_sysinfo_machine_id ON system_info(machine_id);
CREATE INDEX IF NOT EXISTS idx_sessions_acc_tag ON sessions(acc_tag);
CREATE INDEX IF NOT EXISTS idx_transactions_acc_tag ON transactions(acc_tag);
CREATE INDEX IF NOT EXISTS idx_tickets_acc_tag ON tickets(acc_tag);
CREATE INDEX IF NOT EXISTS idx_sessions_device_id ON sessions(device_id);
CREATE INDEX IF NOT EXISTS idx_tickets_claimant_tag ON tickets(claimant_tag);