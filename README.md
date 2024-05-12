# mon-dell-me4012

## Описание

Простой мониторинг СХД Dell ME4012. Скорее всего, работать будет и с другими СХД линейки ME4, не проверялось.
На данный момент реализован сбор следующих данных:
- дисковые группы
- сенсоры
- блоки питания
- вольюмы

Для каждого типа дарнных реализован дискаверинг.

Пример использования в zabbix-agent:
```
UserParameter=me4012.discover[*], cd /opt/mon && ./mon-dell-me4012 -discovery -discovery_name $1
UserParameter=me4012.metric[*], cd /opt/mon && ./mon-dell-me4012 --metric_group $1 -entity_name $2 -metric_name $3
```

## Зависимости

- Redis
