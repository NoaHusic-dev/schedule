# Raspored sati za 2. RT

Svelte aplikacija prikazuje dnevni raspored samo za razred 2. RT. Raspored i
podaci o predmetima, nastavnicima i učionicama dohvaćaju se iz API-ja projekta
`tsd-timetable` na `/api/v1/schedules/current`.

## Pokretanje

Pokreni backend projekta na `http://localhost:3000`, zatim u ovom direktoriju:

```sh
bun install
bun run dev
```

Vite prosljeđuje `/api` zahtjeve backendu. Za provjeru tipova i produkcijsku
izgradnju pokreni `bun run check` i `bun run build`.
