-- name: GetUserFromEmail :one
SELECT * FROM "User"
WHERE email = $1 LIMIT 1;

-- name: GetUserFromUsername :one
SELECT * FROM "User"
WHERE username = $1 LIMIT 1;

-- name: GetLiveClassFromId :one
SELECT * FROM "LiveClass"
WHERE id = $1 LIMIT 1;

-- name: GetLiveClassFromEmail :many
SELECT * FROM "LiveClass"
WHERE email = $1;

-- name: GetAllCourses :many
SELECT 
    c.id, 
    c.name, 
    c.language, 
    c.description, 
    c.is_public, 
    c.img_url, 
    c.uid, 
    u.id AS user_id, 
    u.username, 
    u.name AS user_name, 
    u.email, 
    u.picture,
    GREATEST(
        similarity(c.language, $4), 
        similarity(u.username, $4), 
        similarity(u.name, $4), 
        similarity(c.name, $4)
    ) AS relevance
FROM "Course" c
JOIN "User" u ON c.uid = u.id
WHERE c.is_public = $1
AND (
    $4 IS NULL OR $4 = '' -- If no search term, don't filter
    OR similarity(c.language, $4) > 0.2
    OR similarity(u.username, $4) > 0.2
    OR similarity(u.name, $4) > 0.2
    OR similarity(c.name, $4) > 0.2
)
ORDER BY 
    CASE WHEN $5 = 'asc' THEN c.id END ASC,
    CASE WHEN $5 = 'desc' THEN c.id END DESC,
    relevance DESC
LIMIT $2 OFFSET $3;



-- name: GetCourse :one
SELECT 
    c.*, 
    u.id AS user_id, 
    u.name AS user_name, 
    u.picture AS user_picture
FROM "Course" c
INNER JOIN "User" u ON c.uid = u.id
WHERE c.id = $1;

-- name: GetAllLessons :many
SELECT 
    l.id, 
    l.title, 
    l.body, 
    l.is_public, 
    l.user_id, 
    u.id AS user_id, 
    u.username, 
    u.name AS user_name, 
    u.email, 
    u.picture,
    GREATEST(
        similarity(l.language, $4), 
        similarity(u.username, $4), 
        similarity(u.name, $4), 
        similarity(l.title, $4)
    ) AS relevance
FROM "LessonPost" l
JOIN "User" u ON l.uid = u.id
WHERE l.is_public = $1
AND (
    $4 IS NULL OR $4 = '' -- No filtering if search term is empty
    OR similarity(l.language, $4) > 0.2 
    OR similarity(u.username, $4) > 0.2
    OR similarity(u.name, $4) > 0.2
    OR similarity(l.title, $4) > 0.2
)
ORDER BY 
    CASE WHEN $5 = 'asc' THEN l.id END ASC,
    CASE WHEN $5 = 'desc' THEN l.id END DESC,
    relevance DESC
LIMIT $2 OFFSET $3;
