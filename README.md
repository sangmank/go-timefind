go-timefind
===========


go-timefind is a search library in time domain, and a library for
regular jobs. It makes regular scheduling, like any other alarm clocks
or cron in *nix, easier to write.

go-timefind only supports conditions that are mostly AND'ed, but you
can implement other logical operations using this module more easily.

Time Queries 
============

Month(), Day(), Hour(), Minute(), Weekday() are for integer numbers,
while Months(), Days(), Hours(), Minutes(), Weekdays() are for a
string.

Buiding time queries
--------------------
Firstly, you need to build a query in time domain:

    # daily schedule at both 7:10 PM and 7:20 PM (logically or'ed)
    q1, err = timefind.Hour(19)
    q1.Minute(10,20)
    # or, in one call like crontab rule
    q1, err = timefind.NewFromString("10,20 19 * * *")

    # Every weekday 7:30 PM
    q2, err = timefind.WeekDays("1-5")
    q2.Hour(19)
    q2.Minute(30)
    # or, in one call
    q2, err = timefind.NewFromString("30 19 * * 1-5")

    # 7PM on the first day of every month 
    q3, _ = timefind.Day(1)
    q3.Hour(19)
    q3.Minute(0)
    # or, in one call
    q3, err = timefind.NewFromString("0 19 1 * *")

    # And'ing two rules.
    q = timefind.And(q1, q2)

Setting limits
---------------
You can set minimum limit for the search:

    # time query after time t
    tq1.After(t)


Enumeration with time queries
=============================

Getting matching times
----------------------
The next entry from now:

    time = tq.Next(time.Now())

Or, a shorthand version:
    time = entry.NextNow()

The next five entries from now:
    times = tq.NextList(time.now(), 5)

Or, a shorthand version:
    times = tq.NextNowList(5)

Getting channels for operations
--------------------------------
This is for the fans of time.After() of Golang.

A channel that fires at the next matching time from now:

    <-tq.WaitNext()


Difference between crontab
==========================

crontab manual says that a command is executed when one of two day
fields (day of month and day of week) matches with the current
time. This means that cron ORs these two conditions.

In go-timefind, every added condition between different fields only
makes the preexisting query more specific, basically doing only ANDs
between conditions. We only support AND for simplicity and
usability. In go, you can build OR'ed condition with WaitNext() pretty
easily with additional goroutines.

Why isn't the name go-cron
==========================

There already is go-cron project, and I just wanted an easier name
than the shortened version of a greek word. Also, `find` command in
*nix is great, and I wish this library could be as effective as the
tool in its job.
