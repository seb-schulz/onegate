# OneGate

OneGate is a legacy-free single-sign on service. It uses rather passkeys or fido2 stick for user authentication than passwords.

## Develop Environment

This project is using [devcontainer](https://containers.dev/). It is using [podman](https://podman.io/) and [podman-compose](https://github.com/containers/podman-compose) to get it up and running.

As database only MariaDB is currently supported. In case you need access while you developing, you could use `podman exec -ti onegatedevcontainer_db_1 /bin/bash -c 'mariadb -h 127.0.0.1 -u root --password=$MARIADB_ROOT_PASSWORD $MARIADB_DATABASE'`.
