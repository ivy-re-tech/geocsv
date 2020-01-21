# GeoCSV
GeoCSV is an easier way to insert latitude longitude pairs into PostGIS from a CSV

## Background

A lot of times, data comes in a CSV that contains a column for latitude, and a column for longitude. This isn't ideal for copying data straight into a PostGIS table, since you'll have to define a lat/lng column,
then run a command like:

```sql
UDPATE table SET location='POINT(' || lat || ' ' || lng || ')'
```

Instead, it'd be nice if we could just run
```shell script
cat csv.csv | psql -c "COPY table FROM STDIN CSV"
```
Or, for you shell pipe purists:
```shell script
psql -c "COPY table FROM STDIN CSV" < csv.csv
```

And have everything get loaded into the Geography/Geometry column

GeoCSV is built for this exact use case.

Imagine a CSV that looks this:

| Address        | Latitude |  Longitude |
|----------------|----------|------------|
| 123 Main St    | 47.1234  | -122.1235  |
| 124 Main St    | 47.1235  | -122.1235  |
| 125 Main St    | 47.1236  | -122.1235  |

Running
```shell script
geocsv --lat 1 --lng 2 csv.csv
```

Would output:

| Address        | Coords                     |
|----------------|----------------------------|
| 123 Main St    | POINT( -122.1235  47.1234) |
| 124 Main St    | POINT( -122.1235  47.1235) |
| 125 Main St    | POINT( -122.1235  47.1236) |

## Using geocsv

geocsv is designed to be used like any standard unix tools, accepting input from `stdin` by default and outputting to `stdout` by default. This property makes it very easy to use this utility within an existing bash pipeline

Contrived Example:
```shell script
cat csv.csv | sed -e "s/\t/,/g" | geocsv --lat 1 --lng 2 | psql -c "COPY table FROM STDIN CSV"
```

## Examples

CSV with header, keep header in output:
```shell script
geocsv --lat 1 --lng 2 --has-header --keep-header in.csv out.csv
```
CSV with header, but don't keep header in output:
```shell script
geocsv --lat 1 --lng 2 --has-header in.csv out.csv
```

Output WKT in a specific column:
```shell script
geocsv --lat 1 --lng 2 --wkt 0
```

Keep Lat/Lng Columns:
```shell script
geocsv --lat 0 --lng 1 --keep-columns
```

Specify input delimiter:
```shell script
geocsv --lat 0 --lat 1 --d '\t'
```


