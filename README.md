# Raspored sati za 2. RT v1.0

Raspored sati namijenjen lakšem dohvaćanju i prikazu
rasporeda razreda 2. RT. Planirana adresa: `https://schedule.noahusicdev.site`, ali trenutno je objavljeno na https://schedule-tsd.netlify.app/.

## Autorstvo i vlasništvo

Ovaj projekt temelji se na [`tsd-timetable`](https://github.com/h00s-go/tsd-timetable),
koji je izradio h00s-go. Izvorni projekt, njegov izvorni kod i podatkovni izvor
pripadaju h00s-go; ovaj repozitorij nije tvrdnja o njihovu autorstvu ili
vlasništvu. Koristim ga samo za lakše dohvaćanje i prikaz rasporeda.

## Pokretanje

Potrebni su [Bun](https://bun.sh) i Node.js.

```bash
cd frontend
bun install
bun run dev
```

Razvojni server (Vite) prosljeđuje `/api` na javni API
`https://tsd-timetable.husak.me`, pa nije potreban vlastiti backend. Za lokalni
Go backend pokreni `API_TARGET=http://localhost:3000 bun run dev`.
