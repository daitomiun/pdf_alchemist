-- +goose Up
create table pdf_metadata (
	id UUID primary key,
	created_at TIMESTAMP NOT NULL,
	file_path varchar(200) NOT NULL,
	group_id UUID NULL
);

-- +goose Down
drop table pdf_metadata;
