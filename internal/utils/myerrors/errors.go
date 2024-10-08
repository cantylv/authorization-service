package myerrors

import "errors"

var (
	// HTTP RESPONSES
	ErrInternal             = errors.New("internal server error, please try again later")
	ErrInvalidData          = errors.New("you has passed invalid data in request data")
	ErrNoRequestIdInContext = errors.New("no request_id in request context")
	// CUSTOM
	ErrOnlyRootCanDeleteUser           = errors.New("only root user can delete user from system")
	ErrOnlyOwnerCanAddUserToGroup      = errors.New("only owner of group can add user to his group")
	ErrOnlyOwnerCanDeleteUserFromGroup = errors.New("only owner of group can delete user from his group")
	ErrOnlyOwnerCanAppointNewOwner     = errors.New("only owner can attain new owner")
	ErrOnlyRootCanBeOwnerOfUsersGroup  = errors.New("only root can be an owner of users group")
	ErrOnlyRootCanChangeBidStatus      = errors.New("only root user can approve or reject bid")
	ErrOnlyRootCanAddAgent             = errors.New("only root user can add server agent")
	ErrOnlyRootCanDeleteAgent          = errors.New("only root user can delete server agent")
	ErrOnlyRootCanGetAgents            = errors.New("only root user can get server agents")
	ErrGetUserAgents                   = errors.New("you can't get user agents")
	ErrCantDeleteRoot                  = errors.New("cant't delete root user")
	ErrUserEmailMustBeDiff             = errors.New("user emails must be different")
	ErrUserAlreadyInGroup              = errors.New("user already in group")
	ErrUserIsNotInGroup                = errors.New("user is not in group")
	ErrUserIsAlreadyOwner              = errors.New("user is already an owner")
	ErrUserIsNotOwner                  = errors.New("user is not an owner")
	ErrUserIsResponsible               = errors.New("user is responsible for group/groups, so root user need to appoint new owner")
	ErrDeleteRootFromGroup             = errors.New("user doesn't have enough rights to delete root user from group")
	// DATABASE
	ErrNoRowsAffected         = errors.New("no rows were affected")
	ErrUserNotExist           = errors.New("user is not exist")
	ErrGroupNotExist          = errors.New("group is not exist")
	ErrAgentNotExist          = errors.New("agent is not exist")
	ErrBidNotExist            = errors.New("user doesn't have bid with this name")
	ErrOwnerCantExitFromGroup = errors.New("to leave a group you need to remove the rights of the group owner")
	ErrUserAlreadyExist       = errors.New("user with this email already exist")
	ErrGroupAlreadyExist      = errors.New("group with this name already exist")
	ErrBidAlreadyExist        = errors.New("bid with this name already exist")
	ErrAgentAlreadyExist      = errors.New("agent with this name already exist")
	ErrGroupAgentAlreadyExist = errors.New("agent with this name already belongs to the selected group")
	ErrUserAgentAlreadyExist  = errors.New("agent with this name already belongs to the selected user")
	ErrGroupAgentNotExist     = errors.New("agent with this name not belongs to the selected group")
	ErrUserAgentNotExist      = errors.New("agent with this name not belongs to the selected user")
	// DTO
	ErrInvalidEmail     = errors.New("incorrect email was sent, correct format is username@domain.extension, e.g.: gref@sber.ru")
	ErrInvalidStatus    = errors.New("status must be in range(approved, rejected)")
	ErrInvalidFirstName = errors.New("incorrect first name was sent, it must start with a capital letter and be between 2 and 50 characters long")
	ErrInvalidLastName  = errors.New("incorrect last name was sent, it must start with a capital letter and be between 2 and 50 characters long")
	ErrPasswordTooLong  = errors.New("password is too long, it must be between 8 and 30 characters long")
	ErrPasswordTooShort = errors.New("password is too short, it must be between 8 and 30 characters long")
	ErrPasswordFormat   = errors.New("password must contain at least one digit and one capital letter")
)
