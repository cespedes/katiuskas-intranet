# Katiuskas API

There is a server running in the host `intranet.katiuskas.es`
which allows access to the SQL Database using a RESTful API.

This API will be the only accepted way to access all the stored data in the
database.

We have a helpers in the to access the API: `katiuskas-api`.

The base URL of the API is `htts://intranet.software.imdea.org/api/v1`.

General guidelines for RESTful APIs: https://opensource.zalando.com/restful-api-guidelines/

This is the version 1 of the API.  If there are any incompatible changes,
we will change the base URL with another number (`/api/v2`, `/api/v3` and so on).

In order to be authenticated, a query must have a token.  This token must sent
using an `Authorization: Bearer xxx` header.

The API will access a SQL database with the following sybsystems:

| name         | Description                                        | Status        |
|--------------|----------------------------------------------------|---------------|
| auth         | Authentication                                     | unimplemented |
| users        | List of users                                      | unfinished    |
| money        | Accounting                                         | unimplemented |

## Subsystems

### Authentication (`POST /auth`)

Authenticates user.

I am not sure yet how to use this endpoint.  It is currently not implemented.

## users

This is the main subsystem in the API.  It is used to retreive, create, update and delete information
about users in the system.

The main endpoints will be:

| Endpoint                            | Description                                     | Status        |
|-------------------------------------|-------------------------------------------------|---------------|
| `GET /user`                         | returns details of current user                 | unfinished    |
| `GET /users`                        | returns an array with the list of all the users | unimplemented |
| `POST /users`                       | creates a user (only for admins)                | unimplemented |
| `GET /users/<user>`                 | returns details about one specific user         | unimplemented |
| `PUT /users/<user>`                 | edits a user (only for admins)                  | unimplemented |
| `DELETE /users/<user>`              | deletes a user (only for admins)                | unimplemented |
|                                     |                                                 |               |
| `GET /users/<user>/roles`           | returns list of roles for a user                | unimplemented |
| `POST /users/<user>/roles/<role>`   | adds a role to a user                           | unimplemented |
| `DELETE /users/<user>/roles/<role>` | deletes a role from a user                      | unimplemented |

Pending:
- board
- socio
- federation
- baja temporal
- tokens

## money

We will have to access:
- accounts
- transactions

The main endpoints will be:

| Endpoint                                   | Description                                       | Status        |
|--------------------------------------------|---------------------------------------------------|---------------|
| `GET /money/accounts`                      | returns list of accounts                          | unimplemented |
| `GET /money/accounts/<account>`            | returns details of one account (its transactions) | unimplemented |
|                                            |                                                   |               |
| `POST /money/accounts`                     | creates an account                                | unimplemented |
| `DELETE /money/accounts/<account>`         | deletes an account                                | unimplemented |
|                                            |                                                   |               |
| `GET /money/transactions/<transaction>`    | returns details of one transaction                | unimplemented |
| `POST /money/transactions`                 | creates a transaction                             | unimplemented |
| `DELETE /money/transactions/<transaction>` | deletes a transaction                             | unimplemented |

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
  person_id integer NOT NULL REFERENCES person(id),
  alta      date NOT NULL,
  baja      date
);

CREATE TABLE federation (
  name      text NOT NULL PRIMARY KEY,
  id        integer
);

CREATE TABLE board (
  person_id integer NOT NULL REFERENCES person(id),
  position  text NOT NULL,
  start     date NOT NULL,
  end       date,
  CHECK (start < end)
);
```
