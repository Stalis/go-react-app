-- +migrate Up
create table public.users
(
	id bigserial not null
		constraint users_pk
			primary key,
	username text not null 
        constraint users_username_uindex 
            unique,
	password_hash text not null,
	created_date timestamp with time zone default now()
);

-- +migrate Down
drop table public.users cascade;
