package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/omnicofig/cli/pkg/profile"
	"github.com/spf13/cobra"
)

func NewProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage your OmniConfig user profile",
		Long: `Manage your OmniConfig user profile.

The profile stores your name, business, and contact info. It is used
as metadata when making config changes (for audit logging).

Commands:
  profile init  - Interactive setup
  profile show  - Display current profile
  profile set   - Update a profile field`,
	}

	cmd.AddCommand(NewProfileInitCmd())
	cmd.AddCommand(NewProfileShowCmd())
	cmd.AddCommand(NewProfileSetCmd())

	return cmd
}

func NewProfileInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Interactive profile setup",
		Long:  `Walk through an interactive setup to create your OmniConfig profile.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := profile.Load()
			if err != nil {
				return fmt.Errorf("failed to load profile: %w", err)
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Println("OmniConfig Profile Setup")
			fmt.Println("========================")
			fmt.Println()

			current := ""
			if p.Name != "" {
				current = fmt.Sprintf(" [%s]", p.Name)
			}
			fmt.Printf("Your name%s: ", current)
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			if name != "" {
				p.Name = name
			}

			if p.Business != "" {
				current = fmt.Sprintf(" [%s]", p.Business)
			} else {
				current = ""
			}
			fmt.Printf("Business/company name%s: ", current)
			business, _ := reader.ReadString('\n')
			business = strings.TrimSpace(business)
			if business != "" {
				p.Business = business
			}

			if p.Email != "" {
				current = fmt.Sprintf(" [%s]", p.Email)
			} else {
				current = ""
			}
			fmt.Printf("Email (optional)%s: ", current)
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)
			if email != "" {
				p.Email = email
			}

			if p.Phone != "" {
				current = fmt.Sprintf(" [%s]", p.Phone)
			} else {
				current = ""
			}
			fmt.Printf("Phone (optional)%s: ", current)
			phone, _ := reader.ReadString('\n')
			phone = strings.TrimSpace(phone)
			if phone != "" {
				p.Phone = phone
			}

			if err := profile.Save(p); err != nil {
				return fmt.Errorf("failed to save profile: %w", err)
			}

			fmt.Println()
			fmt.Println("Profile saved:")
			fmt.Println(p.Display())
			return nil
		},
	}
}

func NewProfileShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Display your profile",
		Long:  `Show the current OmniConfig user profile.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := profile.Load()
			if err != nil {
				return fmt.Errorf("failed to load profile: %w", err)
			}

			if !p.IsComplete() {
				fmt.Println("No profile set. Run 'omniconfig profile init' to create one.")
				return nil
			}

			fmt.Println(p.Display())
			return nil
		},
	}
}

func NewProfileSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <field> <value>",
		Short: "Update a profile field",
		Long: `Update a single profile field.

Fields: name, business, email, phone

Examples:
  omniconfig profile set name "Jane Doe"
  omniconfig profile set business "Acme Corp"`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			field := strings.ToLower(args[0])
			value := args[1]

			p, err := profile.Load()
			if err != nil {
				return fmt.Errorf("failed to load profile: %w", err)
			}

			switch field {
			case "name":
				p.Name = value
			case "business":
				p.Business = value
			case "email":
				p.Email = value
			case "phone":
				p.Phone = value
			default:
				return fmt.Errorf("unknown field: %q (use: name, business, email, phone)", field)
			}

			if err := profile.Save(p); err != nil {
				return fmt.Errorf("failed to save profile: %w", err)
			}

			fmt.Printf("Profile updated: %s = %s\n", field, value)
			return nil
		},
	}
}