-- +migrate Up
create extension if not exists "uuid-ossp";

create table public.sessions
(
	id bigserial not null
		constraint sessions_pk
			primary key,
	token uuid default uuid_generate_v4() not null 
        constraint sessions_token_uindex
            unique,
	user_id bigint not null
		constraint sessions_users_id_fk
			references public.users
				on delete cascade,
	created_date timestamp with time zone default now() not null,
	expired_date timestamp with time zone default (now() + '01:00:00'::interval) not null
);

create index sessions_user_id_index
	on public.sessions (user_id);

-- +migrate Down
drop table public.sessions cascade;
drop extension "uuid-ossp";