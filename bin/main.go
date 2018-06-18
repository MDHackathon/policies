package main

import (
	"fmt"
	"github.com/MDHackathon/policies/pkg"
	"github.com/minio/minio-go"
	"log"
	"strings"
	"time"
)

// This will load the policies from a json configuration
// and apply a choosen operation on all object who match
// the policies rules.
//
//	S3 storage:
//		- In  -> the storage source
//		- Out -> the storage destination
//
// Rules:
// 		- match(name) -> allow a reggex on the key
//		- match(data) -> allow to operate an a set of object define by date
//
// Operation:
//		- Copy -> Copy an object from a s3 storage to an other ( actif <-> passif )
//		- Move -> Moce an object from a s3 storage to an other (  hot  <->  cold  )
//

func main() {

	var (
		S3IN  *minio.Client
		S3OUT *minio.Client
		p     policies.Policies
		t     time.Time
		err   error
	)

	if p, err = policies.LoadPolicieFromPath("./config/policies.json"); err != nil {
		log.Println("failed to laod json ")
		return
	}
	if S3IN, err = policies.GetS3Client(p.InKey, p.InSecret, p.InEndpoint, true); err != nil {
		log.Println("failed to init input client", err)
		return
	}
	if S3OUT, err = policies.GetS3Client(p.OutKey, p.OutSecret, p.OutEndpoint, true); err != nil {
		log.Println("failed to init output client", err)
		return
	}

	doneCh := make(chan struct{})
	defer close(doneCh)

	objectCh := S3IN.ListObjects(p.InBucket, p.InPath, true, doneCh)
	for obj := range objectCh {
		if obj.Err != nil {
			fmt.Println("Failed to list: ", obj.Err)
		} else {
			op := false
			// TODO
			// Match should be more generic hidding the logic here avoiding the next part of code
			if strings.Compare(p.RuleType, "match") == 0 {
				if ret, err := policies.MatchName(obj.Key, p.RuleValue); err != nil {
					fmt.Println("somethiong goes wrong check the policies parameters", err)
				} else {
					op = true
				}
			}
			if strings.Compare(p.RuleType, "date") == 0 {
				if t, err = time.Parse(time.RFC3339, p.RuleValue); err != nil {
					fmt.Println("Failed to parse time, stopping policie", p.RuleValue, err)
					break
				}
				if policies.MatchDate(obj.LastModified, p.RuleCmpOp, t) {
					op = true
				}
			}
			if op == true {
				if err = ExecuteOperation(p, S3IN, S3OUT); err != nil {
					fmt.Println("Failed to execute the op {}")
					break
				}
			}
		}
	}
}
