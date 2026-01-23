1) Semantica di CreateShortUrl
- Stessa original_url → stesso code.
- Se la URL esiste già, viene restituito il record esistente.
- Non supportiamo (per ora) più codici diversi per la stessa URL.

2) Generazione del codice
- Codice derivato da un hash della URL (più un salt di servizio).
- L’hash viene poi trasformato in un alfabeto tipo base62 e accorciato a una lunghezza fissa (es. 8 caratteri).
- È deterministico: stessa URL produce lo stesso candidato codice (a parità di salt).

3) Gestione collisioni
- Vincolo di unicità sul codice nel DB.
- Se c’è conflitto sul codice, si rigenera un nuovo codice cambiando un sale (es. un contatore).
- Numero massimo di tentativi, poi errore interno.
- Se c’è conflitto sull’URL, significa che la stessa URL è già stata inserita → si restituisce quella (idempotenza).