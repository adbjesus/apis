create table bing_search (
  timestamp timestamp with time zone,
  date_id integer,
  query text,
  count integer,
  PRIMARY KEY (date_id,query)
);
