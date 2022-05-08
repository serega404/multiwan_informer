# MultiWan Informer
![License](https://img.shields.io/github/license/serega404/multiwan_informer?style=flat-square&)
![Go version](https://img.shields.io/github/go-mod/go-version/serega404/multiwan_informer?style=flat-square&filename=src%2Fgo.mod)

Язык: [Русский](./README_RU.md) | [English](./README.md)

## Для чего эта программа?

Эта программа нужна тем кому нужно следить за подключением к сети нескольких сетевых карт, благодаря этой программе вы сможете получать уведомления в телеграмме об отключении или подключении сети на определённой сетевой карте

## Скриншоты

![Telegram chat](./img/telegram.png)

## Как настроить?

Для начала скачайте артефакт сборки, после чего настрйоки ваши параметры в файле ```conf.json```

### Настройка

<details>
<summary>Пример конфига</summary>

```json
{
    "Interfaces": [{
            "DisplayName": "MainInterface",
            "IpOrInterfaceName": "enp2s0"
        },
        {
            "DisplayName": "SecondInterface",
            "IpOrInterfaceName": "enp3s0"
        }
    ],
    "WaitTimeSec": 15,
    "PingAddr": "8.8.8.8",
    "TelegramConf": {
        "BotToken": "Token",
        "ChatID": "Id",
        "SendSilent": "false"
    }
}
```

</details>

### Параметры

* `Interfaces` - массив интерфейсов
  * `DisplayName` - имя выводимое в уведомлении
  * `IpOrInterfaceName` - название или собственный IP адрес интерфейса
* `WaitTimeSec` - сколько секунд ожидать до успешного пинга (```ping -w```)
* `PingAddr` - какой адрес пинговать (желательно не домен)
* `TelegramConf` - настройки телеграм бота
  * `BotToken` - токен бота ([получать тут](https://t.me/BotFather/))
  * `ChatID` - id чата ([узнать](https://t.me/chatIDrobot/))
* `SendSilent` - отправлять сообщения без звука

## Сборка

Cборка под вашу систему:
</br>```build -o ./build/multiwan_informer ./src/main.go```

Пример сборки под MIPS:
</br>```GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -o ./build/multiwan_informer_mipsle ./src/main.go```

## Лицензия

Распространяется под лицензией GPLv3. Дополнительные сведения смотрите в файле [`LICENSE`](./LICENSE).
