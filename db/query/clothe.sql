-- name: GetClothe :one
SELECT * FROM clothes
WHERE id = ? LIMIT 1;

-- name: ListClothes :many
SELECT * FROM clothes
ORDER BY name;

-- name: CreateClothe :execresult
INSERT INTO clothes (
  name, amount, cost, price
) VALUES (
  ?, ?, ?, ?
);

-- name: UpdateClothePrice :exec
UPDATE clothes
  SET price = ?
WHERE id = ?;

-- name: UpdateClotheAmount :exec
UPDATE clothes
  SET amount = ?
WHERE id = ?;

-- name: DeleteClothe :exec
DELETE FROM clothes
WHERE id = ?;



-- CREATE TABLE clothes(
-- 	  id INT PRIMARY KEY AUTO_INCREMENT,
--     name VARCHAR(255) NOT NULL,
--     amount INT NOT NULL,
--     cost INT NOT NULL,
--     price INT NOT NULL,
--     create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );