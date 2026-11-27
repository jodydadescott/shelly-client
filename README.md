# shelly-client

Declarative config for Shelly Gen2+ devices (Plus, Pro, Gen3, Gen4).

## Desired config

Each entry under `shelly.shellyConfigs` uses the device's own `Shelly.GetConfig` shape: component
keys (`sys`, `wifi`, `switch:0`, `cover:0`, ...) mapping to that component's config. `auth` takes
`{enable: true|false}` and uses `shelly.password` as the device password.

For a device, the desired config is `common`, then the entry named after its app (model, as in
`Shelly.GetDeviceInfo`), then the entry named after its device ID, each layered on the one before.
List only the keys you manage; keys the device has but the config doesn't name are left alone.

```yaml
hostnames: [shelly-porch.lan]
shelly:
  password: ...
  shellyConfigs:
    common:
      auth: {enable: true}
      cloud: {enable: false}
      sys: {sntp: {server: time.google.com}}
    S1G3:                      # app (example; see config get or GetDeviceInfo)
      switch:0: {in_mode: follow}
    shelly1g3-aabbccddeeff:    # device ID (example)
      sys: {device: {name: porch}}
      switch:0: {name: Porch light}
```

`shelly-client config get` prints a device's running config in this shape, as a starting point.

## Commands

- `config compare`: prints `~ path: running -> desired` for each difference and `? path` for keys
  the device never reports (write-only, such as `wifi.sta.pass`). Exits 0 if equal, 2 on drift, 1 on error.
- `config set`: sets only the components with drift (`wifi` last, then `auth`), then reboots if a
  change requires it. Write-only keys are not drift, so a changed password is only sent with `--force`,
  which sets every desired component.
