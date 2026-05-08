-- +goose up
create table group_files (
	id uuid primary key,
	type varchar(50) not null,
	file_path varchar(200) not null,
	created_at timestamp not null,
	group_id uuid references pdf_metadata(group_id) on delete cascade
);

-- -goose down
drop table group_files;
