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
