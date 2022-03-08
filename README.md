# 🐙 API статистики пользоватя GitHub 📈
<details open="true">
   <summary> 🇬🇧 <b>English Version</b> 🇬🇧 </summary>
   
   <h2>Commits </h2>
   <h3>Request sample </h3>
  
   ``` Elixir
   GET https://hud0shnikgitapi.herokuapp.com/commits/hud0shnik
   ```
  
   ``` Elixir
   GET https://hud0shnikgitapi.herokuapp.com/commits/hud0shnik/2022-01-21
   ```
   <h3>Response sample </h3>
  
   ``` Json
  {
  "date":     "2022-01-21",
  "username": "hud0shnik",
  "commits":  9,
  "color":    4
  }
   ```
   > ***color*** is color of the cell. There are 5 colors in total: from ***gray (0)*** to ***bright green (4)***
   
   <h2>Info</h2>
   <h3>Request sample </h3>
  
   ``` Elixir
   GET https://hud0shnikgitapi.herokuapp.com/info/hud0shnik
   ```
   <h3>Response sample </h3>
  
   ``` Json
   {
  "username":      "hud0shnik",
  "name":          "Danila Egorov",
  "followers":     "14",
  "following":     "9",
  "repositories":  "19",
  "packages":      "0",
  "stars":         "6",
  "contributions": "974",
  "status": "Thinkin' about new API",
  "avatar":        "https://avatars.githubusercontent.com/u/42404892"
   }
   ```
   
</details>

<!---------------------------------------------- Russian Version ----------------------------------------->

<details>
   <summary> 🇷🇺 <b>Русская версия</b> 🇷🇺 </summary>
   <h2>Коммиты </h2>
   <h3>Семпл запроса </h3>
  
   ``` Elixir
   GET https://hud0shnikgitapi.herokuapp.com/commits/hud0shnik
   ```
  
   ``` Elixir
   GET https://hud0shnikgitapi.herokuapp.com/commits/hud0shnik/2022-01-21
   ```
   <h3>Семпл ответа</h3>
  
   ``` Json
  {
  "date":     "2022-01-21",
  "username": "hud0shnik",
  "commits":  9,
  "color":    4
  }
   ```
> Параметр ***color*** - цвет ячейки. Всего есть 5 цветов: от ***серого (0)*** до ***ярко-зеленого (4)***
   
   <h2>Информация</h2>
   <h3>Семпл запроса </h3>
  
   ``` Elixir
   GET https://hud0shnikgitapi.herokuapp.com/info/hud0shnik
   ```
   <h3>Семпл ответа </h3>
  
   ``` Json
   {
  "username":      "hud0shnik",
  "name":          "Danila Egorov",
  "followers":     "14",
  "following":     "9",
  "repositories":  "19",
  "packages":      "0",
  "stars":         "6",
  "contributions": "974",
  "status": "Thinkin' about new API",
  "avatar":        "https://avatars.githubusercontent.com/u/42404892"
   }
   ```
</details>
