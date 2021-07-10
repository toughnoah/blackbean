package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/toughnoah/blackbean/pkg/es"
	"io"
	"log"
)

const (
	AllRoles         = "all"
	IndexPrivilege   = "index"
	ClusterPrivilege = "cluster"
)

func role(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		command = &cobra.Command{
			Use:               "role [subcommand]",
			Short:             "role operations for cluster",
			Long:              "role operations for cluster ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
		}
	)
	command.AddCommand(getRole(cli, out))
	command.AddCommand(createRole(cli, out))
	command.AddCommand(updateRole(cli, out))
	command.AddCommand(deleteRole(cli, out))
	return command
}

func getRole(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		r       = &Role{Client: cli}
		command = &cobra.Command{
			Use:   "get [role]",
			Short: "get specify role",
			Long:  "get specify role ... wordless",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return r.getAllRoles(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					res *esapi.Response
					err error
				)
				if len(args) == 0 {
					res, err = r.getRoles(AllRoles)
				} else {
					res, err = r.getRoles(args[0])
				}
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	return command
}

func createRole(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		r   = &Role{Client: cli}
		req = new(es.RequestBody)
		i   = Indices{
			client: cli,
		}
		command = &cobra.Command{
			Use:               "create [role]",
			Short:             "create specify user",
			Long:              "create specify user ... wordless",
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: noCompletions,
			RunE: func(cmd *cobra.Command, args []string) error {
				if es.GetFlagValue(cmd, "cluster_privilege") == "" &&
					es.GetFlagValue(cmd, "indices") == "" {
					return errors.New("at least one of flags cluster_privilege and indices should be specified")
				}
				r.Role = args[0]
				res, err := r.createRole(req)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVar(&r.ClusterPrivilege, "cluster_privilege", "", "specify privilege to be assigned, use ',' to split multi privileges.")
	f.StringVar(&r.Indices, "indices", "", "specify indices wanted, use ',' to split multi indices.")
	f.StringVar(&r.IndexPrivilege, "indices_privilege", "read", "specify indices indices_privilege, use ',' to split multi privileges.")
	if err := command.RegisterFlagCompletionFunc("cluster_privilege", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return r.getAllPrivilege(ClusterPrivilege), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	if err := command.RegisterFlagCompletionFunc("indices_privilege", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return r.getAllPrivilege(IndexPrivilege), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	if err := command.RegisterFlagCompletionFunc("indices", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	es.AddRequestBodyFlag(command, req)
	return command
}

func updateRole(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		r   = &Role{Client: cli}
		req = new(es.RequestBody)
		i   = Indices{
			client: cli,
		}
		command = &cobra.Command{
			Use:   "update [role]",
			Short: "update specify user",
			Long:  "update specify user ... wordless",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return r.getAllRoles(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				if es.GetFlagValue(cmd, "cluster_privilege") == "" &&
					es.GetFlagValue(cmd, "indices") == "" {
					return errors.New("at least one of flags cluster_privilege, and indices should be specified")
				}
				r.Role = args[0]
				res, err := r.updateRole(req)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVar(&r.ClusterPrivilege, "cluster_privilege", "", "specify privilege to be assigned, use ',' to split multi privileges.")
	f.StringVar(&r.Indices, "indices", "", "specify indices wanted, use ',' to split multi indices.")
	f.StringVar(&r.IndexPrivilege, "indices_privilege", "read", "specify indices indices_privilege, use ',' to split multi privileges.")
	f.BoolVar(&r.AddOnly, "add_only", true, "add only, not to override")
	if err := command.RegisterFlagCompletionFunc("cluster_privilege", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return r.getAllPrivilege(ClusterPrivilege), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	if err := command.RegisterFlagCompletionFunc("indices_privilege", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return r.getAllPrivilege(IndexPrivilege), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	if err := command.RegisterFlagCompletionFunc("indices", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return i.getAllIndices(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}); err != nil {
		log.Fatal(err)
	}
	es.AddRequestBodyFlag(command, req)
	return command
}

func deleteRole(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		r       = &Role{Client: cli}
		command = &cobra.Command{
			Use:   "delete [role]",
			Short: "delete specify role",
			Long:  "delete specify role ... wordless",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return r.getAllRoles(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := r.deleteRole(args[0])
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	return command
}

type Role struct {
	Client           *elasticsearch.Client
	Role             string
	ClusterPrivilege string
	IndexPrivilege   string
	Indices          string
	AddOnly          bool
}

func (r *Role) getRoles(roles string) (*esapi.Response, error) {
	if roles == AllRoles {
		return r.Client.Security.GetRole(r.Client.Security.GetRole.WithPretty())
	}
	return r.Client.Security.GetRole(r.Client.Security.GetRole.WithName(splitWords(roles)...), r.Client.Security.GetRole.WithPretty())
}

func (r *Role) getAllRoles() []string {
	var (
		resMap   map[string]interface{}
		resSlice []string
	)
	ret, err := r.Client.Security.GetRole()
	if err != nil {
		return nil
	}
	if err = json.NewDecoder(ret.Body).Decode(&resMap); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	for i, _ := range resMap {
		resSlice = append(resSlice, i)
	}
	return resSlice
}

func (r *Role) getAllPrivilege(PrivilegeType string) []string {
	var (
		resMap map[string][]string
	)
	ret, err := r.Client.Security.GetBuiltinPrivileges()
	if err != nil {
		return nil
	}
	if err = json.NewDecoder(ret.Body).Decode(&resMap); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	return resMap[PrivilegeType]
}

func (r *Role) createRole(req *es.RequestBody) (*esapi.Response, error) {
	rawBody, err := es.GetRawRequestBody(req)
	if err != nil {
		return nil, err
	}
	if rawBody != nil {
		return r.Client.Security.PutRole(r.Role, bytes.NewReader(rawBody))
	}
	roleReqBody := &RoleRequestBody{}
	if r.ClusterPrivilege != "" {
		roleReqBody.Cluster = splitWords(r.ClusterPrivilege)
	}
	if r.Indices != "" {
		roleReqBody.RoleIndices = append(roleReqBody.RoleIndices, map[string]interface{}{
			"names":      splitWords(r.Indices),
			"privileges": splitWords(r.IndexPrivilege),
		})
	}
	bytesBody, err := json.Marshal(roleReqBody)
	if err != nil {
		return nil, err
	}
	return r.Client.Security.PutRole(r.Role, bytes.NewReader(bytesBody))
}

type RoleRequestBody struct {
	Cluster     []string                 `json:"cluster,omitempty"`
	RoleIndices []map[string]interface{} `json:"indices,omitempty"`
	MetaData    map[string]int           `json:"meta_data,omitempty"`
}

func (r *Role) updateRole(req *es.RequestBody) (*esapi.Response, error) {
	rawBody, err := es.GetRawRequestBody(req)
	if err != nil {
		return nil, err
	}
	if rawBody != nil {
		return r.Client.Security.PutRole(r.Role, bytes.NewReader(rawBody))
	}
	originRole, err := r.getRole()
	if err != nil {
		return nil, err
	}
	roleReqBody := &RoleRequestBody{}
	if r.ClusterPrivilege != "" && r.AddOnly {
		hasClusterPrivilege := originRole[r.Role]["cluster"].([]interface{})
		mergedRoles := mergeRoles(hasClusterPrivilege, ConvertStrSliceToInterface(splitWords(r.ClusterPrivilege)))
		roleReqBody.Cluster = ConvertInterfaceSliceToStr(mergedRoles)
	} else if r.ClusterPrivilege != "" && !r.AddOnly {
		roleReqBody.Cluster = splitWords(r.ClusterPrivilege)
	}
	if r.Indices != "" {
		roleReqBody.RoleIndices = append(roleReqBody.RoleIndices, map[string]interface{}{
			"names":      splitWords(r.Indices),
			"privileges": splitWords(r.IndexPrivilege),
		})
	}
	if r.AddOnly {
		hasIndicesPrivilege := ConvertInterfaceSliceToMapStr(originRole[r.Role]["indices"].([]interface{}))
		roleReqBody.RoleIndices = append(roleReqBody.RoleIndices, hasIndicesPrivilege...)
	}
	bytesBody, err := json.Marshal(roleReqBody)
	if err != nil {
		return nil, err
	}
	return r.Client.Security.PutRole(r.Role, bytes.NewReader(bytesBody))
}

func (r *Role) getRole() (map[string]map[string]interface{}, error) {
	var (
		resMap map[string]map[string]interface{}
	)
	ret, err := r.Client.Security.GetRole(r.Client.Security.GetRole.WithName(r.Role))
	if err != nil {
		return nil, err
	}
	if err = json.NewDecoder(ret.Body).Decode(&resMap); err != nil {
		return nil, errors.Errorf("Error parsing the response body: %s", err)
	}
	if len(resMap) == 0 {
		return nil, errors.Errorf("cluster has no such role %s", r.Role)
	}
	return resMap, nil

}

func (r *Role) deleteRole(role string) (*esapi.Response, error) {
	return r.Client.Security.DeleteRole(role)

}

func ConvertStrSliceToInterface(s []string) []interface{} {
	i := make([]interface{}, len(s))
	for ss, v := range s {
		i[ss] = v
	}
	return i
}
func ConvertInterfaceSliceToStr(i []interface{}) []string {
	s := make([]string, len(i))
	for ii, v := range i {
		s[ii] = v.(string)
	}
	return s
}
func ConvertInterfaceSliceToMapStr(i []interface{}) []map[string]interface{} {
	m := make([]map[string]interface{}, len(i))
	for ii, v := range i {
		m[ii] = v.(map[string]interface{})
	}
	return m
}
