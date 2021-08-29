# GORM SQLChaos

GORM SQLChaos manipulates DML at program runtime based on [gorm callbacks](https://gorm.io/docs/write_plugins.html)

## Motivation

In Financial Business distributed system, account imbalance problems caused by unstable networks
or human mistakes may cause serious impacts. We built Imbalance Monitor&Analysis System,
so we want to create data imbalance situations between our business systems to verify
if our monitor reports these imbalances timely.
Also, we want this situation to be controllable and runs periodically to ensure the system works fine.

Yep, [Chaos Engineering](https://principlesofchaos.org/) ;).

So I developed SQLChaos and embedded it into our business systems.

NOTE: if you're looking for SQL injection attack or any related tools, SQLChaos is **not** what you want.

## Features

* Easy to embed into your code;
* Modify DML SQL values at program runtime;
* Support `INSERT`, `UPDATE` SQL;

## How it works

SQLChaos registers hooks on gorm `Before("gorm:update")` and `Before("gorm:create")` callbacks.
It will fetch values from *Statement.Dest pointer, which is a staging store before the real operation is performed,
and try to match user defined conditions and apply assigments.

## Using SQLChaos

### Setup

Embed SQLChaos where your gorm.DB setupped.

```go
db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{}, &sqlchaos.Config{
    DBName:     "dummy",
    RuleReader: sqlchaos.WithSimpleHTTPRuleReader(),
})
if err != nil {
    fmt.Fprintf(os.Stderr, "connect db failed:%v", err)
    return
}
```

SQLChaos provides a simple HTTP server for implementing `RuleReader` to enable or disable
chaos rules at program runtime. You can implement yours and replace it as your need.

SimpleHTTPRuleReader listens at `SQLCHAOS_HTTP`.

```bash
> SQLCHAOS_HTTP=127.0.0.1:8081 ./your-program
```

After your program started, you will see `SQLChaos enabled`. And SQLChaos listens at port `8081`.

### Enable Chaos Rule

Try to enable a rule,
```bash
export DBNAME=dummy;export TABLE=users;
# For dabase dummy table users, before INSERT Statement evaluated,
# set balance to 1024 and age to 40 for every record
# which age is greater then or equals to 1 and less then 50.
> curl -XPOST "http://127.0.0.1:8081/$DBNAME/$TABLE" -H'Content-Type: application/json' \
    -d'{"dml":"INSERT","when":"age>=1 AND age<50","then":"balance=1024,age=40"}'
```

After enabled, every record you created using gorm which matches `age>=1` and `age<50` condition, 
the `balance` will be set to `1024` and `age` will be `40` before it inserted into database.

To disable the table chaos rule,
```bash
> curl -XDELETE "http://127.0.0.1:8081/$DBNAME/$TABLE" 
```

For more practical examples, please check `./example`.

## Limits

* gorm `Save` may not envoke hooks, so SQLChaos only be called in `Create`, `Update`, `Updates`;
* `when` condition supports `AND` only, and operators like `= < <= > >=` are supported;
* Only read values from gorm Statement.Dest which is basically same as
    the argument you pass to `Create` `Update` `Updates` functions.
    **Ensure values which `when` needed are present where you call on `Create`, `Update`, `Updates`.**
    **Record will not be matched if values `when` required are absent.**
