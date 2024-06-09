CREATE TABLE account_types (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE subscription_costs (
  id SERIAL PRIMARY KEY,
  subscription_type VARCHAR(50) NOT NULL,
  subscription_day INTEGER NOT NULL,
  cost DECIMAL(10, 2) NOT NULL,
  UNIQUE (subscription_type, subscription_day)
);

CREATE TABLE payment_types (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE status (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  lastname1 VARCHAR(255) NOT NULL,
  lastname2 VARCHAR(255),
  email VARCHAR(255) UNIQUE NOT NULL,
  phone VARCHAR(50),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT email_not_empty CHECK (email <> '')
);

CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  account_id UUID UNIQUE NOT NULL,
  account_type_id INTEGER NOT NULL,
  subscription_id INTEGER,
  status_id INTEGER,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE subscriptions (
  id SERIAL PRIMARY KEY,
  account_id UUID NOT NULL,
  subscription_cost_id INTEGER NOT NULL,
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  status_id INTEGER
);

CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  account_id UUID NOT NULL,
  payment_type_id INTEGER NOT NULL,
  cost DECIMAL(10, 2) NOT NULL,
  status_id INTEGER,
  payment_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (account_id) REFERENCES accounts (account_id),
  FOREIGN KEY (payment_type_id) REFERENCES payment_types (id),
  FOREIGN KEY (status_id) REFERENCES status (id)
);

ALTER TABLE accounts
ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id),
ADD CONSTRAINT fk_account_type_id FOREIGN KEY (account_type_id) REFERENCES account_types (id),
ADD CONSTRAINT fk_subscription_id FOREIGN KEY (subscription_id) REFERENCES subscriptions (id);

ALTER TABLE subscriptions
ADD CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (account_id),
ADD CONSTRAINT fk_subscription_cost_id FOREIGN KEY (subscription_cost_id) REFERENCES subscription_costs (id),
ADD CONSTRAINT fk_status_id FOREIGN KEY (status_id) REFERENCES status (id);