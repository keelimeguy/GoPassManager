package database

import (
    "database/sql"
    "io/ioutil"
    "fmt"

    _ "github.com/lib/pq"
)

var (
    db *sql.DB
)

func Open(user string, password string, name string) error {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, name)
    var err error
    if db, err = sql.Open("postgres", dbinfo); err != nil {
        return err
    }

    if _, err = db.Exec("SELECT 1"); err != nil {
        return err
    }

    return nil
}

func Close() error {
    return db.Close()
}

func Create(user string, password string, name string, file string) error {
    instr_bytes, err := ioutil.ReadFile(file)
    if err != nil {
        return err
    }
    return CreateFromString(user, password, name, string(instr_bytes))
}

func CreateFromString(user string, password string, name string, instr_str string) error {
    var err error
    dbinfo := fmt.Sprintf("user=%s password=%s sslmode=disable", user, password)

    if db, err = sql.Open("postgres", dbinfo); err != nil {
        return err
    }

    if _, err = db.Exec("DROP DATABASE IF EXISTS "+name); err != nil {
        return err
    }

    if _, err = db.Exec("CREATE DATABASE "+name); err != nil {
        return err
    }

    db.Close()

    if err = Open(user, password, name); err != nil {
        Delete(user, password, name)
        return err
    }

    if _, err = db.Exec(instr_str); err != nil {
        db.Close()
        Delete(user, password, name)
        return err
    }

    db.Close()
    return nil
}

func Delete(user string, password string, name string) error {
    var err error
    dbinfo := fmt.Sprintf("user=%s password=%s sslmode=disable", user, password)

    if db, err = sql.Open("postgres", dbinfo); err != nil {
        return err
    }
    defer db.Close()

    if _, err = db.Exec("DROP DATABASE "+name); err != nil {
        return err
    }

    return nil
}
