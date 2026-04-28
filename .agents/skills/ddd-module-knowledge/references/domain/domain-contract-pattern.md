#### Domain Contract pattern

Domain contracts are interfaces that define the contract between the domain layer and the infrastructure layer.

Example:

```go
package contract

import (
	"context"
)

type [contract-name]Repository interface {}
```

Implementation example: see -> `internal/user/domain/contract/user_repository.go`
