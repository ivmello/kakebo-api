CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL
);
CREATE UNIQUE INDEX users_email_idx ON users USING btree (email);

CREATE TABLE transactions (
	id uuid PRIMARY KEY,
	user_id int4 NOT NULL,
	description varchar(255) NOT NULL,
	amount int8 DEFAULT 0 NOT NULL,
	transaction_type varchar(20) NOT NULL,
	"date" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);