# Drop Ins
Drop ins are noah's way of emulating an actual standalone postgres database. 
When the client sends a query for something like `current_database()` then noah will 
replace the value of the function with the value from noah's current state. Because
the database name will vary significantly depending on the data node that serves the
query; so by replacing or "dropping in" a value on the fly, so that the result will
be consistent regardless of the coordinator receiving the query or the data node 
processing the query.

For example, in the following query:
```postgresql
SELECT current_database();
```

The drop in system will edit the query before it is sent to the data node, resulting
in the following query:
```postgresql
SELECT 'noah';
```

## Table Data
The drop in system will eventually extend beyond just function calls.
When the client queries certain system tables noah might substitute the data from 
those tables. 


# Current Drop In Functions
These are all of the current functions that noah will drop in it's own value and
some details about each function.

> ### `current_database()` returns `text`
> Will return `noah` at the moment. In the future there are plans for adding support
for multiple databases in noah, each with separate settings and possibly even
separate postgresql configurations. But at the moment all of noah's queries will
behave as if they are in a single database.

> ### `current_schema()` returns `text`
> Will return `public` at the moment. There is no support for multiple schemas at
the time of writing this. 