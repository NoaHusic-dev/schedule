import tailwindcss from "@tailwindcss/vite";
import adapter from "@sveltejs/adapter-static";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit({
      // SPA: every route falls back to the client shell, which the Go server serves.
      adapter: adapter({ fallback: "index.html" }),
      // The TV runs for weeks without anyone touching it; the root layout reloads it when a new
      // build is deployed.
      version: { pollInterval: 5 * 60 * 1000 },
      compilerOptions: {
        runes: ({ filename }) => (filename.split(/[/\\]/).includes("node_modules") ? undefined : true),
      },
    }),
  ],
  server: {
    // Development is same-origin like production: the browser talks only to Vite, which forwards
    // /api to the API. By default that is the public tsd-timetable instance; to use a local Go
    // backend instead, run with API_TARGET=http://localhost:3000.
    proxy: {
      "/api": {
        target: process.env.API_TARGET ?? "https://tsd-timetable.husak.me",
        changeOrigin: true,
      },
    },
  },
});
