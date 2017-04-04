create table fb_info (
  timestamp timestamp with time zone,
  date_id integer,
  account text,
  name text,
  likes integer,
  talking_about_count integer,
  PRIMARY KEY (date_id,account)
);

create table fb_posts (
  timestamp timestamp with time zone,
  account text,
  id text,
  message text,
  type text,
  likes_count integer,
  comments_count integer,
  shares_count integer,
  PRIMARY KEY (id)
);
