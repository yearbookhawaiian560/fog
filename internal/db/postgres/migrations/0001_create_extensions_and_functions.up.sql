CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE  FUNCTION update_uat()
RETURNS TRIGGER AS $$
BEGIN
    NEW.uat = now();
    RETURN NEW;
END;
$$ language 'plpgsql';


CREATE EXTENSION IF NOT EXISTS vector;
