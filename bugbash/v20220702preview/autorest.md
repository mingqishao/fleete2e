### AutoRest Configuration

> see https://aka.ms/autorest

> NOTE: 07-02-preview has no changes, so we point at 06-02-preview to generate the code

```yaml
version: 3.9.1
go: true
track2: true
debug: true
input-file:
- ../specs/specification/containerservice/resource-manager/Microsoft.ContainerService/preview/2022-06-02-preview/fleets.json
- ./fleetMemberships.json
file-prefix: zz_generated_
azure-arm: true
module-version: 0.0.1
use-extension:
  "@autorest/modelerfour": "~4.24.0"
  "@autorest/go": "4.0.0-preview.43"
directive:
  - from: swagger-document
    where: $..properties[?(@.readOnly)]
    debug: true
    transform: |
      delete $["readOnly"]
```
