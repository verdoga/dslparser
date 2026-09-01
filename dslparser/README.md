# dslparser

`dslparser` — библиотека Go для разбора учебных сценариев DSL 1.1 в универсальное неизменяемое AST. Библиотека сохраняет исходный регистр имён и нормализованные фрагменты, вычисляет позиции в Unicode-кодовых точках, восстанавливается после локальных ошибок и не выполняет ввода-вывода, кроме чтения переданного `io.Reader`.

## Границы

Библиотека разбирает только версию 1.1, не рендерит сценарий, не исправляет вход и не предоставляет CLI. Фатальная ошибка означает, что `Document` создать нельзя. Локальные ошибки доступны как `Diagnostic` у документа, узла или элемента; при них возвращается частичное дерево и `error == nil`.

## Установка

Модуль использует Go 1.26 и только стандартную библиотеку:

```bash
go get dslparser
```

При разработке в этом репозитории импорт `dslparser` разрешается локальным `go.mod`.

## Минимальный пример

```go
package main

import (
    "fmt"
    "log"

    "dslparser"
)

func main() {
    doc, err := dslparser.ParseString("@dsl-version 1.1\n# Урок\n@question Ответьте")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(doc.Version(), len(doc.Roots()), len(doc.Diagnostics()))
}
```

## Публичный API

* `Parse`, `ParseString` и `ParseReader` используют общий конвейер и возвращают `*Document` либо типизированный `*FatalError`.
* `Document` предоставляет версию, диапазон директивы, корни и общий список диагностик.
* `Node`, `Element`, `BlockInfo`, `Position`, `Span` и `Diagnostic` предоставляют только методы чтения; методы, возвращающие срезы, создают защитные копии.
* `Walk` выполняет прямой обход. Возврат `false` посетителем пропускает потомков текущего узла, но не останавливает весь обход.

Точный список идентификаторов показывает `go doc dslparser`. Архитектура описана в [docs/architecture.md](docs/architecture.md), покрытие грамматики — в [docs/grammar-coverage.md](docs/grammar-coverage.md), диагностики — в [docs/diagnostics.md](docs/diagnostics.md), расширение — в [docs/extending.md](docs/extending.md).

## Проверка

Из корня модуля:

```bash
gofmt -w .
go test ./...
go test -race ./...
go vet ./...
go list ./...
go doc dslparser
```

Проверка отсутствия внешних модулей:

```bash
go list -m all
```

Единственной строкой вывода должна быть `dslparser`.
