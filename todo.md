# TODO / README – gRPC URL Shortener + Analytics

Questo documento descrive **l’intero flusso di lavoro del progetto**, dall’inizializzazione fino ai criteri di completamento. È pensato come guida operativa e checklist condivisa.

---

## 1. Scopo del progetto

Costruire un microservizio backend che esponga **solo API gRPC** per:

- creazione di short URL
- risoluzione di short URL
- consultazione statistiche di utilizzo

Il progetto è volutamente focalizzato su:

- progettazione architetturale
- correttezza sotto concorrenza
- gestione professionale del database
- deploy riproducibile locale + cloud

Nessun frontend.

---

## 2. Stack e vincoli

### Stack tecnologico

- Linguaggio: Go
- Comunicazione: gRPC
- Database: PostgreSQL
- Ambiente locale: Docker + Docker Compose
- Cloud: Render (Postgres gestito)

### Vincoli chiave

- Nessun accesso HTTP/REST
- Nessuna modifica manuale al DB
- Tutte le modifiche DB tramite migrazioni versionate
- Incrementi atomici direttamente nel DB
- Error handling gRPC coerente e documentato
- Codice organizzato per responsabilità

---

## 3. Struttura del repository (concettuale)

- /proto → contratti gRPC
- /cmd → entrypoint del server
- /internal
  - api → layer gRPC
  - service → business logic
  - repository → accesso DB

- /migrations → migrazioni SQL versionate
- /docker → Dockerfile e compose
- README.md
- DECISIONS.md

---

## 4. Workflow di sviluppo

### Regole di lavoro

- Commit piccoli e frequenti
- Ogni modifica significativa passa da PR
- PR revisionata anche se si lavora in due
- Ogni step deve produrre un risultato verificabile

---

## 5. Sprint 1 – Fondazioni

### 5.1 Decision log iniziale

Creare `DECISIONS.md` con almeno:

- strategia generazione short code
- semantica di Create (idempotente o no)
- strategia di incremento atomico

Ogni decisione deve includere:

- alternative scartate
- motivazione tecnica

---

### 5.2 Database locale

Checklist:

- Docker Compose con Postgres
- Variabili ambiente per connessione
- DB avviabile con un comando

Verifica:

- connessione al DB funzionante

---

### 5.3 Migrazioni iniziali

Creare prima migrazione:

- tabella `short_urls`
- colonne:
  - id
  - code (UNIQUE)
  - original_url
  - created_at
  - access_count (default 0)
  - last_access_at (nullable)

Requisiti:

- indice su `code`
- migrazione applicabile più volte senza errori

Verifica:

- migrazione applicata su DB locale
- schema verificato manualmente

---

### 5.4 Database su Render

Checklist:

- creare Postgres su Render
- recuperare credenziali
- configurare SSL se richiesto
- applicare **le stesse migrazioni** del locale

Verifica:

- versione schema identica a locale

---

### 5.5 gRPC contract

Definire service con RPC:

- CreateShortUrl
- ResolveShortUrl
- GetStats

Linee guida:

- request/response minimali ma complete
- tipi coerenti
- campi obbligatori chiari

Error policy:

- INVALID_ARGUMENT → input non valido
- NOT_FOUND → codice inesistente
- ALREADY_EXISTS → collisione
- INTERNAL / UNAVAILABLE → DB

Verifica:

- proto compilabile
- contract stabile e documentato

---

## 6. Sprint 2 – Implementazione MVP

### 6.1 CreateShortUrl

Flusso:

- validazione URL
- generazione codice
- tentativo insert
- gestione collisioni (retry o errore)

Checklist:

- nessuna query non necessaria
- collisioni gestite in modo esplicito

Verifica:

- stessa URL → comportamento documentato
- collisioni simulate

---

### 6.2 ResolveShortUrl (punto critico)

Flusso:

- lookup per code
- se non esiste → NOT_FOUND
- incremento `access_count` atomico
- update `last_access_at`

Vincoli:

- incremento fatto **nel DB**
- nessun read-modify-write in memoria

Verifica:

- test concorrenti
- nessuna perdita di incrementi

---

### 6.3 GetStats

Flusso:

- lookup per code
- ritorno metadati

Vincoli:

- nessuna mutazione

Verifica:

- dati coerenti dopo accessi multipli

---

## 7. Docker e avvio progetto

Checklist:

- Dockerfile per API
- Docker Compose con API + Postgres
- variabili ambiente documentate

Comandi documentati:

- avvio stack
- applicazione migrazioni

Verifica:

- servizio avviabile da zero

---

## 8. Test manuali e concorrenti

- Create → Resolve → Stats
- Resolve concorrenti sullo stesso codice
- verifica contatori

---

## 9. Criteri di completamento

Il progetto è considerato completo se:

- stack avviabile via Docker Compose
- migrazioni applicate locale e Render
- tutte le RPC funzionanti
- incremento accessi corretto sotto concorrenza
- README sufficiente per terzi

---

## 10. Regola finale

Ogni scelta deve poter essere spiegata.
Se non è spiegabile, probabilmente è fragile.
