CREATE TABLE IF NOT EXISTS urls(
    url_id SERIAL PRIMARY KEY,
    url_string VARCHAR (500) UNIQUE NOT NULL,
    url_method VARCHAR (50) NOT NULL,
    time_interval INT NOT NULL,
    unix_time_added INT NOT NULL
);

CREATE TABLE IF NOT EXISTS checks
(
    check_id SERIAL PRIMARY KEY,
    url_id   SERIAL,
    status_code INT NOT NULL,
    unix_time_added INT NOT NULL,
    FOREIGN KEY (url_id) REFERENCES urls(url_id)
);