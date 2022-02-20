# 🐙 API статистики пользоватя GitHub 📈
<details open="true">
   <summary> 🇬🇧 <b>English Version</b> 🇬🇧 </summary>
   <h3>Request sample </h3>
  
   ``` Elixir
   GET https://hud0shnikgitapi.herokuapp.com/hud0shnik
   ```
  
   ``` Elixir
   GET https://hud0shnikgitapi.herokuapp.com/hud0shnik/2022-01-21
   ```
   <h3>Response sample </h3>
  
   ``` Json
  {
  "date":     "2022-01-21",
  "username": "hud0shnik",
  "name":     "Danila Egorov",
  "avatar":   "https://avatars.githubusercontent.com/u/42404892",
  "commits":  9,
  "color":    4,
  "stars":    6
  }
   ```
   > ***color*** is color of the cell. There are 5 colors in total: from ***gray (0)*** to ***bright green (4)***
</details>

<!---------------------------------------------- Russian Version ----------------------------------------->

<details>
   <summary> 🇷🇺 <b>Русская версия</b> 🇷🇺 </summary>
<h3>Семпл реквеста </h3>

``` Elixir
GET https://hud0shnikgitapi.herokuapp.com/hud0shnik
```

``` Elixir
GET https://hud0shnikgitapi.herokuapp.com/hud0shnik/2022-01-21
```
<h3>Семпл респонса </h3>

``` Json
  {
  "date":     "2022-01-21",
  "username": "hud0shnik",
  "name":     "Danila Egorov",
  "avatar":   "https://avatars.githubusercontent.com/u/42404892",
  "commits":  9,
  "color":    4,
  "stars":    6
  }
```
> Параметр ***color*** - цвет ячейки. Всего есть 5 цветов: от ***серого (0)*** до ***ярко-зеленого (4)***
</details>
