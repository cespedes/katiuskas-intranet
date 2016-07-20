* Different pages for these users:
  + NoUser
  + NoSocio
  + ExSocio
  + SocioBajaTemporal
  + SocioActivo
  + SocioJunta
* Pages
  + root (main menu)
  + my info
  + lista de socios
  + junta directiva
  + money

```
CREATE TABLE place_type (
  id          SERIAL PRIMARY KEY,
  type        TEXT
);

CREATE TABLE place (
  id            SERIAL PRIMARY KEY,
  place_type_id INTEGER REFERENES place_type(id),
  name          TEXT
);

CREATE TABLE activity (
  id          SERIAL PRIMARY KEY,
  organizer   INTEGER REFERENCES person(id),
  date_begin  DATE,
  date_end    DATE,
  state       INTEGER,
  title       TEXT
);

CREATE TABLE activity_place (
  activity_id INTEGER NOT NULL REFERENCES activity(id),
  place_id    INTEGER NOT NULL REFERENCES place(id)
);

CREATE TABLE activity_person (
  activity_id INTEGER NOT NULL REFERENCES activity(id),
  person_id   INTEGER NOT NULL REFERENCES person(id)
);

CREATE TABLE equipment (
  id   SERIAL PRIMARY KEY,
  name TEXT,
  cost NUMERIC
);

CREATE TABLE activity_equipment (
  id           SERIAL PRIMARY KEY,
  activity_id  INTEGER NOT NULL REFERENCES activity(id),
  equipment_id INTEGER NOT NULL REFERENCES equipment(id),
  quantity     INTEGER NOT NULL DEFAULT 1,
  date_out     DATE NOT NULL DEFAULT now(),
  date_in      DATE
);

CREATE TABLE activity_report (
  id          SERIAL PRIMARY KEY,
  activity_id INTEGER NOT NULL REFERENCES activity(id),
  report      TEXT
);

CREATE TABLE activity_file (
  id          SERIAL PRIMARY KEY,
  activity_id INTEGER NOT NULL REFERENCES activity(id),
  filename    TEXT
);
```

* Última alta/baja:
  `select * from socio natural join (select id_person,max(alta) as alta from socio group by id_person) a;`
