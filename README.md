![Go 1.16](https://img.shields.io/badge/Go-v1.16-blue)
[![codecov](https://codecov.io/gh/toughnoah/blackbean/branch/master/graph/badge.svg?token=4UUTYZ6NCF)](https://codecov.io/gh/toughnoah/blackbean)
[![CI Workflow](https://github.com/toughnoah/blackbean/actions/workflows/test-coverage.yaml/badge.svg)](https://github.com/toughnoah/blackbean/actions/workflows/test-coverage.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/toughnoah/blackbean)](https://goreportcard.com/report/github.com/toughnoah/blackbean)
# blackbean
The blackbean is a command tool for elasticsearch operations by using cobra. Besides, blackbean is the name of my lovely French bulldog.
![image](https://github.com/toughnoah/blackbean/blob/45d1aca63307d79b9b3a56028494219bf5ba25b2/img/blackbean.png)
## Configuration
Define your config file firstly, and specify `current` as one of your cluster.
```
cluster:
  default:
    url: https://a.es.com:9200
    username: Noah
    password: abc
  backup:
    url: https://b.es.com:9200
    username: Noah
    password: abc
current: default
```

## Shell completion
All commands have fulfilled necessary completion, including flags. Enjoy yourself with blackbean!
```console
[root@noah ~]# echo "source <(blackbean completion bash)" >> ~/.bashrc
```
```console
[root@noah ~]# blackbean [tab][tab]
alias       cat         current-es  help        repo        role        use         watcher
apply       completion  explain     index       reroute     snapshot    user

```
```console
[root@docker ~]# blackbean cat [tab][tab]
allocations   cachemem      health        largeindices  nodes         segmem        threadpoo
```

## Common Flag
You can use `-d` or `--data` to specify raw request body like `'{"query":"match_all":{}"}'`.
Also, you can directly read from file using `-f` or `--filename`. Both json and yaml are supported.
```console
[root@noah ~]# blackbean index search test-* --
--config     --config=    --data       --data=      --filename   --filename=
```
```console
[root@noah ~]# cat query.json
{"size":1,
  "query":{
"match_all": {}}}

```
## Command
```console
[root@noah ~]# blackbean
blackbean command provides a set of commands to talk with es via cli.
Besides, blackbean is the name of my favorite french bulldog.

Usage:
  blackbean [command]

Available Commands:
  alias       alias index
  apply       apply cluster changes
  cat         cat allocation/nodes/health/nodes/threadpool/cache memory/segments memory/large indices.
  completion  Generate completion script
  current-es  show current cluster context
  explain     explain index allocation
  help        Help about any command
  index       index operations
  repo        repo operations
  reroute     reroute for cluster
  role        role operations for cluster
  snapshot    snapshot operations
  use         change current cluster context
  user        user for cluster
  watcher     operate watcher

Flags:
      --config string   config file (default is $HOME/.blackbean.yaml)
  -h, --help            help for blackbean
  -t, --toggle          Help message for toggle

Use "blackbean [command] --help" for more information about a command.
```
### Use
```console
[root@noah ~]# blackbean current-es
current using cluster: qa
```
```console
[root@noah ~]# blackbean use [tab][tab]
prod  qa
[root@noah ~]# blackbean use qa
change to cluster: qa

```
### Cat
```console
[root@noah ~]# blackbean cat health 
[200 OK] epoch      timestamp cluster       status node.total node.data shards  pri relo init unassign pending_tasks max_task_wait_time active_shards_percent
1624371902 14:25:02  black-cluster green          12         9   9304 4652    0    0        0             0                  -                100.0%
```

### Apply
```console
[root@noah ~]# blackbean apply settings -e
null       primaries
```
```console
[root@noah ~]# blackbean apply
apply cluster changes ... wordless

Usage:
  blackbean apply [command]

Available Commands:
  clearCache  apply indices to clear cache
  flush       apply indices to flush
  settings    apply cluster settings change

Flags:
  -h, --help   help for apply

Global Flags:
      --config string   config file (default is $HOME/.blackbean.yaml)

Use "blackbean apply [command] --help" for more information about a command.
```
```console
[root@noah ~]# blackbean apply settings -h
apply cluster settings ... wordless

Usage:
  blackbean apply [resource] [flags]

Flags:
  -e, --allocation_enable string                   to set allocation enable value, primaries or null
  -f, --breaker_fielddata string                   to set breaker_fielddata value, such as 60%
  -r, --breaker_request string                     to set breaker_request value, such as 60%
  -t, --breaker_total string                       to set breaker_total value, such as 60%
  -a, --cluster_concurrent_rebalanced string       to set cluster_concurrent_rebalanced value, such as 10
  -h, --help                                       help for apply
  -b, --max_bytes_per_sec string                   to set indices recovery max_bytes_per_sec, default 40
  -m, --max_compilations_rate string               to set max_compilations_rate value, such as 75/5m
  -s, --max_shards_per_node string                 to set max_shards_per_node value, such as 1000
  -n, --node_concurrent_recoveries string          to set node_concurrent_recoveries value, such as 10
  -i, --node_initial_primaries_recoveries string   to set node_initial_primaries_recoveries value, such as 10
  -w, --watermark_high string                      to set watermark_high value, such as 85%
  -l, --watermark_low string                       to set watermark_low value, such as 90%

Global Flags:
  -c, --cluster string   to specify a es cluster (default "default")
      --config string    config file (default is $HOME/.blackbean.yaml)
```
### Repo
```console
[root@noah ~]# blackbean repo
repo operations ... wordless

Usage:
  blackbean repo [command]

Available Commands:
  create      create specific snapshots
  delete      delete specific snapshots
  get         get specific repository

Flags:
  -h, --help   help for repo

Global Flags:
      --config string   config file (default is $HOME/.blackbean.yaml)

Use "blackbean repo [command] --help" for more information about a command.
```
### Snapshot
```console
[root@noah ~]# blackbean snapshot
snapshot operations ... wordless

Usage:
  blackbean snapshot [command]

Available Commands:
  create      create specific snapshots
  delete      delete specific snapshots
  get         get specific snapshots
  restore     get specific index to restore

Flags:
  -h, --help   help for snapshot

Global Flags:
      --config string   config file (default is $HOME/.blackbean.yaml)

Use "blackbean snapshot [command] --help" for more information about a command.
```

### Index
```console
[root@noah ~]# blackbean index
index operations ... wordless

Usage:
  blackbean index [command]

Available Commands:
  create      create index from command
  delete      delete index from command
  get         get index from cluster
  reindex     do reindex
  search      search index from cluster

...

```
```console
[root@noah ~]# blackbean index search test-* -f query.json
[200 OK] {
  "took" : 172,
  "timed_out" : false,
  ...
```
```console
[root@noah ~]# blackbean index search test-* -f query.yaml
[200 OK] {
  "took" : 172,
  "timed_out" : false,
  ...
```
```console
[root@noah ~]# blackbean index search test-* -d '{"query":{"match_all":{}}}'
[200 OK] {
  "took" : 172,
  "timed_out" : false,
  ...
```
```console
[root@noah ~]# blackbean index get test-2021.06
[200 OK] {
  "test-2021.06" : {
    ...
```

### Alias
```console
[root@noah ~]# blackbean alias
alias index ... wordless

Usage:
  blackbean alias [command]

Available Commands:
  create      create alias for index
  delete      delete alias for index
  get         get alias for index or get alias list

...
```

### Reroute
```console
[root@noah ~]# blackbean reroute
reroute for cluster ... wordless

Usage:
  blackbean reroute [command]

Available Commands:
  allocateReplicas allocate replicas index
  cancel           cancel allocating index
  failed           retry failed allocation
  move             move index
...
```

### User

```console
[root@noah ~]# blackbean user
user for cluster ... wordless

Usage:
  blackbean user [command]

Available Commands:
  create      create specify user
  delete      delete specify user
  get         get specify user
  update      update specify user
...
```

### Role
```console
[root@noah ~]# blackbean role
role operations for cluster ... wordless

Usage:
  blackbean role [command]

Available Commands:
  create      create specify user
  delete      delete specify role
  get         get specify role
  update      update specify user
...


```
### Explain
```console
[root@noah ~]# blackbean explain -h
explain index allocation ... wordless

Usage:
  blackbean explain [index] [flags]

Flags:
      --current_node string   specifies the node ID or the name of the node to only explain a shard that is currently located on the specified node.
  -h, --help                  help for explain
      --primary               if true, returns explanation for the primary shard for the given shard ID.
      --shard string          specifies the ID of the shard that you would like an explanation for.
...
```

### Watcher
```console
[root@noah ~]# blackbean watcher
operate watcherr ... wordless

Usage:
  blackbean watcher [command]

Available Commands:
  start       start watcher
  stats       get watcher stats
  stop        stop watcher
...
```


## Contact Me
Any advice is welcome! Please email to toughnoah@163.com
