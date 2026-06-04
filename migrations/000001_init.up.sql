CREATE TABLE users (
    user_id character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    username character varying(100) NOT NULL,
    password_hash character varying(255) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (user_id)
);

CREATE UNIQUE INDEX idx_users_username ON users (username);

CREATE TABLE tags (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    description character varying(500),
    user_id character varying(255),
    CONSTRAINT tags_pkey PRIMARY KEY (id)
);

CREATE SEQUENCE tags_id_seq
    AS bigint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE tags_id_seq OWNED BY tags.id;
ALTER TABLE ONLY tags ALTER COLUMN id SET DEFAULT nextval('tags_id_seq'::regclass);

CREATE UNIQUE INDEX idx_tags_user_name ON tags (user_id, name);

CREATE TABLE investment (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    balance double precision DEFAULT 0 NOT NULL,
    created_at timestamp with time zone NOT NULL,
    tags text[] DEFAULT '{}'::text[] NOT NULL,
    initial_balance double precision NOT NULL,
    user_id character varying(255),
    CONSTRAINT investment_pkey PRIMARY KEY (id)
);

CREATE SEQUENCE investment_id_seq
    AS bigint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE investment_id_seq OWNED BY investment.id;
ALTER TABLE ONLY investment ALTER COLUMN id SET DEFAULT nextval('investment_id_seq'::regclass);

CREATE INDEX idx_investment_user_id ON investment (user_id);

CREATE TABLE investment_history (
    id bigint NOT NULL,
    investment_id bigint NOT NULL,
    date timestamp with time zone NOT NULL,
    amount double precision NOT NULL,
    movement_type character varying(50) NOT NULL,
    balance_after double precision NOT NULL,
    CONSTRAINT investment_history_pkey PRIMARY KEY (id),
    CONSTRAINT fk_investment_history FOREIGN KEY (investment_id) REFERENCES investment (id)
);

CREATE SEQUENCE investment_history_id_seq
    AS bigint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE investment_history_id_seq OWNED BY investment_history.id;
ALTER TABLE ONLY investment_history ALTER COLUMN id SET DEFAULT nextval('investment_history_id_seq'::regclass);

CREATE INDEX idx_investment_history_investment_id ON investment_history (investment_id);
