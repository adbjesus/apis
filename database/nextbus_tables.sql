create table stops (
  timestamp timestamp with time zone,
  agency text,
  route text,
  tag text,
  title text,
  shorttitle text,
  lat float,
  lon float,
  stopid text,
  PRIMARY KEY (tag)
);

create table vehicle_locations (
  timestamp timestamp with time zone,
  agency text,
  route text,
  bus_id text,
  direction text,
  lat float,
  lon float,
  predictable boolean
);

create table predictions (
  timestamp timestamp with time zone,
  agency text,
  route text,
  bus_id text,
  stop text,
  direction text,
  seconds integer,
  FOREIGN KEY (stop) REFERENCES stops(tag)
)
