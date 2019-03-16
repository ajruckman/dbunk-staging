// Copyright (c) A.J. Ruckman 2019

package genrules

import (
    "bufio"
    "fmt"
    "net/http"
    "regexp"
    "sort"
    "strings"
    "time"

    "gopkg.in/cheggaaa/pb.v1"

    "github.com/ajruckman/dbunk-staging/internal/common"
    "github.com/ajruckman/dbunk-staging/internal/config"
    "github.com/ajruckman/dbunk-staging/internal/dbunkdb"
    "github.com/ajruckman/dbunk-staging/internal/log"
    "github.com/ajruckman/dbunk-staging/pkg/pqext"
    "github.com/ajruckman/dbunk-staging/pkg/schema"
)

func GenRules() {
    begin := time.Now()

    log.Info("Resetting blacklist")
    reset()

    log.Info("Reading hosts from sources")
    hosts := readSources()

    log.Info("Sorting hosts")
    sortHosts(hosts)

    log.Info("Converting hosts to rules")
    rules := hostsToRules(hosts)

    log.Info("Saving rules")
    saveRules(rules)

    log.Info("Rule generator completed", log.F{
        "Duration": time.Since(begin),
    })
}

func reset() {
    _, err := dbunkdb.DB.DB.Exec(`TRUNCATE TABLE blacklist`)
    log.Error(err)
}

func mapExistingRules() (seen map[string]map[string]bool) {
    seen = map[string]map[string]bool{}

    existing, err := dbunkdb.Blacklist()
    log.Error(err)

    for _, v := range existing {
        seen[v.Source][*v.Host] = true
    }

    return
}

func readSources() (hosts []schema.Blacklist) {
    var (
        seen  = mapExistingRules()
        skips int
    )

    for _, url := range config.Conf.Sources {
        insSource(url)

        resp, err := http.Get(url)
        log.Error(err)

        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            row := matchHost.FindStringSubmatch(scanner.Text())
            if len(row) < 1 {
                continue
            }

            for h := 1; h < len(row); h++ {
                name := row[h]
                name = strings.TrimSuffix(name, `.`)
                name = strings.TrimSuffix(name, `\`)
                name = strings.TrimPrefix(name, `www.`)

                if strings.Count(name, ".") == 0 {
                    goto skip
                }

                if _, ok := seen[url]; !ok {
                    seen[url] = map[string]bool{}
                    goto add
                }

                if _, ok := seen[url][name]; !ok {
                    goto add
                }

            skip:
                skips++
                continue

            add:
                hosts = append(hosts, schema.Blacklist{
                    Host:   &name,
                    Source: url,
                })
                seen[url][name] = true

            }
        }
    }

    log.Info(fmt.Sprintf("Skips: %d", skips))

    return
}

func insSource(source string) {
    _, err := dbunkdb.DB.DB.Exec(`
	  INSERT INTO sources (source, enabled)
	  VALUES ($1, TRUE)
	  ON CONFLICT DO NOTHING`,
        source)
    log.Error(err)
}

func sortHosts(hosts []schema.Blacklist) {
    sort.Slice(hosts, func(i, j int) bool {
        return len(*hosts[i].Host) < len(*hosts[j].Host)
    })
}

func hostsToRules(hosts []schema.Blacklist) (rules []schema.Blacklist) {
    var (
        bar       = pb.StartNew(len(hosts))
        n         schema.Blacklist
        seen      = map[string][]schema.Blacklist{}
        seenFirst bool
        skips     int
    )

    bar.Start()

    for i, v := range hosts {
        bar.Increment()

        key := common.HostToKey(*v.Host)
        sld := common.HostToRawSLD(*v.Host)

        if len(sld) == 0 {
            skips++
            goto next
        } else if len(key) == 0 {
            key = sld
        }

        if _, ok := seen[sld]; ok {
            for _, rule := range seen[sld] {
                if rule.Source != v.Source {
                    continue
                } else if rule.Rule.MatchString(*v.Host) {
                    skips++
                    goto next
                }
            }
        } else {
            seenFirst = true
        }

        n = schema.Blacklist{
            Key:    key,
            Rule:   pqext.Regexp{Regexp: regexp.MustCompile(`(?:^|.*\.)` + strings.Replace(*v.Host, `.`, `\.`, -1) + `$`)},
            Host:   hosts[i].Host,
            Source: v.Source,
        }

        if seenFirst {
            seen[sld] = []schema.Blacklist{n}
        } else {
            seen[sld] = append(seen[sld], n)
        }

        rules = append(rules, n)

    next:
    }

    bar.Finish()
    log.Info(fmt.Sprintf("Skips: %d", skips))

    return
}

func saveRules(rules []schema.Blacklist) {
    bar := pb.New(len(rules))
    bar.Start()

    for _, v := range rules {
        bar.Increment()
        //fmt.Printf("%-10d %-150s %s\n", i, v.Rule.String(), v.Source)

        _, err := dbunkdb.DB.DB.Exec(`        
          INSERT INTO blacklist (key, rule, host, source)
          VALUES ($1, $2, $3, $4)`,
            v.Key,
            v.Rule.String(),
            v.Host,
            v.Source)
        log.Error(err)
    }

    bar.Finish()
}

var (
    matchHost = regexp.MustCompile(`^(?:0\.0\.0\.0|127\.0\.0\.1|::|::0|::1)\s+([^\s]+).*`)
)
