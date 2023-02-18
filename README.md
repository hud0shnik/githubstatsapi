# 🐙 My GitHub API 📈

[![License - BSD 3-Clause](https://img.shields.io/static/v1?label=License&message=BSD+3-Clause&color=%239a68af&style=for-the-badge)](/LICENSE)

<details open="true">
   <summary> 🇬🇧 <b>English Version</b> 🇬🇧 </summary>
   
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
   
</details>

<!---------------------------------------------- Russian Version ----------------------------------------->

<details>
   <summary> 🇷🇺 <b>Русская версия</b> 🇷🇺 </summary>
   <h2>Коммиты</h2>
   <h3>Семпл запроса</h3>
  
   ``` Elixir
   https://githubstatsapi.vercel.app/api/commits?id=hud0shnik
   ```
  
   ``` Elixir
   https://githubstatsapi.vercel.app/api/commits?id=hud0shnik&date=2022-01-21
   ```
   <h3>Семпл ответа</h3>
  
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
> Параметр ***color*** - цвет ячейки. Всего есть 5 цветов: от ***серого (0)*** до ***ярко-зеленого (4)***
   
   <h2>Пользователь</h2>
   <h3>Семпл запроса</h3>
  
   ``` Elixir
   https://githubstatsapi.vercel.app/api/user?id=hud0shnik
   ```
   <h3>Семпл ответа</h3>
  
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
   
   <h2>Репозиторий</h2>
   <h3>Семпл запроса</h3>
  
   ``` Elixir
   https://githubstatsapi.vercel.app/api/repo?username=hud0shnik&reponame=OsuStatsApi
   ```
   <h3>Семпл ответа</h3>
  
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
</details>
