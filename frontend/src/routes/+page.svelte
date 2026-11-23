<script lang="ts">
  import logo from "#lib/assets/tsd-logo-dark.png";
  import { formatClock, formatDate, isoDate } from "#lib/helpers/format.ts";
  import { formatPeriods, periodIds } from "#lib/helpers/timetable.ts";
  import { clock } from "#lib/stores/clock.svelte.ts";
  import { schedules } from "#lib/stores/schedules.svelte.ts";

  const classSchedule = $derived(schedules.get("2 RT"));
  const isToday = $derived(schedules.date === isoDate(clock.now));
  const slots = $derived(isToday ? (classSchedule?.schedule ?? []) : []);
  const isFinished = $derived(
    isToday &&
      slots.length > 0 &&
      clock.now.getTime() > new Date(slots.at(-1)?.end ?? 0).getTime(),
  );
  const activeIndex = $derived(
    isToday ? slots.findIndex((slot) => {
      const now = clock.now.getTime();
      return now >= new Date(slot.start).getTime() && now < new Date(slot.end).getTime();
    }) : -1,
  );
  const nextIndex = $derived(
    isToday ? slots.findIndex((slot) => new Date(slot.start).getTime() > clock.now.getTime()) : -1,
  );


  $effect(() => clock.subscribe());
  $effect(() => schedules.subscribe());
</script>

<svelte:head>
  <title>2. RT · Raspored sati</title>
  <meta
    name="description"
    content="Dnevni raspored sati za razred 2. RT Tehničke škole Daruvar."
  />
</svelte:head>

<div class="page-shell">
  <header class="topbar">
    <a class="brand" href="/" aria-label="Raspored sati za 2. RT">
      <img src={logo} alt="" />
      <span>Tehnička škola Daruvar</span>
    </a>
    <div class="live-clock" aria-label="Trenutno vrijeme">
      <time>{formatClock(clock.now)}</time>
    </div>
  </header>

  <main>
    <section class="hero" aria-labelledby="page-title">
      <div class="hero-copy">
        <p class="eyebrow">Raspored nastave</p>
        <h1 id="page-title">2. <span>RT</span></h1>
        <p class="date">{formatDate(clock.now)}</p>
      </div>
      <div class="class-mark" aria-hidden="true">RT</div>
    </section>

    <section class="schedule-section" aria-labelledby="schedule-title">
      <div class="section-heading">
        <div>
          <p class="eyebrow">Dnevni pregled</p>
          <h2 id="schedule-title">Raspored za danas</h2>
        </div>
        {#if schedules.loaded && isToday && isFinished}
          <span class="slot-count finished">Završena nastava</span>
        {/if}
      </div>

      {#if !schedules.loaded}
        <div class="message-card" role="status">
          <span class="loader" aria-hidden="true"></span>
          <p>Učitavam raspored…</p>
        </div>
      {:else if !isToday}
        <div class="message-card" role="status">
          <span class="message-icon" aria-hidden="true">↻</span>
          <p>Danas raspored još nije dostupan.</p>
          <span class="message-detail">Podaci se automatski osvježavaju.</span>
        </div>
      {:else if slots.length === 0}
        <div class="message-card" role="status">
          <span class="message-icon" aria-hidden="true">✓</span>
          <p>Nema upisanih sati za 2. RT.</p>
          <span class="message-detail">Odabrani dan je slobodan.</span>
        </div>
      {:else}
        <ol class="timeline">
          {#each slots as slot, index (`${slot.start}-${slot.end}`)}
            {@const periods = formatPeriods(periodIds(slot.start, slot.end))}
            {@const state = index === activeIndex ? "U tijeku" : index === nextIndex ? "Sljedeći" : ""}
            <li class:active={index === activeIndex} class="timeline-item">
              <div class="time-column">
                <time datetime={slot.start}>{new Date(slot.start).toLocaleTimeString("hr-HR", { hour: "2-digit", minute: "2-digit" })}</time>
                <time class="end-time" datetime={slot.end}>{new Date(slot.end).toLocaleTimeString("hr-HR", { hour: "2-digit", minute: "2-digit" })}</time>
              </div>
              <div class="timeline-rail" aria-hidden="true">
                <span class="rail-dot"></span>
              </div>
              <article class="lesson-card">
                <div class="lesson-heading">
                  <div class="lesson-labels">
                    {#if periods}
                      <span class="period-label">Sat {periods}</span>
                    {/if}
                    {#if state}
                      <span class:current={state === "U tijeku"} class="state-label">{state}</span>
                    {/if}
                  </div>
                  {#if slot.lessons.length > 1}
                    <span class="group-count">{slot.lessons.length} grupe</span>
                  {/if}
                </div>
                <div class="lesson-list">
                  {#each slot.lessons as lesson, lessonIndex (`${lesson.subject.short}-${lesson.groups?.join("-") ?? lessonIndex}`)}
                    <div class="lesson-detail" class:group-a={lesson.groups?.some((group) => group.trim().toUpperCase() === "A")}>
                      <div class="subject">
                        {#if lesson.groups?.length}
                          <span class="group-name">{lesson.groups.join(" · ")}</span>
                        {/if}
                        <h3>{lesson.subject.name}</h3>
                      </div>
                      <div class="lesson-meta">
                        {#if lesson.teacher.short}
                          <span class="teacher">
                            <span class="meta-icon" aria-hidden="true">◉</span>
                            {lesson.teacher.short}
                          </span>
                        {/if}
                        {#if lesson.room}
                          <span class="room">
                            {lesson.room}
                          </span>
                        {/if}
                      </div>
                    </div>
                  {/each}
                </div>
              </article>
            </li>
          {/each}
        </ol>
      {/if}
    </section>

    <footer class="page-footer">
      <span>
        Projekat napravio <a href="https://noahusicdev.site" target="_blank" rel="noreferrer">Noa Husić</a>
        <span class="footer-divider">·</span>
        <a href="https://github.com/noahusic-dev" target="_blank" rel="noreferrer">GitHub</a>
      </span>
      <span>
        Napravljeno pomoću projekta
        <a href="https://github.com/h00s" target="_blank" rel="noreferrer">h00s</a>
      </span>
    </footer>
  </main>
</div>

<style>
  .page-shell {
    min-height: 100vh;
    min-height: 100dvh;
    overflow: hidden;
    background:
      radial-gradient(ellipse at 82% 4%, rgb(91 153 255 / 12%), transparent 29rem),
      #f4f6fb;
  }

  .topbar,
  main {
    width: min(100% - 3rem, 1024px);
    margin-inline: auto;
  }

  .topbar {
    display: flex;
    min-height: 76px;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid #e4e9f2;
  }

  .brand,
  .live-clock {
    display: inline-flex;
    align-items: center;
  }

  .brand {
    gap: 0.75rem;
    color: #263248;
    font-size: 0.9rem;
    font-weight: 700;
    text-decoration: none;
  }

  .brand img {
    width: 38px;
    height: 38px;
    object-fit: contain;
  }

  .live-clock {
    gap: 0.55rem;
    color: #536178;
    font-size: 0.9rem;
    font-variant-numeric: tabular-nums;
    font-weight: 650;
  }

  main {
    padding-top: clamp(2rem, 6vw, 4.5rem);
  }

  .hero {
    position: relative;
    display: flex;
    min-height: 180px;
    align-items: center;
    justify-content: space-between;
    overflow: hidden;
    padding: clamp(1.25rem, 4vw, 2.5rem);
    border-radius: 26px;
    background: linear-gradient(118deg, #172c4d 0%, #234b80 65%, #3279bd 100%);
    box-shadow: 0 20px 50px rgb(25 54 94 / 16%);
    color: white;
  }

  .hero::after {
    position: absolute;
    right: 12%;
    bottom: -180px;
    width: 370px;
    height: 370px;
    border: 1px solid rgb(255 255 255 / 13%);
    border-radius: 50%;
    box-shadow: 0 0 0 38px rgb(255 255 255 / 4%), 0 0 0 78px rgb(255 255 255 / 4%);
    content: "";
  }

  .hero-copy {
    z-index: 1;
  }

  .eyebrow {
    margin: 0 0 0.5rem;
    color: #74839b;
    font-size: 0.72rem;
    font-weight: 800;
    letter-spacing: 0.14em;
    text-transform: uppercase;
  }

  .hero .eyebrow {
    color: #b6d6ff;
  }

  h1 {
    margin: 0;
    font-size: clamp(2.8rem, 7vw, 4.5rem);
    font-weight: 800;
    letter-spacing: -0.07em;
    line-height: 1;
  }

  h1 span {
    color: #92c7ff;
  }

  .date {
    margin: 1rem 0 0;
    color: #d3e2f5;
    font-size: 1rem;
    font-weight: 500;
    text-transform: capitalize;
  }

  .class-mark {
    z-index: 1;
    display: grid;
    width: clamp(88px, 14vw, 120px);
    aspect-ratio: 1;
    place-items: center;
    border: 1px solid rgb(255 255 255 / 23%);
    border-radius: 32px;
    background: rgb(255 255 255 / 10%);
    color: rgb(255 255 255 / 85%);
    font-size: clamp(2rem, 5vw, 3.2rem);
    font-weight: 800;
    letter-spacing: -0.06em;
    transform: rotate(7deg);
  }

  .schedule-section {
    padding-top: clamp(2.5rem, 6vw, 4rem);
  }

  .section-heading {
    display: flex;
    align-items: end;
    justify-content: space-between;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  h2 {
    margin: 0;
    color: #202c41;
    font-size: clamp(1.45rem, 3vw, 1.9rem);
    font-weight: 750;
    letter-spacing: -0.04em;
  }

  .slot-count {
    padding: 0.45rem 0.75rem;
    border: 1px solid #e3e9f2;
    border-radius: 999px;
    background: white;
    color: #63718a;
    font-size: 0.8rem;
    font-weight: 650;
    white-space: nowrap;
  }

  .slot-count.finished {
    border-color: #cde8da;
    background: #eff8f3;
    color: #2f7656;
  }

  .timeline {
    display: grid;
    gap: 0.85rem;
    margin: 0;
    padding: 0;
    list-style: none;
  }

  .timeline-item {
    display: grid;
    grid-template-columns: 76px 24px minmax(0, 1fr);
    gap: 0.85rem;
    align-items: stretch;
  }

  .time-column {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    padding-top: 1.25rem;
    color: #34425a;
    font-size: 0.9rem;
    font-variant-numeric: tabular-nums;
    font-weight: 750;
  }

  .end-time {
    margin-top: 0.25rem;
    color: #8a96a9;
    font-size: 0.78rem;
    font-weight: 550;
  }

  .timeline-rail {
    position: relative;
    display: flex;
    justify-content: center;
  }

  .timeline-item:not(:last-child) .timeline-rail::after {
    position: absolute;
    top: 26px;
    bottom: -18px;
    width: 2px;
    background: #dce3ee;
    content: "";
  }

  .rail-dot {
    z-index: 1;
    width: 12px;
    height: 12px;
    margin-top: 1.45rem;
    border: 3px solid #f4f6fb;
    border-radius: 50%;
    background: #a9b5c8;
    box-shadow: 0 0 0 1px #dce3ee;
  }

  .active .rail-dot {
    background: #3279bd;
    box-shadow: 0 0 0 4px rgb(50 121 189 / 16%);
  }

  .lesson-card {
    min-width: 0;
    overflow: hidden;
    border: 1px solid #e6eaf1;
    border-radius: 16px;
    background: white;
    box-shadow: 0 5px 16px rgb(31 51 84 / 4%);
  }

  .active .lesson-card {
    border-color: #9dc8f0;
    box-shadow: 0 8px 24px rgb(50 121 189 / 10%);
  }

  .lesson-heading {
    display: flex;
    min-height: 45px;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    padding: 0.55rem 1rem;
    border-bottom: 1px solid #f0f2f6;
  }

  .lesson-labels,
  .lesson-meta {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.5rem;
  }

  .period-label,
  .state-label,
  .group-count {
    color: #718099;
    font-size: 0.72rem;
    font-weight: 750;
  }

  .period-label {
    text-transform: uppercase;
  }

  .state-label {
    padding: 0.2rem 0.48rem;
    border-radius: 999px;
    background: #f0f3f8;
    color: #5c6b81;
  }

  .state-label.current {
    background: #e7f3ff;
    color: #2265a3;
  }

  .group-count {
    color: #818da0;
  }

  .lesson-detail {
    display: flex;
    min-height: 70px;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.85rem 1rem;
  }

  .lesson-detail.group-a {
    border-left: 4px solid #3279bd;
    background: linear-gradient(90deg, rgb(50 121 189 / 8%), white 32%);
  }

  .lesson-detail.group-a .group-name {
    display: inline-block;
    padding: 0.18rem 0.48rem;
    border-radius: 999px;
    background: #3279bd;
    color: white;
  }

  .lesson-detail + .lesson-detail {
    border-top: 1px solid #f0f2f6;
  }

  .subject {
    min-width: 0;
  }

  .subject h3 {
    overflow: hidden;
    margin: 0;
    color: #26344b;
    font-size: 1rem;
    font-weight: 700;
    text-overflow: ellipsis;
  }

  .group-name {
    display: block;
    margin-bottom: 0.2rem;
    color: #3279bd;
    font-size: 0.72rem;
    font-weight: 700;
  }

  .lesson-meta {
    justify-content: flex-end;
    color: #68768b;
    font-size: 0.8rem;
  }

  .teacher,
  .room {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
  }

  .room {
    padding: 0.35rem 0.55rem;
    border-radius: 7px;
    background: #f1f5fa;
    color: #4a5b74;
    font-weight: 650;
  }

  .meta-icon {
    color: #91a0b5;
    font-size: 0.76rem;
  }

  .message-card {
    display: flex;
    min-height: 170px;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.55rem;
    padding: 2rem;
    border: 1px solid #e6eaf1;
    border-radius: 18px;
    background: white;
    color: #34425a;
    text-align: center;
  }

  .message-card p {
    margin: 0;
    font-size: 1rem;
    font-weight: 700;
  }

  .message-detail {
    color: #7d899d;
    font-size: 0.85rem;
  }

  .message-icon {
    display: grid;
    width: 36px;
    height: 36px;
    margin-bottom: 0.25rem;
    place-items: center;
    border-radius: 12px;
    background: #edf4fc;
    color: #3279bd;
    font-size: 1.2rem;
  }

  .loader {
    width: 24px;
    height: 24px;
    border: 3px solid #dfe9f5;
    border-top-color: #3279bd;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .page-footer {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    padding: 2.5rem 0 1.5rem;
    color: #8490a2;
    font-size: 0.75rem;
  }

  .page-footer a {
    color: #64748b;
    font-weight: 700;
    text-decoration: none;
  }

  .page-footer a:hover {
    color: #3279bd;
  }

  .footer-divider {
    padding-inline: 0.2rem;
    color: #bac3d0;
  }

  @media (max-width: 600px) {
    .topbar,
    main {
      width: min(100% - 1.5rem, 1024px);
    }

    .topbar {
      min-height: 64px;
    }

    .brand {
      gap: 0.55rem;
      font-size: 0.76rem;
    }

    .brand img {
      width: 32px;
      height: 32px;
    }

    .live-clock {
      gap: 0.45rem;
      font-size: 0.8rem;
    }

    main {
      padding-top: 1.5rem;
    }

    .hero {
      min-height: 180px;
      padding: 1.5rem;
      border-radius: 20px;
    }

    .hero::after {
      right: -120px;
      bottom: -210px;
    }

    .class-mark {
      width: 82px;
      border-radius: 22px;
    }

    .section-heading {
      align-items: center;
      margin-bottom: 1.1rem;
    }

    .timeline {
      gap: 0.75rem;
    }

    .timeline-item {
      grid-template-columns: 52px 14px minmax(0, 1fr);
      gap: 0.5rem;
    }

    .time-column {
      padding-top: 1.1rem;
      font-size: 0.77rem;
    }

    .end-time {
      font-size: 0.7rem;
    }

    .timeline-item:not(:last-child) .timeline-rail::after {
      top: 24px;
      bottom: -16px;
    }

    .rail-dot {
      width: 10px;
      height: 10px;
      margin-top: 1.3rem;
    }

    .lesson-heading {
      min-height: 40px;
      padding-inline: 0.75rem;
    }

    .lesson-detail {
      min-height: 0;
      flex-direction: column;
      align-items: stretch;
      gap: 0.65rem;
      padding: 0.8rem 0.75rem;
    }

    .subject h3 {
      font-size: 0.94rem;
      white-space: normal;
    }

    .lesson-meta {
      justify-content: flex-start;
      gap: 0.65rem;
      font-size: 0.74rem;
    }

    .page-footer {
      flex-direction: column;
      gap: 0.45rem;
      padding-top: 2rem;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .loader {
      animation-duration: 2s;
    }
  }
</style>
