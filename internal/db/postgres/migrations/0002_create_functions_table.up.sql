CREATE TABLE functions (
    function_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    function_name VARCHAR(255) NOT NULL UNIQUE,
    embedding vector(1536), -- text-embedding-3-small model is 1536 dimensions,
    js TEXT NOT NULL -- actual ECMAScript 5.1 (ES5) code
);