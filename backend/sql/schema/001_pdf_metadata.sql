-- +goose up
create table pdf_metadata (
	id UUID primary key,
	user_id UUID NULL,
	created_at TIMESTAMP NOT NULL,
	file_path varchar(200) NOT NULL,
	group_id UUID NULL references pdf_groups(ID) on delete cascade,
	user_id UUID NULL references users(ID) on delete cascade,
);

create table pdf_groups (
	id UUID primary key,
);

create table group_files (
	id UUID primary key,
	type varchar(50) not null,
	file_path varchar(200) not null,
	group_id UUID NULL references pdf_groups(ID) on delete cascade,
	created_at TIMESTAMP NOT NULL,
);

-- +goose down
DROP TABLE pdf_metadata;
