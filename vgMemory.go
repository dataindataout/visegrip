// visegrip memory checker
// Does the OOM killer crash MySQL? Bet your per-thread memory is set too high. Let's find out!

// Valerie Parham-Thompson for sharing 2016
// keep it free and open

// this is also my go hello world :)

package main

import (
"database/sql"
"fmt"
"log"
"strings"
_ "github.com/go-sql-driver/mysql"
)

func main() {

// connect to the mysql instance

    db, err := sql.Open("mysql", "valerie:password@tcp(192.168.56.68:3306)/")

    if err != nil {
        log.Fatal(err)
    }

// get conns variable/status

    //configured max connections
    var varMaxConn string
    err = db.QueryRow("select variable_value from information_schema.global_variables where variable_name='MAX_CONNECTIONS'").Scan(&varMaxConn)
    if err != nil {
        log.Fatal(err)
    }

    //observed max connections
    var statMaxConn string
    err = db.QueryRow("select variable_value from information_schema.global_status where variable_name like 'MAX_USED_CONNECTIONS'").Scan(&statMaxConn)
    if err != nil {
        log.Fatal(err)
    }

// do you want to use max conns, max used conns, or something else?

    //TODO: add all current memory buckets

    var chosenConnectionThreshold float32

    //reader := bufio.NewReader(os.Stdin)
    s := []string{"Would you like to calculate memory on your configured max connections (", varMaxConn, "), on the max connections seen on this server (", statMaxConn, "), or some other number? Enter your desired number of connections here:", "\n>>>"}

    fmt.Printf(strings.Join(s, ""))

    var i int
    fmt.Scan(&i)
    chosenConnectionThreshold = float32(i)
    //TODO: handle NaN

// get global memory

    fmt.Println("\nGlobal memory settings (in MB)\n==============================")

    var globalName string
    var globalValue float32
    var globalValueMB float32
    var totalGlobalValue float32
    var totalGlobalValueMB float32

    //TODO: add all current memory buckets
    //and handle errors for any missing ones amongst versions

    rows, err := db.Query("select variable_name, variable_value from information_schema.global_variables where variable_name in ('KEY_BUFFER_SIZE', 'QUERY_CACHE_SIZE', 'TMP_TABLE_SIZE', 'INNODB_BUFFER_POOL_SIZE', 'INNODB_ADDITIONAL_MEM_POOL_SIZE', 'INNODB_LOG_BUFFER_SIZE')")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(&globalName, &globalValue)
        if err != nil {
            log.Fatal(err)
        }
        globalValueMB = globalValue / 1024 / 1024 //convert to MB
        totalGlobalValue += globalValue //running global total
        fmt.Println(globalName, globalValueMB)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    totalGlobalValueMB = totalGlobalValue / 1024 / 1024
    fmt.Println("Total Global =", totalGlobalValueMB)

// get per-thread memory

    fmt.Println("\nPer-thread memory settings (in MB)\n==================================")

    var perThreadName string
    var perThreadValue float32
    var perThreadValueMB float32
    var totalPerThreadValue float32
    var totalPerThreadValueMB float32

    //TODO: add all current memory buckets
    //and handle errors for any missing ones amongst versions

    rows, err = db.Query("select variable_name, variable_value from information_schema.global_variables where variable_name in ('SORT_BUFFER_SIZE', 'READ_BUFFER_SIZE', 'READ_RND_BUFFER_SIZE', 'JOIN_BUFFER_SIZE', 'THREAD_STACK', 'BINLOG_CACHE_SIZE')")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&perThreadName, &perThreadValue)
        if err != nil {
            log.Fatal(err)
        }
        perThreadValueMB = perThreadValue / 1024 / 1024 //convert to MB
        totalPerThreadValue += perThreadValue //running per-thread total
        fmt.Println(perThreadName, perThreadValueMB)
    }
    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    totalPerThreadValueMB = totalPerThreadValue / 1024 / 1024
    fmt.Println("Total Per-Thread =", totalPerThreadValueMB)

// print total memory usage

    var totalMemoryUsage float32
    totalMemoryUsage = (totalGlobalValue + (chosenConnectionThreshold * totalPerThreadValue) ) / 1024 / 1024

    fmt.Println("\nTotal Memory Usage at", chosenConnectionThreshold, "Connections (in MB)\n==============================================")

    fmt.Println("At", chosenConnectionThreshold, "connections, your memory usage will be", totalMemoryUsage, "MB.")

    //TODO: print your available memory for comparison

    defer db.Close()

} //end main
