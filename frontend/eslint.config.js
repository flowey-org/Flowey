import eslint from "@eslint/js";
import pluginStylistic from "@stylistic/eslint-plugin";
import pluginImport from "eslint-plugin-import";
import pluginVue from "eslint-plugin-vue";
import pluginTypeScript from "typescript-eslint";
import parserVue from "vue-eslint-parser";

const sources = {
  files: ["**/*.js", "**/*.ts", "**/*.vue"],
  ignores: ["dist/*", "node_modules/*"],
};

export default [
  ...pluginTypeScript.config(
    {
      ...sources,
      extends: [
        eslint.configs.recommended,
        ...pluginTypeScript.configs.strictTypeChecked,
        ...pluginTypeScript.configs.stylisticTypeChecked,
      ],
      rules: {
        "sort-imports": ["error", { ignoreDeclarationSort: true }],
        "@typescript-eslint/no-non-null-assertion": "off",
        "@typescript-eslint/restrict-template-expressions": ["error", {
          allowBoolean: true,
        }],
      },
      languageOptions: {
        parser: parserVue,
        parserOptions: {
          parser: pluginTypeScript.parser,
          project: true,
          tsconfigRootDir: import.meta.dirname,
          extraFileExtensions: [".vue"],
        },
      },
    },
    {
      files: ["**/*.js"],
      ...pluginTypeScript.configs.disableTypeChecked,
    },
  ),
  {
    ...sources,
    plugins: {
      "@stylistic": pluginStylistic,
    },
    rules: {
      ...pluginStylistic.configs.customize({
        braceStyle: "1tbs",
        quotes: "double",
        semi: true,
      }).rules,
      "object-curly-newline": ["error", { multiline: true, consistent: true }],
      "object-property-newline": ["error", { allowAllPropertiesOnSameLine: true }],
    },
  },
  {
    ...sources,
    plugins: {
      import: pluginImport,
    },
    rules: {
      ...pluginImport.configs.recommended.rules,
      ...pluginImport.configs.typescript.rules,
      "import/order": [
        "error", {
          "alphabetize": {
            order: "asc",
            caseInsensitive: true,
          },
          "newlines-between": "always",
        },
      ],
      "import/newline-after-import": "error",
    },
    settings: {
      "import/resolver": {
        typescript: true,
        node: true,
      },
    },
  },
  ...pluginVue.configs["flat/recommended"],
];
