package policies

import (
	"fmt"
	"github.com/minio/minio-go"
	"log"
	"strings"
	"time"
)

func Execute(data []byte) {
	var (
		S3IN  *minio.Client
		S3OUT *minio.Client
		p     Policies
		t     time.Time
		err   error
	)

	if p, err = LoadPolicieFromByte(data); err != nil {
		log.Println("failed to laod json ")
		return
	}

	if S3IN, err = GetS3Client(p.InKey, p.InSecret, p.InEndpoint, true); err != nil {
		log.Println("failed to init input client", err)
		return
	}
	if S3OUT, err = GetS3Client(p.OutKey, p.OutSecret, p.OutEndpoint, true); err != nil {
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
				if ret, err := MatchName(obj.Key, p.RuleValue); err != nil {
					fmt.Println("somethiong goes wrong check the policies parameters", err)
				} else if ret {
					fmt.Println("match : ", obj.Key)
					op = true
				} else {
					fmt.Println("not match : ", obj.Key)
				}
			}
			if strings.Compare(p.RuleType, "date") == 0 {
				if t, err = time.Parse(time.RFC3339, p.RuleValue); err != nil {
					fmt.Println("Failed to parse time, stopping policie", p.RuleValue, err)
					break
				}
				if MatchDate(obj.LastModified, p.RuleCmpOp, t) {
					op = true
				}
			}
			if op == true {
				if err = ExecuteOperation(p, S3IN, S3OUT, obj); err != nil {
					fmt.Println("Failed to execute the op ", err)
					break
				}
			}
		}
	}
}
