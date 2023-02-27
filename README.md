# 🐙 GithubStatsAPI 📈

GithubStatsAPI provides fast access to GitHub user statistics and commits

<h2>Commits</h2>
<h3>Request sample</h3>

``` Elixir
https://githubstatsapi.vercel.app/api/commits?id=hud0shnik
```

``` Elixir
https://githubstatsapi.vercel.app/api/commits?id=hud0shnik&date=2022-01-21
```
<h3>Response sample</h3>

``` Json
{
"success":  true,
"error":    "",
"date":     "2022-01-21",
"username": "hud0shnik",
"commits":  9,
"color":    4
}
```
> ***color*** is color of the cell. There are 5 colors in total: from ***gray (0)*** to ***bright green (4)***

<h2>User</h2>
<h3>Request sample</h3>

``` Elixir
https://githubstatsapi.vercel.app/api/user?id=hud0shnik
```
<h3>Response sample</h3>

``` Json
{
"success":       true,
"error":         "",
"username":      "hud0shnik",
"name":          "Danila Egorov",
"followers":     "59",
"following":     "0",
"repositories":  "25",
"packages":      "0",
"stars":         "4",
"contributions": "1,980",
"status":        "Drawin'",
"avatar":        "https://avatars.githubusercontent.com/u/42404892"
}
```
<h2>Repo</h2>
<h3>Request sample</h3>

``` Elixir
https://githubstatsapi.vercel.app/api/repo?username=hud0shnik&reponame=OsuStatsApi
```
<h3>Response sample</h3>

``` Json
{
"success":  true,
"error":    "",
"username": "hud0shnik",
"reponame": "OsuStatsApi",
"commits":  "411",
"branches": "2",
"tags":     "0",
"stars":    "2",
"watching": "1",
"forks":    "1"
}
```

<img src="https://wakatime.com/badge/user/ee2709af-fc5f-498b-aaa1-3ea47bf12a00/project/a706f6cd-74fe-4b0f-ab24-a030f4bb3191.svg?style=for-the-badge">

[![License - BSD 3-Clause](https://img.shields.io/static/v1?label=License&message=BSD+3-Clause&color=%239a68af&style=for-the-badge)](/LICENSE)
