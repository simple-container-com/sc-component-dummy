# sc-component-dummy

The template every SC **component** clones from, AND the throwaway smoke
subject of slice 1a of the [sc-components-platform](https://github.com/simple-container-com/forge/blob/main/docs/roadmap/items/sc-components-platform.md)
milestone.

An **SC component** is a new class of SC plugin whose resource provisions a
cloud resource AND a full forge-fleet-style microservice (go backend +
React SPA via `go:embed`) into the operator's chosen template (lambda /
cloud-compose / k8s). When a deploying service declares `uses:
[my-component]` in its `client.yaml`, the plugin's `ComputeProcess` sidecar
pushes a tenant registration to the deployed component service via a
JWT-authed HTTP call — tenants persist in the component service's own DB,
not in Pulumi state.

The design doc is at [`forge/docs/design/2026-11-10/sc-components-platform/`](https://github.com/simple-container-com/forge/tree/main/docs/design/2026-11-10/sc-components-platform/).
Read it before working on this repo. §P1 explains the Path A sidecar
wrapper (Schema-advertised component types + `.Raw`-reading, NOT Go
interface assertion — see `CLAUDE.md` and `TRAPS.md` for why). §P3 is the
push-register model this repo's plugin implements.

Five real components will clone this repo: **Backup → Monitoring →
Notifier → Credentials rotation → Incident Management**. `sc-component-dummy`
is the forcing function that proves the primitives work end-to-end, with
the smallest possible payload (a tenants list in Mongo, no real work).

## Three modules

```
sc-component-dummy/
├── api/         ← config tier. Published as its own Go module, referenced
│                  by the plugin binary (NOT by SC core — see CLAUDE.md on
│                  why Schema-advertised types + `.Raw` reader, not Go
│                  interface assertion).
├── plugin/      ← provisioning tier. The hashicorp/go-plugin subprocess
│                  binary SC core spawns via the `plugins:` block.
│                  Advertises `sc-dummy` under `component_resource_types`
│                  in its Schema. Provision is a no-op aside from
│                  republishing the service URL dispatched by Path A.
│                  ComputeProcess is where the push-register call lives.
└── service/     ← the microservice. go-aws-lambda-sdk + MongoDB + embedded
                   React SPA. Published to Dockerhub by the operator
                   (NOT by this repo's CI — see CLAUDE.md on the Dockerhub
                   publish gate).
```

## Running tests

```
make test       # fans out to all three modules
```

Both `api/` and `plugin/` tests must stay green on `main`. `service/`
tests cover the HTTP handlers + JWT middleware; they need Go 1.25+ and
no external services.

## Status

- Repo scaffolded 2026-11-11.
- Slice 1a (the SC core-side changes this plugin needs — Schema+Raw
  dispatch + `ComponentConfig` reader) is still owed from
  `simple-container-com/api`. Plugin binary builds + schema stubs are
  ready; the runtime handshake with SC core won't do anything useful
  until the api PR lands.
- Smoke rig TBD pending host-3 allocation from sc-ru-clouds's VPS fleet
  (coord via Smecsia). Plan: byo-VM with `ingress.enable: false` + high
  port HTTP; tar.gz plugin delivery; locally-built service image loaded
  on the host via `docker load`.
- Service image Dockerhub publish gated on operator sign-off (CLAUDE.md
  "package publishing, registries or release pipelines" ratchet).

## License

MIT. See [LICENSE](LICENSE).
