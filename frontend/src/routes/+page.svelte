<script lang="ts">
  import { onMount } from "svelte";
  import logo from "#lib/assets/tsd-logo-dark.png";
  import { formatClock, formatDate, isoDate } from "#lib/helpers/format.ts";
  import { formatPeriods, periodIds } from "#lib/helpers/timetable.ts";
  import { clock } from "#lib/stores/clock.svelte.ts";
  import { schedules } from "#lib/stores/schedules.svelte.ts";

  const APP_VERSION = "1.0.0";
  let selectedClass = $state("2 RT");
  const classSchedule = $derived(schedules.get(selectedClass));
  let showSettings = $state(false);
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
  let theme = $state<"dark" | "light">("dark");
  let selectedGroup = $state<"A" | "B">("A");
  const lastRefreshedText = $derived(
    schedules.lastRefreshed
      ? new Date(schedules.lastRefreshed).toLocaleTimeString("hr-HR", { hour: "2-digit", minute: "2-digit" })
      : "",
  );

  function applyTheme(nextTheme: "dark" | "light") {
    theme = nextTheme;
    document.documentElement.dataset.theme = nextTheme;
  }

  function toggleTheme() {
    const nextTheme = theme === "dark" ? "light" : "dark";
    applyTheme(nextTheme);
    localStorage.setItem("schedule-theme", nextTheme);
  }

  function selectGroup(group: "A" | "B") {
    selectedGroup = group;
    localStorage.setItem("schedule-group", group);
  }

  function selectClass(className: string) {
    selectedClass = className;
    localStorage.setItem("schedule-class", className);
    closeSettings();
  }

  function closeSettings() {
    showSettings = false;
  }

  onMount(() => {
    const savedTheme = localStorage.getItem("schedule-theme");
    applyTheme(savedTheme === "light" ? "light" : "dark");
    const savedGroup = localStorage.getItem("schedule-group");
    selectedGroup = savedGroup === "B" ? "B" : "A";
    const savedClass = localStorage.getItem("schedule-class");
    if (savedClass) selectedClass = savedClass;

  });

  $effect(() => clock.subscribe());
  $effect(() => schedules.subscribe());
</script>

<svelte:head>
  <title>{selectedClass} · Raspored sati</title>
  <meta
    name="description"
    content={`Dnevni raspored sati za razred ${selectedClass} Tehničke škole Daruvar.`}
  />
</svelte:head>

<div class="page-shell">
  <header class="topbar">
    <a class="brand" href="/" aria-label={`Raspored sati za ${selectedClass}`}>
      <img src={logo} alt="" />
      <span>Tehnička škola Daruvar</span>
    </a>
    <div class="live-clock" aria-label="Trenutno vrijeme">
      <time>{formatClock(clock.now)}</time>
      <button
        class="settings-toggle"
        type="button"
        aria-label="Postavke"
        title="Postavke"
        onclick={() => (showSettings = true)}
      >⚙</button>
      <button
        class="theme-toggle"
        type="button"
        aria-label={theme === "dark" ? "Promijeni na svijetlu temu" : "Promijeni na tamnu temu"}
        title={theme === "dark" ? "Svijetla tema" : "Tamna tema"}
        onclick={toggleTheme}
      >
        {theme === "dark" ? "☀" : "☾"}
      </button>
    </div>
  </header>

  <main>
    <section class="hero" aria-labelledby="page-title">
      <div class="hero-copy">
        <p class="eyebrow">Raspored nastave</p>
        <h1 id="page-title">{selectedClass.split(" ")[0]} <span>{selectedClass.split(" ").slice(1).join(" ")}</span></h1>
        <p class="date">{formatDate(clock.now)}</p>
      </div>
      <div class="class-mark" aria-hidden="true">{selectedClass.split(" ").slice(1).join("") || selectedClass}</div>
    </section>

    <section class="schedule-section" aria-labelledby="schedule-title">
      <div class="section-heading">
        <div>
          <p class="eyebrow">Dnevni pregled</p>
          <h2 id="schedule-title">Raspored za danas</h2>
        </div>
        <div class="heading-actions">
          <div class="group-picker" aria-label="Istaknuta grupa">
            <button
              type="button"
              class:selected={selectedGroup === "A"}
              aria-pressed={selectedGroup === "A"}
              onclick={() => selectGroup("A")}
            >A grupa</button>
            <button
              type="button"
              class:selected={selectedGroup === "B"}
              aria-pressed={selectedGroup === "B"}
              onclick={() => selectGroup("B")}
            >B grupa</button>
          </div>
          {#if schedules.loaded && isToday && isFinished}
            <span class="slot-count finished">Završena nastava</span>
          {/if}
          <button
            class="refresh-button"
            type="button"
            aria-label="Osvježi raspored"
            title="Osvježi raspored"
            disabled={schedules.refreshing}
            onclick={() => schedules.refresh()}
          >
            ↻
          </button>
        </div>
      </div>

      {#if schedules.loaded}
        <div class="refresh-status" role="status">
          {#if schedules.offline}
            <span>Offline · prikazuju se zadnji spremljeni podaci</span>
          {:else if lastRefreshedText}
            <span>Zadnje osvježeno u {lastRefreshedText}</span>
          {/if}
        </div>
      {/if}

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
          <p>Nema upisanih sati za {selectedClass}.</p>
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
                    <div
                      class="lesson-detail"
                      class:group-highlighted={lesson.groups?.some((group) => group.trim().toUpperCase() === selectedGroup)}
                    >
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
        Projekt napravio <a href="https://noahusicdev.site" target="_blank" rel="noreferrer">Noa Husić</a>
        <span class="footer-divider">·</span>
        <a href="https://github.com/noahusic-dev" target="_blank" rel="noreferrer">GitHub</a>
      </span>
      <span>
        Napravljeno pomoću projekta
        <a href="https://github.com/h00s" target="_blank" rel="noreferrer">h00s</a>
        <span class="footer-divider">·</span>
        <a href="https://kruno.husak.me" target="_blank" rel="noreferrer">Krunoslav Husak</a>
        <span class="footer-divider">·</span>
        Verzija {APP_VERSION}
      </span>
    </footer>
  </main>
</div>

{#if showSettings}
  <div
    class="settings-backdrop"
    role="presentation"
    onclick={(event) => {
      if (event.target === event.currentTarget) closeSettings();
    }}
  >
    <div
      class="settings-panel"
      role="dialog"
      aria-modal="true"
      aria-labelledby="settings-title"
    >
      <div class="settings-heading">
        <div>
          <p class="eyebrow">Prikaz rasporeda</p>
          <h2 id="settings-title">Postavke</h2>
        </div>
        <button class="settings-close" type="button" aria-label="Zatvori postavke" onclick={closeSettings}>×</button>
      </div>

      <label for="class-select">Razred</label>
      <select
        id="class-select"
        bind:value={selectedClass}
        onchange={(event) => selectClass(event.currentTarget.value)}
      >
        {#if schedules.classes.length === 0}
          <option value={selectedClass}>{selectedClass}</option>
        {:else}
          {#each schedules.classes as classScheduleOption (classScheduleOption.class)}
            <option value={classScheduleOption.class}>{classScheduleOption.class}</option>
          {/each}
        {/if}
      </select>

      <div class="settings-theme">
        <span>Tema</span>
        <button type="button" onclick={toggleTheme}>
          {theme === "dark" ? "Svijetla tema" : "Tamna tema"}
        </button>
      </div>

      <section class="usage-guide" aria-labelledby="usage-title">
        <h3 id="usage-title">Kako koristiti</h3>
        <ol>
          <li>Odaberi svoj razred u postavkama.</li>
          <li>Izaberi A ili B grupu za isticanje predmeta.</li>
          <li>Svijetlo plavi sat označuje nastavu koja je trenutno u tijeku.</li>
        </ol>
      </section>
    </div>
  </div>
{/if}

<style>
  .page-shell {
    display: flex;
    min-height: 100dvh;
    flex-direction: column;
    overflow-x: hidden;
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
    padding-top: env(safe-area-inset-top);
    padding-inline: env(safe-area-inset-left) env(safe-area-inset-right);
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

  .theme-toggle {
    display: grid;
    width: 44px;
    height: 44px;
    margin-left: 0.15rem;
    place-items: center;
    border: 1px solid #dfe5ef;
    border-radius: 50%;
    background: white;
    color: #536178;
    font-size: 1rem;
    cursor: pointer;
  }

  .theme-toggle:focus-visible {
    outline: 3px solid rgb(50 121 189 / 25%);
    outline-offset: 2px;
  }

  .settings-toggle {
    display: grid;
    width: 44px;
    height: 44px;
    margin-left: 0.15rem;
    place-items: center;
    border: 1px solid #dfe5ef;
    border-radius: 50%;
    background: white;
    color: #536178;
    cursor: pointer;
  }

  main {
    display: flex;
    flex: 1;
    flex-direction: column;
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

  .heading-actions {
    display: flex;
    align-items: center;
    gap: 0.65rem;
  }

  .group-picker {
    display: inline-flex;
    padding: 3px;
    border: 1px solid #dde4ef;
    border-radius: 999px;
    background: white;
  }

  .group-picker button {
    min-width: 70px;
    min-height: 44px;
    padding: 0.38rem 0.65rem;
    border: 0;
    border-radius: 999px;
    background: transparent;
    color: #68768b;
    font-size: 0.72rem;
    font-weight: 750;
    cursor: pointer;
  }

  .group-picker button.selected {
    background: #3279bd;
    color: white;
  }

  .group-picker button:focus-visible {
    outline: 3px solid rgb(50 121 189 / 25%);
    outline-offset: 2px;
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

  .refresh-button {
    display: grid;
    width: 44px;
    height: 44px;
    flex: 0 0 44px;
    place-items: center;
    border: 1px solid #dde4ef;
    border-radius: 50%;
    background: white;
    color: #536178;
    font-size: 1.15rem;
    cursor: pointer;
  }

  .refresh-button:disabled {
    opacity: 0.55;
    cursor: wait;
  }

  .refresh-status {
    min-height: 1.25rem;
    margin: -0.7rem 0 1rem;
    color: #7d899d;
    font-size: 0.72rem;
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
    width: 14px;
    height: 14px;
    border: 3px solid #bde2ff;
    background: #77bdf8;
    box-shadow: 0 0 0 5px rgb(119 189 248 / 18%);
    animation: current-pulse 1.8s ease-in-out infinite;
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
    border-color: #77bdf8;
    background: linear-gradient(135deg, #f7fbff 0%, white 55%);
    box-shadow: 0 8px 28px rgb(119 189 248 / 22%), 0 0 0 1px rgb(119 189 248 / 18%);
  }

  @keyframes current-pulse {
    0%, 100% {
      box-shadow: 0 0 0 5px rgb(119 189 248 / 18%);
    }
    50% {
      box-shadow: 0 0 0 9px rgb(119 189 248 / 8%);
    }
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

  .lesson-detail.group-highlighted {
    border-left: 4px solid #3279bd;
    background: linear-gradient(90deg, rgb(50 121 189 / 8%), white 32%);
  }

  .lesson-detail.group-highlighted .group-name {
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
    margin-top: auto;
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

  .settings-backdrop {
    position: fixed;
    z-index: 80;
    inset: 0;
    display: grid;
    place-items: center;
    padding: 1rem;
    background: rgb(10 20 35 / 42%);
    backdrop-filter: blur(6px);
  }

  .settings-panel {
    width: min(390px, 100%);
    padding: 1.25rem;
    border: 1px solid #dde4ef;
    border-radius: 18px;
    background: white;
    box-shadow: 0 20px 60px rgb(20 38 66 / 24%);
  }

  .settings-heading {
    display: flex;
    align-items: start;
    justify-content: space-between;
    margin-bottom: 1.25rem;
  }

  .settings-heading h2 {
    font-size: 1.35rem;
  }

  .settings-close {
    display: grid;
    width: 32px;
    height: 32px;
    place-items: center;
    border: 1px solid #dde4ef;
    border-radius: 50%;
    background: #f7f9fc;
    color: #69778d;
    font-size: 1.1rem;
    cursor: pointer;
  }

  .settings-panel label,
  .settings-theme span {
    display: block;
    margin-bottom: 0.45rem;
    color: #68768b;
    font-size: 0.78rem;
    font-weight: 750;
  }

  .settings-panel select {
    width: 100%;
    height: 46px;
    margin-bottom: 1.1rem;
    padding: 0 0.75rem;
    border: 1px solid #dfe5ef;
    border-radius: 10px;
    background: white;
    color: #26344b;
  }

  .settings-theme {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .usage-guide {
    margin-top: 1.25rem;
    padding-top: 1rem;
    border-top: 1px solid #e8edf5;
  }

  .usage-guide h3 {
    margin: 0 0 0.55rem;
    font-size: 0.95rem;
  }

  .usage-guide ol {
    display: grid;
    margin: 0;
    padding-left: 1.15rem;
    gap: 0.35rem;
    color: #68768b;
    font-size: 0.78rem;
  }

  .settings-theme span {
    margin: 0;
  }

  .settings-theme button {
    padding: 0.48rem 0.7rem;
    border: 1px solid #dfe5ef;
    border-radius: 999px;
    background: #f7f9fc;
    color: #536178;
    font-size: 0.75rem;
    font-weight: 700;
    cursor: pointer;
  }

  :global(html[data-theme="dark"]) .page-shell {
    background:
      radial-gradient(ellipse at 82% 4%, rgb(91 153 255 / 10%), transparent 29rem),
      #0d1626;
  }

  :global(html[data-theme="dark"]) .topbar {
    border-bottom-color: #26344a;
  }

  :global(html[data-theme="dark"]) .brand,
  :global(html[data-theme="dark"]) .live-clock {
    color: #d6dfed;
  }

  :global(html[data-theme="dark"]) .theme-toggle {
    border-color: #35465f;
    background: #17243a;
    color: #d6dfed;
  }

  :global(html[data-theme="dark"]) .settings-toggle {
    border-color: #35465f;
    background: #17243a;
    color: #d6dfed;
  }

  :global(html[data-theme="dark"]) .refresh-button {
    border-color: #35465f;
    background: #17243a;
    color: #d6dfed;
  }

  :global(html[data-theme="dark"]) .refresh-status {
    color: #aab9cc;
  }

  :global(html[data-theme="dark"]) .settings-panel {
    border-color: #304159;
    background: #111d2f;
  }

  :global(html[data-theme="dark"]) .settings-heading h2,
  :global(html[data-theme="dark"]) .settings-panel select {
    color: #edf3fb;
  }

  :global(html[data-theme="dark"]) .settings-close,
  :global(html[data-theme="dark"]) .settings-theme button {
    border-color: #304159;
    background: #17243a;
    color: #bac7d8;
  }

  :global(html[data-theme="dark"]) .settings-panel select {
    border-color: #304159;
    background: #17243a;
  }

  :global(html[data-theme="dark"]) .schedule-section h2,
  :global(html[data-theme="dark"]) .subject h3,
  :global(html[data-theme="dark"]) .time-column {
    color: #edf3fb;
  }

  :global(html[data-theme="dark"]) .eyebrow,
  :global(html[data-theme="dark"]) .end-time,
  :global(html[data-theme="dark"]) .lesson-meta,
  :global(html[data-theme="dark"]) .page-footer,
  :global(html[data-theme="dark"]) .message-detail {
    color: #9aabc2;
  }

  :global(html[data-theme="dark"]) .slot-count,
  :global(html[data-theme="dark"]) .lesson-card,
  :global(html[data-theme="dark"]) .message-card {
    border-color: #26344a;
    background: #111d2f;
  }

  :global(html[data-theme="dark"]) .slot-count {
    color: #bcc8d8;
  }

  :global(html[data-theme="dark"]) .group-picker {
    border-color: #304159;
    background: #17243a;
  }

  :global(html[data-theme="dark"]) .group-picker button {
    color: #aab9cc;
  }

  :global(html[data-theme="dark"]) .lesson-heading,
  :global(html[data-theme="dark"]) .lesson-detail + .lesson-detail {
    border-color: #26344a;
  }

  :global(html[data-theme="dark"]) .lesson-detail.group-highlighted {
    background: linear-gradient(90deg, rgb(50 121 189 / 18%), #111d2f 32%);
  }

  :global(html[data-theme="dark"]) .room {
    background: #1b2a40;
    color: #cbd6e5;
  }

  :global(html[data-theme="dark"]) .timeline-item:not(:last-child) .timeline-rail::after {
    background: #2b3b53;
  }

  :global(html[data-theme="dark"]) .rail-dot {
    border-color: #0d1626;
  }

  :global(html[data-theme="dark"]) .page-footer a {
    color: #aab9cc;
  }

  :global(html[data-theme="dark"]) .usage-guide {
    border-top-color: #304159;
  }

  :global(html[data-theme="dark"]) .usage-guide h3,
  :global(html[data-theme="dark"]) .usage-guide ol {
    color: #edf3fb;
  }

  :global(html[data-theme="dark"]) .usage-guide ol {
    color: #bac7d8;
  }

  :global(html[data-theme="dark"]) .active .lesson-card {
    border-color: #77bdf8;
    background: linear-gradient(135deg, rgb(119 189 248 / 14%), #111d2f 55%);
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
      gap: 0.5rem;
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

    .theme-toggle {
      width: 44px;
      height: 44px;
      font-size: 0.88rem;
    }

    .settings-toggle {
      width: 44px;
      height: 44px;
      font-size: 0.88rem;
    }

    main {
      padding-top: 1.5rem;
    }

    .hero {
      min-height: 150px;
      padding: 1.25rem;
      border-radius: 20px;
    }

    h1 {
      font-size: 2.8rem;
    }

    .date {
      margin-top: 0.65rem;
      font-size: 0.86rem;
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

    .heading-actions {
      flex-direction: column;
      align-items: flex-end;
      gap: 0.35rem;
    }

    .group-picker button {
      min-width: 58px;
      padding-inline: 0.45rem;
      font-size: 0.66rem;
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
      width: 100%;
      max-width: 100%;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 0.35rem;
      margin-top: auto;
      padding:
        1.75rem env(safe-area-inset-right, 0)
        calc(1.25rem + env(safe-area-inset-bottom, 0))
        env(safe-area-inset-left, 0);
      font-size: clamp(0.68rem, 2.8vw, 0.75rem);
      text-align: center;
      overflow-wrap: anywhere;
      white-space: normal;
    }

  }

  @media (prefers-reduced-motion: reduce) {
    .loader {
      animation-duration: 2s;
    }

    .active .rail-dot {
      animation: none;
    }
  }
</style>
