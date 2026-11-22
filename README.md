# Raspored sati za 2. RT

Privatna instanca rasporeda sati, namijenjena lakšem dohvaćanju i prikazu
rasporeda razreda 2. RT. Planirana adresa: `https://schedule.noahusicdev.site`.

## Autorstvo i vlasništvo

Ovaj projekt temelji se na [`tsd-timetable`](https://github.com/h00s-go/tsd-timetable),
koji je izradio h00s-go. Izvorni projekt, njegov izvorni kod i podatkovni izvor
pripadaju h00s-go; ovaj repozitorij nije tvrdnja o njihovu autorstvu ili
vlasništvu. Koristim ga samo za lakše dohvaćanje i prikaz rasporeda.

Na zahtjev h00s-go uklonit ću ovu kopiju i prestati je koristiti.

## Pokretanje frontenda

Frontend je u direktoriju `frontend/` i dohvaća podatke iz
`/api/v1/schedules/current`. Za lokalni razvoj pokreni backend na
`http://localhost:3000`, a zatim u `frontend/` pokreni:

```sh
bun install
bun run dev
```

Više detalja o frontend provjerama i izgradnji nalazi se u
[frontend/README.md](./frontend/README.md).
