package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/toughnoah/blackbean/pkg/es"
	"golang.org/x/term"
	"io"
	"strings"
)

const (
	PassRolesCheck = "ignoreEmptyRoles"
	AllUser        = "all"
)

func user(cli *elasticsearch.Client, out io.Writer, in io.ReadWriter, fd int) *cobra.Command {
	var (
		command = &cobra.Command{
			Use:               "user [subcommand]",
			Short:             "user for cluster",
			Long:              "user for cluster ... wordless",
			Args:              cobra.NoArgs,
			ValidArgsFunction: noCompletions,
		}
	)
	command.AddCommand(getUser(cli, out))
	command.AddCommand(createUser(cli, out, in, fd))
	command.AddCommand(updateUser(cli, out, in, fd))
	command.AddCommand(deleteUser(cli, out))
	return command
}

func getUser(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		u       = &User{Client: cli}
		command = &cobra.Command{
			Use:   "get [user]",
			Short: "get specify user",
			Long:  "get specify user ... wordless",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return u.getAllUser(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					res *esapi.Response
					err error
				)
				if len(args) == 0 {
					res, err = u.getUser(AllRoles)
				} else {
					res, err = u.getUser(args[0])
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

func createUser(cli *elasticsearch.Client, out io.Writer, in io.ReadWriter, fd int) *cobra.Command {
	var (
		u = &User{
			Client: cli,
			In:     in,
			Fd:     fd,
		}
		req     = &es.RequestBody{}
		command = &cobra.Command{
			Use:               "create [user]",
			Short:             "create specify user",
			Long:              "create specify user ... wordless",
			Args:              cobra.ExactArgs(1),
			ValidArgsFunction: noCompletions,
			RunE: func(cmd *cobra.Command, args []string) error {
				u.Username = args[0]
				res, err := u.createUser(req)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVarP(&u.Roles, "roles", "", "", "specify roles to be assigned, use ',' to split multi roles.")
	f.StringVarP(&u.Email, "email", "", "", "specify email to be assigned")
	f.StringVarP(&u.FullName, "full_name", "", "", "specify full_name to be assigned")
	_ = command.MarkFlagRequired("roles")
	es.AddRequestBodyFlag(command, req)
	return command
}

func updateUser(cli *elasticsearch.Client, out io.Writer, in io.ReadWriter, fd int) *cobra.Command {
	var (
		u = &User{
			Client: cli,
			In:     in,
			Fd:     fd,
		}
		req     = &es.RequestBody{}
		command = &cobra.Command{
			Use:   "update [user]",
			Short: "update specify user",
			Long:  "update specify user ... wordless",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return u.getAllUser(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				if es.GetFlagValue(cmd, "roles") == "none" && es.GetFlagValue(cmd, "change_password") == "false" {
					return errors.New("at least to specify one of 'roles' and 'change_password' flag")
				}
				u.Username = args[0]
				res, err := u.updateUser(req)
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	f := command.Flags()
	f.StringVarP(&u.Roles, "roles", "", PassRolesCheck, "specify roles to be assigned, use ',' to split multi roles.")
	f.StringVarP(&u.Email, "email", "", "", "specify email to be assigned")
	f.StringVarP(&u.FullName, "full_name", "", "", "specify full_name to be assigned")
	f.BoolVar(&u.ChangePassword, "change_password", false, "to change password")
	f.BoolVar(&u.AddOnly, "add_only", true, "add only, not to override")
	es.AddRequestBodyFlag(command, req)
	return command
}

func deleteUser(cli *elasticsearch.Client, out io.Writer) *cobra.Command {
	var (
		u       = &User{Client: cli}
		command = &cobra.Command{
			Use:   "delete [user]",
			Short: "delete specify user",
			Long:  "delete specify user ... wordless",
			Args:  cobra.ExactArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				if len(args) != 0 {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return u.getAllUser(), cobra.ShellCompDirectiveNoFileComp
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				res, err := u.deleteUser(args[0])
				if err == nil {
					fmt.Fprintln(out, res)
				}
				return err
			},
		}
	)
	return command
}

type User struct {
	Client         *elasticsearch.Client
	Body           *createUserBody
	In             io.ReadWriter
	Out            io.Writer
	Username       string
	Fd             int
	Roles          string
	Email          string
	FullName       string
	AddOnly        bool
	ChangePassword bool
}

func (u *User) getAllUser() []string {
	var (
		resMap   map[string]interface{}
		resSlice []string
	)
	ret, _ := u.Client.Security.GetUser()

	_ = json.NewDecoder(ret.Body).Decode(&resMap)
	for i := range resMap {
		resSlice = append(resSlice, i)
	}
	return resSlice
}

func (u *User) createUser(req *es.RequestBody) (*esapi.Response, error) {
	rawBody, err := es.GetRawRequestBody(req)
	if err != nil {
		return nil, err
	}
	if rawBody != nil {
		return u.Client.Security.PutUser(u.Username, bytes.NewReader(rawBody))
	}
	exists, err := u.getExistRoles()
	if err != nil {
		return nil, err
	}
	pwd, err := u.getPasswordFromTerminal()
	if err != nil {
		return nil, err
	}
	body := &createUserBody{
		Password: pwd,
		Roles:    exists,
		Email:    u.Email,
		FullName: u.FullName,
		Metadata: intelligence{
			Intelligence: 7,
		},
	}
	bytesBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return u.Client.Security.PutUser(u.Username, bytes.NewReader(bytesBody))
}

func (u *User) getUser(user string) (*esapi.Response, error) {
	if user == AllUser {
		return u.Client.Security.GetUser(u.Client.Security.GetUser.WithPretty())
	}
	return u.Client.Security.GetUser(u.Client.Security.GetUser.WithUsername(splitWords(user)...), u.Client.Security.GetUser.WithPretty())
}

func (u *User) getPasswordFromTerminal() (pwd string, err error) {
	oldState, err := term.MakeRaw(u.Fd)
	if err != nil {
		return
	}
	ss := term.NewTerminal(u.In, "> ")
	pwd, err = ss.ReadPassword("password: ")
	if err != nil {
		return
	}
	confirmPwd, err := ss.ReadPassword("confirm password: ")
	if err != nil {
		return
	}
	if pwd != confirmPwd {
		err = errors.New("two input password must be consistent")
		return
	}
	if err = term.Restore(u.Fd, oldState); err != nil {
		return
	}
	return
}

type createUserBody struct {
	Password string        `json:"password,omitempty"`
	Roles    []interface{} `json:"roles,omitempty"`
	FullName string        `json:"full_name,omitempty"`
	Email    string        `json:"email,omitempty"`
	Metadata intelligence  `json:"metadata,omitempty"`
}

type intelligence struct {
	Intelligence int `json:"intelligence"`
}

func (u *User) getExistRoles() (exist []interface{}, err error) {
	if u.Roles == PassRolesCheck {
		return
	}
	var (
		rolesMap map[string]interface{}
		notExist []string
	)
	newRolesSlice := splitWords(u.Roles)
	res, err := u.Client.Security.GetRole(u.Client.Security.GetRole.WithName(newRolesSlice...))
	if err != nil {
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&rolesMap); err != nil {
		return nil, errors.Errorf("Error parsing the response body: %s", err)
	}
	for _, newRole := range newRolesSlice {
		if rolesMap[newRole] != nil {
			exist = append(exist, newRole)
		} else {
			notExist = append(notExist, newRole)
		}
	}
	if len(notExist) == len(newRolesSlice) && len(exist) == 0 {
		err = errors.Errorf("role: %s  does not exist\n", strings.Join(notExist, ","))
	} else if len(notExist) != 0 {
		fmt.Fprintf(u.Out, "role: %s  does not exist\n", strings.Join(notExist, ","))
	}
	return
}

func (u *User) updateUser(req *es.RequestBody) (*esapi.Response, error) {
	rawBody, err := es.GetRawRequestBody(req)
	if err != nil {
		return nil, err
	}
	if rawBody != nil {
		return u.Client.Security.PutUser(u.Username, bytes.NewReader(rawBody))
	}

	var body = &createUserBody{}
	user, err := u.getSpecificUser(u.Username)
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, errors.Errorf("User: %s  does not exist\n", u.Username)
	}
	exists, err := u.getExistRoles()
	if err != nil {
		return nil, err
	}
	if u.AddOnly {
		hasRoles := user[u.Username]["roles"].([]interface{})
		newRoles := mergeRoles(exists, hasRoles)
		body.Roles = newRoles
	} else {
		body.Roles = exists
	}
	body.Email = u.Email
	body.FullName = u.FullName
	if u.ChangePassword {
		pwd, err := u.getPasswordFromTerminal()
		if err != nil {
			return nil, err
		}
		body.Password = pwd
	}
	bytesBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return u.Client.Security.PutUser(u.Username, bytes.NewReader(bytesBody))
}

func mergeRoles(roles, newRoles []interface{}) []interface{} {
	mergedRoles := append(roles, newRoles...)
	mergedRoleSet := mapset.NewSetFromSlice(mergedRoles)
	mergedRoles = mergedRoleSet.ToSlice()
	return mergedRoles
}

func (u *User) getSpecificUser(user string) (userMap map[string]map[string]interface{}, err error) {
	res, err := u.Client.Security.GetUser(u.Client.Security.GetUser.WithUsername(user))
	if err != nil {
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&userMap); err != nil {
		return nil, errors.Errorf("Error parsing the response body: %s", err)
	}
	return
}

func (u *User) deleteUser(user string) (*esapi.Response, error) {
	return u.Client.Security.DeleteUser(user)
}
