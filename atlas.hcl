variable "src" {
    type = string
    default = "file://script/init.sql"
}

variable "migrations" {
    type = string
    default = "file://migrations"
}

variable "simulation" {
    type = string
    default ="docker://postgres/15.4/dev?search_path=public"
}

variable "postgres_user" {
    type = string
    default = "admin"
}

variable "postgres_password" {
    type = string
    default = "1"
}

variable "postgres_host" {
    type = string
    default = "localhost"
}

variable "postgres_db" {
    type = string
    default = "ewallet"
}

env "local" {
    // definition of desired schema
    src = var.src

    // target database
    url = format(
        "postgres://%s:%s@%s/%s?search_path=public&sslmode=disable",
        var.postgres_user, var.postgres_password, var.postgres_host, var.postgres_db,
    )

    // simulation database
    dev = var.simulation

    // migration setting
    migration {
        dir = var.migrations
    }
}

env "script" {
    url = "file://script/init.sql"
    dev = var.simulation
}

env "ci" {
    // definition of desired schema
    src = var.src

    // simulation database
    dev = var.simulation

    // migration setting
    migration {
        dir = var.migrations
    }
}

env "cd" {
    migration {
        dir = var.migrations
    }

    url = format(
        "postgres://%s:%s@%s/%s?search_path=public&sslmode=disable",
        var.postgres_user, var.postgres_password, var.postgres_host, var.postgres_db,
    )
}
