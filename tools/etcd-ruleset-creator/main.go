package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/heetch/rules-engine/rule"
)

// CREATE YOUR RULESET IN THIS FUNCTION BODY
func getRuleset() (*rule.Ruleset, error) {
	// EXAMPLE:
	return rule.NewStringRuleset(
		rule.New(
			rule.Eq(
				rule.StringParam("product-id"),
				rule.StringValue("fr-paris"),
			),
			rule.ReturnsString("matched"),
		),
	)
}

// NO NEED TO TOUCH THIS
func main() {
	addr := flag.String("addr", "", "etcd addr")
	namespace := flag.String("namespace", "", "prefix to use for namespacing")
	name := flag.String("name", "", "name of the ruleset")
	flag.Parse()

	if *addr == "" || *namespace == "" || *name == "" {
		flag.Usage()
		return
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	rs, err := getRuleset()
	if err != nil {
		log.Fatal(err)
	}

	raw, err := json.Marshal(rs)
	if err != nil {
		log.Fatal(err)
	}

	keyPrefix := path.Join(strings.TrimLeft(*namespace, "/"), "/")
	rsName := strings.Trim(*name, "/")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Put(ctx, path.Join(keyPrefix, rsName), string(raw))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Ruleset \"%s\" successfully saved.\n", rsName)
}