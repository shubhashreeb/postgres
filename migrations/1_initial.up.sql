CREATE TABLE mytable (
    id SERIAL PRIMARY KEY,
    msg character varying(3000),
    data JSONB
);

INSERT INTO mytable (id, msg, data) VALUES (1, 'This is Aaron message1', '{"Name": "Aaron","age": 6}');
INSERT INTO mytable (id, msg, data) VALUES (2, 'This is Aaron message2', '{"Name": "Aaron","age": 6}');

CREATE OR REPLACE FUNCTION mytable_changes()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('mytable_changes', TG_OP || '|' || NEW.id || '|' || NEW.msg || '|' || to_jsonb(NEW));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER mytable_trigger
AFTER INSERT OR UPDATE OR DELETE ON mytable
FOR EACH ROW
EXECUTE FUNCTION mytable_changes();


