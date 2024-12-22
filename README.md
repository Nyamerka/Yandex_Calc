# Yandex-калькулятор.

## Made by Бутер Бродский aka Nyamerka)

Данное api было написано в рамках прохождения курса разработки на языке Go от Яндекса.

### Как запустить?

```
cd cmd
go run main.go
```

После ввода этих команд Вы должны получить:

```
Server is running on port :8080...
```

### Примеры использования:

Примеры использования с корректными входными данными

Пример 1:
```
 curl -X POST -H "Content-Type: application/json" -d '{"expression": "3 + 2"}' http://localhost:8080/api/v1/calculate
```

Результат:
```
{"result":"5"}
```

Пример 2:
```
curl -X POST -H "Content-Type: application/json" -d '{"expression": "3 + 2 * 5"}' http://localhost:8080/api/v1/calculate
```

Результат:
```
{"result":"13"}
```

Пример 3 (более сложный):
```
curl -X POST -H "Content-Type: application/json" -d '{"expression": "(1+1)*2-5*(3+2)+(-6+4)"}' http://localhost:8080/api/v1/calculate
```

Результат:
```
{"result":"-23"}
```

Примеры некорректного использования:

Пример 1 (не хватает закрывающей скобки):
```
curl -X POST -H "Content-Type: application/json" -d '{"expression": "-1+2--1*((2+3)"}' http://localhost:8080/api/v1/calculate
```

Результат:
```
{"error":"Expression is not valid"}
```