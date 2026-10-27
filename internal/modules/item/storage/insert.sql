INSERT INTO item_tag_relations 
SELECT 2, v.id 
FROM item_tag_values AS v 
JOIN item_tag_keys AS k ON v.tag_key_id = k.id 
JOIN (
    VALUES 
        ($1, $2),
        ($1, $3),
        ($4, $5)
) AS tags(key_name, value_name)
ON tags.key_name = k.name 
AND tags.value_name = v.name;
