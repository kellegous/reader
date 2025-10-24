import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  root: "ui",
  base: "/ui/",
  plugins: [react()],
  build: {
    outDir: "../internal/ui/assets",
    assetsDir: ".",
    emptyOutDir: true,
  },
});
