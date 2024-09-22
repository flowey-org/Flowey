import eslint from "@eslint/js";
import stylistic from "@stylistic/eslint-plugin";
import globals from "globals";

export default [
  {
    files: ["**/*.js"],
    languageOptions: {
      globals: {
        ...globals.browser,
      },
    },
  },
  eslint.configs.recommended,
  stylistic.configs.customize({
    braceStyle: "1tbs",
    quotes: "double",
    semi: true,
  }),
];
