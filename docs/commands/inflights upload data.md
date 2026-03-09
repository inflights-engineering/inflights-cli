# inflights upload data

Upload processed data (deliverables) for a given flight. Can upload a single deliverable or an entire directory.

## Usage

```bash
inflights upload data <flightId> <path> [options]
```

## Options

| Flag             | Description                                            |
| ---------------- | ------------------------------------------------------ |
| `--type <type>`  | Deliverable type (required when uploading a single file) |

### Available deliverable types

| Type                  | Description                                             |
| --------------------- | ------------------------------------------------------- |
| `point_cloud`         | Point Cloud from Photogrammetry - 3D                    |
| `point_cloud_lidar`   | Point Cloud from LIDAR - 3D                             |
| `textured_mesh`       | Textured Mesh - 3D                                      |
| `orthomosaic`         | Orthomosaic - 2D                                        |
| `dsm`                 | Digital Surface Model - 2.5D                            |
| `cad_roof`            | CAD Model Roof - 3D                                     |
| `cad_terrain`         | CAD Model Terrain - 3D                                  |
| `cad_power_pole`      | CAD Model - Power Pole                                  |
| `surface_xml`         | Surface data XML - 3D                                   |
| `surface_points`      | Surface Points (Landxml/CSV)                            |
| `survey_report`       | Survey Report                                           |
| `stockpile_report`    | Stockpile Volumetric report                             |
| `thermal_report`      | Thermal Report                                          |
| `solar_report`        | Solar Site Infrared Inspection report                   |
| `ndvi`                | NDVI maps (index maps)                                  |
| `collision_analysis`  | Collision analysis between power grid and vegetation    |
| `multispectral`       | Multispectral Images                                    |
| `aerial_media`        | Aerial pictures and movies                              |

## Examples

Upload a single deliverable:

```bash
inflights upload data FL-1042 ./orthomosaic.tif --type orthomosaic
# → Uploading orthomosaic (450 MB)…
# → ████████████████████ 100%
# → FL-1042 orthomosaic uploaded.
```

Upload all processed data at once:

```bash
inflights upload data FL-1042 ./processed-output/
# → Uploading processed data (1.8 GB)…
# → ████████████████████ 100%
# → FL-1042 processed data uploaded. Status → completed.
```

## API

```
POST /flights/:flightId/data
Content-Type: multipart/form-data
```

## Roles

Processor.
