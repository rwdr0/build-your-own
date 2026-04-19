# SQLite clone built with Go's awesome concurrency model.

[![progress-banner](https://backend.codecrafters.io/progress/sqlite/82fa30cf-d84f-4044-9c89-329855cc0f5d)](https://app.codecrafters.io/courses/sqlite/overview)

## Features Implemented

- Print page size, number of tables and table names
- Count rows in a table
- Read data from a single column
- Read data from multiple columns
- Filter data with a WHERE clause
- Retrieve data using a full-table scan
- Scan database concurrently, resulting in high performance

## Architecture

![Architecture](./sqlite_arch.svg)

Each page on disk is deserialized and printed as soon as they are read, resulting in lower query time
