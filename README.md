# 🐙 GithubStatsAPI 📈

<i><b>GithubStatsAPI</i></b> provides fast access to GitHub user statistics and commits

## Overview

- [Overview](#overview)
- [Commits](#commits)
- [User](#user)
- [Repo](#repo)
- [Samples](#samples)

## Commits

### Request

``` Elixir
https://githubstatsapi.vercel.app/api/commits
```

Parameter       | Value type | Required | Description   
----------------|------------|----------|------------
id              |   string   |    Yes   | username
date            |   string   |    No    | date (like 2022-01-21)

### Structures

#### UserCommits

Field                       |       Type         | Description
----------------------------|--------------------|------------
success                     |        bool        |
error                       |       string       | api error response (default value= "")
date                        |       string       | date (like 2022-01-21)
username                    |       string       |
commits                     |        int         |
color                       |        int         |


> ***color*** is color of the cell. There are 5 colors in total: from ***gray (0)*** to ***bright green (4)***

## User

### Request

``` Elixir
https://githubstatsapi.vercel.app/api/user
```

Parameter       | Value type | Required | Description   
----------------|------------|----------|------------
id              |   string   |    Yes   | username
type            |   string   |    No    | response type (like "string")

### Structures

#### UserInfo

Field                       |       Type         | Description
----------------------------|--------------------|------------
success                     |        bool        |
error                       |       string       | api error response (default value= "")
username                    |       string       |
name                        |       string       |
followers                   |       string       |
following                   |       string       |
repositories                |       string       |
packages                    |       string       |
stars                       |       string       |
contributions               |       string       |
status                      |       string       |
avatar                      |       string       | avatar url


## Repo

### Request

``` Elixir
https://githubstatsapi.vercel.app/api/repo
```

Parameter       | Value type | Required | Description   
----------------|------------|----------|-------------
username        |   string   |    Yes   | username
reponame        |   string   |    Yes   | reponame
type            |   string   |    No    | response type (like "string")

### Structures

#### RepoInfo

Field                       |       Type         | Description
----------------------------|--------------------|------------
success                     |        bool        |
error                       |       string       | api error response (default value= "")
username                    |       string       |
reponame                    |       string       |
commits                     |       string       |
branches                    |       string       |
tags                        |       string       |
stars                       |       string       |
watching                    |       string       |
forks                       |       string       |


### Samples

#### Request

``` Elixir
https://githubstatsapi.vercel.app/api/commits?id=hud0shnik&date=2023-03-05
``` 

#### Response

``` Json
{
  "success": true,
  "error": "",
  "date": "2023-03-05",
  "username": "hud0shnik",
  "commits": 6,
  "color": 3
}
```

#### Request

``` Elixir
https://githubstatsapi.vercel.app/api/user?id=hud0shnik
``` 

#### Response

``` Json
{
  "success": true,
  "error": "",
  "username": "hud0shnik",
  "name": "Danila Egorov",
  "followers": 61,
  "following": 0,
  "repositories": 29,
  "packages": 0,
  "stars": 7,
  "contributions": 0,
  "status": "Drawin'",
  "avatar": "https://avatars.githubusercontent.com/u/42404892?v=4"
}
```

#### Request

``` Elixir
https://githubstatsapi.vercel.app/api/repo?username=hud0shnik&reponame=osustatsapi
``` 

#### Response

``` Json
{
  "success": true,
  "error": "",
  "username": "hud0shnik",
  "reponame": "osustatsapi",
  "commits": 658,
  "branches": 1,
  "tags": 0,
  "stars": 2,
  "watching": 1,
  "forks": 1
}
```


<img src="https://wakatime.com/badge/user/ee2709af-fc5f-498b-aaa1-3ea47bf12a00/project/a706f6cd-74fe-4b0f-ab24-a030f4bb3191.svg?style=for-the-badge">

[![License - BSD 3-Clause](https://img.shields.io/static/v1?label=License&message=BSD+3-Clause&color=%239a68af&style=for-the-badge)](/LICENSE)
