CREATE TABLE IF NOT EXISTS sensors (
                                       id SERIAL PRIMARY KEY,
                                       name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    value FLOAT DEFAULT 0,
    unit VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_sensors_type ON sensors(type);
CREATE INDEX IF NOT EXISTS idx_sensors_location ON sensors(location);
CREATE INDEX IF NOT EXISTS idx_sensors_status ON sensors(status);

INSERT INTO sensors (name, type, location, value, unit, status)
VALUES ('Temperature Sensor 1', 'temperature', 'Living Room', 22.5, '°C', 'active')
    ON CONFLICT DO NOTHING;

GRANT ALL ON SCHEMA public TO smarthome;
GRANT ALL ON ALL TABLES IN SCHEMA public TO smarthome;