-- +goose Up
-- +goose StatementBegin
CREATE TABLE albums (
    id int NOT NULL UNIQUE AUTO_INCREMENT,
    title varchar(255) NOT NULL,
    artist varchar (255) NOT NULL,
    price float NOT NULL,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE albums;
-- +goose StatementEnd
