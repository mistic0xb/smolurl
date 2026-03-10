-- Write your migrate up statements here
CREATE OR REPLACE FUNCTION camel(input_row anyelement)
    RETURNS jsonb
    LANGUAGE plpgsql
    AS $$
DECLARE
    result jsonb := '{}';
    rec record;
BEGIN
    FOR rec IN
    SELECT
        lower(substring(regexp_replace(initcap(regexp_replace(key, '_', ' ', 'g')), '\s', '', 'g'), 1, 1)) || substring(regexp_replace(initcap(regexp_replace(key, '_', ' ', 'g')), '\s', '', 'g'), 2) AS camel_key,
        value
    FROM
        jsonb_each(to_jsonb(input_row))
        LOOP
            result := result || jsonb_build_object(rec.camel_key, rec.value);
        END LOOP;
    RETURN result;
END;
$$;

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
