create table google_search (
  timestamp timestamp with time zone,
  date_id integer,
  query text,
  count integer,
  PRIMARY KEY (date_id,query)
);

create table google_trends (
	timestamp timestamp with time zone,
	date_id integer,
	query text,
	result text,
	PRIMARY KEY (date_id,query)
);

create table google_news (
  timestamp timestamp with time zone,
  date_id integer,
  query text,
  result text,
  PRIMARY KEY (date_id,query)
);

create table google_first_page (
  timestamp timestamp with time zone,
  query text,
  link text,
  page text,
  PRIMARY KEY (query,link)
);