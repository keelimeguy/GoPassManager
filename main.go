package main

import (
    "os/exec"
    "strconv"
    "bufio"
    "fmt"
    "os"

    db "pass_manager/database"
)

var (
    db_CREATE_COMMAND string
    pub_key string
    prv_key string
    prv_pass string
    database_user string
    database_password string
    database_name string
)

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Panic:",r)

            reader := bufio.NewReader(os.Stdin)
            fmt.Print("\nEnter to exit.")
            reader.ReadString('\n')
        }
    }()

    if len(os.Args) != 7 {
        for i, arg := range os.Args {
            fmt.Println(i,arg)
        }
        panic("Usage: "+os.Args[0]+" <pub_key> <prv_key> <prv_pass> <database_user> <database_password> <database_name>")
    }
    pub_key = os.Args[1]
    prv_key = os.Args[2]
    prv_pass = os.Args[3]
    database_user = os.Args[4]
    database_password = os.Args[5]
    database_name = os.Args[6]

    if pub_key == "" || prv_key == "" {
        panic("Unauthorized.")
    }

    if err := db.Open(database_user, database_password, database_name); err != nil {
        if err.Error() == "pq: database \""+database_name+"\" does not exist" {
            fmt.Println("Creating database.")
            err = db.CreateFromString(database_user, database_password, database_name, db_CREATE_COMMAND)
            check(err)
            err = db.Open(database_user, database_password, database_name)
        }
        check(err)
    }
    defer db.Close()

    mainLoop()
}

func clear(msg string) {
    cmd := exec.Command("/bin/bash", "-c", "clear;")
    cmd.Stdout = os.Stdout
    check(cmd.Run())
    if msg != "" {
        fmt.Println(msg)
        fmt.Println()
    }
}

func mainLoop() {
    id, err := db.GetKeyId(pub_key)
    check(err)

    reader := bufio.NewReader(os.Stdin)
    domains, err := db.GetDomainList(pub_key)
    check(err)

    num_per_page := 15
    num_pages := int(len(domains)/num_per_page)+1
    cur_page := 1

    done := false
    for !done {
        header := ""

        fmt.Println("Active Key:", id)
        fmt.Println("Page", cur_page, "/", num_pages)

        if len(domains) > 0 {
            for i, domain := range domains {
                if i < num_per_page*(cur_page-1) {
                    continue
                }
                if i >= num_per_page*cur_page {
                    break
                }
                fmt.Println((i+1),":",domain)
            }
        } else {
            fmt.Println("(no domains)")
        }

        fmt.Println(">Enter domain number to read domain password.")
        fmt.Println(">Enter '>' and '<' to navigate pages.")
        fmt.Println(">Enter 'd<num>' to delete selected domain.")
        fmt.Println(">Enter 'n' to create new domain.")
        fmt.Println(">Enter 'q' to exit.")
        fmt.Print(": ")
        choice, err := reader.ReadString('\n')
        check(err)
        choice = choice[:len(choice)-1]

        switch choice {
            case "<":
                if cur_page > 1 {
                    cur_page--
                }

            case ">":
                if cur_page < num_pages {
                    cur_page++
                }

            case "n":
                fmt.Print("Domain name: ")
                domain, err := reader.ReadString('\n')
                check(err)
                domain = domain[:len(domain)-1]

                fmt.Print("Domain username: ")
                cmd := exec.Command("/bin/bash", "-c", "read -s -p '' user; echo $user;")
                cmd.Stdin = os.Stdin
                out, err := cmd.Output()
                check(err)
                fmt.Println()

                username := string(out)
                username = username[:len(username)-1]

                fmt.Print("Domain password: ")
                cmd = exec.Command("/bin/bash", "-c", "read -s -p '' pass; echo $pass;")
                cmd.Stdin = os.Stdin
                out, err = cmd.Output()
                check(err)
                fmt.Println()

                password := string(out)
                password = password[:len(password)-1]

                if err := db.CreateDomain(domain, username, password, pub_key); err != nil {
                    if err.Error() == "pq: duplicate key value violates unique constraint \"domains_domain_key_id_key\"" {
                        fmt.Println("Domain exists, replace?")
                        fmt.Print("(y)es else no: ")
                        var input string
                        input, err = reader.ReadString('\n')
                        check(err)
                        if input == "y\n" {
                            header = "Replaced "+domain+".."
                            err = db.UpdateDomain(domain, username, password, pub_key)
                        } else {
                            header = "Did not replace "+domain+".."
                        }
                    }
                    check(err)
                }

                domains, err = db.GetDomainList(pub_key)
                check(err)

            case "q":
                done = true

            default:
                if len(choice)>1 && choice[0] == 'd' {
                    if i, err := strconv.Atoi(choice[1:]); err == nil && i > 0 && i <= len(domains) {
                        err := db.DeleteDomain(domains[i-1], pub_key)
                        check(err)
                        header = "Removed: "+domains[i-1]

                        domains, err = db.GetDomainList(pub_key)
                        check(err)

                    } else {
                        header = "Unknown input: "+choice
                    }
                } else if i, err := strconv.Atoi(choice); err == nil && i > 0 && i <= len(domains) {
                    username, password, err := db.ReadDomainUserPass(domains[i-1], pub_key, prv_key, prv_pass)
                    check(err)
                    clear(username)
                    fmt.Print("(PRESS ENTER)")
                    reader.ReadString('\n')
                    clear(password)
                    fmt.Print("(PRESS ENTER)")
                    reader.ReadString('\n')
                } else {
                    header = "Unknown input: "+choice
                }
        }

        clear(header)
    }
}
