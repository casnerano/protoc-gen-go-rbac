# protoc-gen-go-rbac

Пакет предоставляет возможности разграничения доступа к grpc-методам на основе ролей пользователя прямо в proto-контрактах.

**В комплекте**:
- плагин для `protoc` и protobuf-расширения
- интерцептор для проверки прав доступа по ролям

## Основные возможности
 - Настройка правил доступа к сервису в целом
 - Настройка правил доступа к методам сервиса
 - Настройки правил сервиса наследуются для методов (метод может переопределить)

### Пример

#### Настройка правил в proto-контракте
```protobuf
import "github.com/casnerano/protoc-gen-go-rbac/proto/rbac.proto";

// Сервис по умолчанию закрыт для всех.
service ExampleService {
  // Открываем сервис для всех.
  option (protoc_gen_go_rbac.service_rules) = {
    access_level: ACCESS_LEVEL_PUBLIC,
  };

  // Метод наследует правила сервиса.
  rpc Stats(google.protobuf.Empty) returns (google.protobuf.Empty);

  // Метод переопределяет уровень доступа сервиса,
  // и устанавливает его открытым для определенных ролей.
  rpc Update(google.protobuf.Empty) returns (google.protobuf.Empty) {
    option (protoc_gen_go_rbac.method_rules) = {
      access_level: ACCESS_LEVEL_PRIVATE,
      allowed_roles: ["manager", "director"]
    };
  }
}
```

### Подключение

Скачать плагин
```bash
GOBIN=${PWD}/bin go install github.com/casnerano/protoc-gen-go-rbac/cmd/protoc-gen-go-rbac@latest
```

Добавить в `buf` новый плагин
```yaml
- name: rbac
  out: .
  opt:
  - paths=source_relative
  path: bin/protoc-gen-go-rbac
  strategy: directory
```

и добавить в опции плагина `go` следующую опцию
```yaml
- Mgithub.com/casnerano/protoc-gen-go-rbac/proto/rbac.proto=github.com/casnerano/protoc-gen-go-rbac/proto
```

например,
```yaml
- name: go
  out: .
  opt:
    - paths=source_relative
    - Mgithub.com/casnerano/protoc-gen-go-rbac/proto/rbac.proto=github.com/casnerano/protoc-gen-go-rbac/proto
  path: bin/protoc-gen-go
  strategy: directory
```

Подключить grpc-интерцептор
```go
// Todo
```
**Важно:** интерцептор работает по концепции zero-trust,
и после подключения ограничивает доступ ко всем сервисам, если они не переопределены другими правилами.

Выполнить `buf generate`.