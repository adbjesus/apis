create table twitter_search (
  timestamp timestamp with time zone,
  query text,
  id text,
  post text,
  retweetCount integer,
  favoriteCount integer,
  PRIMARY KEY (id,query)
);

create table twitter_timeline (
  timestamp timestamp with time zone,
  screenname text,
  id text,
  post text,
  retweetCount integer,
  favoriteCount integer,
  PRIMARY KEY (id,screenname)
);
