package diskscript

import (
	"context"
	"errors"
	"fmt"
	clientruntime "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/alert"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/alertmanager/pkg/labels"
	promconfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/version"
	"golang.org/x/mod/semver"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/url"
	"os"
	"path"
	"time"
)

var (
	httpConfigFile string
)

const (
	defaultAmHost      = "localhost"
	defaultAmPort      = "9093"
	defaultAmApiv2path = "/api/v2"
)

type AlertAddCmd struct {
	annotations  []string
	generatorURL string
	labels       []string
	start        string
	end          string
}

func parseLabels(inputLabels []string) (models.LabelSet, error) {
	labelSet := make(models.LabelSet, len(inputLabels))

	for _, l := range inputLabels {
		matcher, err := labels.ParseMatcher(l)
		if err != nil {
			return models.LabelSet{}, err
		}
		if matcher.Type != labels.MatchEqual {
			return models.LabelSet{}, errors.New("labels must be specified as key=value pairs")
		}

		labelSet[matcher.Name] = matcher.Value
	}

	return labelSet, nil
}

func (a *AlertAddCmd) addAlert(ctx context.Context) error {
	if len(a.labels) > 0 {
		if _, err := parseLabels([]string{a.labels[0]}); err != nil {
			a.labels[0] = fmt.Sprintf("alertname=%s", a.labels[0])
		}
	}

	labels, err := parseLabels(a.labels)
	if err != nil {
		return err
	}

	annotations, err := parseLabels(a.annotations)
	if err != nil {
		return err
	}

	var startsAt, endsAt time.Time
	if a.start != "" {
		startsAt, err = time.Parse(time.RFC3339, a.start)
		if err != nil {
			return err
		}
	}
	if a.end != "" {
		endsAt, err = time.Parse(time.RFC3339, a.end)
		if err != nil {
			return err
		}
	}

	pa := &models.PostableAlert{
		Alert: models.Alert{
			GeneratorURL: strfmt.URI(a.generatorURL),
			Labels:       labels,
		},
		Annotations: annotations,
		StartsAt:    strfmt.DateTime(startsAt),
		EndsAt:      strfmt.DateTime(endsAt),
	}
	alertParams := alert.NewPostAlertsParams().WithContext(ctx).
		WithAlerts(models.PostableAlerts{pa})

	amclient := NewAlertmanagerClient(AlertmanagerURL)
	_, err = amclient.Alert.PostAlerts(alertParams)
	return err
}

func NewAlertmanagerClient(amURL *url.URL) *client.AlertmanagerAPI {
	address := defaultAmHost + ":" + defaultAmPort
	schemes := []string{"http"}

	if amURL.Host != "" {
		address = amURL.Host
	}
	if amURL.Scheme != "" {
		schemes = []string{amURL.Scheme}
	}

	cr := clientruntime.New(address, path.Join(amURL.Path, defaultAmApiv2path), schemes)

	if amURL.User != nil && httpConfigFile != "" {
		kingpin.Fatalf("basic authentication and http.config.file are mutually exclusive")
	}

	if amURL.User != nil {
		password, _ := amURL.User.Password()
		cr.DefaultAuthentication = clientruntime.BasicAuth(amURL.User.Username(), password)
	}

	if httpConfigFile != "" {
		var err error
		httpConfig, _, err := promconfig.LoadHTTPConfigFile(httpConfigFile)
		if err != nil {
			kingpin.Fatalf("failed to load HTTP config file: %v", err)
		}

		httpclient, err := promconfig.NewClientFromConfig(*httpConfig, "amtool")
		if err != nil {
			kingpin.Fatalf("failed to create a new HTTP client: %v", err)
		}
		cr = clientruntime.NewWithClient(address, path.Join(amURL.Path, defaultAmApiv2path), schemes, httpclient)
	}

	c := client.New(cr, strfmt.Default)

	//if !versionCheck {
	//	return c
	//}

	status, err := c.General.GetStatus(nil)
	if err != nil || status.Payload.VersionInfo == nil || version.Version == "" {
		// We can not get version info, or we do not know our own version. Let amtool continue.
		return c
	}

	if semver.MajorMinor("v"+*status.Payload.VersionInfo.Version) != semver.MajorMinor("v"+version.Version) {
		fmt.Fprintf(os.Stderr, "Warning: amtool version (%s) and alertmanager version (%s) are different.\n", version.Version, *status.Payload.VersionInfo.Version)
	}

	return c
}

func (c *Mount) AlertAdd(ctx context.Context, PercentageOfOccupancy int64) bool {
	summary := fmt.Sprintf("summary=磁盘挂载:%s超过限制：%d%%", c.Mount, c.Threshold)
	description := fmt.Sprintf("description=磁盘挂载:%s超过限制：%d%% 当前占用:%d%%", c.Mount, c.Threshold, PercentageOfOccupancy)
	alertname := fmt.Sprintf("alertname=磁盘挂载:%s超过限制：%d%%", c.Mount, c.Threshold)
	labels := make([]string, 0)
	labels = append(labels, fmt.Sprintf("%s", alertname), fmt.Sprintf("intance=%s", Hostname))
	labels = append(labels, c.Labels...)
	labels = append(labels, Conf.Alert.Labels...)
	if c.Alertname != "" {
		alertname = fmt.Sprintf("alertname=%s", c.Alertname)
	}
	alertcli := AlertAddCmd{
		annotations:  []string{summary, description},
		generatorURL: "",
		labels:       labels,
		start:        "",
		end:          "",
	}
	err := alertcli.addAlert(ctx)
	if err != nil {
		log.Printf("%s err %s", alertname, err)
		return false
	}
	log.Printf("%s success", alertname)
	return true
}

func (c *Directory) AlertAdd(ctx context.Context, filesizeformat string) bool {
	summary := fmt.Sprintf("summary=目录:%s超过限制：%s", c.Directory, c.Threshold)
	description := fmt.Sprintf("description=目录:%s超过限制：%s 当前占用:%s", c.Directory, c.Threshold, filesizeformat)
	alertname := fmt.Sprintf("alertname=目录:%s超过限制：%s", c.Directory, c.Threshold)
	labels := make([]string, 0)
	labels = append(labels, fmt.Sprintf("%s", alertname), fmt.Sprintf("intance=%s", Hostname))
	labels = append(labels, c.Labels...)
	labels = append(labels, Conf.Alert.Labels...)
	if c.Alertname != "" {
		alertname = fmt.Sprintf("alertname=%s", c.Alertname)
	}
	alertcli := AlertAddCmd{
		annotations:  []string{summary, description},
		generatorURL: "",
		labels:       labels,
		start:        "",
		end:          "",
	}
	err := alertcli.addAlert(ctx)
	if err != nil {
		log.Printf("%s err %s", alertname, err)
		return false
	}
	log.Printf("%s success", alertname)
	return true
}

func (c *File) AlertAdd(ctx context.Context, filesizeformat string) bool {
	summary := fmt.Sprintf("summary=文件:%s超过限制：%s", c.File, c.Threshold)
	description := fmt.Sprintf("description=文件:%s超过限制：%s 当前占用:%s", c.File, c.Threshold, filesizeformat)
	alertname := fmt.Sprintf("alertname=文件:%s超过限制：%s", c.File, c.Threshold)
	labels := make([]string, 0)
	labels = append(labels, fmt.Sprintf("%s", alertname), fmt.Sprintf("intance=%s", Hostname))
	labels = append(labels, c.Labels...)
	labels = append(labels, Conf.Alert.Labels...)
	if c.Alertname != "" {
		alertname = fmt.Sprintf("alertname=%s", c.Alertname)
	}
	alertcli := AlertAddCmd{
		annotations:  []string{summary, description},
		generatorURL: "",
		labels:       labels,
		start:        "",
		end:          "",
	}
	err := alertcli.addAlert(ctx)
	if err != nil {
		log.Printf("%s 发送失败 %s", alertname, err)
		return false
	}
	log.Printf("%s 发送成功", alertname)
	return true
}
