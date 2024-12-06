# pg_explain

A terminal tool for exploring postgres query plans.

## USAGE

Read in a json formatted query plan from STDIN

```
> cat explain_plan.json | pg_explain
```

Execute a sql file passed in as an argument

```
> pg_explain exec my_query.sql
```

Show a previously executed plan

```
> pg_explain
```

## Storing plans

Each time a query is executed `pg_explain` stores the resulting query plan
along with the relevant settings at time of execution into a `.pgex` file
stored in a `_pgex` directory in the location where pg_explain was executed.

## Examples

Navigate nodes with j and k

![CleanShot 2024-12-06 at 11 28 18](https://github.com/user-attachments/assets/46dda840-7246-42c4-88ee-250a7c98f1a0)

Navigate node stats with [ and ]

* Rows
* Buffers
* Cost 
* Time

![CleanShot 2024-12-06 at 11 33 57](https://github.com/user-attachments/assets/a5826afc-d355-48f3-8f93-685906a0226b)
