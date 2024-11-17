package osrelease

import (
	"fmt"
	"strings"
)

type Info struct {
	// identification
	Name       string `json:"name,omitempty" field:"NAME"`
	ID         string `json:"id" field:"ID"` // required
	CPEName    string `json:"cpe_name,omitempty" field:"CPE_NAME"`
	PrettyName string `json:"pretty_name,omitempty" field:"PRETTY_NAME"`

	// variant identification
	IDLike        []string `json:"id_like,omitempty" field:"ID_LIKE"`
	VariantID     string   `json:"variant_id,omitempty" field:"VARIANT_ID"`
	OstreeVersion string   `json:"ostree_version,omitempty" field:"OSTREE_VERSION"`

	// links
	HomeURL          string `json:"home_url,omitempty" field:"HOME_URL"`
	SupportURL       string `json:"support_url,omitempty" field:"SUPPORT_URL"`
	BugReportURL     string `json:"bug_report_url,omitempty" field:"BUG_REPORT_URL"`
	DocumentationURL string `json:"documentation_url,omitempty" field:"DOCUMENTATION_URL"`
	PrivacyPolicyURL string `json:"privacy_policy_url,omitempty" field:"PRIVACY_POLICY_URL"`

	// lifecycle related
	SupportEnd   string `json:"support_end,omitempty" field:"SUPPORT_END"`
	Discontinued bool   `json:"discontinued,omitempty"`

	// version info
	VersionID    string `json:"version_id" field:"VERSION_ID"` // required
	Version      string `json:"version,omitempty" field:"VERSION"`
	MajorVersion string `json:"major_version,omitempty"` // calculated
	MinorVersion string `json:"minor_version,omitempty"` // calculated
	BuildID      string `json:"build_id,omitempty" field:"BUILD_ID"`

	// codename support
	VersionCodename string `json:"version_codename,omitempty" field:"VERSION_CODENAME"`
	DebianCodename  string `json:"debian_codename,omitempty" field:"DEBIAN_CODENAME"`
	UbuntuCodename  string `json:"ubuntu_codename,omitempty" field:"UBUNTU_CODENAME"`
}

func (i Info) String() string {
	var name string
	if i.PrettyName == "" {
		name = fmt.Sprintf("%q @ %s", i.Name, i.VersionID)
	} else {
		name = fmt.Sprintf("%q", i.PrettyName)
	}

	if !strings.Contains(name, i.VersionID) {
		name = fmt.Sprintf("%s @ %s", name, i.VersionID)
	}

	var postFix []string

	if i.VariantID != "" { // && !strings.Contains(strings.ToLower(name), strings.ToLower(i.VariantID)) {
		postFix = append(postFix, fmt.Sprintf("[variant=%s]", i.VariantID))
	}

	if len(i.IDLike) > 0 {
		postFix = append(postFix, fmt.Sprintf("[like=%s]", strings.Join(i.IDLike, ",")))
	}

	postFix = append(postFix, fmt.Sprintf("[id=%s]", i.ID))

	return fmt.Sprintf("%s %s", name, strings.Join(postFix, " "))
}
