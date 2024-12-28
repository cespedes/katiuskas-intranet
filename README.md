# Katiuskas Intranet

There is a server running in the host `intranet.katiuskas.es`
which serves a few pages to administer Katiuskas internal pages.

This is a work in progress.  Some things are not implemented yet,
and the pages are still living next to the legacy version, which
is a server-side HTML pages generator written in Go.

From now on I will talk about the *next* version (not yet finished).

There are 3 kind of pages:
- static files (HTML, CSS, JavaScript)
- RESTful JSON API
- (maybe) some helper pages for authentication:
  - Google Authentication
  - Electronic certificate
  - Hash as sent via e-mail in mail+phone (or Telegram, or WhatsApp).
  - 2FA...

# Login

When the user first access the web page, and if she does not have a valid token,
a welcome page is shown with several login methods:

- Google Authentication

  When clicked, the page is redirected to Google's OAuth authorization endpoint,
  which authenticates the user and redirects to redirect URL
  (`/auth/google`).  That URL will get the `code` using JavaScript ahd feed it
  to a call to the API, which will in turn create a token. The code in JavaScript
  will then set a cookie with that token.

- Input with e-mail + phone.  If that e-mail and phone is registered
  in the database, an e-mail is sent with a URL (`/auth/hash`) including a hash.
  In that URL, the hash will be get using JavaScript and it will be sent to a call
  to the API, which will check it and create a token.

- Electronic certificate

  When clicked, the page is redirected to another server which will ask for
  a client certificate; if valid, it will create a valid hash and redirect
  to the "hash" auth page (`/auth/hash`).  This is the only place where a
  server-side page with root access to the API will need to be used.

# Katiuskas API

The server in the host `intranet.katiuskas.es`
which allows access to the SQL Database using a RESTful API.

This API will be the only accepted way to access all the stored data in the
database.

We have a helper command to access the API: `katiuskas-api`.

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
