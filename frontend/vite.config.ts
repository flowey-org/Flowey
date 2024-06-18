import fs from "fs";
import { fileURLToPath } from "node:url";

import { defaultAssetName } from "@vite-pwa/assets-generator/config";
import vue from "@vitejs/plugin-vue";
import { RollupOptions } from "rollup";
import { defineConfig } from "vite";
import { analyzer } from "vite-bundle-analyzer";
import { VitePWA, VitePWAOptions } from "vite-plugin-pwa";

await fs.promises.mkdir(`${import.meta.dirname}/public/assets/pwa`, { recursive: true });

const rollupOptions: RollupOptions = {
  output: {
    assetFileNames: (assetInfo) => {
      const name = assetInfo.name!;
      const extType = name.substring(name.lastIndexOf(".") + 1);
      return `assets/${extType}/[name].[hash:20][extname]`;
    },
    chunkFileNames: `assets/js/[name].[hash:20].js`,
    entryFileNames: `assets/js/[name].[hash:20].js`,
    hashCharacters: "hex",
    manualChunks: (id) => {
      if (id.includes("node_modules")) {
        return "vendor";
      }
      return null;
    },
  },
};

const vitePWAOptions: Partial<VitePWAOptions> = {
  filename: "service-worker.ts",
  injectRegister: false,
  injectManifest: {
    globPatterns: [
      "**/*.{js,css,html,woff2}",
      "favicon.{ico,svg}",
    ],
  },
  manifest: {
    name: "Flowey",
    short_name: "Flowey",
    lang: "en",
    display: "standalone",
    start_url: ".",
    theme_color: "#ffffff",
  },
  pwaAssets: {
    preset: {
      assetName: (type, size) => {
        return `assets/pwa/${defaultAssetName(type, size)}`;
      },
      transparent: {
        sizes: [64, 192, 512],
        favicons: [[48, "favicon.ico"]],
        padding: 0,
      },
      maskable: {
        sizes: [512],
        padding: 0,
      },
      apple: {
        sizes: [180],
        padding: 0,
      },
    },
  },
  srcDir: "src",
  strategies: "injectManifest",
};

export default defineConfig(({ mode }) => {
  const productionMode = mode === "production";
  const analyzeMode = mode === "analyze";
  return {
    build: {
      minify: productionMode && "esbuild",
      sourcemap: !productionMode,
      rollupOptions,
    },
    define: {
      DEV_MODE: !productionMode,
    },
    esbuild: {
      supported: {
        "top-level-await": true,
      },
    },
    plugins: [
      ...analyzeMode ? [analyzer()] : [],
      VitePWA(vitePWAOptions),
      vue(),
    ],
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),
      },
    },
  };
});
