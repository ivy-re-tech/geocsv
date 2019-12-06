# GeoCSV
GeoCSV is an easier way to insert latitude longitude pairs into PostGIS from a CSV

A lot of times, data comes in a CSV that contains a column for latitude, and a column for longitude
this isn't ideal for copying data straight into a PostGIS table, since you'll have to define a lat/lng column,
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
geocsv --lat 1 --lng 2 csv.csv | psql -c "COPY table FROM STDIN CSV"
```

Would output:

| Address        | Coords                     |
|----------------|----------------------------|
| 123 Main St    | POINT( -122.1235  47.1234) |
| 124 Main St    | POINT( -122.1235  47.1235) |
| 125 Main St    | POINT( -122.1235  47.1236) |
