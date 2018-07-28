package database

import (
    _ "github.com/lib/pq"
)

func CreateDomain(domain string, username string, password string, pub_key string) error {
    var id string
    return db.QueryRow("INSERT INTO Domains (domain, username, key_id, passhash) VALUES ($1, $2, pgp_key_id(dearmor($4)), pgp_pub_encrypt($3, dearmor($4))) RETURNING id;", domain, username, password, pub_key).Scan(&id)
}

func UpdateDomain(domain string, username string, password string, pub_key string) error {
    var id string
    return db.QueryRow("UPDATE Domains SET username = $2, passhash = pgp_pub_encrypt($3, dearmor($4)) WHERE domain = $1 AND key_id = pgp_key_id(dearmor($4)) RETURNING id;", domain, username, password, pub_key).Scan(&id)
}

func DeleteDomain(domain string, pub_key string) error {
    var id string
    return db.QueryRow("DELETE FROM Domains WHERE domain = $1 AND key_id = pgp_key_id(dearmor($2)) RETURNING id;", domain, pub_key).Scan(&id)
}

func ReadDomainUserPass(domain string, pub_key string, prv_key string, prv_pass string) (string, string, error) {
    var username, password string
    if err := db.QueryRow("SELECT username, pgp_pub_decrypt(passhash, dearmor($3), $4) As password FROM Domains WHERE domain = $1 AND key_id = pgp_key_id(dearmor($2));", domain, pub_key, prv_key, prv_pass).Scan(&username, &password); err != nil {
        return "", "", err
    }
    return username, password, nil
}

func GetKeyId(pub_key string) (string, error) {
    var key_id string
    if err := db.QueryRow("SELECT pgp_key_id(dearmor($1));", pub_key).Scan(&key_id); err != nil {
        return "", err
    }
    return key_id, nil
}

func GetDomainList(pub_key string) ([]string, error) {
    rows, err := db.Query("SELECT domain FROM Domains WHERE key_id = pgp_key_id(dearmor($1));", pub_key)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var domains []string

    for rows.Next() {
        var domain string
        if err := rows.Scan(&domain); err != nil {
            return nil, err
        }
        domains = append(domains, domain)
    }

    if rows.Err() != nil {
        return nil, rows.Err()
    }

    return domains, nil
}
