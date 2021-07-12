![Go 1.16](https://img.shields.io/badge/Go-v1.16-blue)
[![codecov](https://codecov.io/gh/toughnoah/blackbean/branch/master/graph/badge.svg?token=4UUTYZ6NCF)](https://codecov.io/gh/toughnoah/blackbean)
[![CI Workflow](https://github.com/toughnoah/blackbean/actions/workflows/test-coverage.yaml/badge.svg)](https://github.com/toughnoah/blackbean/actions/workflows/test-coverage.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/toughnoah/blackbean)](https://goreportcard.com/report/github.com/toughnoah/blackbean)
# blackbean
The blackbean is a command tool for elasticsearch operations by using cobra. Besides, blackbean is the name of my lovely French bulldog.

![avatar](img/blackbean.png)

<!-- vscode-markdown-toc -->
* 1. [Start to enjoy blackbean](#Starttoenjoyblackbean)
* 2. [Configuration](#Configuration)
* 3. [Shell completion](#Shellcompletion)
* 4. [Common Flag](#CommonFlag)
* 5. [Command](#Command)
	* 5.1. [Current es](#Currentes)
	* 5.2. [Use](#Use)
	* 5.3. [Cat](#Cat)
	* 5.4. [Apply](#Apply)
	* 5.5. [Repo](#Repo)
	* 5.6. [Snapshot](#Snapshot)
	* 5.7. [Index](#Index)
	* 5.8. [Alias](#Alias)
	* 5.9. [Reroute](#Reroute)
	* 5.10. [User](#User)
	* 5.11. [Role](#Role)
	* 5.12. [Explain](#Explain)
	* 5.13. [Template](#Template)
	* 5.14. [Watcher](#Watcher)
* 6. [Contact Me](#ContactMe)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->
##  1. <a name='Starttoenjoyblackbean'></a>Start to enjoy blackbean
```
curl -LO https://github.com/toughnoah/blackbean/releases/download/v1.0.2/blackbean
chmod +x ./blackbean
mv ./blackbean /usr/local/bin
```
```console
git clone -b master git@github.com:toughnoah/blackbean.git
cd blackbean
make
```

##  2. <a name='Configuration'></a>Configuration
Define your config `.blackbean.yaml` file firstly and put it into your home directory, do not forget specify `current` as one of your cluster.
```
cluster:
  prod:
    url: https://a.es.com:9200
    username: Noah
    password: abc
  qa:
    url: https://b.es.com:9200
    username: Noah
    password: abc
current: qa
```

##  3. <a name='Shellcompletion'></a>Shell completion
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

##  4. <a name='CommonFlag'></a>Common Flag
You can use `-d` or `--data` to specify raw request body like `'{"query":"match_all":{}}'`.
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
##  5. <a name='Command'></a>Command
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
###  5.1. <a name='Currentes'></a>Current es
```console
[root@noah ~]# blackbean current-es
current using cluster: qa
```

###  5.2. <a name='Use'></a>Use
```console
[root@noah ~]# blackbean use [tab][tab]
prod  qa
[root@noah ~]# blackbean use qa
change to cluster: qa

```
###  5.3. <a name='Cat'></a>Cat
```console
[root@noah ~]# blackbean cat health 
[200 OK] epoch      timestamp cluster       status node.total node.data shards  pri relo init unassign pending_tasks max_task_wait_time active_shards_percent
1624371902 14:25:02  black-cluster green          12         9   9304 4652    0    0        0             0                  -                100.0%
```

###  5.4. <a name='Apply'></a>Apply
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
  --allocation_enable string                   to set allocation enable value, primaries or null
  --breaker_fielddata string                   to set breaker_fielddata value, such as 60%
  --breaker_request string                     to set breaker_request value, such as 60%
  --breaker_total string                       to set breaker_total value, such as 60%
  --cluster_concurrent_rebalanced string       to set cluster_concurrent_rebalanced value, such as 10
  --help                                       help for apply
  --max_bytes_per_sec string                   to set indices recovery max_bytes_per_sec, default 40
  --max_compilations_rate string               to set max_compilations_rate value, such as 75/5m
  --max_shards_per_node string                 to set max_shards_per_node value, such as 1000
  --node_concurrent_recoveries string          to set node_concurrent_recoveries value, such as 10
  --node_initial_primaries_recoveries string   to set node_initial_primaries_recoveries value, such as 10
  --watermark_high string                      to set watermark_high value, such as 85%
  --watermark_low string                       to set watermark_low value, such as 90%

Global Flags:
  -c, --cluster string   to specify a es cluster (default "default")
      --config string    config file (default is $HOME/.blackbean.yaml)
```
###  5.5. <a name='Repo'></a>Repo
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
###  5.6. <a name='Snapshot'></a>Snapshot
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

###  5.7. <a name='Index'></a>Index
```console
[root@noah ~]# blackbean index
index operations ... wordless

Usage:
  blackbean index [command]

Available Commands:
  bulk        send bulk request
  create      create index from command
  delete      delete index from command
  get         get index from cluster
  msearch     send msearch request
  reindex     do reindex
  search      search index from cluster
  write       write index from command
...
```
```console
[root@noah ~]# blackbean index msearch -h
send msearch request ... wordless

Usage:
  blackbean index msearch [flags]

Flags:
  -d, --data string                         the path to raw file with request body (default "{}")
  -h, --help                                help for msearch
      --max_concurrent_searches int         ID of the pipeline to use to preprocess incoming documents.
      --max_concurrent_shard_requests int   if true, the request’s actions must target an index alias.
      --raw_file string                     the path to raw file with request body
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

###  5.8. <a name='Alias'></a>Alias
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

###  5.9. <a name='Reroute'></a>Reroute
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

###  5.10. <a name='User'></a>User

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

###  5.11. <a name='Role'></a>Role
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
###  5.12. <a name='Explain'></a>Explain
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

###  5.13. <a name='Template'></a>Template
```console
[root@noah ~]# blackbean template
template operations ... wordless

Usage:
  blackbean template [command]

Available Commands:
  apply       create or update template
  delete      delete or update template
  get         get template
...
```

###  5.14. <a name='Watcher'></a>Watcher
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


##  6. <a name='ContactMe'></a>Contact Me
Any advice is welcome! Please email to toughnoah@163.com
