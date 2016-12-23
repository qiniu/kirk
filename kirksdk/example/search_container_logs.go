package example

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"qiniupkg.com/kirk/kirksdk"
)

func SearchContainerLogsExample() {
	client := kirksdk.NewQcosClient(kirksdk.QcosConfig{
		AccessKey: "test_fake_ak",
		SecretKey: "test_fake_sk",
		Host:      "test_fake_host",
	})

	err := SearchContainerLogs(client, "access", "", time.Now().Add(time.Hour*time.Duration(-72)), time.Now(), true, 25000, os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
}

func SearchContainerLogs(client kirksdk.QcosClient, repoType, queryString string, fromTime, toTime time.Time, tail bool, size int, out io.Writer) (err error) {

	args := kirksdk.SearchContainerLogsArgs{
		RepoType: repoType,
	}

	var query string
	if tail {
		query = addClause(query, fmt.Sprintf("collectedAtNano:[%d TO %%d]", fromTime.UnixNano()))
	} else {
		query = addClause(query, fmt.Sprintf("collectedAtNano:[%%d TO %d]", toTime.UnixNano()))
	}

	if len(queryString) > 0 {
		// we will use fmt.Sprintf to format the query string later,
		// so the raw % should be replaced into %% to prevent BADPREC MISSING error from fmt.Sprintf
		query = addClause(query, strings.NewReplacer("%", "%%").Replace(queryString))
	}

	if tail {
		args.Sort = "collectedAtNano:desc"
		args.Query = fmt.Sprintf(query, toTime.UnixNano())
	} else {
		args.Sort = "collectedAtNano:asc"
		args.Query = fmt.Sprintf(query, fromTime.UnixNano())
	}

	expectedSize := size
	maxSizePerRequest := 10000
	if expectedSize > maxSizePerRequest {
		args.Size = maxSizePerRequest
	} else {
		args.Size = expectedSize
	}

	result, err := client.SearchContainerLogs(context.TODO(), args)
	if err != nil {
		return err
	}

	total := result.Total
	logs := make([]string, 0)
	logs = append(logs, formatLogs(result.Data, repoType)...)

	// We cannot get logs more than total.
	if total < expectedSize {
		expectedSize = total
	}

	for {
		expectedSize -= len(result.Data)
		if expectedSize <= 0 {
			break
		}

		if expectedSize > maxSizePerRequest {
			args.Size = maxSizePerRequest
		} else {
			args.Size = expectedSize
		}
		if tail {
			args.Query = fmt.Sprintf(query, result.Data[len(result.Data)-1].CollectedAtNano-1)
		} else {
			args.Query = fmt.Sprintf(query, result.Data[len(result.Data)-1].CollectedAtNano+1)
		}
		result, err = client.SearchContainerLogs(context.TODO(), args)
		if err != nil {
			return err
		}
		if len(result.Data) <= 0 {
			break
		}
		logs = append(logs, formatLogs(result.Data, repoType)...)
	}

	count := len(logs)

	if err != nil {
		fmt.Printf("search container logs err: %s\n", err)
	} else {
		if tail {
			for i := len(logs) - 1; i >= 0; i-- {
				fmt.Fprintln(out, logs[i])
			}
		} else {
			for _, log := range logs {
				fmt.Fprintln(out, log)
			}
		}
		fmt.Fprintf(out, "\nSummary: %d/%d\n", count, total)
	}

	return
}

func formatLogs(hits []kirksdk.Hit, repoType string) (logs []string) {
	var log string
	logs = make([]string, 0)
	for _, hit := range hits {
		if repoType == "access" {
			log = formatAccessLogField(hit)
		} else {
			log = hit.Log
		}
		logs = append(logs, log)
	}
	return
}

func formatAccessLogField(hit kirksdk.Hit) string {
	if hit.Method == "" {
		hit.Method = "-"
	}
	if hit.Url == "" {
		hit.Url = "-"
	}
	if hit.RequestHeader == "" {
		hit.RequestHeader = "-"
	}
	if hit.RequestParams == "" {
		hit.RequestParams = "-"
	}
	if hit.RequestBody == "" {
		hit.RequestBody = "-"
	}
	if hit.ResponseHeader == "" {
		hit.ResponseHeader = "-"
	}
	if hit.ResponseBody == "" {
		hit.ResponseBody = "-"
	}
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%s\t%s\t%gms", hit.StartAt.Format("2006-01-02 15:04:05.000000Z07:00"), hit.Method, hit.Url, hit.RequestHeader, hit.RequestParams, hit.RequestBody, hit.StatusCode, hit.ResponseHeader, hit.ResponseBody, float64(hit.ElapsedNano)/1e6)
}

func addClause(q string, clause string) string {
	if q == "" {
		return clause
	}
	if clause == "" {
		return q
	}
	return fmt.Sprintf("%s AND %s", q, clause)
}
