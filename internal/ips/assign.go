package ips

import (
	"github.com/cherryservers/cherryctl/internal/utils"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func (c *Client) Assign() *cobra.Command {
	var (
		ipID           string
		projectID      int
		targetIPID     string
		targetID       int
		targetHostname string
	)
	ipAssignCmd := &cobra.Command{
		Use:     `assign -i <ip_address_id> {--target-hostname | --target-id | --target-ip-id} [-p <project_id>]`,
		Aliases: []string{"attach"},
		Short:   "Assign an IP address to a specified server or other IP address.",
		Long:    "Assign an IP address to a specified server or another IP address. IP address assignment to another IP is possible only if routed IP type is floating and target IP is subnet or primary-ip type.",
		Example: `  # Assign an IP address to a server:
  cherryctl ip assign -i 30c15082-a06e-4c43-bfc3-252616b46eba --server-id 12345`,

		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			request := &cherrygo.AssignIPAddress{}

			if targetIPID != "" && IsValidUUID(targetIPID) {
				request.IpID = targetIPID
			} else if targetHostname != "" {
				srvID, err := utils.ServerHostnameToID(targetHostname, projectID, c.ServerService)
				if err != nil {
					return errors.Wrap(err, "Could not find a target by hostname")
				}
				request.ServerID = srvID
			} else if targetID != 0 {
				request.ServerID = targetID
			}

			if request.ServerID == 0 && request.IpID == "" {
				return errors.New("Could not find a target")
			}

			i, _, err := c.Service.Assign(ipID, request)
			if err != nil {
				return errors.Wrap(err, "Could not assign IP address")
			}

			header := []string{"ID", "Address", "Cidr", "Type", "Region"}
			data := make([][]string, 1)
			data[0] = []string{i.ID, i.Address, i.Cidr, i.Type, i.Region.Name}

			return c.Out.Output(i, header, &data)
		},
	}

	ipAssignCmd.Flags().StringVarP(&ipID, "ip-address-id", "i", "", "The ID of a IP address.")
	ipAssignCmd.Flags().IntVarP(&projectID, "project-id", "p", 0, "The project's ID.")

	ipAssignCmd.Flags().StringVarP(&targetHostname, "target-hostname", "", "", "The hostname of the server to assign IP to.")
	ipAssignCmd.Flags().IntVarP(&targetID, "target-id", "", 0, "The ID of the server to assign IP to.")
	ipAssignCmd.Flags().StringVarP(&targetIPID, "target-ip-id", "", "", "Subnet or primary-ip type IP ID to route IP to.")

	ipAssignCmd.MarkFlagsMutuallyExclusive("target-hostname", "target-id", "target-ip-id")
	ipAssignCmd.MarkFlagsRequiredTogether("target-hostname", "project-id")
	ipAssignCmd.MarkFlagRequired("ip-address-id")

	return ipAssignCmd
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
