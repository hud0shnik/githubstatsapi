# 🐙 GithubStatsAPI 📈

<i><b>GithubStatsAPI</i></b> provides fast access to GitHub user statistics and commits

<h2>/commits/</h2>
<h3>Request</h3>

``` Elixir
https://githubstatsapi.vercel.app/api/commits
```

Parameter       | Value type | Required | Description   
----------------|------------|----------|------------
id              |   string   |    Yes   | username
date            |   string   |    No    | date (like 2022-01-21)

<h3>Structures</h3>

<h4>UserCommits</h4>

Field                       |       Type         | Description
----------------------------|--------------------|------------
success                     |        bool        |
error                       |       string       | api error response (default value= "")
date                        |       string       | date (like 2022-01-21)
username                    |       string       |
commits                     |        int         |
color                       |        int         |


> ***color*** is color of the cell. There are 5 colors in total: from ***gray (0)*** to ***bright green (4)***

<h2>/user/</h2>
<h3>Request</h3>

``` Elixir
https://githubstatsapi.vercel.app/api/user
```

Parameter       | Value type | Required | Description   
----------------|------------|----------|------------
id              |   string   |    Yes   | username
type            |   string   |    No    | response type (like "string")

<h3>Structures</h3>

<h4>UserInfo</h4>

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


<h2>/repo/</h2>
<h3>Request</h3>

``` Elixir
https://githubstatsapi.vercel.app/api/repo
```

Parameter       | Value type | Required | Description   
----------------|------------|----------|-------------
username        |   string   |    Yes   | username
reponame        |   string   |    Yes   | reponame
type            |   string   |    No    | response type (like "string")

<h3>Structures</h3>

<h4>RepoInfo</h4>

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


<img src="https://wakatime.com/badge/user/ee2709af-fc5f-498b-aaa1-3ea47bf12a00/project/a706f6cd-74fe-4b0f-ab24-a030f4bb3191.svg?style=for-the-badge">

[![License - BSD 3-Clause](https://img.shields.io/static/v1?label=License&message=BSD+3-Clause&color=%239a68af&style=for-the-badge)](/LICENSE)
