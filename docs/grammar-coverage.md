# Покрытие грамматики DSL 1.1

Таблица отражает фактический реестр `internal/registry/v1_1.go`. «Свободный текст» означает один именованный элемент без токенизации; «токены» поддерживают двойные кавычки и экранирование. Неизвестный inline-тег сохраняется с универсальными токенами, неизвестный block-тег — с непрозрачным телом.

| Тег | Форма | Аргументы | Тело | Обработчик / контекст |
|---|---|---|---|---|
| `dsl-version` | inline | version, токены | нет | version |
| `task` | inline | id, токены | логическая область | task start |
| `endtask` | inline | нет | нет | task end |
| `header` | inline | title, свободный текст | нет | task boundary |
| `newpage` | inline | нет | нет | step boundary |
| `step` | inline | title, свободный текст | логическая область | step start |
| `editor` | inline/block | content, свободный текст | opaque raw | raw body |
| `speaking` | inline | нет | нет | speaking context |
| `media` | inline | media tokens | нет | media |
| `example` | inline/block | content, свободный текст | opaque shaped | example |
| `wordlist` | inline/block | content, свободный текст | opaque shaped | wordlist |
| `table` | block | optional title | opaque raw | raw body |
| `script` | block | optional title | opaque raw | raw body |
| `text` | block | optional title | opaque raw, пустые строки сохраняются | raw body |
| `key` | block | optional title | opaque raw | raw body |
| `instr` | inline/block | content, свободный текст | opaque raw | raw body |
| `note` | inline/block | content, свободный текст | opaque raw | raw body |
| `alt` | inline/block | content, свободный текст | opaque raw | raw body |
| `question` | inline | content, свободный текст | нет | free text |
| `multifill` | inline/block | optional instruction | opaque raw, пустые строки сохраняются | raw body |
| `choice` | block | optional instruction | opaque shaped | items |
| `multichoice` | block | optional instruction | opaque shaped | items |
| `matching` | block | optional instruction | opaque shaped | matching |
| `ordering` | block | optional instruction | opaque shaped | ordering |
| `variants` | block | нет | structural | variants context |
| `variant` | inline | name, свободный текст | logical | variant context |

Помимо тегов реализованы заголовки произвольного положительного уровня, текстовые строки, пятисимвольный placeholder `_____`, фигурные границы, логические задания/шаги и прямой обход AST. `table`, `script`, `text`, `key` и остальные raw-тела намеренно не разбирают внутреннюю DSL-разметку.
