## Fly

Create new flyway scripts under `src/main/migration` and open them in VsCode with one command.

### Installation

Install Go, then run:

```bash
go install github.com/setlog/fly
```

### Usage

```bash
fly # creates src/main/migration/VXXX.XXX__migration.sql
fly my_change # creates src/main/migration/VXXX.XXX__my_change.sql
```

The file is then opened in VsCode if you have it installed.
