# Tankerkoenig to InfluxDB

## Build

```bash
make
```

## Configuration

Create a configuration file and save it to `config.json`.

```json
{
  "influxDB": {
    "serverUrl": "https://...:8086",
    "token": "...",
    "bucket": "...",
    "org": "...",
    "measurement": "..."
  },
  "stations": [
    {
      "id": "<uuid of the station>",
      "brand": "...",
      "city": "...",
      "street": "..."
    }
    // ...
  ]
}
```

## Start

```bash
find ../tankerkoenig-data/prices/2022/06/* -type f -exec bin/tankerkoenig-to-influxdb -c config.json {} +
```
