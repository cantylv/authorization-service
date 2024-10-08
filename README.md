# HTTP-микросервис авторизации
Необходимо реализовать микросервис для работы с привилегиями пользователей (создание/
удаление пользователя, добавить/убрать права пользователя, проверка прав пользователя). Сервис 
должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON.

## Содержание
- [Технологии](#технологии)
- [Предметная область](#предметная-область)
- [ER диаграммы](#er-диаграммы)
- [Концептуальные схемы](#концептуальные-схемы)
- [CI/CD](#ci/cd)
- [Использование](#использование)
- [API](#api)


## Технологии
- [Golang 1.23.1](https://go.dev/dl/)
- [Docker 4.31.0](https://docs.docker.com/engine/install/)
- [Viper 1.19.0](https://github.com/spf13/viper)

## Предметная область
Для микросервиса привелегий/прав пользователя я решил сделать группы с правами на использование агентов. Все права группы
наследуются ее участниками, т.е. если у группы 'devs' есть доступ к агенту с именем 'pools', то у всех участников группы 
'devs' будет доступ к агенту 'pools'. Это позволяет экономить память при хранении прав пользователей, но немного проигрывает
в cpu, так как требует join. Однако для более гибкой настройки прав пользователя была создана ручка, которая добавляет 
агента пользователю, т.е. у пользователя есть два хранилища прав - его собственное и унаследованное от всех групп, в которые он входит. В итоге доступ к агенту будет, если хотя бы в одном из хранилищ есть агента.  

В такой системе должен быть root пользователь, который будет иметь абсолютные права на все группы, агенты. Именно он создает агенты, которые используются в task_manager. Например, я написал archive_manager, к которому в начальный момент времени не имеет доступ никто, потому что информации о нем в микросервисе привелегии/прав еще нет. root пользователь должен добавить агента archive, чтобы к нему можно было обращаться. После добавления, только root имеет доступ.  

Пользователь может получить доступ к агенту 3 способами:
1) root дает прямые права пользователю на пользование агентом
2) пользователь А создает заявку на создание группы, root принимает заявку, вследствие чего создается группа с ответственным в лице пользователя А. root пользователь наделяет группу правами пользованиями услугами агента, следовательно пользователь получает доступ к агенту.
3) ответственный за группу может добавить в нее пользователя, после этого он получит права группы.

## ER диаграммы

### Микросервис прав пользователя 
```mermaid
erDiagram
    "user" {
        UUID id PK "DEFAULT gen_random_uuid()"
        TEXT(6-50) email UK "NOT NULL"
        TEXT(145) password "NOT NULL"
        TEXT(2-50) first_name "NOT NULL"
        TEXT(2-50) last_name "NOT NULL"
        TIMESTAMPTZ created_at "DEFAULT now()"
        TIMESTAMPTZ updated_at "DEFAULT now()"
    }

    "group" {
        INT id PK "GENERATED ALWAYS AS IDENTITY"
        TEXT(2-30) name UK "NOT NULL"
        UUID owner_id FK "ON DELETE RESTRICT"
        TIMESTAMPTZ created_at "DEFAULT now()"
        TIMESTAMPTZ updated_at "DEFAULT now()"
    }

    agent {
        INT id PK "GENERATED ALWAYS AS IDENTITY"
        TEXT(2-50) name UK "NOT NULL"
        TIMESTAMPTZ created_at "DEFAULT now()"
    }

    bid {
        INT id PK "GENERATED ALWAYS AS IDENTITY"
        TEXT group_name
        UUID user_id FK "ON DELETE CASCADE"
        status_type status
        TIMESTAMPTZ created_at "DEFAULT now()"
        TIMESTAMPTZ updated_at "DEFAULT now()"
    }

    privelege_group {
        INT id PK "GENERATED ALWAYS AS IDENTITY"
        INT agent_id FK "ON DELETE CASCADE"
        INT group_id FK "ON DELETE CASCADE"
        TIMESTAMPTZ created_at "DEFAULT now()"
    }

    privelege_user {
        INT id PK "GENERATED ALWAYS AS IDENTITY"
        INT agent_id FK "ON DELETE CASCADE"
        INT user_id FK "ON DELETE CASCADE"
        TIMESTAMPTZ created_at "DEFAULT now()"
    }

    participation {
        INT id PK "GENERATED ALWAYS AS IDENTITY"
        UUID user_id FK "ON DELETE CASCADE"
        INT group_id FK "ON DELETE CASCADE"
    }

    "user" ||--o{ bid : "has"
    "user" ||--o{ participation : "participates in"
    "user" ||--o{ privelege_user : "has access to"
    "group" ||--o{ bid : "has"
    "group" ||--o{ participation : "has"
    "group" ||--o{ privelege_group : "has access to"
    agent ||--o{ privelege_group : "is accessible by"
    agent ||--o{ privelege_user : "is accessible by"
```

### Archive manager  
```mermaid
erDiagram
    record {
        INT id PK "GENERATED ALWAYS AS IDENTITY"
        TEXT(10-500) text "NOT NULL"
    }
```

## Концептуальная схема
![Концептуальная схема работы системы](./src/scheme.png)


## CI/CD
Настроен только CI: используется staticcheck и линтер, также идет проверка тестов, но тесты не успел написать, хотя опыт в написании unit-тестов есть (testify + gomock).

## Использование
Для того запустить у себя проект необходимо:
1) Установить go и docker нужных версий, также должен быть установлен клиент git.
2) Склонировать репозиторий с кодом и запустить следующие команды:
```
git clone git@github.com:cantylv/authorization-service.git
make init 
make start
```
После выполнения этих команд вы можете делать запросы, пример запросов будет ниже.

## API
Вы можете посмотреть OpenAPI [здесь](src/open-api.yaml).