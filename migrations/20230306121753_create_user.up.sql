CREATE TABLE account(
  id uuid DEFAULT public.gen_random_uuid() NOT NULL,
  firstName character varying NOT NULL,
  lastName character varying NOT NULL,
  balance float DEFAULT 10000.00
)