# Golang Middleware for 1C http service

# Описание
    В связи с тем, что Платформа 1С до сих порт (8.3.18.*) с сеансами работает довольно прохо, 
    а так же нет возможности авторизации с помощью токенов, был разработан этот микросервис.

    Основная идея:
        1. Проксировавть любые запросы от любых клиентов в 1С.
        2. Для авторизации получать в методе auth получать логин и пароль, пытаться с их помощью 
            поднять сессию на стороне 1С, при успешной автотризации генерирвать токен и возвращаьт его клиенту
        3. Для последующих запросов использовать ранее выданный токен
        4. Для получения инфоррмации о причине отказа в авторизации на шаге 1 предполагается использовать
            метод check-login, для него в конфигурации присутствует логин и пароль
        5. При необходимости можно добавить и другие внутренние ресурсы и использовать этот микросервис как точку входа.
        6. При большой нагрузке возможно поднять несколько копий данного сервиса

# Требования
    - KeyDB - для хранения сессий пользователей
    - http сервис на 1С - куда будут перенаправляться все запросы
    - jaeger - не обязательно. Можно не включать.
    
    - На стороне 1С включить переиспользование сеансов для сервиса. Оптимально - 200 сек.

# Использование
    Необходимо создать файл настроек. Пример находится в файле - config_sample.yml

    1. Компиляция и использование как сервис. В корне есть MakeFile. А так же есть файл сервиса в каталоге - systemd.
        Требуется дополнительно установить и настроить KeyDB сервер.
    2. Использование в Docker. В корне есть docker-compose.yaml.
    3. Запуск без компиляции - make run
