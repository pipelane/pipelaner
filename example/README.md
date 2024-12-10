# Использование pkl во внешних проектах
1) В проекте создать директорию pkl, в ней создать файлик PklProject, со следующим контентом:
```
dependencies {
  ["pipelaner"] {
    uri = "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@x.x.x"
  }
}
```
2) Объявить корневой файл с конфигурацией в заголовке файла объявить:
```
amends "@pipelaner/Pipelaner.pkl" // нужно разобраться чтобы работало так

amends "package://pkg.pkl-lang.org/github.com/pipelane/pipelaner/pipelaner@0.0.4#/Pipelaner.pkl" // сейчас получилось только без алиаса
```
3) Конфигурируем сами пайплайны, пример в pkl/config.pkl
4) При необходимости создаем свои имплементации компонент pkl/components.pkl
5) В случае если были созданы свои компоненты необходимо сгенерировать код для их использования для файла pkl/components.pkl, используется следующая команда:
```
pkl-gen-go pkl/example.pkl
```
6) Имплементируем созданные компоненты и регистрируем их в source pipelaner'a пример в main.go