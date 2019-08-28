
/* http://big-elephants.com/2013-09/exploring-query-locks-in-postgres/ */
/* https://www.postgresql.org/docs/9.3/view-pg-locks.html */
SELECT locktype,
       virtualtransaction AS vxid,
       tuple              AS tup,
       relation::regclass,
       mode,
       virtualxid         AS target_vxid,
       transactionid      AS target_xid,
       l.granted
FROM pg_catalog.pg_locks l
WHERE NOT pid = pg_backend_pid() AND
      (NOT locktype = 'relation' OR granted = false)
ORDER BY vxid;

/* Translation table */
SELECT l.locktype,
       l.mode,
       l.pid,
       l.virtualtransaction AS my_vxid,
       l.virtualxid         AS target_vxid,
       l.transactionid      AS target_xid
FROM pg_locks l
WHERE
       l.mode = 'ExclusiveLock' AND
       l.locktype = 'transactionid';


/* give me all non-table locks and all not-granted locks */
/* http://big-elephants.com/2013-09/exploring-query-locks-in-postgres/ */
/* https://www.postgresql.org/docs/9.3/view-pg-locks.html */
SELECT l1.locktype,
       l1.pid,
       l1.virtualtransaction AS vxid,
       l2.transactionid      AS xid,
       l1.tuple              AS tup,
       l1.relation::regclass,
       l1.mode,
       l1.virtualxid         AS target_vxid,
       l1.transactionid      AS target_xid,
       l1.granted
FROM pg_catalog.pg_locks l1
JOIN pg_catalog.pg_locks l2
  ON l1.virtualtransaction = l2.virtualtransaction AND
     l2.locktype = 'transactionid' AND
     /* l1.granted = false AND */
     l2.mode = 'ExclusiveLock'
WHERE NOT l1.pid = pg_backend_pid() AND
      (NOT l1.locktype = 'relation' OR l1.granted = false)
ORDER BY vxid;


/* See what queries we're stuck on */
SELECT t.pid, l.transactionid, t.state, t.query, t.waiting
FROM pg_stat_activity t
JOIN pg_locks l ON (
  l.mode = 'ExclusiveLock' AND
  l.locktype = 'transactionid' AND
  l.pid = t.pid
)
WHERE NOT t.pid = pg_backend_pid() AND
      NOT t.query LIKE '%django_migrations%';


/******************************************************************************************/
/******************************************************************************************/
/******************************************************************************************/

/* SET application_name='%your_logical_name%'; */

/* https://wiki.postgresql.org/wiki/Lock_Monitoring#Logging_for_later_analysis */
SELECT blocked_locks.pid     AS blocked_pid,
         blocked_activity.usename  AS blocked_user,
         blocking_locks.pid     AS blocking_pid,
         blocking_activity.usename AS blocking_user,
         blocked_activity.query    AS blocked_statement,
         blocking_activity.query   AS current_statement_in_blocking_process,
         blocked_activity.application_name AS blocked_application,
         blocking_activity.application_name AS blocking_application
   FROM  pg_catalog.pg_locks         blocked_locks
    JOIN pg_catalog.pg_stat_activity blocked_activity  ON blocked_activity.pid = blocked_locks.pid
    JOIN pg_catalog.pg_locks         blocking_locks
        ON blocking_locks.locktype = blocked_locks.locktype
        AND blocking_locks.DATABASE IS NOT DISTINCT FROM blocked_locks.DATABASE
        AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
        AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
        AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
        AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
        AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
        AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
        AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
        AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
        AND blocking_locks.pid != blocked_locks.pid
    JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
   WHERE NOT blocked_locks.GRANTED;


SELECT bl.pid     AS blocked_pid,
     a.usename  AS blocked_user,
     a.query    AS blocked_statement
FROM  pg_catalog.pg_locks         bl
 JOIN pg_catalog.pg_stat_activity a  ON a.pid = bl.pid
WHERE NOT bl.granted;

/* https://www.endpoint.com/blog/2014/11/12/dear-postgresql-where-are-my-logs */
show log_destination; /* stderr */
show logging_collector; /* off */


/* example from http://elioxman.blogspot.com/2013/02/postgres-deadlock.html?m=1 */
CREATE TABLE parent (
  id integer PRIMARY KEY,
  name text
);

CREATE TABLE child (
  id integer PRIMARY KEY,
  parent_id integer REFERENCES parent(id),
  name text
);

BEGIN;
INSERT INTO child VALUES (3, 1, 'CHILD_A');
SELECT * from pg_sleep(10);
UPDATE parent SET name='Parent_B' WHERE id=1;
END;

BEGIN;
INSERT INTO child VALUES (4, 1, 'CHILD_B');
SELECT * from pg_sleep(10);
UPDATE parent SET name='Parent_C' WHERE id=1;
END;

BEGIN;
UPDATE parent SET name=MD5(random()::text) WHERE id=1;
INSERT INTO child VALUES (6, 1, 'CHILD_C');
SELECT * from pg_sleep(60 * 60 * 24);
END;

SELECT * FROM parent;

/* https://www.heatware.net/databases/how-view-see-table-row-locks-postgres/ */
/* actually gives me something */
SELECT t.relname, l.locktype, l.pid, l.mode, l.granted,
       l.page, l.tuple, l.virtualxid, l.transactionid
FROM pg_locks l
JOIN pg_stat_all_tables AS t
  ON l.relation=t.relid
WHERE t.relname NOT IN ('pg_class', 'pg_index', 'pg_namespace')
ORDER BY l.pid, l.relation ASC;


/* http://big-elephants.com/2013-09/exploring-query-locks-in-postgres/ */
SELECT blockeda.pid AS blocked_pid, blockeda.query as blocked_query,
      blockinga.pid AS blocking_pid, blockinga.query as blocking_query
FROM pg_catalog.pg_locks blockedl
JOIN pg_stat_activity blockeda ON blockedl.pid = blockeda.pid
JOIN pg_catalog.pg_locks blockingl ON
  (blockingl.transactionid=blockedl.transactionid
   AND blockedl.pid != blockingl.pid)
JOIN pg_stat_activity blockinga ON blockingl.pid = blockinga.pid
WHERE NOT blockedl.granted;


/* SET deadlock_timeout = 60000; /1* 1m. Timeout is in milliseconds. *1/ */
/* SHOW deadlock_timeout; */
