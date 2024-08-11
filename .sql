-- Active: 1723171055346@@127.0.0.1@5432@postgres@public
DROP TABLE IF EXISTS news;
DROP TABLE IF EXISTS newscategories;
CREATE TABLE IF NOT EXISTS News (
    id SERIAL PRIMARY KEY NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS NewsCategories (
    categoryId BIGINT NOT NULL PRIMARY KEY,
    newsId BIGINT NOT NULL REFERENCES News(Id)
);

CREATE TABLE IF NOT EXISTS NewNewsCategories (
    categoryId BIGINT NOT NULL,
    newsId BIGINT NOT NULL REFERENCES News(Id),
    PRIMARY KEY (categoryId, newsId)
);
INSERT INTO
    NewNewsCategories
SELECT
FROM
    NewsCategories;
DROP TABLE IF EXISTS newscategories;
ALTER TABLE NewNewsCategories RENAME TO newscategories;

INSERT INTO newscategories VALUES (34, 344);