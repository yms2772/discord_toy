# Discord Toy Bot
[bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)를 사용한 간단한 디스코드 봇

# Docker
## Self compile
```bash
git clone https://github.com/yms2772/discord_toy.git

cd discord_toy

docker build-t discordbot .

docker run -d -e TZ=Asia/Seoul -p 8080:80 discordbot
```

## Options
|env|description|value|
|------|---|---|
|TZ|set timezone|default: Asia/Seoul|

## Feature
+ ~dday
    + add {date: YYMMDD} {comment}
    + del {target id}
