#### Domain Contract pattern

Domain contracts are interfaces that define the contract between the domain layer and the infrastructure layer.
They are defined in `domain/contract/` and **the infrastructure layer is responsible for implementing them** — no other layer may provide implementations.

Contracts are grouped by type into files, the file name indicates the contract type. For example: Repository.go, datasource.go...

Use this example for:
- Repository contracts struct name -> sufix: Repository ex: userRepository
- Data source contracts struct name -> sufix: DataSource ex: userDataSource
- Façade contracts struct name -> sufix: Facade ex: BillingFacade
- Gateway contracts struct name -> sufix: Gateway ex: RestGateway

Example:

```go
package contract

import (
	"context"
)

type [contract-name]Repository interface {}
```

Implementation example: see ->  [reference](internal/user/domain/contract/repository.go)
