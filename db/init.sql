CREATE TABLE video (
    id bigint GENERATED ALWAYS AS IDENTITY,
    title text,
    fileSize bigint,
    chunks bigint DEFAULT 0,
    totalChunks bigint,
    isCompiled bool DEFAULT false,
    isSegmented bool DEFAULT false
);
