CREATE TABLE config (
	access_token varchar,
	room varchar,
	session_cookie varchar,
	leaderboard_url varchar
);
CREATE TABLE IF NOT EXISTS "days" (
	user_id varchar references users(id),
	day int,
	part int,
	get_star_ts varchar,
	primary key (user_id, day, part)
);
CREATE TABLE IF NOT EXISTS "users" (
	id varchar primary key,
	name varchar,
	local_score integer not null default 0
);
