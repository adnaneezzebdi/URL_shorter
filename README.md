Scopo del progetto: "URL shortener gRPC con analytics base".
Stack: Go, gRPC, PostgreSQL, Docker.

Sezione "MVP scope":
- creare short URL
- risolvere short URL
- leggere statistiche (hit count, created_at, last_access)

Setup locale (DB + migrazioni):
- copia `.env.example` in `.env` e compila i valori
- installa la CLI `migrate` (golang-migrate) per Windows o Linux
- `make migrate-up` per applicare lo schema
- `make migrate-status` per verificare la versione
