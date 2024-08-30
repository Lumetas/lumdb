# lumdb
Эта поделка - моя собственная субд написанная на go в рамках изучения языка go и как давняя идея лежавшая на антресолях.

При сборке пожалуйста укажите стандартный путь и разделитель пути в константах "app_path" и "path_sepirator" для windows и linux/unix они будут отличаться. Так же рядом советую вам указать порт.

## Использование
После сборки и запуска вы можете отправлять на адрес сервер и указанный вами порт пост запросы для управления базой. Я приведу примеры действий. В формате поле1=значение1 & поле2=значение2

Инициализация: init=password

Создание базы данных: password=password & create_data_base=db_name

Создание таблицы: password=password & db=db_name & create_table=table_name

Удаление таблицы: password=password & db=db_name & delete_table=table_name

Удаление базы: password=password & delete_db=db_name

Добавление записи в таблицу: password=password & db=db_name & table=table_name & data=string

Поиск по таблице: password=password & db=db_name & table=table_name & select=query & count=count |& index=true (count-количество возвращаемых значений. index-необязательный параметр, если равен "true", то вернутся идентификаторы записей)

Удаление записи: password=password & db=db_name & table=table_name & remove=id

Изменение записи: password=password & db=db_name & table=table_name & edit=id & data=string

Просмотр списка таблиц в базе: password=password & db=db_name & show_tables=1

Просмотр списка баз данных: password=password & show_databases=1
