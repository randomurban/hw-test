create table event (
    id serial primary key,
    title varchar(100),
    start_time timestamp with time zone,
    end_time timestamp with time zone,
    owner int,
    description varchar(200),
    notice_time timestamp with time zone
);