package response

type LaunchPlan struct {
}

type AllInstance struct {
	Result []Instance `json:"list"`
}

type AllInstanceNames struct {
	Result []string `json:"result"`
}

type List struct {
}

type Instance struct {
	Domain                 string        `json:"domain"`
	Placement_requirements []string      `json:"placement_requirements"`
	Ip                     string        `json:"ip"`
	Site                   string        `json:"site,omitempty"`
	Shape                  string        `json:"shape"`
	Imagelist              interface{}   `json:"imagelist,omitempty"` // null problem
	Image_format           string        `json:"image_format"`
	Relationships          interface{}   `json:"relationships,omitempty"` //array of what?
	Availability_domain    string        `json:"availability_domain"`
	Networking             Networking    `json:"networking"`
	Storage_attachments    Storage       `json:"storage_attachments"`
	Hostname               string        `json:"hostname"`
	Quota_reservation      interface{}   `json:"quota_reservation,omitempty"` //null problem
	Disk_attach            string        `json:"disk_attach,omitempty"`
	Label                  string        `json:"label"`
	Priority               string        `json:"priority"`
	Platform               string        `json:"platform"`
	State                  string        `json:"state"`
	Virtio                 interface{}   `json:"virtio,omitempty"` // null problem
	Vnc                    string        `json:"vnc"`
	Desired_state          string        `json:"desired_state"`
	Tags                   []string      `json:"tags,omitempty"`
	Start_time             string        `json:"start_time"`
	Quota                  string        `json:"quota"`
	Entry                  interface{}   `json:"entry,omitempty"` // null problem
	Error_reason           string        `json:"error_reason"`
	SSHKeys                []AllSSHNames `json:"sshkeys"`
	Resolvers              interface{}   `json:"resolvers,omitempty"` // null problem
	Account                string        `json:"account"`
	Name                   string        `json:"name"`
	Vcanble_id             string        `json:"vcable_id"`
	Hypervisor             Hypervisor    `json:"hypervisor"`
	Uri                    string        `json:"uri"`
	Reverse_dns            bool          `json:"reverse_dns"`
	Attributes             Attributes    `json:"attributes"`
	Boot_order             []int         `json:"boot_order"`
}

type Networking struct {
	Eth0 Nic `json:"eth0"`
}

type Nic struct {
	Model     string   `json:"model"`
	Seclists  []string `json:"seclists"`
	Dns       []string `json:"dns"`
	Vethernet string   `json:"vethernet"`
	Nat       string   `json:"nat"`
}

type Storage struct {
	Index               string `json:"index"`
	Storage_volume_name string `json:"storage_volume_name"`
	Name                string `json:"name"`
}

type Hypervisor struct {
	Mode string `json:"mode"`
}

type Attributes struct {
	SSHKeys []string `json:"sshkeys"`
	Network Network  `json:"network"`
	Dns     Dns      `json:"dns"`
}

type Dns struct {
	Domain      string `json:"domain"`
	Hostname    string `json:"hostname"`
	Vcable_eth0 string `json:"nimbula_vcable-eth0"`
}

type Network struct {
	Vcable_eth0    Vcable   `json:"nimbula_vcable-eth0"`
	Model          string   `json:"model,omitempty"`
	Vethernet_type string   `json:"vethernet_type"`
	Id             string   `json:"id"`
	Dhcp_options   []string `json"dhcp_options,omitempty"`
}

type Vcable struct {
	Vethernet_id string   `json:"vethernet_id"`
	Address      []string `json:"address"`
}
