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
CREATE TABLE activity (
  id          SERIAL PRIMARY KEY,
  organizer   INTEGER REFERENCES person(id),
  date        DATE,
  type        TEXT,
  place       TEXT,
  description TEXT,
);

CREATE TABLE activity_person (
  id_activity INTEGER NOT NULL REFERENCES activity(id),
  id_person INTEGER NOT NULL REFERENCES person(id),
);
```
