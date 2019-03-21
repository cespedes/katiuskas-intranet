# PostgreSQL schema

```
CREATE DOMAIN gender AS character(1)
  CHECK (VALUE = ANY (ARRAY['', 'M', 'F']));

CREATE TABLE person (
  id            serial PRIMARY KEY,
  name          text DEFAULT '---' NOT NULL,
  surname       text DEFAULT '---' NOT NULL,
  dni           text DEFAULT '' NOT NULL,
  birth         date,
  address       text DEFAULT '' NOT NULL,
  zip           text DEFAULT '' NOT NULL,
  city          text DEFAULT '' NOT NULL,
  province      text DEFAULT '' NOT NULL,
  gender        gender,
  emerg_contact text DEFAULT '' NOT NULL,
  lopd          date
);

CREATE TABLE socio (
  id        serial PRIMARY KEY,
  id_person integer NOT NULL REFERENCES person(id),
  alta      date NOT NULL,
  baja      date
);

CREATE TABLE federation (
  name      text NOT NULL PRIMARY KEY,
  id        integer
);

CREATE TABLE board2 (
  id_person integer NOT NULL REFERENCES person(id),
  position  text NOT NULL,
  start     date NOT NULL,
  end       date,
  CHECK (start < end)
);
```
