import eslintJS from '@eslint/js'
import typescriptParser from '@typescript-eslint/parser'
import pluginEslintImport from 'eslint-plugin-import'
import pluginEslintPrettier from 'eslint-plugin-prettier/recommended'
import pluginEslintTypescript from 'typescript-eslint'

export default [
  {
    files: ['**/*.ts'],
  },
  eslintJS.configs.recommended,
  ...pluginEslintTypescript.configs.strict,
  pluginEslintImport.configs.typescript,
  pluginEslintPrettier,
  {
    plugins: {
      import: pluginEslintImport,
    },

    languageOptions: {
      parser: typescriptParser,
      ecmaVersion: 2022,
      sourceType: 'module',
      parserOptions: {
        project: './tsconfig.json',
      },
    },

    rules: {
      indent: ['error', 2, { SwitchCase: 1 }],
      quotes: ['error', 'single', { avoidEscape: true }],
      semi: [2, 'never'],
      curly: 'warn',
      eqeqeq: 'warn',
      'no-throw-literal': 'warn',
      'import/order': [
        'error',
        {
          'newlines-between': 'always',
          alphabetize: { order: 'asc' },
        },
      ],
      'sort-imports': [
        'error',
        {
          ignoreDeclarationSort: true,
          ignoreCase: true,
        },
      ],
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          ignoreRestSiblings: true,
          destructuredArrayIgnorePattern: '^_',
          varsIgnorePattern: '^_',
        },
      ],
    },
  },
]
